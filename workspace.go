// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/logrusorgru/aurora/v4"
	"github.com/stainless-api/stainless-api-go"
	"github.com/urfave/cli/v3"
)

var initWorkspaceCommand = cli.Command{
	Name:  "init",
	Usage: "Initialize stainless workspace configuration in current directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "project-name",
			Usage: "Project name to use for this workspace",
		},
	},
	Action:          handleInitWorkspace,
	HideHelpCommand: true,
}

func handleInitWorkspace(ctx context.Context, cmd *cli.Command) error {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11"))

	fmt.Printf("%s\n", titleStyle.Render("workspace init"))

	// Check for existing workspace configuration
	existingConfig, existingPath, err := FindWorkspaceConfig()
	if err == nil && existingConfig != nil {
		fmt.Printf("Existing workspace detected: %s (project: %s)\n", aurora.Bold(existingPath), existingConfig.ProjectName)
	}

	// Get current directory and show where the file will be written
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	configPath := filepath.Join(dir, "stainless-workspace.json")
	fmt.Printf("Writing workspace config to: %s\n", aurora.Bold(configPath))

	fmt.Println()

	// Set up form theme and keymap for use in multiple forms
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Input.AcceptSuggestion = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "complete"),
	)
	keyMap.Input.Next = key.NewBinding(
		key.WithKeys("tab", "down", "enter"),
		key.WithHelp("tab/↓/enter", "next"),
	)
	keyMap.Input.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "previous"),
	)

	// Create custom theme with bullet point cursor and no borders
	theme := huh.ThemeBase()
	theme.Focused.Base = theme.Focused.Base.BorderStyle(lipgloss.NormalBorder())
	theme.Focused.Title = theme.Focused.Title.Bold(true)

	// Get project name from flag or prepare for interactive prompt
	projectName := cmd.String("project-name")
	var openAPISpec, stainlessConfig string

	// If project name wasn't provided via flag, prompt for all fields interactively
	if projectName == "" {
		projectInfoMap := fetchUserProjects(ctx)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("project name").
					Value(&projectName).
					Suggestions(slices.Collect(maps.Keys(projectInfoMap))).
					Description("Enter the stainless project for this workspace").
					Validate(createProjectValidator(projectInfoMap)),
				huh.NewInput().
					Title("OpenAPI spec path (optional)").
					Description("Relative path to your OpenAPI spec file").
					Placeholder("openapi.yml").
					Value(&openAPISpec),
				huh.NewInput().
					Title("Stainless config path (optional)").
					Description("Relative path to your Stainless config file").
					Placeholder("openapi.stainless.yml").
					Value(&stainlessConfig),
			),
		).WithTheme(theme).WithKeyMap(keyMap)

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

	if err := InitWorkspaceConfig(projectName, openAPISpec, stainlessConfig); err != nil {
		return fmt.Errorf("failed to initialize workspace: %v", err)
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
	client := stainlessv0.NewClient(getClientOptions(ctx, nil)...)
	params := stainlessv0.ProjectListParams{}

	res, err := client.Projects.List(context.TODO(), params)
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

// WorkspaceConfig stores workspace-level configuration
type WorkspaceConfig struct {
	ProjectName     string `json:"projectName"`
	OpenAPISpec     string `json:"openapi_spec,omitempty"`
	StainlessConfig string `json:"stainless_config,omitempty"`
}

// FindWorkspaceConfig searches for a stainless-workspace.json file starting from the current directory
// and moving up to parent directories until found or root is reached
func FindWorkspaceConfig() (*WorkspaceConfig, string, error) {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	for {
		configPath := filepath.Join(dir, "stainless-workspace.json")
		if _, err := os.Stat(configPath); err == nil {
			// Found config file
			config, err := LoadWorkspaceConfig(configPath)
			return config, configPath, err
		}

		// Move up to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			return nil, "", nil
		}
		dir = parent
	}
}

// LoadWorkspaceConfig loads the workspace config from the specified path
func LoadWorkspaceConfig(configPath string) (*WorkspaceConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var config WorkspaceConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveWorkspaceConfig saves the workspace config to the specified path
func SaveWorkspaceConfig(configPath string, config *WorkspaceConfig) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// GetProjectNameFromConfig returns the project name from workspace config if available
func GetProjectNameFromConfig() string {
	config, _, err := FindWorkspaceConfig()
	if err != nil || config == nil || config.ProjectName == "" {
		return ""
	}
	return config.ProjectName
}

// InitWorkspaceConfig initializes a new workspace config in the current directory
func InitWorkspaceConfig(projectName, openAPISpec, stainlessConfig string) error {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(dir, "stainless-workspace.json")
	config := WorkspaceConfig{
		ProjectName:     projectName,
		OpenAPISpec:     openAPISpec,
		StainlessConfig: stainlessConfig,
	}

	return SaveWorkspaceConfig(configPath, &config)
}
