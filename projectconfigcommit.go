// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
)

func createProjectsConfigCommitsCreateSubcommand(initialBody []byte) Subcommand {
	var projectName *string = nil
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("projects.config.commits.create", flag.ExitOnError)

	flagSet.Func(
		"project-name",
		"",
		func(string string) error {
			projectName = &string
			return nil
		},
	)

	flagSet.Func(
		"branch",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "branch", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"commit-message",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "commit_message", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"allow-empty",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "allow_empty", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"openapi-spec",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "openapi_spec", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"stainless-config",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "stainless_config", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Projects.Config.Commits.New(
				context.TODO(),
				*projectName,
				stainlessv0.ProjectConfigCommitNewParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = serializeQuery(query).Encode()
					r.Header = serializeHeader(header)
					return mn(r)
				}),
				option.WithRequestBody("application/json", body),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", res.JSON.RawJSON())
		},
	}
}
