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

var projectsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new project.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "display-name",
			Usage:    "Human-readable project name",
			BodyPath: "display_name",
		},
		&requestflag.Flag[string]{
			Name:     "org",
			Usage:    "Organization name",
			BodyPath: "org",
		},
		&requestflag.Flag[any]{
			Name:     "revision",
			Usage:    "File contents to commit",
			BodyPath: "revision",
		},
		&requestflag.Flag[string]{
			Name:     "slug",
			Usage:    "Project name/slug",
			BodyPath: "slug",
		},
		&requestflag.Flag[[]string]{
			Name:     "target",
			Usage:    "Targets to generate for",
			BodyPath: "targets",
		},
	},
	Before:          before,
	Action:          handleProjectsCreate,
	HideHelpCommand: true,
}

var projectsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project by name.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name: "project",
		},
	},
	Before:          before,
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "Update a project's properties.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name: "project",
		},
		&requestflag.Flag[string]{
			Name:     "display-name",
			BodyPath: "display_name",
		},
	},
	Before:          before,
	Action:          handleProjectsUpdate,
	HideHelpCommand: true,
}

var projectsList = cli.Command{
	Name:  "list",
	Usage: "List projects in an organization, from oldest to newest.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from a previous response",
			QueryPath: "cursor",
		},
		&requestflag.Flag[float64]{
			Name:        "limit",
			Usage:       "Maximum number of projects to return, defaults to 10 (maximum: 100).",
			Default:       10,
			DefaultText: "10",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "org",
			QueryPath: "org",
		},
	},
	Before:          before,
	Action:          handleProjectsList,
	HideHelpCommand: true,
}

func handleProjectsCreate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectNewParams{}

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
	_, err = client.Projects.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects create", obj, format, transform)
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectGetParams{
		Project: stainless.Opt(cmd.Value("project").(string)),
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
	_, err = client.Projects.Get(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects retrieve", obj, format, transform)
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectUpdateParams{
		Project: stainless.Opt(cmd.Value("project").(string)),
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
	_, err = client.Projects.Update(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects update", obj, format, transform)
}

func handleProjectsList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := stainless.ProjectListParams{}

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
		_, err = client.Projects.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "projects list", obj, format, transform)
	} else {
		iter := client.Projects.ListAutoPaging(ctx, params, options...)
		return streamOutput("projects list", func(w *os.File) error {
			for iter.Next() {
				item := iter.Current()
				obj := gjson.Parse(item.RawJSON())
				if err := ShowJSON(w, "projects list", obj, format, transform); err != nil {
					return err
				}
			}
			return iter.Err()
		})
	}
}
