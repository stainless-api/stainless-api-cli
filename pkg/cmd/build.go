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

var buildsCreate = requestflag.WithInnerFlags(cli.Command{
	Name:    "create",
	Usage:   "Create a build, on top of a project branch, against a given input revision.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "project",
			Usage:    "Project name",
			Required: true,
			BodyPath: "project",
		},
		&requestflag.Flag[any]{
			Name:     "revision",
			Usage:    "Specifies what to build: a branch name, commit SHA, merge command\n(\"base..head\"), or file contents.",
			Required: true,
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
		&requestflag.Flag[map[string]any]{
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
}, map[string][]requestflag.HasOuterFlag{
	"target-commit-messages": {
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.cli",
			InnerField: "cli",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.csharp",
			InnerField: "csharp",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.go",
			InnerField: "go",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.java",
			InnerField: "java",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.kotlin",
			InnerField: "kotlin",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.node",
			InnerField: "node",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.openapi",
			InnerField: "openapi",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.php",
			InnerField: "php",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.python",
			InnerField: "python",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.ruby",
			InnerField: "ruby",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.sql",
			InnerField: "sql",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.terraform",
			InnerField: "terraform",
		},
		&requestflag.InnerFlag[string]{
			Name:       "target-commit-messages.typescript",
			InnerField: "typescript",
		},
	},
})

var buildsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Retrieve a build by its ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "build-id",
			Usage:    "Build ID",
			Required: true,
		},
	},
	Action:          handleBuildsRetrieve,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:    "list",
	Usage:   "List user-triggered builds for a given project.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "project",
			Usage:     "Project name",
			Required:  true,
			QueryPath: "project",
		},
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

var buildsCompare = requestflag.WithInnerFlags(cli.Command{
	Name:    "compare",
	Usage:   "Create two builds whose outputs can be directly compared with each other.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[map[string]any]{
			Name:     "base",
			Usage:    "Parameters for the base build",
			Required: true,
			BodyPath: "base",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "head",
			Usage:    "Parameters for the head build",
			Required: true,
			BodyPath: "head",
		},
		&requestflag.Flag[string]{
			Name:     "project",
			Usage:    "Project name",
			Required: true,
			BodyPath: "project",
		},
		&requestflag.Flag[[]string]{
			Name:     "target",
			Usage:    "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			BodyPath: "targets",
		},
	},
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"base": {
		&requestflag.InnerFlag[string]{
			Name:       "base.branch",
			Usage:      "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			InnerField: "branch",
		},
		&requestflag.InnerFlag[any]{
			Name:       "base.revision",
			Usage:      "Specifies what to build: a branch name, a commit SHA, or file contents.",
			InnerField: "revision",
		},
		&requestflag.InnerFlag[string]{
			Name:       "base.commit-message",
			Usage:      "Optional commit message to use when creating a new commit.",
			InnerField: "commit_message",
		},
	},
	"head": {
		&requestflag.InnerFlag[string]{
			Name:       "head.branch",
			Usage:      "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			InnerField: "branch",
		},
		&requestflag.InnerFlag[any]{
			Name:       "head.revision",
			Usage:      "Specifies what to build: a branch name, a commit SHA, or file contents.",
			InnerField: "revision",
		},
		&requestflag.InnerFlag[string]{
			Name:       "head.commit-message",
			Usage:      "Optional commit message to use when creating a new commit.",
			InnerField: "commit_message",
		},
	},
})

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
		return ShowJSONIterator(os.Stdout, "builds list", iter, format, transform)
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
