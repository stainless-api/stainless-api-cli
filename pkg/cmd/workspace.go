// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-go"
	"github.com/urfave/cli/v3"
)

var workspaceInit = cli.Command{
	Name:  "init",
	Usage: "Initialize stainless workspace configuration in current directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "project",
			Usage: "Project name to use for this workspace",
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
			Name:  "download-config",
			Usage: "Download Stainless config to workspace",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "download-targets",
			Usage: "Download configured targets after build completion",
			Value: true,
		},
	},
	Action:          handleWorkspaceInit,
	HideHelpCommand: true,
}

var workspaceStatus = cli.Command{
	Name:            "status",
	Usage:           "Show workspace configuration status",
	Action:          handleWorkspaceStatus,
	HideHelpCommand: true,
}

func handleWorkspaceInit(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	// Check for existing workspace configuration
	var existingConfig WorkspaceConfig
	found, err := existingConfig.Find()
	if err == nil && found {
		Info("Existing workspace detected: %s (project: %s)", existingConfig.ConfigPath, existingConfig.Project)
	}

	// Get current directory and show where the file will be written
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	configPath := filepath.Join(dir, "stainless-workspace.json")
	Info("Writing workspace config to: %s", configPath)

	Spacer()

	// Get values from flags or prepare for interactive prompt
	projectName := cmd.String("project")
	openAPISpec := cmd.String("openapi-spec")
	stainlessConfig := cmd.String("stainless-config")

	// Pre-fill OpenAPI spec and Stainless config if found and not provided via flags
	if openAPISpec == "" {
		openAPISpec = findOpenAPISpec()
	}
	if stainlessConfig == "" {
		stainlessConfig = findStainlessConfig()
	}

	// Skip interactive form if all values are provided via flags or auto-detected
	// Project name is required, but openAPISpec and stainlessConfig are optional
	allValuesProvided := projectName != "" &&
		(cmd.IsSet("openapi-spec") || openAPISpec != "") &&
		(cmd.IsSet("stainless-config") || stainlessConfig != "")

	if !allValuesProvided {
		projectInfoMap := fetchUserProjects(ctx, cc.client)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("project").
					Value(&projectName).
					Suggestions(slices.Collect(maps.Keys(projectInfoMap))).
					Description("Enter the stainless project for this workspace"),
				huh.NewInput().
					Title("openapi_spec (optional)").
					Description("Relative path to your OpenAPI spec file").
					Placeholder("openapi.yml").
					Value(&openAPISpec),
				huh.NewInput().
					Title("stainless_config (optional)").
					Description("Relative path to your Stainless config file").
					Placeholder("openapi.stainless.yml").
					Value(&stainlessConfig),
			),
		).WithTheme(GetFormTheme(0)).WithKeyMap(GetFormKeyMap())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get workspace configuration: %v", err)
		}
	}

	group := Info("Initializing workspace...")
	group.Property("project_name", projectName)
	if openAPISpec != "" {
		group.Property("openapi_spec", openAPISpec)
	}
	if stainlessConfig != "" {
		group.Property("stainless_config", stainlessConfig)
	}

	config, err := NewWorkspaceConfig(projectName, openAPISpec, stainlessConfig)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %v", err)
	}
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save workspace config: %v", err)
	}
	group.Success("Workspace initialized")

	Spacer()

	if stainlessConfig == "" {
		downloadConfig, err := Confirm(cmd, "download-config",
			"Download Stainless config to workspace?",
			"Manages Stainless config as part of your source code instead of in the cloud",
			true)
		if err != nil {
			return fmt.Errorf("failed to get stainless config preference: %v", err)
		}

		if downloadConfig {
			if err := downloadStainlessConfig(ctx, cc.client, projectName, &config); err != nil {
				return fmt.Errorf("config download failed: %v", err)
			}
		}
	} else {
		Info("Using existing Stainless config: %s", stainlessConfig)
	}

	Spacer()

	configureTargetsFlag, err := Confirm(cmd, "download-targets",
		"Configure targets?",
		"Set up output paths for SDK generation targets",
		true)
	if err != nil {
		return fmt.Errorf("failed to get target configuration preference: %v", err)
	}
	Info("Configuring targets...")
	var selectedTargets []string
	if configureTargetsFlag {
		// Get available targets from project's latest build with workspace config for defaults
		targetInfo := getAvailableTargetInfo(ctx, cc.client, projectName, config)
		availableTargets := targetInfoToOptions(targetInfo)

		targetForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("targets").
					Description("Select target languages for code generation").
					Options(availableTargets...).
					Value(&selectedTargets),
			),
		).WithTheme(GetFormTheme(1)).WithKeyMap(GetFormKeyMap())

		if err := targetForm.Run(); err != nil {
			return fmt.Errorf("failed to get target selection: %v", err)
		}

		if len(selectedTargets) > 0 {
			if err := configureTargets(projectName, selectedTargets, &config); err != nil {
				return fmt.Errorf("target configuration failed: %v", err)
			}
		}
	}

	if config.Targets != nil && len(config.Targets) > 0 {
		Spacer()

		build, err := waitForLatestBuild(ctx, cc.client, projectName)
		if err != nil {
			return fmt.Errorf("build wait failed: %v", err)
		}

		Spacer()

		if err := pullConfiguredTargets(ctx, cc.client, *build, config); err != nil {
			return fmt.Errorf("target download failed: %v", err)
		}
	}

	return nil
}

type projectInfo struct {
	Name string
	Org  string
}

// fetchUserProjects retrieves the list of projects the user has access to
func fetchUserProjects(ctx context.Context, client stainless.Client) map[string]projectInfo {
	params := stainless.ProjectListParams{}

	res, err := client.Projects.List(ctx, params)
	if err != nil {
		// Return empty map if we can't fetch projects
		return map[string]projectInfo{}
	}

	projectInfoMap := make(map[string]projectInfo)
	for _, project := range res.Data {
		if project.Slug != "" {
			projectInfoMap[project.Slug] = projectInfo{
				Name: project.Slug,
				Org:  project.Org,
			}
		}
	}

	return projectInfoMap
}

// findOpenAPISpec searches for common OpenAPI spec files in the current directory
func findOpenAPISpec() string {
	commonOpenAPIFiles := []string{
		"openapi.json",
		"openapi.yml",
		"openapi.yaml",
		"api.yml",
		"api.yaml",
		"spec.yml",
		"spec.yaml",
	}

	for _, filename := range commonOpenAPIFiles {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}
	return ""
}

// findStainlessConfig searches for common Stainless config files in the current directory
func findStainlessConfig() string {
	commonStainlessFiles := []string{
		"openapi.stainless.yml",
		"openapi.stainless.yaml",
		"stainless.yml",
		"stainless.yaml",
	}

	for _, filename := range commonStainlessFiles {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}
	return ""
}

func handleWorkspaceStatus(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	if cc.workspaceConfig.ConfigPath == "" {
		group := Warn("No workspace configuration found")
		group.Info("Run 'stl workspace init' to initialize a workspace in this directory.")
		return nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Get relative path from cwd to config file
	relPath, err := filepath.Rel(cwd, cc.workspaceConfig.ConfigPath)
	if err != nil {
		relPath = cc.workspaceConfig.ConfigPath // fallback to absolute path
	}

	group := Success("Workspace configuration found")
	group.Property("path", relPath)
	group.Property("project", cc.workspaceConfig.Project)

	if cc.workspaceConfig.OpenAPISpec != "" {
		// Check if OpenAPI spec file exists
		configDir := filepath.Dir(cc.workspaceConfig.ConfigPath)
		specPath := filepath.Join(configDir, cc.workspaceConfig.OpenAPISpec)
		if _, err := os.Stat(specPath); err == nil {
			group.Property("openapi_spec", cc.workspaceConfig.OpenAPISpec)
		} else {
			group.Property("openapi_spec", cc.workspaceConfig.OpenAPISpec+" (not found)")
		}
	} else {
		group.Property("openapi_spec", "(not configured)")
	}

	if cc.workspaceConfig.StainlessConfig != "" {
		// Check if Stainless config file exists
		configDir := filepath.Dir(cc.workspaceConfig.ConfigPath)
		stainlessPath := filepath.Join(configDir, cc.workspaceConfig.StainlessConfig)
		if _, err := os.Stat(stainlessPath); err == nil {
			group.Property("stainless_config", cc.workspaceConfig.StainlessConfig)
		} else {
			group.Property("stainless_config", cc.workspaceConfig.StainlessConfig+" (not found)")
		}
	} else {
		group.Property("stainless_config", "(not configured)")
	}

	return nil
}
