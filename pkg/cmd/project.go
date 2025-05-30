// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

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
	Usage: "Retrieve a project by name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "Update a project's properties",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
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

var projectsList = cli.Command{
	Name:  "list",
	Usage: "List projects in an organization",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "org",
			Action: getAPIFlagAction[string]("query", "org"),
		},
		&cli.StringFlag{
			Name:   "cursor",
			Action: getAPIFlagAction[string]("query", "cursor"),
		},
		&cli.FloatFlag{
			Name:   "limit",
			Action: getAPIFlagAction[float64]("query", "limit"),
		},
	},
	Before:          initAPICommand,
	Action:          handleProjectsList,
	HideHelpCommand: true,
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := stainlessv0.ProjectGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainlessv0.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Get(
		context.TODO(),
		params,
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
	params := stainlessv0.ProjectUpdateParams{}
	if cmd.IsSet("project") {
		params.Project = stainlessv0.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Update(
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

func handleProjectsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	params := stainlessv0.ProjectListParams{}
	res, err := cc.client.Projects.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
