// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
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
		&requestflag.Flag[string]{
			Name:      "build-id",
			Usage:     "Build ID",
			QueryPath: "build_id",
		},
		&cli.StringFlag{
			Name:  "project",
			Usage: "Project name (required when build-id is not provided)",
		},
		&requestflag.Flag[string]{
			Name:    "branch",
			Usage:   "Branch name (defaults to main if not provided)",
			Default: "main",
		},
		&requestflag.Flag[string]{
			Name:      "target",
			Usage:     "SDK language target name(s). Can be specified multiple times.",
			QueryPath: "target",
		},
		&requestflag.Flag[string]{
			Name:      "type",
			QueryPath: "type",
		},
		&requestflag.Flag[string]{
			Name:        "output",
			Usage:       "Output format: url (download URL) or git (temporary access token).",
			DefaultText: "url",
			HideDefault: true,
			QueryPath:   "output",
		},
	},
	Before: before,
	Action: handleBuildsTargetOutputsRetrieve,
}

func handleBuildsTargetOutputsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildTargetOutputGetParams{}
	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Builds.TargetOutputs.Get(ctx, params, options...)
	if err != nil {
		return err
	}

	buildID := cmd.Value("build-id").(string)
	if buildID == "" {
		latestBuild, err := getLatestBuild(ctx, client, cmd.String("project"), cmd.String("branch"))
		if err != nil {
			return fmt.Errorf("failed to get latest build: %v", err)
		}
		buildID = latestBuild.ID
	}

	wc := getWorkspace(ctx)
	downloadPaths, targets, _ := parseTargetPaths(wc, cmd.StringSlice("target"))

	if len(targets) == 0 {
		return fmt.Errorf("at least one target must be specified")
	}

	outputType := cmd.String("type")
	outputFormat := cmd.String("output")
	isPull := cmd.Bool("pull")

	for _, target := range targets {
		params := stainless.BuildTargetOutputGetParams{
			BuildID: buildID,
			Target:  stainless.BuildTargetOutputGetParamsTarget(target),
			Type:    stainless.BuildTargetOutputGetParamsType(outputType),
			Output:  stainless.BuildTargetOutputGetParamsOutput(outputFormat),
		}
		res, err := client.Builds.TargetOutputs.Get(
			ctx,
			params,
			debugMiddlewareOption,
		)
		if err != nil {
			return fmt.Errorf("failed to get output for target %s: %v", target, err)
		}

		if !isPull {
			json := gjson.Parse(res.RawJSON())
			format := cmd.Root().String("format")
			transform := cmd.Root().String("transform")
			if err := ShowJSON(os.Stdout, "builds:target_outputs retrieve", json, format, transform); err != nil {
				return err
			}
		}

		if isPull {
			group := console.Info("Downloading %s", target)

			// Get target output path from downloadPaths (which includes workspace config)
			targetDir := downloadPaths[target]

			if err := build.PullOutput(res.Output, res.URL, res.Ref, cmd.String("branch"), targetDir, group); err != nil {
				return fmt.Errorf("failed to pull %s: %v", target, err)
			}
		}
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
