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

func createProjectsBranchesCreateSubcommand(initialBody []byte) Subcommand {
	var project *string = nil
	query := []byte("{}")
	header := []byte("{}")
	body := initialBody
	var flagSet = flag.NewFlagSet("projects.branches.create", flag.ExitOnError)

	flagSet.Func(
		"project",
		"",
		func(string string) error {
			project = &string
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
		"branch-from",
		"",
		func(string string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "branch_from", string)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	flagSet.BoolFunc(
		"force",
		"",
		func(_ string) error {
			var jsonErr error
			body, jsonErr = jsonSet(body, "force", true)
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Projects.Branches.New(
				context.TODO(),
				*project,
				stainlessv0.ProjectBranchNewParams{},
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

func createProjectsBranchesRetrieveSubcommand() Subcommand {
	var project *string = nil
	var branch *string = nil
	query := []byte("{}")
	header := []byte("{}")
	var flagSet = flag.NewFlagSet("projects.branches.retrieve", flag.ExitOnError)

	flagSet.Func(
		"project",
		"",
		func(string string) error {
			project = &string
			return nil
		},
	)

	flagSet.Func(
		"branch",
		"",
		func(string string) error {
			branch = &string
			return nil
		},
	)

	return Subcommand{
		flagSet: flagSet,
		handle: func(client *stainlessv0.Client) {
			res, err := client.Projects.Branches.Get(
				context.TODO(),
				*project,
				*branch,
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
