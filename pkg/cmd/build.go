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

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
	Flags: []cli.Flag{
		&requestflag.YAMLFlag{
			Name:  "revision",
			Usage: "Specifies what to build: a branch name, commit SHA, merge command\n(\"base..head\"), or file contents.",
			Config: requestflag.RequestConfig{
				BodyPath: "revision",
			},
		},
		&requestflag.BoolFlag{
			Name:  "allow-empty",
			Usage: "Whether to allow empty commits (no changes). Defaults to false.",
			Config: requestflag.RequestConfig{
				BodyPath: "allow_empty",
			},
		},
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
			Config: requestflag.RequestConfig{
				BodyPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name:  "commit-message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: requestflag.RequestConfig{
				BodyPath: "commit_message",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: requestflag.RequestConfig{
				BodyPath: "targets",
			},
		},
	},
	Action:          handleBuildsCreate,
	HideHelpCommand: true,
}

var buildsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a build by its ID.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
		},
	},
	Action:          handleBuildsRetrieve,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List user-triggered builds for a given project.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: requestflag.RequestConfig{
				QueryPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response.",
			Config: requestflag.RequestConfig{
				QueryPath: "cursor",
			},
		},
		&requestflag.FloatFlag{
			Name:  "limit",
			Usage: "Maximum number of builds to return, defaults to 10 (maximum: 100).",
			Value: requestflag.Value[float64](10),
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.YAMLFlag{
			Name:  "revision",
			Usage: "A config commit SHA used for the build",
			Config: requestflag.RequestConfig{
				QueryPath: "revision",
			},
		},
	},
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Create two builds whose outputs can be directly compared with each other.",
	Flags: []cli.Flag{
		&requestflag.YAMLFlag{
			Name:  "base",
			Usage: "Parameters for the base build",
			Config: requestflag.RequestConfig{
				BodyPath: "base",
			},
		},
		&requestflag.YAMLFlag{
			Name:  "head",
			Usage: "Parameters for the head build",
			Config: requestflag.RequestConfig{
				BodyPath: "head",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: requestflag.RequestConfig{
				BodyPath: "targets",
			},
		},
	},
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildNewParams{}

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
	_, err = client.Builds.New(
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
	return ShowJSON("builds create", json, format, transform)
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
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
	_, err = client.Builds.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "build-id"),
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds retrieve", json, format, transform)
}

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildListParams{}

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
	_, err = client.Builds.List(
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
	return ShowJSON("builds list", json, format, transform)
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildCompareParams{}

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
	_, err = client.Builds.Compare(
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
	return ShowJSON("builds compare", json, format, transform)
}
