// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/pkg/console"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var buildsTargetOutputsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a method to download an output for a given build target.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "pull",
		},
		&jsonflag.JSONStringFlag{
			Name:  "build-id",
			Usage: "Build ID",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "build_id",
			},
		},
		&cli.StringFlag{
			Name:  "project",
			Usage: "Project name (required when build-id is not provided)",
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "Branch name (defaults to main if not provided)",
			Value: "main",
		},
		&jsonflag.JSONStringFlag{
			Name:  "target",
			Usage: "SDK language target name",
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
			Name:  "output",
			Usage: "Output format: url (download URL) or git (temporary access token).",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "output",
			},
			Value: "url",
		},
	},
	Action: handleBuildsTargetOutputsRetrieve,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	buildID := cmd.String("build-id")
	if buildID == "" {
		latestBuild, err := getLatestBuild(ctx, cc.client, cmd.String("project"), cmd.String("branch"))
		if err != nil {
			return fmt.Errorf("failed to get latest build: %v", err)
		}
		buildID = latestBuild.ID
	}

	params := stainless.BuildTargetOutputGetParams{
		BuildID: buildID,
	}
	var resBytes []byte
	res, err := cc.client.Builds.TargetOutputs.Get(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&resBytes),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(resBytes))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if err := ShowJSON("builds:target_outputs retrieve", json, format, transform); err != nil {
		return err
	}

	group := console.Info("Downloading output")
	if cmd.Bool("pull") {
		return pullOutput(res.Output, res.URL, res.Ref, "", &group)
	}

	return nil
}

func getLatestBuild(ctx context.Context, client stainless.Client, project, branch string) (*stainless.Build, error) {
	if project == "" {
		return nil, fmt.Errorf("project is required when build-id is not provided")
	}

	params := stainless.BuildListParams{
		Project: stainless.String(project),
		Limit:   stainless.Float(1.0),
	}
	if branch != "" {
		params.Branch = stainless.String(branch)
	}

	res, err := client.Builds.List(
		ctx,
		params,
	)
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("no builds found for project %s", project)
	}

	return &res.Data[0], nil
}
