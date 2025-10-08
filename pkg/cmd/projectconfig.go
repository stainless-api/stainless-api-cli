// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
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
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: `Branch name, defaults to "main".`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "branch",
			},
			Value: "main",
		},
		&jsonflag.JSONStringFlag{
			Name: "include",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "include",
			},
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
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name:  "spec",
			Usage: "OpenAPI spec",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "spec",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
			Value: "main",
		},
	},
	Action:          handleProjectsConfigsGuess,
	HideHelpCommand: true,
}

func handleProjectsConfigsRetrieve(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectConfigGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Configs.Get(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
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

func handleProjectsConfigsGuess(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectConfigGuessParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	var res []byte
	_, err := cc.client.Projects.Configs.Guess(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
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
