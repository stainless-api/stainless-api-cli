// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var buildTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "TODO",
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
	Action:          handleBuildTargetOutputsRetrieve,
	HideHelpCommand: true,
}

func handleBuildTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)

	res, err := cc.client.BuildTargetOutputs.Get(
		context.TODO(),
		stainlessv0.BuildTargetOutputGetParams{},
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
