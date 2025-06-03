// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var buildsTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Download the output of a build target",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "build-id",
			Action: getAPIFlagAction[string]("query", "build_id"),
		},
		&cli.StringFlag{
			Name:   "target",
			Action: getAPIFlagAction[string]("query", "target"),
		},
		&cli.StringFlag{
			Name:   "type",
			Action: getAPIFlagAction[string]("query", "type"),
		},
		&cli.StringFlag{
			Name:   "output",
			Action: getAPIFlagAction[string]("query", "output"),
		},
	},
	Before:          initAPICommand,
	Action:          handleBuildsTargetOutputsRetrieve,
	HideHelpCommand: true,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
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
