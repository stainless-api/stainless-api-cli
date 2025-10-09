package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Initialize a Stainless project",
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

func singleFieldForm(field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).WithTheme(GetFormTheme(0)).Run()
}

func handleInit(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	if err := authenticate(ctx, cmd, false); err != nil {
		return err
	}

	// Check for existing workspace configuration
	var existingConfig WorkspaceConfig
	if found, err := existingConfig.Find(); err == nil && found {
		title := fmt.Sprintf("Existing workspace detected: %s (project: %s)", existingConfig.ConfigPath, existingConfig.Project)
		overwrite, err := Confirm(cmd, "", title, "Do you want to overwrite your existing workplace configuration?", true)
		if err != nil || !overwrite {
			return err
		}
		if err := os.Remove(existingConfig.ConfigPath); err != nil {
			return err
		}
		existingConfig = WorkspaceConfig{}
	}

	orgs := fetchUserOrgs(cc.client, ctx)

	if len(orgs) == 0 {
		signupURL := "https://app.stainless.com/signup?source=cli"
		group := Info("Creating organization for user...")
		group.Property("url", signupURL)

		ok, err := Confirm(cmd, "browser", "Open browser?", "", true)
		if err != nil {
			return err
		} else if ok && browser.OpenURL(signupURL) != nil {
			Info("Opening browser...")
		}

		group.Progress("Waiting for organization to be created...")

		for {
			time.Sleep(5 * time.Second)
			if orgs = fetchUserOrgs(cc.client, ctx); len(orgs) > 0 {
				group.Success("Organization found! Continuing...")
				break
			}
		}

		Spacer()
	}

	// Determine organization
	var org string
	switch {
	case cmd.IsSet("org") && cmd.String("org") != "":
		org = cmd.String("org")
	case len(orgs) == 1:
		org = orgs[0]
	default:
		err := singleFieldForm(huh.NewSelect[string]().
			Title("org").
			Description("Enter the organization for this project").
			Options(huh.NewOptions(orgs...)...).
			Height(len(orgs) + 2).
			Value(&org))
		if err != nil {
			return err
		}
	}

	var projects []stainless.Project
	if projectsResponse, err := cc.client.Projects.List(ctx, stainless.ProjectListParams{}); err == nil {
		projects = projectsResponse.Data
	}

	var targets []stainless.Target
	project := ""

	if cmd.IsSet("project") {
		project = cmd.String("project")
		projectExists := false
		for _, p := range projects {
			// User can specify display name or slug, but we should normalize to slug here:
			if project == p.Slug || project == p.DisplayName {
				project = p.Slug
				projectExists = true
				break
			}
		}

		if !projectExists {
			confirm, err := Confirm(cmd, "", fmt.Sprintf("Project '%s' does not exist", project), "Do you want to create a new project?", true)
			if err != nil {
				return err
			}
			if !confirm {
				return fmt.Errorf("project '%s' does not exist", project)
			}
			if project, targets, err = createProject(ctx, cmd, cc, org, project); err != nil {
				return err
			}
		}
	} else if len(projects) > 0 {
		options := make([]huh.Option[*stainless.Project], 0, len(projects)+1)
		for _, project := range projects {
			options = append(options, huh.NewOption(project.Slug, &project))
		}
		options = append(options, huh.NewOption[*stainless.Project]("New Project", &stainless.Project{}))

		var picked *stainless.Project
		err := singleFieldForm(huh.NewSelect[*stainless.Project]().
			Title("project").
			Description("Choose a project for this workspace").
			Options(options...).
			Value(&picked))
		if err != nil {
			return err
		}
		project = picked.Slug
		targets = picked.Targets
	}

	if project == "" {
		var err error
		if project, targets, err = createProject(ctx, cmd, cc, org, ""); err != nil {
			return err
		}
	}

	return initializeWorkspace(ctx, cmd, cc, project, targets)
}

func createProject(ctx context.Context, cmd *cli.Command, cc *apiCommandContext, org, projectName string) (string, []stainless.Target, error) {
	info := Info("Creating a new project")

	if projectName == "" {
		err := singleFieldForm(huh.NewInput().
			Title("Project Display Name").
			Description("Enter a display name for your new project").
			Value(&projectName).
			DescriptionFunc(func() string {
				if projectName == "" {
					return "Project name, slug will be 'my-project'."
				}
				return fmt.Sprintf("Project name, slug will be '%s'.", nameToSlug(projectName))
			}, &projectName).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("project name is required")
				}
				return nil
			}))
		if err != nil {
			return "", nil, err
		}
	}

	info.Property("project", projectName)

	// Determine targets
	var selectedTargets []stainless.Target
	if cmd.IsSet("targets") {
		for _, target := range strings.Split(cmd.String("targets"), ",") {
			selectedTargets = append(selectedTargets, stainless.Target(strings.TrimSpace(target)))
		}
	} else {
		allTargets := slices.DeleteFunc(getAllTargetInfo(), func(item TargetInfo) bool {
			return item.Name == "node" // Remove node (deprecated option)
		})

		options := make([]huh.Option[stainless.Target], len(allTargets))
		for i, target := range allTargets {
			options[i] = huh.NewOption(target.DisplayName, stainless.Target(target.Name)).Selected(target.DefaultSelected)
		}
		err := singleFieldForm(huh.NewMultiSelect[stainless.Target]().
			Title("targets").
			Description("Select target languages for code generation").
			Options(options...).
			Value(&selectedTargets))
		if err != nil {
			return "", nil, err
		}
	}

	info.Property("targets", fmt.Sprintf("%v", selectedTargets))

	slug := nameToSlug(projectName)
	_, err := cc.client.Projects.New(
		ctx,
		stainless.ProjectNewParams{
			DisplayName: projectName,
			Org:         org,
			Slug:        slug,
			Targets:     selectedTargets,
		},
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return "", nil, err
	}

	info.Success("Project created successfully")
	return slug, selectedTargets, nil
}

func initializeWorkspace(ctx context.Context, cmd *cli.Command, cc *apiCommandContext, projectSlug string, targets []stainless.Target) error {
	info := Info("Initializing workspace...")

	var openAPISpecPath string
	if cmd.IsSet("openapi-spec") {
		openAPISpecPath = cmd.String("openapi-spec")
	} else if path, err := chooseOpenAPISpecLocation(); err != nil {
		return err
	} else {
		openAPISpecPath = path
	}

	var stainlessConfigPath string
	if cmd.IsSet("stainless-config") {
		stainlessConfigPath = cmd.String("stainless-config")
	} else if path, err := chooseStainlessConfigLocation(); err != nil {
		return err
	} else {
		stainlessConfigPath = path
	}

	config, err := NewWorkspaceConfig(projectSlug, openAPISpecPath, stainlessConfigPath)
	if err != nil {
		info.Error("Failed to create workspace config: %v", err)
		return fmt.Errorf("project created but workspace initialization failed: %v", err)
	}

	if err = config.Save(); err != nil {
		info.Error("Failed to save workspace config: %v", err)
		return fmt.Errorf("project created but workspace initialization failed: %v", err)
	}

	info.Success("Workspace initialized at %s", config.ConfigPath)

	Spacer()

	if err := downloadConfigFiles(ctx, cc.client, config); err != nil {
		return fmt.Errorf("project created but downloading configuration files failed: %v", err)
	}

	Spacer()

	if len(targets) > 0 {
		if err := configureTargets(projectSlug, targets, &config); err != nil {
			return fmt.Errorf("target configuration failed: %v", err)
		}
	}

	Spacer()

	// Wait for build and pull outputs
	build, err := waitForLatestBuild(ctx, cc.client, projectSlug)
	if err != nil {
		return fmt.Errorf("build wait failed: %v", err)
	}

	if len(config.Targets) > 0 {
		Spacer()
		if err := pullConfiguredTargets(ctx, cc.client, *build, config); err != nil {
			return fmt.Errorf("pull targets failed: %v", err)
		}
	}

	Spacer()

	// Get terminal width or use a sensible default
	width, _, err := term.GetSize(os.Stderr.Fd())
	if err != nil {
		width = 100
	} else if width > 100 {
		width = 100
	}

	fmt.Fprintf(
		os.Stderr,
		"%s\n",
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(width-2).
			Render(
				"Next steps:\n\n"+
					"  * See "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("https://www.stainless.com/docs/guides/configure")+" to learn how to customize your SDKs\n\n"+
					"  * Use "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("stl builds create")+" to create more builds\n"+
					"  * Use "+lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render("stl dev")+" to launch a development server that helps you build and see output locally.",
			),
	)

	return nil
}

// nameToSlug converts a project name to a URL-friendly slug
func nameToSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and common punctuation with hyphens
	replacements := []string{" ", "_", ".", "/", "\\"}
	for _, r := range replacements {
		slug = strings.ReplaceAll(slug, r, "-")
	}

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
	return strings.Trim(slug, "-")
}

func findFile(name string) string {
	var foundPath string
	_ = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == name {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})
	return foundPath
}

func chooseOpenAPISpecLocation() (string, error) {
	commonOpenAPIFiles := []string{
		"openapi.json", "openapi.yml", "openapi.yaml",
		"api.yml", "api.yaml", "spec.yml", "spec.yaml",
	}

	suggestion := filepath.Join(".stainless", "openapi.json")
	for _, filename := range commonOpenAPIFiles {
		if path := findFile(filename); path != "" {
			suggestion = path
			break
		}
	}

	path := suggestion
	err := singleFieldForm(huh.NewInput().
		Title("OpenAPI Specification Location").
		Description("Path where the OpenAPI specification file should be stored").
		Value(&path).
		Placeholder(suggestion))

	if err != nil {
		return "", err
	}

	Property("OpenAPI file", path)
	return path, nil
}

func chooseStainlessConfigLocation() (string, error) {
	commonStainlessFiles := []string{
		"openapi.stainless.yml", "openapi.stainless.yaml",
		"stainless.yml", "stainless.yaml",
	}

	suggestion := filepath.Join(".stainless", "stainless.yml")
	for _, filename := range commonStainlessFiles {
		if path := findFile(filename); path != "" {
			suggestion = path
			break
		}
	}

	path := suggestion
	err := singleFieldForm(huh.NewInput().
		Title("Stainless Config Location").
		Description("Path where the Stainless configuration file should be stored").
		Value(&path).
		Placeholder(suggestion))

	if err != nil {
		return "", err
	}

	Property("Stainless configuration file", path)
	return path, nil
}

func downloadConfigFiles(ctx context.Context, client stainless.Client, config WorkspaceConfig) error {
	if config.StainlessConfig == "" {
		return fmt.Errorf("No destination for the stainless configuration file")
	}
	if config.OpenAPISpec == "" {
		return fmt.Errorf("No destination for the OpenAPI specification file")
	}

	group := Info("Downloading Stainless config...")
	params := stainless.ProjectConfigGetParams{
		Project: stainless.String(config.Project),
		Include: stainless.String("openapi"),
	}

	configRes, err := client.Projects.Configs.Get(ctx, params)
	if err != nil {
		return fmt.Errorf("config download failed: %v", err)
	}

	group.Property("Available config files", "")
	for key := range *configRes {
		group.Property("- ", key)
	}

	// Helper function to write a file with confirmation if it exists
	writeFileWithConfirm := func(path string, content []byte, description string) error {
		// Create parent directories if they don't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for config file: %w", err)
		}

		// If the file exists and is nonempty, ask for confirmation
		if fileInfo, err := os.Stat(path); err == nil && fileInfo.Size() > 0 {
			shouldOverwrite, err := Confirm(nil, "", fmt.Sprintf("File %s already exists", path), "Do you want to overwrite it?", true)
			if err != nil {
				return fmt.Errorf("failed to confirm file overwrite: %w", err)
			}
			if !shouldOverwrite {
				group.Property("Note", fmt.Sprintf("Keeping existing file at %s", path))
				return nil
			}
		}

		if err := os.WriteFile(path, content, 0644); err != nil {
			group.Error("Failed to save project config to %s: %v", path, err)
			return fmt.Errorf("%s could not write to file: %v", description, err)
		}

		group.Success("%s downloaded to %s", description, path)
		return nil
	}

	// Handle Stainless config file
	{
		var stainlessConfig string
		for _, name := range []string{"stainless.yml", "openapi.stainless.yml"} {
			if try, ok := (*configRes)[name]; ok {
				stainlessConfig = try.Content
				break
			}
		}

		if err := writeFileWithConfirm(config.StainlessConfig, []byte(stainlessConfig), "Stainless configuration"); err != nil {
			return err
		}
	}

	// Handle OpenAPI spec file
	{
		var openAPISpec string
		for _, name := range []string{"openapi.json", "openapi.yml", "openapi.yaml"} {
			if try, ok := (*configRes)[name]; ok {
				openAPISpec = try.Content
				break
			}
		}

		// TODO: we should warn or confirm if the downloaded file has a different file extension than the destination filename
		if err := writeFileWithConfirm(config.OpenAPISpec, []byte(openAPISpec), "OpenAPI specification"); err != nil {
			return err
		}
	}

	return nil
}

// configureTargets prompts user for target output paths and saves them to workspace config
func configureTargets(slug string, targets []stainless.Target, config *WorkspaceConfig) error {
	if len(targets) == 0 {
		return nil
	}

	group := Info("Configuring targets...")

	// Initialize target configs with default paths
	targetConfigs := make(map[string]*TargetConfig, len(targets))
	for _, target := range targets {
		defaultPath := filepath.Join("sdks", fmt.Sprintf("%s-%s", slug, target))
		targetConfigs[string(target)] = &TargetConfig{OutputPath: defaultPath}
	}

	// Create form fields for each target
	pathVars := make(map[stainless.Target]*string, len(targets))
	fields := make([]huh.Field, 0, len(targets))

	for _, target := range targets {
		pathVar := targetConfigs[string(target)].OutputPath
		pathVars[target] = &pathVar
		fields = append(fields, huh.NewInput().
			Title(fmt.Sprintf("%s output path", target)).
			Value(pathVars[target]))
	}

	// Run the form
	form := huh.NewForm(huh.NewGroup(fields...)).
		WithTheme(GetFormTheme(1)).
		WithKeyMap(GetFormKeyMap())
	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get target output paths: %v", err)
	}

	// Update config with user-provided paths
	for target, pathVar := range pathVars {
		if path := strings.TrimSpace(*pathVar); path != "" {
			targetConfigs[string(target)] = &TargetConfig{OutputPath: path}
		} else {
			delete(targetConfigs, string(target))
		}
	}

	// Save updated config
	config.Targets = targetConfigs
	if err := config.Save(); err != nil {
		group.Error("Failed to update workspace config with target paths: %v", err)
		return fmt.Errorf("workspace config update failed: %v", err)
	}

	for target, targetConfig := range targetConfigs {
		group.Property(target+".output_path", targetConfig.OutputPath)
	}

	group.Success("Targets configured to output locally")
	return nil
}

// waitForLatestBuild waits for the latest build to complete
func waitForLatestBuild(ctx context.Context, client stainless.Client, slug string) (*stainless.Build, error) {
	waitGroup := Info("Waiting for build to complete...")

	// Try to get the latest build for this project (which should have been created automatically)
	build, err := getLatestBuild(ctx, client, slug, "main")
	if err != nil {
		return nil, fmt.Errorf("expected build to exist after project creation, but none found: %v", err)
	}

	waitGroup.Property("build_id", build.ID)
	return waitForBuildCompletion(ctx, client, build, &waitGroup)
}

// pullConfiguredTargets pulls build outputs for configured targets
func pullConfiguredTargets(ctx context.Context, client stainless.Client, build stainless.Build, config WorkspaceConfig) error {
	if len(config.Targets) == 0 {
		return nil
	}

	pullGroup := Info("Pulling build outputs...")

	// Create target paths map from workspace config
	targetPaths := make(map[string]string, len(config.Targets))
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

func isValidTarget(targetInfos []TargetInfo, name string) bool {
	for _, info := range targetInfos {
		if info.Name == name {
			return true
		}
	}
	return false
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

	// Mark targets from config as selected
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
			if target.Name == "typescript" || target.Name == "python" {
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
