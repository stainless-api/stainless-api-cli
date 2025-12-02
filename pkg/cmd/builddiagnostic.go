// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/stainless-api/stainless-api-go/shared"
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
		&cli.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response",
		},
		&cli.Float64Flag{
			Name:  "limit",
			Usage: "Maximum number of diagnostics to return, defaults to 100 (maximum: 100)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:  "severity",
			Usage: "Includes the given severity and above (fatal > error > warning > note).",
		},
		&cli.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of language targets to filter diagnostics by",
		},
	},
	Action:          handleBuildsDiagnosticsList,
	HideHelpCommand: true,
}

func handleBuildsDiagnosticsList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildDiagnosticListParams{
		Cursor:   stainless.String(cmd.Value("cursor").(string)),
		Severity: cmd.Value("severity").(stainless.BuildDiagnosticListParamsSeverity),
		Targets:  cmd.Value("target").([]shared.Target),
	}
	if cmd.IsSet("limit") {
		params.Limit = stainless.Opt(cmd.Value("limit").(float64))
	}
	var res []byte
	_, err := client.Builds.Diagnostics.List(
		ctx,
		cmd.Value("build-id").(string),
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
	return ShowJSON("builds:diagnostics list", json, format, transform)
}
