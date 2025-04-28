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

var projectsBranchesCreate = cli.Command{
	Name:  "create",
	Usage: "TODO",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name:   "branch",
			Action: getAPIFlagAction[string]("body", "branch"),
		},
		&cli.StringFlag{
			Name:   "branch-from",
			Action: getAPIFlagAction[string]("body", "branch_from"),
		},
		&cli.BoolFlag{
			Name:   "force",
			Action: getAPIFlagAction[bool]("body", "force"),
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsBranchesCreate,
	HideHelpCommand: true,
}

var projectsBranchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "TODO",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsBranchesRetrieve,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Projects.Branches.New(
		context.TODO(),
		cmd.Value("project").(string),
		stainlessv0.ProjectBranchNewParams{},
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithRequestBody("application/json", cc.body),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleProjectsBranchesRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.Projects.Branches.Get(
		context.TODO(),
		cmd.Value("project").(string),
		cmd.Value("branch").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
