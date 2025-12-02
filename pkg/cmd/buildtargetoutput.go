// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
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
		&cli.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
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
		&cli.StringFlag{
			Name:  "target",
			Usage: "SDK language target name",
		},
		&cli.StringFlag{
			Name: "type",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output format: url (download URL) or git (temporary access token).",
			Value: "url",
		},
	},
	Action: handleBuildsTargetOutputsRetrieve,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
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
	res, err := client.Builds.TargetOutputs.Get(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
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
		return build.PullOutput(res.Output, res.URL, res.Ref, "", &group)
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
