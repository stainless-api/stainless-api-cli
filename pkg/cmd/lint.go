package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var lintCommand = cli.Command{
	Name:  "lint",
	Usage: "Lint your stainless configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "Project name to use for the build",
		},
		&cli.StringFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Usage:   "Path to OpenAPI spec file",
		},
		&cli.StringFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Usage:   "Path to Stainless config file",
		},
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Watch for files to change and re-run linting",
		},
	},
	Action: runLinter,
}

type GenerateSpecParams struct {
	Project string `json:"project"`
	Source  struct {
		Type            string `json:"type"`
		OpenAPISpec     string `json:"openapi_spec"`
		StainlessConfig string `json:"stainless_config"`
	} `json:"source"`
}

func runLinter(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	for {
		diagnostics, err := getDiagnostics(ctx, cmd, cc)
		if err != nil {
			return err
		}

		if cmd.IsSet("format") {
			rawJson, err := json.Marshal(diagnostics)
			if err != nil {
				return err
			}
			jsonObj := gjson.Parse(string(rawJson))
			if err := ShowJSON("Diagnostics", jsonObj, cmd.String("format"), ""); err != nil {
				return err
			}
		} else {
			fmt.Println(ViewDiagnosticsPrint(diagnostics, -1))
		}

		if cmd.Bool("watch") {
			fmt.Println("\nDiagnostic checks will re-run once you edit your configuration files...")
			if err := waitTillConfigChanges(ctx, cmd, cc); err != nil {
				return err
			}
		} else {
			if hasBlockingDiagnostic(diagnostics) {
				os.Exit(1)
			}
			break
		}
	}

	return nil
}
