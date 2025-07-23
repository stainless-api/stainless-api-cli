// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"github.com/urfave/cli/v3"
)

var Command = cli.Command{
	Name:    "stl",
	Usage:   "CLI for the stainless API",
	Version: Version,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
		&cli.StringFlag{
			Name:  "base-url",
			Usage: "Override the base URL for API requests",
		},
	},
	Commands: []*cli.Command{
		{
			Name: "auth",
			Commands: []*cli.Command{
				&authLogin,
				&authLogout,
				&authStatus,
			},
		},

		{
			Name: "workspace",
			Commands: []*cli.Command{
				&workspaceInit,
				&workspaceStatus,
			},
		},

		{
			Name:     "projects",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&projectsCreate,
				&projectsRetrieve,
				&projectsUpdate,
				&projectsList,
			},
		},

		{
			Name:     "projects:branches",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&projectsBranchesCreate,
				&projectsBranchesRetrieve,
				&projectsBranchesList,
				&projectsBranchesDelete,
			},
		},

		{
			Name:     "projects:configs",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&projectsConfigsRetrieve,
				&projectsConfigsGuess,
			},
		},

		{
			Name:     "builds",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&buildsCreate,
				&buildsRetrieve,
				&buildsList,
				&buildsCompare,
			},
		},

		{
			Name:     "builds:diagnostics",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&buildsDiagnosticsList,
			},
		},

		{
			Name:     "builds:target_outputs",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&buildsTargetOutputsRetrieve,
			},
		},

		{
			Name:     "orgs",
			Category: "API RESOURCE",
			Commands: []*cli.Command{
				&orgsRetrieve,
				&orgsList,
			},
		},

		&mcpCommand,
		&devCommand,
	},
	EnableShellCompletion: true,
	HideHelpCommand:       true,
}
