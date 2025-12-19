// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

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
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
			Name:     "branch",
			Usage:    "Branch name",
			BodyPath: "branch",
		},
		&requestflag.Flag[string]{
			Name:     "branch-from",
			Usage:    "Branch or commit SHA to branch from",
			BodyPath: "branch_from",
		},
		&requestflag.Flag[bool]{
			Name:     "force",
			Usage:    "Whether to throw an error if the branch already exists. Defaults to false.",
			BodyPath: "force",
		},
	},
	Action:          handleProjectsBranchesCreate,
	HideHelpCommand: true,
}

var projectsBranchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project branch by name.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
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
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from a previous response",
			QueryPath: "cursor",
		},
		&requestflag.Flag[float64]{
			Name:      "limit",
			Usage:     "Maximum number of items to return, defaults to 10 (maximum: 100).",
			Default:   10,
			QueryPath: "limit",
		},
	},
	Action:          handleProjectsBranchesList,
	HideHelpCommand: true,
}

var projectsBranchesDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete a project branch by name.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
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
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
			Name: "branch",
		},
		&requestflag.Flag[string]{
			Name:      "base",
			Usage:     `The branch or commit SHA to rebase onto. Defaults to "main".`,
			Default:   "main",
			QueryPath: "base",
		},
	},
	Action:          handleProjectsBranchesRebase,
	HideHelpCommand: true,
}

var projectsBranchesReset = cli.Command{
	Name:  "reset",
	Usage: "Reset a project branch.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
			Name: "branch",
		},
		&requestflag.Flag[string]{
			Name:      "target-config-sha",
			Usage:     "The commit SHA to reset the main branch to. Required if resetting the main branch; disallowed otherwise.",
			QueryPath: "target_config_sha",
		},
	},
	Action:          handleProjectsBranchesReset,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectBranchNewParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:branches create", obj, format, transform)
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

	params := stainless.ProjectBranchGetParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Get(
		ctx,
		cmd.Value("branch").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:branches retrieve", obj, format, transform)
}

func handleProjectsBranchesList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectBranchListParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Projects.Branches.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "projects:branches list", obj, format, transform)
	} else {
		iter := client.Projects.Branches.ListAutoPaging(ctx, params, options...)
		return streamOutput("projects:branches list", func(w *os.File) error {
			for iter.Next() {
				item := iter.Current()
				obj := gjson.Parse(item.RawJSON())
				if err := ShowJSON(w, "projects:branches list", obj, format, transform); err != nil {
					return err
				}
			}
			return iter.Err()
		})
	}
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

	params := stainless.ProjectBranchDeleteParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Delete(
		ctx,
		cmd.Value("branch").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:branches delete", obj, format, transform)
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

	params := stainless.ProjectBranchRebaseParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Rebase(
		ctx,
		cmd.Value("branch").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:branches rebase", obj, format, transform)
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

	params := stainless.ProjectBranchResetParams{
		Project: stainless.String(cmd.Value("project").(string)),
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Branches.Reset(
		ctx,
		cmd.Value("branch").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:branches reset", obj, format, transform)
}
