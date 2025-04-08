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

func createBuildsCreateSubcommand(initialBody []byte) Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("builds.create", flag.ExitOnError)

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
		"config-commit",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "config_commit", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"project",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "project", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"targets",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "targets.#", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"+target",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "targets.-1", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Builds.New(
				context.TODO(),
				stainlessv0.BuildNewParams{},
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

func createBuildsRetrieveSubcommand() Subcommand {
	var buildID *string = nil
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("builds.retrieve", flag.ExitOnError)

	flagSet.Func(
		"build-id",
		"",
		func(string string) error {
			buildID = &string
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Builds.Get(
				context.TODO(),
				*buildID,
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
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", res.JSON.RawJSON())
		},
	}
}
