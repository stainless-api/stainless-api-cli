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

func createProjectsSnippetsCreateRequestSubcommand(initialBody []byte) Subcommand {
	var projectName *string = nil
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("projects.snippets.create_request", flag.ExitOnError)

	flagSet.Func(
		"project-name",
		"",
		func(string string) error {
			projectName = &string
			return nil
		},
	)

	flagSet.Func(
		"language",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "language", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"request.method",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "request.method", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"request.parameters.in",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "request.parameters.#.in", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"request.parameters.name",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "request.parameters.#.name", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"request.+parameter",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "request.parameters.-1", map[string]interface{}{})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"request.path",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "request.path", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"version",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "version", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Projects.Snippets.NewRequest(
				context.TODO(),
				*projectName,
				stainlessv0.ProjectSnippetNewRequestParams{},
				option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
					q := r.URL.Query()
					for key, values := range serializeQuery(query) {
						for _, value := range values {
							q.Add(key, value)
						}
					}
					r.URL.RawQuery = q.Encode()

					for key, values := range serializeHeader(header) {
						for _, value := range values {
							r.Header.Add(key, value)
						}
					}

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
