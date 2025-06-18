// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/logrusorgru/aurora/v4"
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
	// Check for existing workspace configuration
	var existingConfig WorkspaceConfig
	found, err := existingConfig.Find()
	if err == nil && found {
		fmt.Printf("Existing workspace detected: %s (project: %s)\n", aurora.Bold(existingConfig.ConfigPath), existingConfig.Project)
	}

	// Get current directory and show where the file will be written
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	configPath := filepath.Join(dir, "stainless-workspace.json")
	fmt.Printf("Writing workspace config to: %s\n", aurora.Bold(configPath))
	fmt.Println()

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
		projectInfoMap := fetchUserProjects(ctx)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("project").
					Value(&projectName).
					Suggestions(slices.Collect(maps.Keys(projectInfoMap))).
					Description("Enter the stainless project for this workspace").
					Validate(createProjectValidator(projectInfoMap)),
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
		).WithTheme(GetFormTheme()).WithKeyMap(GetFormKeyMap())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get workspace configuration: %v", err)
		}

		fmt.Printf("%s project name: %s\n", aurora.Bold("✱"), projectName)
		if openAPISpec != "" {
			fmt.Printf("%s openapi spec: %s\n", aurora.Bold("✱"), openAPISpec)
		}
		if stainlessConfig != "" {
			fmt.Printf("%s stainless config: %s\n", aurora.Bold("✱"), stainlessConfig)
		}
	}

	config, err := NewWorkspaceConfig(projectName, openAPISpec, stainlessConfig)
	if err != nil {
		return fmt.Errorf("failed to create workspace config: %v", err)
	}

	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save workspace config: %v", err)
	}

	fmt.Printf("%s %s\n", aurora.BrightGreen("✱"), fmt.Sprintf("Workspace initialized"))
	return nil
}

type projectInfo struct {
	Name string
	Org  string
}

// fetchUserProjects retrieves the list of projects the user has access to
func fetchUserProjects(ctx context.Context) map[string]projectInfo {
	client := stainlessv0.NewClient(getClientOptions()...)
	params := stainlessv0.ProjectListParams{}

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

func createProjectValidator(projectInfoMap map[string]projectInfo) func(string) error {
	attemptCount := 0
	lastProjectName := ""

	return func(projectName string) error {
		if projectName != lastProjectName {
			attemptCount = 0
			lastProjectName = projectName
		}
		if strings.TrimSpace(projectName) == "" {
			return fmt.Errorf("project name is required")
		}
		if _, exists := projectInfoMap[projectName]; exists {
			return nil
		}

		attemptCount++
		if attemptCount == 1 {
			return fmt.Errorf("project '%s' not found in accessible projects (press Enter again to proceed anyway)", projectName)
		}
		// Allow bypass on second attempt
		return nil
	}
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
	// Look for workspace configuration
	var config WorkspaceConfig
	found, err := config.Find()
	if err != nil {
		return fmt.Errorf("error searching for workspace config: %v", err)
	}

	if !found {
		fmt.Printf("%s No workspace configuration found\n", aurora.Yellow("!"))
		fmt.Printf("Run 'stl workspace init' to initialize a workspace in this directory.\n")
		return nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Get relative path from cwd to config file
	relPath, err := filepath.Rel(cwd, config.ConfigPath)
	if err != nil {
		relPath = config.ConfigPath // fallback to absolute path
	}

	fmt.Printf("%s Workspace configuration found\n", aurora.BrightGreen("✓"))
	fmt.Printf("  Config file: %s\n", aurora.Bold(relPath))
	fmt.Printf("  Project: %s\n", aurora.Bold(config.Project))

	if config.OpenAPISpec != "" {
		// Check if OpenAPI spec file exists
		configDir := filepath.Dir(config.ConfigPath)
		specPath := filepath.Join(configDir, config.OpenAPISpec)
		if _, err := os.Stat(specPath); err == nil {
			fmt.Printf("  OpenAPI spec: %s %s\n", aurora.Bold(config.OpenAPISpec), aurora.BrightGreen("✓"))
		} else {
			fmt.Printf("  OpenAPI spec: %s %s\n", aurora.Bold(config.OpenAPISpec), aurora.BrightRed("✗ (not found)"))
		}
	} else {
		fmt.Printf("  OpenAPI spec: %s\n", aurora.Faint("(not configured)"))
	}

	if config.StainlessConfig != "" {
		// Check if Stainless config file exists
		configDir := filepath.Dir(config.ConfigPath)
		stainlessPath := filepath.Join(configDir, config.StainlessConfig)
		if _, err := os.Stat(stainlessPath); err == nil {
			fmt.Printf("  Stainless config: %s %s\n", aurora.Bold(config.StainlessConfig), aurora.BrightGreen("✓"))
		} else {
			fmt.Printf("  Stainless config: %s %s\n", aurora.Bold(config.StainlessConfig), aurora.BrightRed("✗ (not found)"))
		}
	} else {
		fmt.Printf("  Stainless config: %s\n", aurora.Faint("(not configured)"))
	}

	return nil
}
