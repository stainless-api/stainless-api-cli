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

var generateCreateSpec = cli.Command{
	Name:  "create_spec",
	Usage: "Perform create_spec operation",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "project",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "source.revision",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "source.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "source.type",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "source.type",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "source.openapi_spec",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "source.openapi_spec",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "source.stainless_config",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "source.stainless_config",
			},
		},
	},
	Action:          handleGenerateCreateSpec,
	HideHelpCommand: true,
}

func handleGenerateCreateSpec(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainlessv0.GenerateNewSpecParams{}
	res, err := cc.client.Generate.NewSpec(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
