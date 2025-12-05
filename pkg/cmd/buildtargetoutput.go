// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var buildsTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a method to download an output for a given build target.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
			Config: requestflag.RequestConfig{
				QueryPath: "build_id",
			},
		},
		&requestflag.StringFlag{
			Name:  "target",
			Usage: "SDK language target name",
			Config: requestflag.RequestConfig{
				QueryPath: "target",
			},
		},
		&requestflag.StringFlag{
			Name: "type",
			Config: requestflag.RequestConfig{
				QueryPath: "type",
			},
		},
		&requestflag.StringFlag{
			Name:  "output",
			Usage: "Output format: url (download URL) or git (temporary access token).",
			Value: requestflag.Value[stainless.BuildTargetOutputGetParamsOutput]("url"),
			Config: requestflag.RequestConfig{
				QueryPath: "output",
			},
		},
	},
	Action:          handleBuildsTargetOutputsRetrieve,
	HideHelpCommand: true,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildTargetOutputGetParams{}

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
	_, err = client.Builds.TargetOutputs.Get(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds:target-outputs retrieve", json, format, transform)
}
