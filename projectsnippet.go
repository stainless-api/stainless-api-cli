// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var projectsSnippetsCreateRequest = cli.Command{
	Name:  "create_request",
	Usage: "Perform create_request operation",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project-name",
		},
		&cli.StringFlag{
			Name:   "language",
			Action: getAPIFlagAction[string]("body", "language"),
		},
		&cli.StringFlag{
			Name:   "request.method",
			Action: getAPIFlagAction[string]("body", "request.method"),
		},
		&cli.StringFlag{
			Name:   "request.parameters.in",
			Action: getAPIFlagAction[string]("body", "request.parameters.#.in"),
		},
		&cli.StringFlag{
			Name:   "request.parameters.name",
			Action: getAPIFlagAction[string]("body", "request.parameters.#.name"),
		},
		&cli.BoolFlag{
			Name:   "request.+parameter",
			Action: getAPIFlagActionWithValue[bool]("body", "request.parameters.-1", map[string]interface{}{}),
		},
		&cli.StringFlag{
			Name:   "request.path",
			Action: getAPIFlagAction[string]("body", "request.path"),
		},
		&cli.StringFlag{
			Name:   "request.body.fileParam",
			Action: getAPIFlagAction[string]("body", "request.body.fileParam"),
		},
		&cli.StringFlag{
			Name:   "request.body.filePath",
			Action: getAPIFlagAction[string]("body", "request.body.filePath"),
		},
		&cli.StringFlag{
			Name:   "version",
			Action: getAPIFlagAction[string]("body", "version"),
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsSnippetsCreateRequest,
	HideHelpCommand: true,
}

func handleProjectsSnippetsCreateRequest(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Projects.Snippets.NewRequest(
		context.TODO(),
		cmd.Value("project-name").(string),
		stainlessv0.ProjectSnippetNewRequestParams{},
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
