// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var projectsConfigsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve the configuration files for a given project.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: `Branch name, defaults to "main".`,
			Value: requestflag.Value[string]("main"),
			Config: requestflag.RequestConfig{
				QueryPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name: "include",
			Config: requestflag.RequestConfig{
				QueryPath: "include",
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
		&requestflag.StringFlag{
			Name:  "spec",
			Usage: "OpenAPI spec",
			Config: requestflag.RequestConfig{
				BodyPath: "spec",
			},
		},
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Value: requestflag.Value[string]("main"),
			Config: requestflag.RequestConfig{
				BodyPath: "branch",
			},
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
	params := stainless.ProjectConfigGetParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Configs.Get(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:configs retrieve", obj, format, transform)
}

func handleProjectsConfigsGuess(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.ProjectConfigGuessParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Projects.Configs.Guess(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "projects:configs guess", obj, format, transform)
}
