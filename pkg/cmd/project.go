// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var projectsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new project.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "display-name",
			Usage: "Human-readable project name",
		},
		&cli.StringFlag{
			Name:  "org",
			Usage: "Organization name",
		},
		&cli.StringFlag{
			Name:  "slug",
			Usage: "Project name/slug",
		},
		&cli.StringSliceFlag{
			Name:  "target",
			Usage: "Targets to generate for",
		},
	},
	Action:          handleProjectsCreate,
	HideHelpCommand: true,
}

var projectsRetrieve = cli.Command{
	Name:            "retrieve",
	Usage:           "Retrieve a project by name.",
	Flags:           []cli.Flag{},
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "Update a project's properties.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "display-name",
		},
	},
	Action:          handleProjectsUpdate,
	HideHelpCommand: true,
}

var projectsList = cli.Command{
	Name:  "list",
	Usage: "List projects in an organization, from oldest to newest.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
		},
		&cli.Float64Flag{
			Name:  "limit",
			Usage: "Maximum number of projects to return, defaults to 10 (maximum: 100).",
			Value: 10,
		},
		&cli.StringFlag{
			Name: "org",
		},
	},
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
	if err := unmarshalStdinWithFlags(cmd, map[string]string{
		"display-name": "display_name",
		"org":          "org",
		"slug":         "slug",
		"targets":      "targets",
	}, &params); err != nil {
		return err
	}
	var res []byte
	_, err := client.Projects.New(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects create", json, format, transform)
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectGetParams{}
	var res []byte
	_, err := client.Projects.Get(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects retrieve", json, format, transform)
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectUpdateParams{}
	if err := unmarshalStdinWithFlags(cmd, map[string]string{
		"display-name": "display_name",
	}, &params); err != nil {
		return err
	}
	var res []byte
	_, err := client.Projects.Update(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects update", json, format, transform)
}

func handleProjectsList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectListParams{
		Cursor: stainless.String(cmd.Value("cursor").(string)),
		Org:    stainless.String(cmd.Value("org").(string)),
	}
	if cmd.IsSet("limit") {
		params.Limit = stainless.Opt(cmd.Value("limit").(float64))
	}
	var res []byte
	_, err := client.Projects.List(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects list", json, format, transform)
}
