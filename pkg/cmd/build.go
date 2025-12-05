// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goccy/go-yaml"
	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	cbuild "github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
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
		&requestflag.YAMLFlag{
			Name:  "revision",
			Usage: "Specifies what to build: a branch name, commit SHA, merge command\n(\"base..head\"), or file contents.",
			Config: requestflag.RequestConfig{
				BodyPath: "revision",
			},
		},
		&requestflag.BoolFlag{
			Name:  "allow-empty",
			Usage: "Whether to allow empty commits (no changes). Defaults to false.",
			Config: requestflag.RequestConfig{
				BodyPath: "allow_empty",
			},
		},
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
			Config: requestflag.RequestConfig{
				BodyPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name:  "commit-message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: requestflag.RequestConfig{
				BodyPath: "commit_message",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: requestflag.RequestConfig{
				BodyPath: "targets",
			},
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
		&requestflag.StringFlag{
			Name:  "build-id",
			Usage: "Build ID",
		},
	},
	Action:          handleBuildsRetrieve,
	Before:          before,
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List user-triggered builds for a given project.",
	Flags: []cli.Flag{
		&requestflag.StringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: requestflag.RequestConfig{
				QueryPath: "branch",
			},
		},
		&requestflag.StringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response.",
			Config: requestflag.RequestConfig{
				QueryPath: "cursor",
			},
		},
		&requestflag.FloatFlag{
			Name:        "limit",
			Usage:       "Maximum number of builds to return, defaults to 10 (maximum: 100).",
			Value:       requestflag.Value[float64](10),
			DefaultText: "10",
			Config: requestflag.RequestConfig{
				QueryPath: "limit",
			},
		},
		&requestflag.YAMLFlag{
			Name:  "revision",
			Usage: "A config commit SHA used for the build",
			Config: requestflag.RequestConfig{
				QueryPath: "revision",
			},
		},
	},
	Action:          handleBuildsList,
	HideHelpCommand: true,
}

var buildsCompare = cli.Command{
	Name:  "compare",
	Usage: "Create two builds whose outputs can be directly compared with each other.",
	Flags: []cli.Flag{
		&requestflag.YAMLFlag{
			Name:  "base",
			Usage: "Parameters for the base build",
			Config: requestflag.RequestConfig{
				BodyPath: "base",
			},
		},
		&requestflag.YAMLFlag{
			Name:  "head",
			Usage: "Parameters for the head build",
			Config: requestflag.RequestConfig{
				BodyPath: "head",
			},
		},
		&requestflag.StringSliceFlag{
			Name:  "target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: requestflag.RequestConfig{
				BodyPath: "targets",
			},
		},
	},
	Action:          handleBuildsCompare,
	Before:          before,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	var err error
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	wc := getWorkspace(ctx)

	buildGroup := console.Info("Creating build...")

	var revision map[string]map[string]map[string][]byte
	if cmd.IsSet("revision") {
		var ok bool
		revision, ok = requestflag.CommandRequestValue[any](cmd, "revision").(map[string]map[string]map[string][]byte)
		if !ok {
			revision = make(map[string]map[string]map[string][]byte)
		}
	} else {
		revision = make(map[string]map[string]map[string][]byte)
	}
	var exists bool
	var fileInputMap map[string]map[string][]byte
	if fileInputMap, exists = revision["file_input_map"]; !exists {
		fileInputMap = make(map[string]map[string][]byte)
		revision["file_input_map"] = fileInputMap
	}

	if name, oas, err := convertFileFlag(cmd, "openapi-spec"); err != nil {
		return err
	} else if oas != nil {
		fileInputMap["openapi"+path.Ext(name)] = map[string][]byte{
			"content": oas,
		}
	}

	if name, config, err := convertFileFlag(cmd, "stainless-config"); err != nil {
		return err
	} else if config != nil {
		revision["file_input_map"]["stainless"+path.Ext(name)] = map[string][]byte{
			"content": config,
		}
	}

	var revisionYAML []byte
	if revisionYAML, err = yaml.Marshal(revision); err != nil {
		return err
	}
	cmd.Set("revision", string(revisionYAML))

	downloadPaths, targets, specifiedPath := parseTargetPaths(wc, cmd.StringSlice("target"))

	shouldDownload := specifiedPath || cmd.Bool("pull")
	if !shouldDownload {
		downloadPaths = make(map[stainless.Target]string)
	}

	for _, t := range targets {
		cmd.Set("target", string(t))
	}

	params := stainless.BuildNewParams{}
	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
	)
	if err != nil {
		return err
	}

	build, err := client.Builds.New(
		ctx,
		params,
		options...,
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
			console.Warn("%s", err.Error())
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
	_, err = client.Builds.Get(
		ctx,
		requestflag.CommandRequestValue[string](cmd, "build-id"),
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds retrieve", json, format, transform)
}

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildListParams{}

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
	_, err = client.Builds.List(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds list", json, format, transform)
}

func handleBuildsCompare(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildCompareParams{}

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
	_, err = client.Builds.Compare(
		ctx,
		params,
		options...,
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("builds compare", json, format, transform)
}
