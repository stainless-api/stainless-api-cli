package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/components/dev"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/workspace"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/stainless-api/stainless-api-go/shared"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var devCommand = cli.Command{
	Name:    "preview",
	Aliases: []string{"dev"},
	Usage:   "Development mode with interactive build monitoring",
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
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Which branch to use",
		},
		&cli.StringSliceFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "The target build language(s)",
		},
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Value:   false,
			Usage:   "Run in 'watch' mode to loop and rebuild when files change.",
		},
	},
	Before: before,
	Action: runPreview,
}

func runPreview(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("watch") {
		// Clear the screen and move the cursor to the top
		fmt.Print("\033[2J\033[H")
		os.Stdout.Sync()
	}

	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	wc := getWorkspace(ctx)

	gitUser, err := getGitUsername()
	if err != nil {
		console.Warn("Couldn't get a git user: %s", err)
		gitUser = "user"
	}

	var selectedBranch string
	if cmd.IsSet("branch") {
		selectedBranch = cmd.String("branch")
	} else {
		selectedBranch, err = chooseBranch(gitUser)
		if err != nil {
			return err
		}
	}
	console.Property("branch", selectedBranch)

	// Phase 2: Language selection
	var selectedTargets []string
	targetInfos := getAvailableTargetInfo(ctx, client, cmd.String("project"), wc)
	if cmd.IsSet("target") {
		selectedTargets = cmd.StringSlice("target")
		for _, target := range selectedTargets {
			if !isValidTarget(targetInfos, stainless.Target(target)) {
				return fmt.Errorf("invalid language target: %s", target)
			}
		}
	} else {
		selectedTargets, err = chooseSelectedTargets(targetInfos)
	}

	if len(selectedTargets) == 0 {
		return fmt.Errorf("no languages selected")
	}

	console.Property("targets", strings.Join(selectedTargets, ", "))

	// Convert string targets to stainless.Target
	targets := make([]stainless.Target, len(selectedTargets))
	for i, target := range selectedTargets {
		targets[i] = stainless.Target(target)
	}

	// Phase 3: Start build and monitor progress in a loop
	for {
		// Start the build process
		if err := runDevBuild(ctx, client, wc, cmd, selectedBranch, targets); err != nil {
			if errors.Is(err, build.ErrUserCancelled) {
				return nil
			}
			return err
		}

		if !cmd.Bool("watch") {
			break
		}

		// Clear the screen and move the cursor to the top
		fmt.Print("\nRebuilding...\n\n\033[2J\033[H")
		os.Stdout.Sync()
		console.Property("branch", selectedBranch)
		console.Property("targets", strings.Join(selectedTargets, ", "))
	}
	return nil
}

func chooseBranch(gitUser string) (string, error) {
	now := time.Now()
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	randomSuffix := base64.RawURLEncoding.EncodeToString(randomBytes)
	randomBranch := fmt.Sprintf("%s/%d%02d%02d-%s", gitUser, now.Year(), now.Month(), now.Day(), randomSuffix)

	branchOptions := []huh.Option[string]{}
	if currentBranch, err := getCurrentGitBranch(); err == nil && currentBranch != "main" && currentBranch != "master" {
		branchOptions = append(branchOptions,
			huh.NewOption(currentBranch, currentBranch),
		)
	}
	branchOptions = append(branchOptions,
		huh.NewOption(fmt.Sprintf("%s/dev", gitUser), fmt.Sprintf("%s/dev", gitUser)),
		huh.NewOption(fmt.Sprintf("%s/<random>", gitUser), randomBranch),
	)

	var selectedBranch string
	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("branch").
				Description("Select a Stainless project branch to use for development").
				Options(branchOptions...).
				Value(&selectedBranch),
		),
	).WithTheme(console.GetFormTheme(0))

	if err := branchForm.Run(); err != nil {
		return selectedBranch, fmt.Errorf("branch selection failed: %v", err)
	}

	return selectedBranch, nil
}

func chooseSelectedTargets(targetInfos []TargetInfo) ([]string, error) {
	targetOptions := targetInfoToOptions(targetInfos)

	var selectedTargets []string
	targetForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("targets").
				Description("Select targets to generate (space to select, enter to confirm, select none to select all):").
				Options(targetOptions...).
				Value(&selectedTargets),
		),
	).WithTheme(console.GetFormTheme(0))

	if err := targetForm.Run(); err != nil {
		return nil, fmt.Errorf("target selection failed: %v", err)
	}
	return selectedTargets, nil
}

func runDevBuild(ctx context.Context, client stainless.Client, wc workspace.Config, cmd *cli.Command, branch string, languages []stainless.Target) error {
	projectName := cmd.String("project")
	buildReq := stainless.BuildNewParams{
		Project:    stainless.String(projectName),
		Branch:     stainless.String(branch),
		Targets:    languages,
		AllowEmpty: stainless.Bool(true),
	}

	if name, oas, err := convertFileFlag(cmd, "openapi-spec"); err != nil {
		return err
	} else if oas != nil {
		if buildReq.Revision.OfFileInputMap == nil {
			buildReq.Revision.OfFileInputMap = make(map[string]shared.FileInputUnionParam)
		}
		buildReq.Revision.OfFileInputMap["openapi"+path.Ext(name)] = shared.FileInputParamOfFileInputContent(string(oas))
	}

	if name, config, err := convertFileFlag(cmd, "stainless-config"); err != nil {
		return err
	} else if config != nil {
		if buildReq.Revision.OfFileInputMap == nil {
			buildReq.Revision.OfFileInputMap = make(map[string]shared.FileInputUnionParam)
		}
		buildReq.Revision.OfFileInputMap["stainless"+path.Ext(name)] = shared.FileInputParamOfFileInputContent(string(config))
	}

	downloads := make(map[stainless.Target]string)
	for targetName, targetConfig := range wc.Targets {
		downloads[stainless.Target(targetName)] = targetConfig.OutputPath
	}

	model := dev.NewModel(
		client,
		ctx,
		branch,
		func() (*stainless.Build, error) {
			options := []option.RequestOption{}
			if cmd.Bool("debug") {
				options = append(options, debugMiddlewareOption)
			}
			build, err := client.Builds.New(
				ctx,
				buildReq,
				options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create build: %v", err)
			}
			return build, err
		},
		downloads,
		cmd.Bool("watch"),
	)

	p := console.NewProgram(model)
	finalModel, err := p.Run()

	if err != nil {
		return fmt.Errorf("failed to run TUI: %v", err)
	}
	if buildModel, ok := finalModel.(dev.Model); ok {
		return buildModel.Err
	}
	return nil
}

func getGitUsername() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	username := strings.TrimSpace(string(output))
	if username == "" {
		return "", fmt.Errorf("git username not configured")
	}

	// Convert to lowercase and replace spaces with hyphens for branch name
	return strings.ToLower(strings.ReplaceAll(username, " ", "-")), nil
}

func getCurrentGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", fmt.Errorf("could not determine current git branch")
	}

	return branch, nil
}

type GenerateSpecParams struct {
	Project string `json:"project"`
	Source  struct {
		Type            string `json:"type"`
		OpenAPISpec     string `json:"openapi_spec"`
		StainlessConfig string `json:"stainless_config"`
	} `json:"source"`
}

func getDiagnostics(ctx context.Context, cmd *cli.Command, client stainless.Client, wc workspace.Config) ([]stainless.BuildDiagnostic, error) {
	var specParams GenerateSpecParams
	if cmd.IsSet("project") {
		specParams.Project = cmd.String("project")
	} else {
		specParams.Project = wc.Project
	}
	specParams.Source.Type = "upload"

	configPath := wc.StainlessConfig
	if cmd.IsSet("stainless-config") {
		configPath = cmd.String("stainless-config")
	} else if configPath == "" {
		return nil, fmt.Errorf("You must provide a stainless configuration file with `--config /path/to/stainless.yml` or run this command from an initialized workspace.")
	}

	stainlessConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read your stainless configuration file:\n%w", err)
	}
	specParams.Source.StainlessConfig = string(stainlessConfig)

	oasPath := wc.OpenAPISpec
	if cmd.IsSet("openapi-spec") {
		oasPath = cmd.String("openapi-spec")
	} else if oasPath == "" {
		return nil, fmt.Errorf("You must provide an OpenAPI specification with `--oas /path/to/openapi.json` or run this command from an initialized workspace.")
	}

	openAPISpec, err := os.ReadFile(oasPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read your stainless configuration file:\n%w", err)
	}
	specParams.Source.OpenAPISpec = string(openAPISpec)

	options := []option.RequestOption{}
	if cmd.Bool("debug") {
		options = append(options, debugMiddlewareOption)
	}
	var result []byte
	err = client.Post(
		ctx,
		"api/generate/spec",
		specParams,
		&result,
		options...,
	)
	if err != nil {
		return nil, err
	}

	transform := "spec.diagnostics.@values.@flatten.#(ignored==false)#"
	jsonObj := gjson.Parse(string(result)).Get(transform)
	var diagnostics []stainless.BuildDiagnostic
	json.Unmarshal([]byte(jsonObj.Raw), &diagnostics)
	return diagnostics, nil
}

func hasBlockingDiagnostic(diagnostics []stainless.BuildDiagnostic) bool {
	for _, d := range diagnostics {
		if !d.Ignored {
			switch d.Level {
			case stainless.BuildDiagnosticLevelFatal:
			case stainless.BuildDiagnosticLevelError:
			case stainless.BuildDiagnosticLevelWarning:
				return true
			case stainless.BuildDiagnosticLevelNote:
				continue
			}
		}
	}
	return false
}
