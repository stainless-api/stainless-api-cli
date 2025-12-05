// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var projectsBranchesCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new branch for a project.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: requestflag.RequestConfig{
				BodyPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name:  "branch-from",
			Usage: "Branch or commit SHA to branch from",
			Config: requestflag.RequestConfig{
				BodyPath: "branch_from",
			},
		},
		&requestflag.BoolFlag{
			Name:  "force",
			Usage: "Whether to throw an error if the branch already exists. Defaults to false.",
			Config: requestflag.RequestConfig{
				BodyPath: "force",
			},
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesCreate,
	HideHelpCommand: true,
}

var projectsBranchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project branch by name.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name: "branch",
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesRetrieve,
	HideHelpCommand: true,
}

var projectsBranchesList = cli.Command{
	Name:  "list",
	Usage: "Retrieve a project branch by name.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
			Config: requestflag.RequestConfig{
				QueryPath: "cursor",
			},
		},
		&requestflag.FloatFlag{
			Name:        "limit",
			Usage:       "Maximum number of items to return, defaults to 10 (maximum: 100).",
			Value:       requestflag.Value[float64](10),
			DefaultText: "10",
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesList,
	HideHelpCommand: true,
}

var projectsBranchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a project branch by name.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name: "branch",
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesDelete,
	HideHelpCommand: true,
}

var projectsBranchesRebase = cli.Command{
	Name:  "rebase",
	Usage: "Rebase a project branch.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name: "branch",
		},
		&requestflag.StringFlag{
			Name:        "base",
			Usage:       `The branch or commit SHA to rebase onto. Defaults to "main".`,
			Value:       requestflag.Value[string]("main"),
			DefaultText: "main",
			Config: requestflag.RequestConfig{
				QueryPath: "base",
			},
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesRebase,
	HideHelpCommand: true,
}

var projectsBranchesReset = cli.Command{
	Name:  "reset",
	Usage: "Reset a project branch.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name: "branch",
		},
		&requestflag.StringFlag{
			Name:  "target-config-sha",
			Usage: "The commit SHA to reset the main branch to. Required if resetting the main branch; disallowed otherwise.",
			Config: requestflag.RequestConfig{
				QueryPath: "target_config_sha",
			},
		},
	},
	Before:          before,
	Action:          handleProjectsBranchesReset,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchNewParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.New(
		ctx,
		params,
		options...,
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
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchGetParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "branch"),
		params,
		options...,
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
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.List(
		ctx,
		params,
		options...,
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
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchDeleteParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Delete(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "branch"),
		params,
		options...,
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
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchRebaseParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Rebase(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "branch"),
		params,
		options...,
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
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("branch") && len(unusedArgs) > 0 {
		cmd.Set("branch", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectBranchResetParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Reset(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "branch"),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:branches reset", json, format, transform)
}
