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

	buildID := cmd.String("build-id")
	if buildID == "" {
		latestBuildID, err := getLatestBuildID(ctx, cc.client, cmd.String("project"), cmd.String("branch"))
		if err != nil {
			return fmt.Errorf("failed to get latest build: %v", err)
		}
		buildID = latestBuildID
	}

	params := stainless.BuildTargetOutputGetParams{
		BuildID: buildID,
	}
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
		return pullOutput(res.Output, res.URL, res.Ref, "")
	}

	return nil
}

func getLatestBuildID(ctx context.Context, client stainless.Client, project, branch string) (string, error) {
	if project == "" {
		return "", fmt.Errorf("project is required when build-id is not provided")
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
		return "", err
	}

	if len(res.Data) == 0 {
		return "", fmt.Errorf("no builds found for project %s", project)
	}

	return res.Data[0].ID, nil
}
