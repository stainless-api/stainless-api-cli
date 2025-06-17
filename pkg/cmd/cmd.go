// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"github.com/urfave/cli/v3"
)

var Command = cli.Command{
	Name:  "stl",
	Usage: "CLI for the stainless API",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	},
	Commands: []*cli.Command{
		{
			Name: "projects",
			Commands: []*cli.Command{
				&projectsCreate,
				&projectsRetrieve,
				&projectsUpdate,
				&projectsList,
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
				&buildsCompare,
			},
		},

		{
			Name: "builds:target_outputs",
			Commands: []*cli.Command{
				&buildsTargetOutputsRetrieve,
			},
		},

		{
			Name: "orgs",
			Commands: []*cli.Command{
				&orgsRetrieve,
				&orgsList,
			},
		},

		{
			Name: "generate",
			Commands: []*cli.Command{
				&generateCreateSpec,
			},
		},
	},
	EnableShellCompletion: true,
	HideHelpCommand:       true,
}
