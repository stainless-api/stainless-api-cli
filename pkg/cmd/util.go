// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"slices"
	"strings"
	"golang.org/x/term"
	"github.com/logrusorgru/aurora/v4"
	"reflect"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonview"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/itchyny/json2yaml"
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

type fileReader struct {
	workspaceConfig WorkspaceConfig
	Value         io.Reader
	Base64Encoded bool
}

func (f *fileReader) Set(filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filename, err)
	}
	f.Value = reader
	return nil
}

func (f *fileReader) String() string {
	if f.Value == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(f.Value)
	if f.Base64Encoded {
		return base64.StdEncoding.EncodeToString(buf.Bytes())
	}
	return buf.String()
}

func (f *fileReader) Get() any {
	return f.String()
}

func unmarshalWithReaders(data []byte, v any) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		jsonKey := ft.Tag.Get("json")
		if jsonKey == "" {
			jsonKey = ft.Name
		} else if idx := strings.Index(jsonKey, ","); idx != -1 {
			jsonKey = jsonKey[:idx]
		}

		rawVal, ok := fields[jsonKey]
		if !ok {
			continue
		}

		if ft.Type == reflect.TypeOf((*io.Reader)(nil)).Elem() {
			var s string
			if err := json.Unmarshal(rawVal, &s); err != nil {
				return fmt.Errorf("field %s: %w", ft.Name, err)
			}
			fv.Set(reflect.ValueOf(strings.NewReader(s)))
		} else {
			ptr := fv.Addr().Interface()
			if err := json.Unmarshal(rawVal, ptr); err != nil {
				return fmt.Errorf("field %s: %w", ft.Name, err)
			}
		}
	}

	return nil
}

func unmarshalStdinWithFlags(cmd *cli.Command, flags map[string]string, target any) error {
	var data []byte
	if isInputPiped() {
		var err error
		if data, err = io.ReadAll(os.Stdin); err != nil {
			return err
		}
	}

	// Merge CLI flags into the body
	for flag, path := range flags {
		if cmd.IsSet(flag) {
			var err error
			data, err = sjson.SetBytes(data, path, cmd.Value(flag))
			if err != nil {
				return err
			}
		}
	}

	if data != nil {
		if err := unmarshalWithReaders(data, target); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	return nil
}

func debugMiddleware(debug bool) option.Middleware {
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

		if slices.Contains(names, "openapi-spec") && !cmd.IsSet("openapi-spec") && !cmd.IsSet("revision") && config.OpenAPISpec != "" {
			// OpenAPI spec path is already absolute
			cmd.Set("openapi-spec", config.OpenAPISpec)
		}

		if slices.Contains(names, "stainless-config") && !cmd.IsSet("stainless-config") && !cmd.IsSet("revision") && config.StainlessConfig != "" {
			// Stainless config path is already absolute
			cmd.Set("stainless-config", config.StainlessConfig)
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
func (c *apiCommandContext) GetWorkspaceTargetPaths() map[stainless.Target]string {
	targetPaths := make(map[stainless.Target]string)
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
