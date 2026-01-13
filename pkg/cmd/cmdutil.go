package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/stainless-api/stainless-api-cli/internal/jsonview"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/itchyny/json2yaml"
	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

func getDefaultRequestOptions(cmd *cli.Command) []option.RequestOption {
	opts := []option.RequestOption{
		option.WithHeader("User-Agent", fmt.Sprintf("Stainless/CLI %s", Version)),
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Package-Version", Version),
		option.WithHeader("X-Stainless-Runtime", "cli"),
		option.WithHeader("X-Stainless-CLI-Command", cmd.FullName()),
	}

	// Override base URL if the --base-url flag is provided
	if baseURL := cmd.String("base-url"); baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// Set environment if the --environment flag is provided
	if environment := cmd.String("environment"); environment != "" {
		switch environment {
		case "production":
			opts = append(opts, option.WithEnvironmentProduction())
		case "staging":
			opts = append(opts, option.WithEnvironmentStaging())
		default:
			log.Fatalf("Unknown environment: %s. Valid environments are %s", environment, "production, staging")
		}
	}

	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey == "" {
		config := &AuthConfig{}
		if found, err := config.Find(); err == nil && found && config.AccessToken != "" {
			opts = append(opts, option.WithAPIKey(config.AccessToken))
		}
	}

	if project := os.Getenv("STAINLESS_PROJECT"); project == "" {
		workspaceConfig := WorkspaceConfig{}
		found, err := workspaceConfig.Find()
		if err == nil && found && workspaceConfig.Project != "" {
			cmd.Set("project", workspaceConfig.Project)
		}
	}

	return opts
}

var debugMiddlewareOption = option.WithMiddleware(
	func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		logger := log.Default()

		if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(r)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	},
)

// convertFileFlag reads a file from a flag and mutates the flag's contents to have the file contents rather
// than the file values.
func convertFileFlag(cmd *cli.Command, flagName string) (string, []byte, error) {
	filePath := cmd.String(flagName)
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return path.Base(filePath), nil, fmt.Errorf("failed to read %s file: %v", flagName, err)
		}
		return path.Base(filePath), content, nil
	}
	return "", nil, nil
}

func isInputPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

var au *aurora.Aurora

func init() {
	au = aurora.New(aurora.WithColors(shouldUseColors(os.Stdout)))
}

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	// For non-tty output (probably a pipe), write directly to stdout
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// When streaming output on Unix-like systems, there's a special trick involving creating two socket pairs
	// that we prefer because it supports small buffer sizes which results in less pagination per buffer. The
	// constructs needed to run it don't exist on Windows builds, so we have this function broken up into
	// OS-specific files with conditional build comments. Under Windows (and in case our fancy constructs fail
	// on Unix), we fall back to using pipes (`streamToPagerWithPipe`), which are OS agnostic.
	//
	// Defined in either cmdutil_unix.go or cmdutil_windows.go.
	return streamOutputOSSpecific(label, generateOutput)
}

func streamToPagerWithPipe(label string, generateOutput func(w *os.File) error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	if _, err := exec.LookPath(pagerProgram); err != nil {
		return err
	}

	cmd := exec.Command(pagerProgram)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"LESS=-r -P "+label,
		"MORE=-r -P "+label,
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	// If we would be streaming to a terminal and aren't forcing color one way
	// or the other, we should configure things to use color so the pager gets
	// colorized input.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(w); err != nil && !strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	w.Close()
	return cmd.Wait()
}

func streamToStdout(generateOutput func(w *os.File) error) error {
	signal.Ignore(syscall.SIGPIPE)
	err := generateOutput(os.Stdout)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		return nil
	}
	return err
}

func shouldUseColors(w io.Writer) bool {
	// Check if NO_COLOR environment variable is set
	if _, noColor := os.LookupEnv("NO_COLOR"); noColor {
		return false
	}

	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

// Display JSON to the user in various different formats
func ShowJSON(out *os.File, title string, res gjson.Result, format string, transform string) error {
	if format != "raw" && transform != "" {
		transformed := res.Get(transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	switch strings.ToLower(format) {
	case "auto":
		return ShowJSON(out, title, res, "json", "")
	case "explore":
		return jsonview.ExploreJSON(title, res)
	case "pretty":
		_, err := out.WriteString(jsonview.RenderJSON(title, res) + "\n")
		return err
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(out) {
			_, err := out.Write(pretty.Color(prettyJSON, pretty.TerminalStyle))
			return err
		} else {
			_, err := out.Write(prettyJSON)
			return err
		}
	case "jsonl":
		// @ugly is gjson syntax for "no whitespace", so it fits on one line
		oneLineJSON := res.Get("@ugly").Raw
		if shouldUseColors(out) {
			bytes := append(pretty.Color([]byte(oneLineJSON), pretty.TerminalStyle), '\n')
			_, err := out.Write(bytes)
			return err
		} else {
			_, err := out.Write([]byte(oneLineJSON + "\n"))
			return err
		}
	case "raw":
		if _, err := out.Write([]byte(res.Raw + "\n")); err != nil {
			return err
		}
		return nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return err
		}
		_, err := out.Write([]byte(yaml.String()))
		return err
	default:
		return fmt.Errorf("Invalid format: %s, valid formats are: %s", format, strings.Join(OutputFormats, ", "))
	}
}

// For an iterator over different value types, display its values to the user in
// different formats.
func ShowJSONIterator[T any](stdout *os.File, title string, iter jsonview.Iterator[T], format string, transform string) error {
	if format == "explore" {
		return jsonview.ExploreJSONStream(title, iter)
	}
	return streamOutput(title, func(pager *os.File) error {
		for iter.Next() {
			item := iter.Current()
			jsonData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			obj := gjson.ParseBytes(jsonData)
			if err := ShowJSON(pager, title, obj, format, transform); err != nil {
				return err
			}
		}
		return iter.Err()
	})
}
