package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Initialize a new stainless project interactively",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name:  "org",
			Usage: "Organization name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "org",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "display-name",
			Usage: "Project display name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "display_name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "slug",
			Usage: "Project slug",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "slug",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "targets",
			Usage: "Comma-separated list of target languages",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name:  "+target",
			Usage: "Add a single target language",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.-1",
			},
		},
		&cli.StringFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Usage:   "Path to OpenAPI spec file",
		},
		&cli.BoolFlag{
			Name:  "workspace-init",
			Usage: "Initialize workspace configuration",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "download-config",
			Usage: "Download Stainless config to workspace",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "download-targets",
			Usage: "Download and configure SDK targets",
			Value: true,
		},
	},
	Action:          handleInit,
	HideHelpCommand: true,
}

func handleInit(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	targetInfo := getAvailableTargetInfo(ctx, cc.client, "", WorkspaceConfig{})
	targetInfo = slices.DeleteFunc(targetInfo, func(item TargetInfo) bool {
		// Remove node since it's a deprecated option
		return item.Name == "node"
	})
	availableTargets := targetInfoToOptions(targetInfo)

	orgs := fetchUserOrgs(cc.client, ctx)

	if len(orgs) == 0 {
		signupURL := "https://app.stainless.com/signup?source=cli"
		group := Info("Creating organization for user...")
		group.Property("url", signupURL)

		ok, err := Confirm(cmd, "browser", "Open browser?", "", true)
		if err != nil {
			return err
		}
		if ok {
			if err := browser.OpenURL(signupURL); err != nil {
				Info("Opening browser...")
			}
		}

		group.Progress("Waiting for organization to be created...")

		for {
			time.Sleep(5 * time.Second)
			orgs = fetchUserOrgs(cc.client, ctx)
			if len(orgs) > 0 {
				group.Success("Organization found! Continuing...")
				break
			}
		}

		Spacer()
	}

	org := cmd.String("org")
	if org == "" && len(orgs) > 0 {
		org = orgs[0]
	}

	projectName := cmd.String("display-name")
	if projectName == "" {
		projectName = cmd.String("slug")
	}
	targetsFlag := cmd.String("targets")
	openAPISpec := cmd.String("openapi-spec")

	// Convert comma-separated targets flag to slice for multi-select
	var selectedTargets []string
	if targetsFlag != "" {
		for _, target := range strings.Split(targetsFlag, ",") {
			selectedTargets = append(selectedTargets, strings.TrimSpace(target))
		}
	}

	if openAPISpec == "" {
		openAPISpec = findOpenAPISpec()
	}

	group := Info("Creating a new project...")

	// Check if all required values are provided via flags
	allValuesProvided := org != "" && projectName != "" && openAPISpec != ""
	if !allValuesProvided {

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("org").
					Value(&org).
					Suggestions(orgs).
					Description("Enter the organization for this project").
					Validate(func(s string) error {
						if strings.TrimSpace(s) == "" {
							return fmt.Errorf("organization is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("project").
					Value(&projectName).
					DescriptionFunc(func() string {
						if projectName == "" {
							return "Project name, slug will be 'my-project'."
						}
						slug := nameToSlug(projectName)
						return fmt.Sprintf("Project name, slug will be '%s'.", slug)
					}, &projectName).
					Placeholder("My Project").
					Validate(func(s string) error {
						if strings.TrimSpace(s) == "" {
							return fmt.Errorf("project name is required")
						}
						return nil
					}),
			),
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("targets").
					Description("Select target languages for code generation").
					Options(availableTargets...).
					Value(&selectedTargets),
				huh.NewInput().
					Title("openapi_spec").
					Description("Relative path to your OpenAPI spec file").
					Value(&openAPISpec).
					Validate(func(s string) error {
						if strings.TrimSpace(s) == "" {
							return fmt.Errorf("OpenAPI spec file is required")
						}
						if _, err := os.Stat(s); os.IsNotExist(err) {
							return fmt.Errorf("file '%s' does not exist", s)
						}
						return nil
					}),
			),
		).WithTheme(GetFormTheme(1)).WithKeyMap(GetFormKeyMap())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get project configuration: %v", err)
		}

		group.Property("organization", org)
		group.Property("project_name", projectName)
		if len(selectedTargets) > 0 {
			group.Property("targets", strings.Join(selectedTargets, ", "))
		}
		if openAPISpec != "" {
			group.Property("openapi_spec", openAPISpec)
		}
	}

	// Generate slug from project name
	slug := nameToSlug(projectName)

	// Set the CLI flags so that the JSONFlag middleware can pick them up
	cmd.Set("org", org)
	cmd.Set("display-name", projectName)
	cmd.Set("slug", slug)
	for _, target := range selectedTargets {
		cmd.Set("+target", target)
	}

	// Inject file contents into the API payload if files are provided or found
	if openAPISpec != "" {
		content, err := os.ReadFile(openAPISpec)
		if err == nil {
			// Inject the actual file content into the project creation payload
			jsonflag.Mutate(jsonflag.Body, "revision.openapi\\.yml.content", string(content))
		}
	}

	params := stainless.ProjectNewParams{}
	_, err := cc.client.Projects.New(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}
	group.Success("Project created successfully")

	Spacer()

	var config WorkspaceConfig
	{
		group := Info("Initializing workspace...")

		// Use the same project name (slug) for workspace initialization
		config, err = NewWorkspaceConfig(slug, openAPISpec, "")
		if err != nil {
			group.Error("Failed to create workspace config: %v", err)
			return fmt.Errorf("project created but workspace initialization failed: %v", err)
		}

		err = config.Save()
		if err != nil {
			group.Error("Failed to save workspace config: %v", err)
			return fmt.Errorf("project created but workspace initialization failed: %v", err)
		}

		group.Success("Workspace initialized at " + config.ConfigPath)
	}

	Spacer()

	{
		config.StainlessConfig = "./stainless.yml"
		err = config.Save()
		if err != nil {
			return fmt.Errorf("workspace update failed: %v", err)
		}
		if err := downloadStainlessConfig(ctx, cc.client, slug, &config); err != nil {
			return fmt.Errorf("project created but config download failed: %v", err)
		}
	}

	Spacer()

	{
		if len(selectedTargets) > 0 {
			if err := configureTargets(slug, selectedTargets, &config); err != nil {
				return fmt.Errorf("target configuration failed: %v", err)
			}
		}
	}

	Spacer()

	// Wait for build and pull outputs if workspace is configured
	{
		build, err := waitForLatestBuild(ctx, cc.client, slug)
		if err != nil {
			return fmt.Errorf("build wait failed: %v", err)
		}

		if len(config.Targets) > 0 {
			Spacer()

			if err := pullConfiguredTargets(ctx, cc.client, *build, config); err != nil {
				return fmt.Errorf("pull targets failed: %v", err)
			}
		}
	}

	Spacer()

	fmt.Fprintf(
		os.Stderr,
		"%s\n",
		lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).Render(
			"Next steps:\n\n"+
				"  * See "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("https://www.stainless.com/docs/guides/configure")+" to learn how to customize your SDKs\n\n"+
				"  * Use "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("stl builds create")+" to create more builds\n"+
				"  * Use "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("stl dev")+" to launch a development server that helps you build and see output locally.",
		),
	)

	return nil
}

// fetchUserOrgs retrieves the list of organizations the user has access to
func fetchUserOrgs(client stainless.Client, ctx context.Context) []string {
	res, err := client.Orgs.List(ctx)
	if err != nil {
		// Return empty slice if we can't fetch orgs
		return []string{}
	}

	var orgs []string
	for _, org := range res.Data {
		if org.Slug != "" {
			orgs = append(orgs, org.Slug)
		}
	}

	return orgs
}

// nameToSlug converts a project name to a URL-friendly slug
func nameToSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and common punctuation with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, ".", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "\\", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// downloadStainlessConfig downloads the stainless config file for a project
func downloadStainlessConfig(ctx context.Context, client stainless.Client, slug string, config *WorkspaceConfig) error {
	stainlessConfig := config.StainlessConfig
	if config.StainlessConfig == "" {
		return nil
	}

	group := Info("Downloading Stainless config...")

	params := stainless.ProjectConfigGetParams{
		Project: stainless.String(slug),
	}

	var configRes *stainless.ProjectConfigGetResponse
	var err error
	maxRetries := 3

	// I'm not sure why, but our endpoint here doesn't work immediately after the project is created, but
	// retrying it reliably fixes it.
	for attempt := 1; attempt <= maxRetries; attempt++ {
		configRes, err = client.Projects.Configs.Get(ctx, params)
		if err == nil {
			break
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("config download failed after %d attempts: %v", maxRetries, err)
	}

	content := ""
	if try, ok := (*configRes)["stainless.yml"]; ok {
		content = try.Content
	}
	if try, ok := (*configRes)["openapi.stainless.yml"]; ok {
		content = try.Content
	}

	err = os.WriteFile(stainlessConfig, []byte(content), 0644)
	if err != nil {
		group.Error("Failed to save project config to %s: %v", stainlessConfig, err)
		return fmt.Errorf("config save failed: %v", err)
	}

	group.Success("Stainless config downloaded to %s", stainlessConfig)
	return nil
}

// configureTargets prompts user for target output paths and saves them to workspace config
func configureTargets(slug string, selectedTargets []string, config *WorkspaceConfig) error {
	if len(selectedTargets) == 0 {
		return nil
	}

	group := Info("Configuring targets...")

	// Collect output paths for each selected target
	targets := map[string]*TargetConfig{}
	for _, target := range selectedTargets {
		defaultPath := fmt.Sprintf("./%s-%s", slug, target)
		targets[target] = &TargetConfig{OutputPath: defaultPath}
	}

	pathVars := make(map[string]*string)
	var fields []huh.Field

	for _, target := range selectedTargets {
		pathVar := targets[target].OutputPath
		pathVars[target] = &pathVar
		input := huh.NewInput().
			Title(fmt.Sprintf("%s output path", target)).
			Value(pathVars[target])
		fields = append(fields, input)
	}

	form := huh.NewForm(huh.NewGroup(fields...)).
		WithTheme(GetFormTheme(1)).
		WithKeyMap(GetFormKeyMap())
	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get target output paths: %v", err)
	}

	// Update the targets map with the final values, skipping empty paths
	for target, pathVar := range pathVars {
		if strings.TrimSpace(*pathVar) != "" {
			targets[target] = &TargetConfig{OutputPath: *pathVar}
		} else {
			// Remove the target if path is empty
			delete(targets, target)
		}
	}

	config.Targets = targets
	err := config.Save()
	if err != nil {
		group.Error("Failed to update workspace config with target paths: %v", err)
		return fmt.Errorf("workspace config update failed: %v", err)
	}
	for target, targetConfig := range targets {
		group.Property(target+".output_path", targetConfig.OutputPath)
	}
	group.Success("Targets configured to output locally")
	return nil
}

// waitForLatestBuild waits for the latest build to complete
func waitForLatestBuild(ctx context.Context, client stainless.Client, slug string) (*stainless.BuildObject, error) {
	waitGroup := Info("Waiting for build to complete...")

	// Try to get the latest build for this project (which should have been created automatically)
	build, err := getLatestBuild(ctx, client, slug, "main")
	if err != nil {
		return nil, fmt.Errorf("expected build to exist after project creation, but none found: %v", err)
	}

	waitGroup.Property("build_id", build.ID)
	build, err = waitForBuildCompletion(ctx, client, build.ID, &waitGroup)
	if err != nil {
		return nil, err
	}

	return build, nil
}

// pullConfiguredTargets pulls build outputs for configured targets
func pullConfiguredTargets(ctx context.Context, client stainless.Client, build stainless.BuildObject, config WorkspaceConfig) error {
	if config.Targets == nil || len(config.Targets) == 0 {
		return nil
	}

	pullGroup := Info("Pulling build outputs...")

	// Create target paths map from workspace config
	targetPaths := make(map[string]string)
	for targetName, targetConfig := range config.Targets {
		targetPaths[targetName] = targetConfig.OutputPath
	}

	if err := pullBuildOutputs(ctx, client, build, targetPaths, &pullGroup); err != nil {
		pullGroup.Error("Failed to pull outputs: %v", err)
		return err
	}

	return nil
}

// TargetInfo represents a target with its display name and default selection
type TargetInfo struct {
	DisplayName     string
	Name            string
	DefaultSelected bool
}

// getAllTargetInfo returns all available targets with their display names
func getAllTargetInfo() []TargetInfo {
	return []TargetInfo{
		{"TypeScript", "typescript", false},
		{"Python", "python", false},
		{"Node.js", "node", false},
		{"Go", "go", false},
		{"Java", "java", false},
		{"Kotlin", "kotlin", false},
		{"Ruby", "ruby", false},
		{"Terraform", "terraform", false},
		{"CLI", "cli", false},
		{"C#", "csharp", false},
		{"PHP", "php", false},
	}
}

// targetInfoToOptions converts TargetInfo slice to huh.Options
func targetInfoToOptions(targets []TargetInfo) []huh.Option[string] {
	options := make([]huh.Option[string], len(targets))
	for i, target := range targets {
		options[i] = huh.NewOption(target.DisplayName, target.Name).Selected(target.DefaultSelected)
	}
	return options
}

// getAvailableTargetInfo gets available targets from the project's latest build with workspace config for default selection
func getAvailableTargetInfo(ctx context.Context, client stainless.Client, projectName string, config WorkspaceConfig) []TargetInfo {
	targetInfo := getAllTargetInfo()

	for targetName := range config.Targets {
		for idx, target := range targetInfo {
			if targetName == target.Name {
				targetInfo[idx].DefaultSelected = true
			}
		}
	}

	// If there is no configured targets, just set python and typescript to be true.
	if len(config.Targets) == 0 {
		for idx, target := range targetInfo {
			if "typescript" == target.Name || "python" == target.Name {
				targetInfo[idx].DefaultSelected = true
			}
		}
	}

	// Try to get targets from latest build
	if projectName == "" {
		return targetInfo
	}
	build, err := getLatestBuild(ctx, client, projectName, "main")
	if err != nil {
		return targetInfo
	}
	buildTargets := getBuildTargetInfo(*build)
	if len(buildTargets) == 0 {
		return targetInfo
	}

	return slices.DeleteFunc(targetInfo, func(item TargetInfo) bool {
		for name := range config.Targets {
			if name == item.Name {
				return false
			}
		}
		for _, target := range buildTargets {
			if target.name == item.Name {
				return false
			}
		}
		return true
	})
}
