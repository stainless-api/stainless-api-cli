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

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:     "revision",
			Usage:    "Specifies what to build: a branch name, commit SHA, merge command\n(\"base..head\"), or file contents.",
			BodyPath: "revision",
		},
		&requestflag.Flag[bool]{
			Name:     "allow-empty",
			Usage:    "Whether to allow empty commits (no changes). Defaults to false.",
			BodyPath: "allow_empty",
		},
		&requestflag.Flag[string]{
			Name:     "branch",
			Usage:    "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
			BodyPath: "branch",
		},
		&requestflag.Flag[string]{
			Name:     "commit-message",
			Usage:    "Optional commit message to use when creating a new commit.",
			BodyPath: "commit_message",
		},
		&requestflag.Flag[any]{
			Name:     "target-commit-messages",
			Usage:    "Optional commit messages to use for each SDK when making a new commit.\nSDKs not represented in this object will fallback to the optional\n`commit_message` parameter, or will fallback further to the default\ncommit message.",
			BodyPath: "target_commit_messages",
		},
		&requestflag.Flag[[]string]{
			Name:     "target",
			Usage:    "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			BodyPath: "targets",
		},
	},
	Action:          handleBuildsCreate,
	HideHelpCommand: true,
}

var buildsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a build by its ID.",
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
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
		&requestflag.Flag[string]{
			Name:      "branch",
			Usage:     "Branch name",
			QueryPath: "branch",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from a previous response.",
			QueryPath: "cursor",
		},
		&requestflag.Flag[float64]{
			Name:      "limit",
			Usage:     "Maximum number of builds to return, defaults to 10 (maximum: 100).",
			Default:   10,
			QueryPath: "limit",
		},
		&requestflag.Flag[any]{
			Name:      "revision",
			Usage:     "A config commit SHA used for the build",
			Default:   map[string]any{},
			QueryPath: "revision",
		},
	},
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Create two builds whose outputs can be directly compared with each other.",
	Flags: []cli.Flag{
		&requestflag.Flag[any]{
			Name:     "base",
			Usage:    "Parameters for the base build",
			BodyPath: "base",
		},
		&requestflag.Flag[any]{
			Name:     "head",
			Usage:    "Parameters for the head build",
			BodyPath: "head",
		},
		&requestflag.Flag[[]string]{
			Name:     "target",
			Usage:    "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			BodyPath: "targets",
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
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Builds.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "builds create", obj, format, transform)
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
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Builds.Get(ctx, cmd.Value("build-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "builds retrieve", obj, format, transform)
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
		_, err = client.Builds.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "builds list", obj, format, transform)
	} else {
		iter := client.Builds.ListAutoPaging(ctx, params, options...)
		return streamOutput("builds list", func(w *os.File) error {
			for iter.Next() {
				item := iter.Current()
				obj := gjson.Parse(item.RawJSON())
				if err := ShowJSON(w, "builds list", obj, format, transform); err != nil {
					return err
				}
			}
			return iter.Err()
		})
	}
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
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Builds.Compare(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "builds compare", obj, format, transform)
}
