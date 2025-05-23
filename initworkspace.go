// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
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

	// If project name wasn't provided via flag, prompt the user interactively
	projectName := cmd.String("project-name")
	if projectName == "" {
		projectInfoMap := fetchUserProjects(ctx)

		keyMap := huh.NewDefaultKeyMap()
		keyMap.Input.AcceptSuggestion = key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "complete"),
		)

		// Create custom theme with bullet point cursor and no borders
		theme := huh.ThemeBase()
		theme.Focused.Base = theme.Focused.Base.BorderStyle(lipgloss.NormalBorder())
		theme.Focused.Title = theme.Focused.Title.Bold(true)
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("project name").
					Value(&projectName).
					Suggestions(slices.Collect(maps.Keys(projectInfoMap))).
					Description("Enter the stainless project for this workspace").
					Validate(createProjectValidator(projectInfoMap)),
			),
		).WithTheme(theme).WithKeyMap(keyMap)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get project name: %v", err)
		}
		fmt.Printf("%s project name: %s\n", aurora.Bold("✱"), projectName)
	}

	if err := InitWorkspaceConfig(projectName); err != nil {
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
