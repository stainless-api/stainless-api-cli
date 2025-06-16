// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var projectsBranchesCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new branch for a project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name: "branch",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "branch-from",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch_from",
			},
		},
		&jsonflag.JSONBoolFlag{
			Name: "force",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "force",
			},
		},
	},
	Action:          handleProjectsBranchesCreate,
	HideHelpCommand: true,
}

var projectsBranchesRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project branch",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&cli.StringFlag{
			Name: "branch",
		},
	},
	Action:          handleProjectsBranchesRetrieve,
	HideHelpCommand: true,
}

func handleProjectsBranchesCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainlessv0.ProjectBranchNewParams{}
	if cmd.IsSet("project") {
		params.Project = stainlessv0.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Branches.New(
		context.TODO(),
		params,
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
	cc := getAPICommandContext(cmd)
	params := stainlessv0.ProjectBranchGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainlessv0.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Branches.Get(
		context.TODO(),
		cmd.Value("branch").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
