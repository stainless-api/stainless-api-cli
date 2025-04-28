// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "stainless-api-cli",
		Usage: "CLI for the stainless-v0 API",
		Commands: []*cli.Command{
			{
				Name: "projects",
				Commands: []*cli.Command{
					&projectsRetrieve,
					&projectsUpdate,
				},
			},

			{
				Name: "projects:branches",
				Commands: []*cli.Command{
					&projectsBranchesCreate,
					&projectsBranchesRetrieve,
				},
			},

			{
				Name: "projects:configs",
				Commands: []*cli.Command{
					&projectsConfigsRetrieve,
					&projectsConfigsGuess,
				},
			},

			{
				Name: "builds",
				Commands: []*cli.Command{
					&buildsCreate,
					&buildsRetrieve,
					&buildsList,
				},
			},

			{
				Name: "build_target_outputs",
				Commands: []*cli.Command{
					&buildTargetOutputsRetrieve,
				},
			},
		},
		EnableShellCompletion: true,
		HideHelpCommand:       true,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
