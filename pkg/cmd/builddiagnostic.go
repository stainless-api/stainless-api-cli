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

var buildsDiagnosticsList = cli.Command{
	Name:  "list",
	Usage: "Get diagnostics for a build",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "build-id",
		},
		&jsonflag.JSONStringFlag{
			Name: "cursor",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "limit",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "severity",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "severity",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "targets",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+target",
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
	params := stainless.BuildDiagnosticListParams{}
	res, err := cc.client.Builds.Diagnostics.List(
		context.TODO(),
		cmd.Value("build-id").(string),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
