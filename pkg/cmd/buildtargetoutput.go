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
		&cli.BoolFlag{
			Name: "pull",
		},
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
	Action: handleBuildsTargetOutputsRetrieve,
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

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))

	if cmd.Bool("pull") {
		build, err := cc.client.Builds.Get(ctx, cmd.String("build-id"))
		if err != nil {
			return err
		}
		targetDir := fmt.Sprintf("%s-%s", build.Project, cmd.String("target"))
		return pullOutput(res.Output, res.URL, res.Ref, targetDir)
	}

	return nil
}
