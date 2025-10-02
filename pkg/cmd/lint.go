package cmd

import (
	"context"
	"os"

	"github.com/stainless-api/stainless-api-go/option"
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
	return ShowJSON("Diagnostics", string(result), "json")
}
