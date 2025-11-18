// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	cbuild "github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"

	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v3"
)

// parseTargetPaths processes target flags to extract target:path syntax with workspace config
// Returns a map of target names to their custom paths
func parseTargetPaths(workspaceConfig WorkspaceConfig) map[stainless.Target]string {
	targetPaths := make(map[stainless.Target]string)

	// Check workspace configuration for target paths if loaded
	for targetName, targetConfig := range workspaceConfig.Targets {
		if targetConfig.OutputPath != "" {
			targetPaths[stainless.Target(targetName)] = targetConfig.OutputPath
		}
	}

	// Get the current JSON body with all mutations applied
	body, _, _, err := jsonflag.ApplyMutations([]byte("{}"), []byte("{}"), []byte("{}"))
	if err != nil {
		return targetPaths // If we can't parse, return map with workspace paths
	}

	// Check if there are any targets in the body
	targetsResult := gjson.GetBytes(body, "targets")
	if !targetsResult.Exists() {
		return targetPaths
	}

	// Process the targets array
	var cleanTargets []string
	for _, targetResult := range targetsResult.Array() {
		target := targetResult.String()
		cleanTarget, path := processSingleTarget(target)
		cleanTargets = append(cleanTargets, cleanTarget)
		if path != "" {
			// Command line target:path overrides workspace configuration
			targetPaths[stainless.Target(cleanTarget)] = path
		}
	}

	// Update the JSON body with cleaned targets
	if len(cleanTargets) > 0 {
		body, err = sjson.SetBytes(body, "targets", cleanTargets)
		if err != nil {
			return targetPaths
		}

		// Clear mutations and re-apply the cleaned JSON
		jsonflag.ClearMutations()

		// Re-register the cleaned body
		bodyObj := gjson.ParseBytes(body)
		bodyObj.ForEach(func(key, value gjson.Result) bool {
			jsonflag.Mutate(jsonflag.Body, key.String(), value.Value())
			return true
		})
	}

	return targetPaths
}

// processSingleTarget extracts path from target:path and returns clean target name and path
func processSingleTarget(target string) (string, string) {
	target = strings.TrimSpace(target)
	if !strings.Contains(target, ":") {
		return target, ""
	}

	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return target, ""
	}

	targetName := strings.TrimSpace(parts[0])
	targetPath := strings.TrimSpace(parts[1])
	return targetName, targetPath
}

var buildsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a build, on top of a project branch, against a given input revision.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "revision",
			},
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
			Name:  "wait",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "pull",
			Usage: "Pull the build outputs after completion (only works with --wait)",
		},
		&jsonflag.JSONBoolFlag{
			Name:  "allow-empty",
			Usage: "Whether to allow empty commits (no changes). Defaults to false.",
			Config: jsonflag.JSONConfig{
				Kind:     jsonflag.Body,
				Path:     "allow_empty",
				SetValue: true,
			},
			Value: false,
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: "The project branch to use for the build. If not specified, the\nbranch is inferred from the `revision`, and will 400 when that\nis not possible.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "commit-message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "target",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
			Hidden: true,
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsCreate,
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
	HideHelpCommand: true,
}

var buildsList = cli.Command{
	Name:  "list",
	Usage: "List user-triggered builds for a given project.",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "branch",
			Usage: "Branch name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "cursor",
			Usage: "Pagination cursor from a previous response.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name:  "limit",
			Usage: "Maximum number of builds to return, defaults to 10 (maximum: 100).",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
			Value: 10,
		},
		&jsonflag.JSONStringFlag{
			Name:  "revision",
			Usage: "A config commit SHA used for the build",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "revision",
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
		&jsonflag.JSONStringFlag{
			Name:  "base.branch",
			Usage: "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "base.revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "base.commit_message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "base.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.branch",
			Usage: "Branch to use. When using a branch name as revision, this must match or be\nomitted.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.branch",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.revision",
			Usage: `A branch name, commit SHA, or merge command in the format "base..head"`,
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.revision",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "head.commit_message",
			Usage: "Optional commit message to use when creating a new commit.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "head.commit_message",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "project",
			Usage: "Project name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "project",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Optional list of SDK targets to build. If not specified, all configured\ntargets will be built.",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
	},
	Action:          handleBuildsCompare,
	HideHelpCommand: true,
}

func handleBuildsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	// Handle file flags by reading files and mutating JSON body
	if err := applyFileFlag(cmd, "openapi-spec", "revision.openapi\\.yml.content"); err != nil {
		return err
	}
	if err := applyFileFlag(cmd, "stainless-config", "revision.openapi\\.stainless\\.yml.content"); err != nil {
		return err
	}

	buildGroup := console.Info("Creating build...")
	params := stainless.BuildNewParams{}
	build, err := cc.client.Builds.New(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	buildGroup.Property("build_id", build.ID)

	downloadPaths := parseTargetPaths(cc.workspaceConfig)

	if cmd.Bool("wait") {
		console.Spacer()
		model := tea.Model(buildCompletionModel{cbuild.NewModel(cc.client, ctx, *build, downloadPaths)})
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

	if stainlessutils.NewBuild(c.Build.Build).IsCompleted() {
		return c, tea.Sequence(
			cmd,
			tea.Quit,
		)
	}

	return c, cmd
}

func (c buildCompletionModel) View() string {
	return c.Build.View()
}

func handleBuildsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("build-id") && len(unusedArgs) > 0 {
		cmd.Set("build-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Builds.Get(
		ctx,
		cmd.Value("build-id").(string),
		option.WithMiddleware(cc.AsMiddleware()),
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

func handleBuildsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildListParams{}
	var res []byte
	_, err := cc.client.Builds.List(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
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
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	params := stainless.BuildCompareParams{}
	var res []byte
	_, err := cc.client.Builds.Compare(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
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
