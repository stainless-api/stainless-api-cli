// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
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
		&jsonflag.JSONStringFlag{
			Name: "language",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "language",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.method",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.method",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.parameters.in",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.parameters.#.in",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.parameters.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.parameters.#.name",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "request.+parameter",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "request.parameters.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.path",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.path",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.body.fileParam",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.body.fileParam",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.body.filePath",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.body.filePath",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.queryString.name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.queryString.#.name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.queryString.value",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.queryString.#.value",
			},
		},
		&jsonflag.JSONAnyFlag{
			Name: "request.+query_string",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "request.queryString.-1",
				SetValue: map[string]interface{}{},
			},
			Value: map[string]interface{}{},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.url",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.url",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.postData.mimeType",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.postData.mimeType",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "request.postData.text",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "request.postData.text",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "version",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "version",
			},
		},
	},
	Action:          handleProjectsSnippetsCreateRequest,
	HideHelpCommand: true,
}

func handleProjectsSnippetsCreateRequest(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainlessv0.ProjectSnippetNewRequestParams{}
	res, err := cc.client.Projects.Snippets.NewRequest(
		context.TODO(),
		cmd.Value("project-name").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
