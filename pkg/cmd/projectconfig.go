// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var projectsConfigsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve the configuration files for a given project.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "branch",
			Usage: `Branch name, defaults to "main".`,
			Value: "main",
		},
		&cli.StringFlag{
			Name: "include",
		},
	},
	Action:          handleProjectsConfigsRetrieve,
	HideHelpCommand: true,
}

var projectsConfigsGuess = cli.Command{
	Name:  "guess",
	Usage: "Generate suggestions for changes to config files based on an OpenAPI spec.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "spec",
			Usage: "OpenAPI spec",
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Value: "main",
		},
	},
	Action:          handleProjectsConfigsGuess,
	HideHelpCommand: true,
}

func handleProjectsConfigsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectConfigGetParams{
		Include: stainless.String(cmd.Value("include").(string)),
	}
	if cmd.IsSet("branch") {
		params.Branch = stainless.Opt(cmd.Value("branch").(string))
	}
	var res []byte
	_, err := client.Projects.Configs.Get(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:configs retrieve", json, format, transform)
}

func handleProjectsConfigsGuess(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectConfigGuessParams{}
	if err := unmarshalStdinWithFlags(cmd, map[string]string{
		"spec":   "spec",
		"branch": "branch",
	}, &params); err != nil {
		return err
	}
	var res []byte
	_, err := client.Projects.Configs.Guess(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("projects:configs guess", json, format, transform)
}
