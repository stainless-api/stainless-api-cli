// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
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
			Name:  "branch",
			Usage: "Branch name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch-from",
			Usage: "Branch or commit SHA to branch from",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch_from",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "force",
			Usage: "Whether to throw an error if the branch already exists. Defaults to false.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "force",
				SetValue: true,
			},
			Value: false,
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
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "limit",
			Usage: "Maximum number of items to return, defaults to 10 (maximum: 100).",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
			Value: 10,
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
			Name:  "base",
			Usage: `The branch or commit SHA to rebase onto. Defaults to "main".`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "base",
			},
			Value: "main",
		},
	},
	Action:          handleProjectsBranchesRebase,
	HideHelpCommand: true,
}

var projectsBranchesReset = cli.Command{
	Name:  "reset",
	Usage: "Reset a project branch.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
		&jsonflag.JSONStringFlag{
			Name:  "target-config-sha",
			Usage: "The commit SHA to reset the main branch to. Required if resetting the main branch; disallowed otherwise.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "target_config_sha",
			},
		},
	},
	Action:          handleProjectsBranchesReset,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchNewParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.New(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches create", json, format, transform)
}

func handleProjectsBranchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Get(
		ctx,
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches retrieve", json, format, transform)
}

func handleProjectsBranchesList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchListParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.List(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches list", json, format, transform)
}

func handleProjectsBranchesDelete(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchDeleteParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Delete(
		ctx,
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches delete", json, format, transform)
}

func handleProjectsBranchesRebase(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchRebaseParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Rebase(
		ctx,
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches rebase", json, format, transform)
}

func handleProjectsBranchesReset(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchResetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Branches.Reset(
		ctx,
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches reset", json, format, transform)
}
