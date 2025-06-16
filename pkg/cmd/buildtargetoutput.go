// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var buildsTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Download the output of a build target",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "build-id",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "build_id",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "target",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "output",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "output",
			},
		},
	},
	Action:          handleBuildsTargetOutputsRetrieve,
	HideHelpCommand: true,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainlessv0.BuildTargetOutputGetParams{}
	res, err := cc.client.Builds.TargetOutputs.Get(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
