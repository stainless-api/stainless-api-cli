// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	cbuild "github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/stainless-api/stainless-api-go/shared"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

// parseTargetPaths processes target flags to extract target:path syntax with workspace config
// Returns a map of target names to their custom paths
func parseTargetPaths(workspaceConfig WorkspaceConfig, targetsSlice []string) (downloadPaths map[stainless.Target]string, targets []stainless.Target, specifiedPath bool) {
	downloadPaths = make(map[stainless.Target]string)

	// Check workspace configuration for target paths if loaded
	for targetName, targetConfig := range workspaceConfig.Targets {
		if targetConfig.OutputPath != "" {
			downloadPaths[stainless.Target(targetName)] = targetConfig.OutputPath
		}
	}

	// Process the targets array from the CLI
	for _, target := range targetsSlice {
		cleanTarget, path := processSingleTarget(target)
		targets = append(targets, cleanTarget)
		if path != "" {
			// Command line target:path overrides workspace configuration
			downloadPaths[stainless.Target(cleanTarget)] = path
			specifiedPath = true
		}
	}

	return downloadPaths, targets, specifiedPath
}

// processSingleTarget extracts path from target:path and returns clean target name and path
func processSingleTarget(target string) (stainless.Target, string) {
	target = strings.TrimSpace(target)
	if !strings.Contains(target, ":") {
		return stainless.Target(target), ""
	}

	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return stainless.Target(target), ""
	}

	targetName := strings.TrimSpace(parts[0])
	targetPath := strings.TrimSpace(parts[1])
	return stainless.Target(targetName), targetPath
}

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
	Flags: []cli.Flag{
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
			Name:  "wait",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "pull",
			Usage: "Pull the build outputs after completion (only works with --wait)",
		},
		&cli.BoolFlag{
			Name:  "allow-empty",
			Usage: "Whether to allow empty commits (no changes). Defaults to false.",
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
		},
		&cli.StringFlag{
			Name:  "commit-message",
			Usage: "Optional commit message to use when creating a new commit.",
		},
		&cli.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
		},
	},
	Action:          handleBuildsCreate,
	Before:          before,
	HideHelpCommand: true,
}

var buildsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a build by its ID.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
		},
	},
	Action:          handleBuildsRetrieve,
	Before:          before,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Create two builds whose outputs can be directly compared with each other.",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
		},
	},
	Action:          handleBuildsCompare,
	Before:          before,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	wc := getWorkspace(ctx)

	buildGroup := console.Info("Creating build...")
	params := stainless.BuildNewParams{}
	if err := unmarshalStdinWithFlags(cmd, map[string]string{
		"allow-empty":    "allow_empty",
		"branch":         "branch",
		"commit-message": "commit_message",
		"targets":        "targets",
	}, &params); err != nil {
		return err
	}

	if name, oas, err := convertFileFlag(cmd, "openapi-spec"); err != nil {
		return err
	} else if oas != nil {
		if params.Revision.OfFileInputMap == nil {
			params.Revision.OfFileInputMap = make(map[string]shared.FileInputUnionParam)
		}
		params.Revision.OfFileInputMap["openapi"+path.Ext(name)] = shared.FileInputParamOfFileInputContent(string(oas))
	}

	if name, config, err := convertFileFlag(cmd, "stainless-config"); err != nil {
		return err
	} else if config != nil {
		if params.Revision.OfFileInputMap == nil {
			params.Revision.OfFileInputMap = make(map[string]shared.FileInputUnionParam)
		}
		params.Revision.OfFileInputMap["stainless"+path.Ext(name)] = shared.FileInputParamOfFileInputContent(string(config))
	}

	downloadPaths, targets, specifiedPath := parseTargetPaths(wc, cmd.StringSlice("target"))

	shouldDownload := specifiedPath || cmd.Bool("pull")
	if !shouldDownload {
		downloadPaths = make(map[stainless.Target]string)
	}

	params.Targets = targets

	build, err := client.Builds.New(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
	)
	if err != nil {
		return err
	}

	buildGroup.Property("build_id", build.ID)

	if cmd.Bool("wait") {
		console.Spacer()
		model := tea.Model(buildCompletionModel{
			Build: cbuild.NewModel(client, ctx, *build, cmd.String("branch"), downloadPaths),
		})
		model, err = tea.NewProgram(model).Run()
		if err != nil {
			console.Warn(err.Error())
		}
		b := model.(buildCompletionModel).Build
		build = &b.Build
		console.Spacer()
	}

	data := gjson.Parse(string(build.RawJSON()))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if err := ShowJSON("builds create", data, format, transform); err != nil {
		return err
	}

	for _, target := range data.Get("targets.@values").Array() {
		if target.Get("status").String() == "not_started" ||
			target.Get("commit.completed.conclusion").String() == "error" ||
			target.Get("lint.completed.conclusion").String() == "error" ||
			target.Get("test.completed.conclusion").String() == "error" {
			buildGroup.Error("Build did not succeed!")
			os.Exit(1)
		}
	}

	return nil
}

type buildCompletionModel struct {
	Build cbuild.Model
}

func (c buildCompletionModel) Init() tea.Cmd {
	return c.Build.Init()
}

func (c buildCompletionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	c.Build, cmd = c.Build.Update(msg)

	if c.IsCompleted() {
		return c, tea.Sequence(
			cmd,
			tea.Quit,
		)
	}

	return c, cmd
}

func (c buildCompletionModel) IsCompleted() bool {
	b := stainlessutils.NewBuild(c.Build.Build)
	for _, target := range b.Languages() {
		buildTarget := b.BuildTarget(target)

		var downloadIsCompleted = true
		if buildTarget.IsCommitCompleted() && stainlessutils.IsGoodCommitConclusion(buildTarget.Commit.Completed.Conclusion) {
			if download, ok := c.Build.Downloads[target]; ok {
				if download.Status != "completed" {
					downloadIsCompleted = false
				}
			}
		}

		if buildTarget == nil ||
			!buildTarget.IsCompleted() ||
			!downloadIsCompleted {
			return false
		}
	}

	return true
}

func (c buildCompletionModel) View() string {
	return c.Build.View()
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := client.Builds.Get(
		ctx,
		cmd.Value("build-id").(string),
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds retrieve", json, format, transform)
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildCompareParams{}
	if err := unmarshalStdinWithFlags(cmd, map[string]string{
		"targets": "targets",
	}, &params); err != nil {
		return err
	}
	var res []byte
	_, err := client.Builds.Compare(
		ctx,
		params,
		option.WithMiddleware(debugMiddleware(cmd.Bool("debug"))),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds compare", json, format, transform)
}
