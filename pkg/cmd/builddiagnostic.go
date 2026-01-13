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

var buildsDiagnosticsList = cli.Command{
	Name:    "list",
	Usage:   "Get the list of diagnostics for a given build.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "build-id",
			Usage:    "Build ID",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from a previous response",
			QueryPath: "cursor",
		},
		&requestflag.Flag[float64]{
			Name:        "limit",
			Usage:       "Maximum number of diagnostics to return, defaults to 100 (maximum: 100)",
			Default:       100,
			DefaultText: "100",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "severity",
			Usage:     "Includes the given severity and above (fatal > error > warning > note).",
			QueryPath: "severity",
		},
		&requestflag.Flag[string]{
			Name:      "targets",
			Usage:     "Optional comma-delimited list of language targets to filter diagnostics by",
			QueryPath: "targets",
		},
	},
	Before:          before,
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

	params := stainless.BuildDiagnosticListParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Builds.Diagnostics.List(
			ctx,
			cmd.Value("build-id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "builds:diagnostics list", obj, format, transform)
	} else {
		iter := client.Builds.Diagnostics.ListAutoPaging(
			ctx,
			cmd.Value("build-id").(string),
			params,
			options...,
		)
		return ShowJSONIterator(os.Stdout, "builds:diagnostics list", iter, format, transform)
	}
}
