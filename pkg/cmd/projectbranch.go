// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var projectsBranchesCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new branch for a project.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name: "branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "branch-from",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch_from",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "force",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "force",
				SetValue: true,
			},
		},
	},
	Action:          handleProjectsBranchesCreate,
	HideHelpCommand: true,
}

var projectsBranchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project branch by name.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
	},
	Action:          handleProjectsBranchesRetrieve,
	HideHelpCommand: true,
}

var projectsBranchesList = cli.Command{
	Name:  "list",
	Usage: "Retrieve a project branch by name.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
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
	},
	Action:          handleProjectsBranchesList,
	HideHelpCommand: true,
}

var projectsBranchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a project branch by name.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
	},
	Action:          handleProjectsBranchesDelete,
	HideHelpCommand: true,
}

var projectsBranchesRebase = cli.Command{
	Name:  "rebase",
	Usage: "Rebase a project branch.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
		&jsonflag.JSONStringFlag{
			Name: "base",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "base",
			},
		},
	},
	Action:          handleProjectsBranchesRebase,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectBranchNewParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects:branches create", string(res), format)
}

func handleProjectsBranchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectBranchGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Get(
		context.TODO(),
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects:branches retrieve", string(res), format)
}

func handleProjectsBranchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectBranchListParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects:branches list", string(res), format)
}

func handleProjectsBranchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectBranchDeleteParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Delete(
		context.TODO(),
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects:branches delete", string(res), format)
}

func handleProjectsBranchesRebase(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectBranchRebaseParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Rebase(
		context.TODO(),
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("projects:branches rebase", string(res), format)
}
