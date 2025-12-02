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

var buildsTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a method to download an output for a given build target.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "SDK language target name",
		},
		&cli.StringFlag{
			Name: "type",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output format: url (download URL) or git (temporary access token).",
			Value: "url",
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
	params := stainless.BuildTargetOutputGetParams{
		BuildID: cmd.Value("build-id").(string),
		Target:  cmd.Value("target").(stainless.BuildTargetOutputGetParamsTarget),
		Type:    cmd.Value("type").(stainless.BuildTargetOutputGetParamsType),
	}
	if cmd.IsSet("output") {
		params.Output = cmd.Value("output").(stainless.BuildTargetOutputGetParamsOutput)
	}
	var res []byte
	_, err := client.Builds.TargetOutputs.Get(
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
	return ShowJSON("builds:target-outputs retrieve", json, format, transform)
}
