package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/tidwall/sjson"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

func getDefaultRequestOptions() []option.RequestOption {
	return []option.RequestOption{
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Runtime", "cli"),
	}
}

type apiCommandContext struct {
	client stainlessv0.Client
	body   []byte
	query  []byte
	header []byte
}

func (c apiCommandContext) AsMiddleware() option.Middleware {
	return func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		q := r.URL.Query()
		for key, values := range serializeQuery(c.query) {
			for _, value := range values {
				q.Set(key, value)
			}
		}
		r.URL.RawQuery = q.Encode()

		for key, values := range serializeHeader(c.header) {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}

		// Handle request body merging if there's a body to process
		if r.Body != nil && len(c.body) > 2 { // More than just "{}"
			// Read the existing request body
			existingBody, err := io.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read existing request body: %v", err)
			}

			// Start with existing body as base (default from API params)
			mergedBody := existingBody
			if len(existingBody) == 0 {
				mergedBody = []byte("{}")
			}

			// Parse command body and merge top-level keys
			commandResult := gjson.ParseBytes(c.body)
			if commandResult.IsObject() {
				commandResult.ForEach(func(key, value gjson.Result) bool {
					// Set each top-level key from command body, overwriting existing values
					var err error
					mergedBody, err = sjson.SetBytes(mergedBody, key.String(), value.Value())
					if err != nil {
						// Continue on error to merge as much as possible
						return true
					}
					return true
				})
			}

			// Set the new body
			r.Body = io.NopCloser(strings.NewReader(string(mergedBody)))
			r.ContentLength = int64(len(mergedBody))
			r.Header.Set("Content-Type", "application/json")
		}

		return mn(r)
	}
}

func getAPICommandContext(cmd *cli.Command) *apiCommandContext {
	client := stainlessv0.NewClient(getDefaultRequestOptions()...)
	body := getStdInput()
	if body == nil {
		body = []byte("{}")
	}
	var query = []byte("{}")
	var header = []byte("{}")

	// Apply JSON flag mutations
	body, query, header, err := jsonflag.Apply(body, query, header)
	if err != nil {
		log.Fatal(err)
	}

	return &apiCommandContext{client, body, query, header}
}

func serializeQuery(params []byte) url.Values {
	serialized := url.Values{}

	var serialize func(value gjson.Result, path string)
	serialize = func(res gjson.Result, path string) {
		if res.IsObject() {
			for key, value := range res.Map() {
				newPath := path
				if len(newPath) == 0 {
					newPath += key
				} else {
					newPath = "[" + key + "]"
				}

				serialize(value, newPath)
			}
		} else if res.IsArray() {
			for _, value := range res.Array() {
				serialize(value, path)
			}
		} else {
			serialized.Add(path, res.String())
		}
	}
	serialize(gjson.GetBytes(params, "@this"), "")

	for key, values := range serialized {
		serialized.Set(key, strings.Join(values, ","))
	}

	return serialized
}

func serializeHeader(params []byte) http.Header {
	serialized := http.Header{}

	var serialize func(value gjson.Result, path string)
	serialize = func(res gjson.Result, path string) {
		if res.IsObject() {
			for key, value := range res.Map() {
				newPath := path
				if len(newPath) > 0 {
					newPath += "."
				}
				newPath += key

				serialize(value, newPath)
			}
		} else if res.IsArray() {
			for _, value := range res.Array() {
				serialize(value, path)
			}
		} else {
			serialized.Add(path, res.String())
		}
	}
	serialize(gjson.GetBytes(params, "@this"), "")

	return serialized
}

func getStdInput() []byte {
	if !isInputPiped() {
		return nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return data
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

	if isTerminal(w) {
		return true
	}
	return false
}

func colorizeJSON(input string, w io.Writer) string {
	if !shouldUseColors(w) {
		return input
	}
	return string(pretty.Color(pretty.Pretty([]byte(input)), nil))
}

// GetProjectName returns the project name from the command line flag or workspace config
func GetProjectName(cmd *cli.Command, flagName string) string {
	// First check if the flag was provided
	projectName := cmd.String(flagName)
	if projectName != "" {
		return projectName
	}

	// Otherwise, try to get from workspace config
	configProjectName := GetProjectNameFromConfig()
	if configProjectName != "" {
		// Log that we're using the workspace config if in interactive mode
		if isTerminal(os.Stdout) {
			fmt.Printf("%s %s\n", au.BrightBlue("i"), fmt.Sprintf("Using project '%s' from workspace config", configProjectName))
		}
	}

	return configProjectName
}

// CheckInteractiveAndInitWorkspace checks if running in interactive mode and prompts to init workspace if needed
func CheckInteractiveAndInitWorkspace(cmd *cli.Command, projectName string) {
	// Only run in interactive mode with a terminal
	if !isTerminal(os.Stdout) {
		return
	}

	// Check if workspace config exists
	config, _, _ := FindWorkspaceConfig()
	if config != nil {
		return
	}

	// Prompt user to initialize workspace
	var answer string
	fmt.Printf("%s %s", au.BrightYellow("?"), fmt.Sprintf("Would you like to initialize a workspace config with project '%s'? [y/N] ", projectName))
	fmt.Scanln(&answer)

	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		if err := InitWorkspaceConfig(projectName, "", ""); err != nil {
			fmt.Printf("%s %s\n", au.BrightRed("✱"), fmt.Sprintf("Failed to initialize workspace: %v", err))
			return
		}
		fmt.Printf("%s %s\n", au.BrightGreen("✱"), fmt.Sprintf("Workspace initialized with project: %s", projectName))
	}
}
