package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strings"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonview"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/itchyny/json2yaml"
	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

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
			log.Fatalf("Unknown environment: %s. Valid environments are: production, staging", environment)
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
			opts = append(opts, option.WithProject(workspaceConfig.Project))
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

func ShowJSON(title string, res gjson.Result, format string, transform string) error {
	if format != "raw" && transform != "" {
		transformed := res.Get(transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	switch strings.ToLower(format) {
	case "auto":
		return ShowJSON(title, res, "json", "")
	case "explore":
		return jsonview.ExploreJSON(title, res)
	case "pretty":
		jsonview.DisplayJSON(title, res)
		return nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(os.Stdout) {
			fmt.Print(string(pretty.Color(prettyJSON, pretty.TerminalStyle)))
		} else {
			fmt.Print(string(prettyJSON))
		}
		return nil
	case "raw":
		fmt.Println(res.Raw)
		return nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return err
		}
		fmt.Print(yaml.String())
		return nil
	default:
		return fmt.Errorf("Invalid format: %s, valid formats are: %s", format, strings.Join(OutputFormats, ", "))
	}
}
