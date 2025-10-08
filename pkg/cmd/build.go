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

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "revision",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name:  "allow-empty",
			Usage: "Whether to allow empty commits (no changes). Defaults to false.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "allow_empty",
				SetValue: true,
			},
			Value: false,
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "commit-message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
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
		&cli.StringFlag{
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
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "limit",
			Usage: "Maximum number of builds to return, defaults to 10 (maximum: 100).",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
			Value: 10,
		},
		&jsonflag.JSONStringFlag{
			Name:  "revision",
			Usage: "A config commit SHA used for the build",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "revision",
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
		&jsonflag.JSONStringFlag{
			Name:  "base.branch",
			Usage: "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "base.revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "base.commit_message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.branch",
			Usage: "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.commit_message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}

func handleBuildsCreate(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildNewParams{}
	var res []byte
	_, err := cc.client.Builds.New(
		context.TODO(),
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
	return ShowJSON("builds create", json, format, transform)
}

func handleBuildsRetrieve(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Builds.Get(
		context.TODO(),
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds retrieve", json, format, transform)
}

func handleBuildsList(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildListParams{}
	var res []byte
	_, err := cc.client.Builds.List(
		context.TODO(),
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
	return ShowJSON("builds list", json, format, transform)
}

func handleBuildsCompare(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildCompareParams{}
	var res []byte
	_, err := cc.client.Builds.Compare(
		context.TODO(),
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
	return ShowJSON("builds compare", json, format, transform)
}
