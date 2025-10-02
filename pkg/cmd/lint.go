package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
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
	var specParams GenerateSpecParams
	specParams.Project = cmd.String("project")
	specParams.Source.Type = "upload"

	stainlessConfig, err := os.ReadFile(cmd.String("stainless-config"))
	if err != nil {
		return err
	}
	specParams.Source.StainlessConfig = string(stainlessConfig)

	openAPISpec, err := os.ReadFile(cmd.String("openapi-spec"))
	if err != nil {
		return err
	}
	specParams.Source.OpenAPISpec = string(openAPISpec)

	var result []byte
	err = cc.client.Post(ctx, "api/generate/spec", specParams, &result, option.WithMiddleware(cc.AsMiddleware()))
	if err != nil {
		return err
	}

	transform := "spec.diagnostics.@values.@flatten.#(ignored==false)#"
	jsonObj := gjson.Parse(string(result)).Get(transform)
	var diagnostics []stainless.BuildDiagnostic
	json.Unmarshal([]byte(jsonObj.Raw), &diagnostics)
	if cmd.IsSet("format") {
		if err := ShowJSON("Diagnostics", jsonObj, cmd.String("format"), ""); err != nil {
			return err
		}
	} else {
		fmt.Println(ViewDiagnosticsPrint(diagnostics, -1))
	}

	for _, d := range diagnostics {
		if !d.Ignored {
			switch d.Level {
			case stainless.BuildDiagnosticLevelFatal:
			case stainless.BuildDiagnosticLevelError:
			case stainless.BuildDiagnosticLevelWarning:
				os.Exit(1)
			case stainless.BuildDiagnosticLevelNote:
				continue
			}
		}
	}
	return nil
}
