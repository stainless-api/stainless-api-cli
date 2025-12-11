package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonview"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/itchyny/json2yaml"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
	"golang.org/x/sys/unix"
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
			log.Fatalf("Unknown environment: %s. Valid environments are: production, staging", environment)
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

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	// For non-tty output (probably a pipe), write directly to stdout
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// Windows lacks UNIX socket APIs, so we fall back to pipes there or if
	// socket creation fails. We prefer sockets when available because they
	// allow for smaller buffer sizes, preventing unnecessary data streaming
	// from the backend. Pipes typically have large buffers but serve as a
	// decent alternative when sockets aren't available.
	if runtime.GOOS == "windows" {
		return streamToPagerWithPipe(label, generateOutput)
	}

	// Try to use socket pair for better buffer control
	pagerInput, pid, err := openSocketPairPager(label)
	if err != nil || pagerInput == nil {
		// Fall back to pipe if socket setup fails
		return streamToPagerWithPipe(label, generateOutput)
	}
	defer pagerInput.Close()

	// If we would be streaming to a terminal and aren't forcing color one way
	// or the other, we should configure things to use color so the pager gets
	// colorized input.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	// If the pager exits before reading all input, then generateOutput() will
	// produce a broken pipe error, which is fine and we don't want to propagate it.
	if err := generateOutput(pagerInput); err != nil &&
		!strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	// Close the file NOW before we wait for the child process to terminate.
	// This way, the child will receive the end-of-file signal and know that
	// there is no more input. Otherwise the child process may block
	// indefinitely waiting for another line (this can happen when streaming
	// less than a screenful of data to a pager).
	pagerInput.Close()

	// Wait for child process to exit
	var wstatus syscall.WaitStatus
	_, err = syscall.Wait4(pid, &wstatus, 0, nil)
	if wstatus.ExitStatus() != 0 {
		return fmt.Errorf("Pager exited with non-zero exit status: %d", wstatus.ExitStatus())
	}
	return err
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

func openSocketPairPager(label string) (*os.File, int, error) {
	fds, err := unix.Socketpair(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, 0, err
	}

	// The child file descriptor will be sent to the child process through
	// ProcAttr and ForkExec(), while the parent process will always close the
	// child file descriptor.
	// The parent file descriptor will be wrapped in an os.File wrapper and
	// returned from this function, or closed if something goes wrong.
	parentFd, childFd := fds[0], fds[1]
	defer unix.Close(childFd)

	// Use small buffer sizes so we don't ask the server for more paginated
	// values than we actually need.
	if err := unix.SetsockoptInt(parentFd, unix.SOL_SOCKET, unix.SO_SNDBUF, 128); err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}
	if err := unix.SetsockoptInt(childFd, unix.SOL_SOCKET, unix.SO_RCVBUF, 128); err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}

	// Set CLOEXEC on the parent file descriptor so it doesn't leak to child
	syscall.CloseOnExec(parentFd)

	parentConn := os.NewFile(uintptr(parentFd), "parent-socket")

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	pagerPath, err := exec.LookPath(pagerProgram)
	if err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}

	env := os.Environ()
	env = append(env, "LESS=-r -P "+label)
	env = append(env, "MORE=-r -P "+label)

	procAttr := &syscall.ProcAttr{
		Dir: "",
		Env: env,
		Files: []uintptr{
			uintptr(childFd),        // stdin (fd 0)
			uintptr(syscall.Stdout), // stdout (fd 1)
			uintptr(syscall.Stderr), // stderr (fd 2)
		},
	}

	pid, err := syscall.ForkExec(pagerPath, []string{pagerProgram}, procAttr)
	if err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}

	return parentConn, pid, nil
}

func shouldUseColors(w io.Writer) bool {
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
