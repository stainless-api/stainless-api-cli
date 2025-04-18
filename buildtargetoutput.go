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

func createBuildTargetOutputsRetrieveSubcommand() Subcommand {
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("build_target_outputs.retrieve", flag.ExitOnError)

	flagSet.Func(
		"build-id",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "build_id", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"target",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "target", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"type",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "type", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.Func(
		"output",
		"",
		func(string string) error {
			var jsonErr error
			query, jsonErr = jsonSet(query, "output", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.BuildTargetOutputs.Get(
				context.TODO(),
				stainlessv0.BuildTargetOutputGetParams{},
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
