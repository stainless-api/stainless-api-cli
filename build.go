// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new build",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "project",
			Action: getAPIFlagAction[string]("body", "project"),
		},
		&cli.StringFlag{
			Name:   "revision",
			Action: getAPIFlagAction[string]("body", "revision"),
		},
		&cli.BoolFlag{
			Name:   "allow-empty",
			Action: getAPIFlagAction[bool]("body", "allow_empty"),
		},
		&cli.StringFlag{
			Name:   "branch",
			Action: getAPIFlagAction[string]("body", "branch"),
		},
		&cli.StringFlag{
			Name:   "commit-message",
			Action: getAPIFlagAction[string]("body", "commit_message"),
		},
		&cli.StringFlag{
			Name:   "targets",
			Action: getAPIFlagAction[string]("body", "targets.#"),
		},
		&cli.StringFlag{
			Name:   "+target",
			Action: getAPIFlagAction[string]("body", "targets.-1"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsCreate,
	HideHelpCommand: true,
}

var buildsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a build by ID",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "build-id",
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsRetrieve,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List builds for a project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "project",
			Action: getAPIFlagAction[string]("query", "project"),
		},
		&cli.StringFlag{
			Name:   "branch",
			Action: getAPIFlagAction[string]("query", "branch"),
		},
		&cli.StringFlag{
			Name:   "cursor",
			Action: getAPIFlagAction[string]("query", "cursor"),
		},
		&cli.FloatFlag{
			Name:   "limit",
			Action: getAPIFlagAction[float64]("query", "limit"),
		},
		&cli.StringFlag{
			Name:   "revision",
			Action: getAPIFlagAction[string]("query", "revision"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Builds.New(
		context.TODO(),
		stainlessv0.BuildNewParams{},
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Builds.Get(
		context.TODO(),
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Builds.List(
		context.TODO(),
		stainlessv0.BuildListParams{},
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
