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

var specRetrieveDecoratedSpec = cli.Command{
	Name:  "retrieve-decorated-spec",
	Usage: "Retrieve the decorated spec for a given application and project.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "client-id",
		},
		&cli.StringFlag{
			Name: "project-name",
		},
	},
	Action:          handleSpecRetrieveDecoratedSpec,
	HideHelpCommand: true,
}

func handleSpecRetrieveDecoratedSpec(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("project-name") && len(unusedArgs) > 0 {
		cmd.Set("project-name", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.SpecGetDecoratedSpecParams{}
	if cmd.IsSet("client-id") {
		params.ClientID = cmd.Value("client-id").(string)
	}
	var res []byte
	_, err := cc.client.Spec.GetDecoratedSpec(
		ctx,
		cmd.Value("project-name").(string),
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
	return ShowJSON("spec retrieve-decorated-spec", json, format, transform)
}
