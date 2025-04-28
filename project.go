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

var projectsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "TODO",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project-name",
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "TODO",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project-name",
		},
		&cli.StringFlag{
			Name:   "display-name",
			Action: getAPIFlagAction[string]("body", "display_name"),
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsUpdate,
	HideHelpCommand: true,
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Projects.Get(
		context.TODO(),
		cmd.Value("project-name").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Projects.Update(
		context.TODO(),
		cmd.Value("project-name").(string),
		stainlessv0.ProjectUpdateParams{},
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
