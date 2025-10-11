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

var buildsDiagnosticsList = cli.Command{
	Name:  "list",
	Usage: "Get the list of diagnostics for a given build.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
		},
		&jsonflag.JSONStringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "limit",
			Usage: "Maximum number of diagnostics to return, defaults to 100 (maximum: 100)",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
			Value: 100,
		},
		&jsonflag.JSONStringFlag{
			Name:  "severity",
			Usage: "Includes the given severity and above (fatal > error > warning > note).",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "severity",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Optional list of language targets to filter diagnostics by",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Optional list of language targets to filter diagnostics by",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsDiagnosticsList,
	HideHelpCommand: true,
}

func handleBuildsDiagnosticsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildDiagnosticListParams{}
	var res []byte
	_, err := cc.client.Builds.Diagnostics.List(
		ctx,
		cmd.Value("build-id").(string),
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
	return ShowJSON("builds:diagnostics list", json, format, transform)
}
