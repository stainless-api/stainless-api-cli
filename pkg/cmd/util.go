// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/term"

	"github.com/logrusorgru/aurora/v4"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v3"
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

type apiCommandContext struct {
	client          stainless.Client
	cmd             *cli.Command
	workspaceConfig WorkspaceConfig
}

func (c apiCommandContext) AsMiddleware() option.Middleware {
	body := getStdInput()
	if body == nil {
		body = []byte("{}")
	}
	var query = []byte("{}")
	var header = []byte("{}")

	// Apply JSON flag mutations
	body, query, header, err := jsonflag.ApplyMutations(body, query, header)
	if err != nil {
		log.Fatal(err)
	}

	debug := c.cmd.Bool("debug")

	return func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		q := r.URL.Query()
		for key, values := range serializeQuery(query) {
			for _, value := range values {
				q.Set(key, value)
			}
		}
		r.URL.RawQuery = q.Encode()

		for key, values := range serializeHeader(header) {
			for _, value := range values {
				r.Header.Set(key, value)
			}
		}

		// Handle request body merging if there's a body to process
		if r.Body != nil || len(body) > 2 { // More than just "{}"
			var existingBody []byte
			var err error

			// Read the existing request body if it exists
			if r.Body != nil {
				existingBody, err = io.ReadAll(r.Body)
				r.Body.Close()
				if err != nil {
					return nil, fmt.Errorf("failed to read existing request body: %v", err)
				}
			}

			// Start with existing body as base (default from API params)
			mergedBody := existingBody
			if len(existingBody) == 0 {
				mergedBody = []byte("{}")
			}

			// Parse command body and merge top-level keys
			commandResult := gjson.ParseBytes(body)
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
			r.Body = io.NopCloser(bytes.NewBuffer(mergedBody))
			r.ContentLength = int64(len(mergedBody))
			r.Header.Set("Content-Type", "application/json")
		}

		// Add debug logging if the --debug flag is set
		if debug {
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
		}

		return mn(r)
	}
}

func getAPICommandContext(cmd *cli.Command) *apiCommandContext {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	var workspaceConfig WorkspaceConfig
	found, _ := workspaceConfig.Find()

	if found {
		names := []string{}
		for _, flag := range cmd.VisibleFlags() {
			names = append(names, flag.Names()...)
		}

		config := workspaceConfig
		// Get the directory containing the workspace config file
		configDir := filepath.Dir(config.ConfigPath)

		if slices.Contains(names, "openapi-spec") && !cmd.IsSet("openapi-spec") && !cmd.IsSet("revision") && config.OpenAPISpec != "" {
			// Set OpenAPI spec path relative to workspace config directory
			openAPIPath := filepath.Join(configDir, config.OpenAPISpec)
			cmd.Set("openapi-spec", openAPIPath)
		}

		if slices.Contains(names, "stainless-config") && !cmd.IsSet("stainless-config") && !cmd.IsSet("revision") && config.StainlessConfig != "" {
			// Set Stainless config path relative to workspace config directory
			stainlessConfigPath := filepath.Join(configDir, config.StainlessConfig)
			cmd.Set("stainless-config", stainlessConfigPath)
		}

		if slices.Contains(names, "project") && !cmd.IsSet("project") && config.Project != "" {
			cmd.Set("project", config.Project)
		}
	}

	return &apiCommandContext{client, cmd, workspaceConfig}
}

// HasWorkspaceTargets returns true if workspace config has configured targets
func (c *apiCommandContext) HasWorkspaceTargets() bool {
	return c.workspaceConfig.ConfigPath != "" && c.workspaceConfig.Targets != nil && len(c.workspaceConfig.Targets) > 0
}

// GetWorkspaceTargetPaths returns a map of target names to their output paths from workspace config
func (c *apiCommandContext) GetWorkspaceTargetPaths() map[string]string {
	targetPaths := make(map[string]string)
	if c.workspaceConfig.ConfigPath != "" && c.workspaceConfig.Targets != nil {
		for targetName, targetConfig := range c.workspaceConfig.Targets {
			if targetConfig.OutputPath != "" {
				targetPaths[targetName] = targetConfig.OutputPath
			}
		}
	}
	return targetPaths
}

// applyFileFlag reads a file from a flag and mutates the JSON body
func applyFileFlag(cmd *cli.Command, flagName, jsonPath string) error {
	filePath := cmd.String(flagName)
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s file: %v", flagName, err)
		}
		jsonflag.Mutate(jsonflag.Body, jsonPath, string(content))
	}
	return nil
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

func ColorizeJSON(input string, w io.Writer) string {
	if !shouldUseColors(w) {
		return input
	}
	return string(pretty.Color(pretty.Pretty([]byte(input)), nil))
}
