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

var projectsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new project.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "display-name",
			Usage: "Human-readable project name",
			Config: requestflag.RequestConfig{
				BodyPath: "display_name",
			},
		},
		&requestflag.StringFlag{
			Name:  "org",
			Usage: "Organization name",
			Config: requestflag.RequestConfig{
				BodyPath: "org",
			},
		},
		&requestflag.YAMLFlag{
			Name:  "revision",
			Usage: "File contents to commit",
			Config: requestflag.RequestConfig{
				BodyPath: "revision",
			},
		},
		&requestflag.StringFlag{
			Name:  "slug",
			Usage: "Project name/slug",
			Config: requestflag.RequestConfig{
				BodyPath: "slug",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "target",
			Usage: "Targets to generate for",
			Config: requestflag.RequestConfig{
				BodyPath: "targets",
			},
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
		&requestflag.StringFlag{
			Name: "display-name",
			Config: requestflag.RequestConfig{
				BodyPath: "display_name",
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
		&requestflag.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
			Config: requestflag.RequestConfig{
				QueryPath: "cursor",
			},
		},
		&requestflag.FloatFlag{
			Name:  "limit",
			Usage: "Maximum number of projects to return, defaults to 10 (maximum: 100).",
			Value: requestflag.Value[float64](10),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.StringFlag{
			Name: "org",
			Config: requestflag.RequestConfig{
				QueryPath: "org",
			},
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
	_, err = client.Projects.New(
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
	return ShowJSON("projects create", json, format, transform)
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectGetParams{}

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
	_, err = client.Projects.Get(
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
	return ShowJSON("projects retrieve", json, format, transform)
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectUpdateParams{}

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
	_, err = client.Projects.Update(
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
	return ShowJSON("projects update", json, format, transform)
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
		ApplicationJSON,
	)
	if err != nil {
		return err
	}
	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.List(
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
	return ShowJSON("projects list", json, format, transform)
}
