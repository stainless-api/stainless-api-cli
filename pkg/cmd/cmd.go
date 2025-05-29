// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"github.com/urfave/cli/v3"
)

var Command = cli.Command{
	Name:  "stainless-api-cli",
	Usage: "CLI for the stainless-v0 API",
	Commands: []*cli.Command{
		{
			Name: "projects",
			Commands: []*cli.Command{
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
			Name: "projects:snippets",
			Commands: []*cli.Command{
				&projectsSnippetsCreateRequest,
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
			Name: "build_target_outputs",
			Commands: []*cli.Command{
				&buildTargetOutputsRetrieve,
			},
		},

		{
			Name: "orgs",
			Commands: []*cli.Command{
				&orgsRetrieve,
				&orgsList,
			},
		},
	},
	EnableShellCompletion: true,
	HideHelpCommand:       true,
}
