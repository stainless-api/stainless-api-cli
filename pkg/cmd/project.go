// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var projectsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new project.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "display-name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "display_name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "org",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "org",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "slug",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "slug",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "targets",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleProjectsCreate,
	HideHelpCommand: true,
}

var projectsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project by name.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
	},
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "Update a project's properties.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name: "display-name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "display_name",
			},
		},
	},
	Action:          handleProjectsUpdate,
	HideHelpCommand: true,
}

var projectsList = cli.Command{
	Name:  "list",
	Usage: "List projects in an organization, from oldest to newest.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "cursor",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "limit",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "org",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "org",
			},
		},
	},
	Action:          handleProjectsList,
	HideHelpCommand: true,
}

func handleProjectsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectNewParams{}
	var res []byte
	_, err := cc.client.Projects.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects create", string(res), format)
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Get(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects retrieve", string(res), format)
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectUpdateParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Update(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects update", string(res), format)
}

func handleProjectsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectListParams{}
	var res []byte
	_, err := cc.client.Projects.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects list", string(res), format)
}
