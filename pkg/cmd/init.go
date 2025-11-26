package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	cbuild "github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

//go:embed example-openapi.yml
var exampleSpecYAML string

//go:embed example-openapi.json
var exampleSpecJSON string

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

func handleInit(ctx context.Context, cmd *cli.Command) error {
	if err := authenticate(ctx, cmd, false); err != nil {
		return err
	}

	if err := ensureExistingWorkspaceIsDeleted(cmd); err != nil {
		return err
	}

	cc := getAPICommandContext(cmd)

	orgs := fetchUserOrgs(cc.client, ctx)
	orgs, err := ensureUserHasOrg(ctx, cmd, cc.client, orgs)
	if err != nil {
		return err
	}

	org, err := askSelectOrganization(cmd, orgs)
	if err != nil {
		return err
	}

	projects := fetchUserProjects(ctx, cc.client, org)

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
			return fmt.Errorf("project '%s' does not exist", project)
		}
	} else {
		project, targets, err = askSelectProject(projects)
		if err != nil {
			return err
		}

		// If project is empty, that means the user selected <New Project>
		if project == "" {
			var err error
			if project, targets, err = askCreateProject(ctx, cmd, cc, org, ""); err != nil {
				return err
			}
		} else {
			console.Property("project", project)
		}
	}

	console.Spacer()

	return initializeWorkspace(ctx, cmd, cc, project, targets)
}

func ensureExistingWorkspaceIsDeleted(cmd *cli.Command) error {
	var existingConfig WorkspaceConfig
	if found, err := existingConfig.Find(); err == nil && found {
		title := fmt.Sprintf("Existing workspace detected: %s (project: %s)", existingConfig.ConfigPath, existingConfig.Project)
		overwrite, err := console.Confirm(cmd, "", title, "Do you want to overwrite your existing workplace configuration?", true)
		if err != nil || !overwrite {
			return err
		}
		if err := os.Remove(existingConfig.ConfigPath); err != nil {
			return err
		}
	}
	return nil
}

// ensureUserHasOrg ensures the user has at least one organization, prompting to create one if needed
func ensureUserHasOrg(ctx context.Context, cmd *cli.Command, client stainless.Client, orgs []string) ([]string, error) {
	if len(orgs) == 0 {
		signupURL := "https://app.stainless.com/signup?source=cli"
		group := console.Info("Creating organization for user...")
		group.Property("url", signupURL)

		ok, err := console.Confirm(cmd, "browser", "Open browser?", "", true)
		if err != nil {
			return nil, err
		} else if ok && browser.OpenURL(signupURL) != nil {
			console.Info("Opening browser...")
		}

		group.Progress("Waiting for organization to be created...")

		for {
			time.Sleep(5 * time.Second)
			if orgs = fetchUserOrgs(client, ctx); len(orgs) > 0 {
				group.Success("Organization found! Continuing...")
				break
			}
		}

		console.Spacer()
	}
	return orgs, nil
}

func askSelectOrganization(cmd *cli.Command, orgs []string) (string, error) {
	var org string
	switch {
	case cmd.IsSet("org") && cmd.String("org") != "":
		org = cmd.String("org")
	case len(orgs) == 1:
		org = orgs[0]
	default:
		err := console.Field(huh.NewSelect[string]().
			Title("org").
			Description("Enter the organization for this project").
			Options(huh.NewOptions(slices.Sorted(slices.Values(orgs))...)...).
			Height(len(orgs) + 2).
			Value(&org))
		if err != nil {
			return "", err
		}
	}
	console.Property("org", org)
	return org, nil
}

func fetchUserProjects(ctx context.Context, client stainless.Client, org string) []stainless.Project {
	var projects []stainless.Project
	if projectsResponse, err := client.Projects.List(ctx, stainless.ProjectListParams{Org: stainless.String(org)}); err == nil {
		projects = projectsResponse.Data
	}
	return projects
}

// askSelectProject prompts the user to select from existing projects or create a new one
func askSelectProject(projects []stainless.Project) (string, []stainless.Target, error) {
	options := make([]huh.Option[*stainless.Project], 0, len(projects)+1)
	options = append(options, huh.NewOption("<New Project>", &stainless.Project{}))
	projects = slices.SortedFunc(slices.Values(projects), func(p1, p2 stainless.Project) int {
		if p1.Slug < p2.Slug {
			return -1
		}
		return 1
	})
	for _, project := range projects {
		options = append(options, huh.NewOption(project.Slug, &project))
	}

	var picked *stainless.Project
	err := console.Field(huh.NewSelect[*stainless.Project]().
		Title("project").
		Description("Choose or create a new project").
		Options(options...).
		Value(&picked))
	if err != nil {
		return "", nil, err
	}
	return picked.Slug, picked.Targets, nil
}

func askCreateProject(ctx context.Context, cmd *cli.Command, cc *apiCommandContext, org, projectName string) (string, []stainless.Target, error) {
	group := console.Property("project", "(new)")

	if projectName == "" {
		err := group.Field(huh.NewInput().
			Title("name").
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
	group.Property("name", projectName)

	// Determine targets
	var selectedTargets []stainless.Target
	if cmd.IsSet("targets") {
		for target := range strings.SplitSeq(cmd.String("targets"), ",") {
			selectedTargets = append(selectedTargets, stainless.Target(strings.TrimSpace(target)))
		}
		if len(selectedTargets) == 0 {
			return "", nil, fmt.Errorf("You must select at least one target!")
		}
	} else {
		allTargets := slices.DeleteFunc(getAllTargetInfo(), func(item TargetInfo) bool {
			return item.Name == "node" // Remove node (deprecated option)
		})

		options := make([]huh.Option[stainless.Target], len(allTargets))
		for i, target := range allTargets {
			options[i] = huh.NewOption(target.DisplayName, stainless.Target(target.Name)).Selected(target.DefaultSelected)
		}
		err := group.Field(huh.NewMultiSelect[stainless.Target]().
			Title("targets").
			Description("Select target languages for code generation").
			Options(options...).
			Validate(func(selected []stainless.Target) error {
				if len(selected) == 0 {
					return fmt.Errorf("You must select at least one target!")
				}
				return nil
			}).
			Value(&selectedTargets))
		if err != nil {
			return "", nil, err
		}
	}

	group.Property("targets", fmt.Sprintf("%v", selectedTargets))

	slug := nameToSlug(projectName)
	params := stainless.ProjectNewParams{
		DisplayName: projectName,
		Org:         org,
		Slug:        slug,
		Targets:     selectedTargets,
		Revision:    make(map[string]stainless.FileInputUnionParam),
	}

	// Get OpenAPI spec content
	oasContent, err := askExistingOpenAPISpec(group)
	if err != nil {
		return "", nil, err
	}

	// Determine format based on content (JSON starts with '{', otherwise assume YAML)
	var oasName string
	trimmedContent := strings.TrimSpace(oasContent)
	if strings.HasPrefix(trimmedContent, "{") {
		oasName = "openapi.json"
	} else {
		oasName = "openapi.yml"
	}

	params.Revision[oasName] = stainless.FileInputUnionParam{
		OfFileInputContent: &stainless.FileInputContentParam{
			Content: oasContent,
		},
	}

	_, err = cc.client.Projects.New(
		ctx,
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return "", nil, err
	}

	group.Success("Project created successfully")
	return slug, selectedTargets, nil
}

// initializeWorkspace sets up the local workspace configuration and downloads files
func initializeWorkspace(ctx context.Context, cmd *cli.Command, cc *apiCommandContext, projectSlug string, targets []stainless.Target) error {
	group := console.Info("Configuring .stainless/workspace.json")

	var openAPISpecPath string
	if cmd.IsSet("openapi-spec") {
		openAPISpecPath = cmd.String("openapi-spec")
		group.Property("openapi_spec", openAPISpecPath)
	} else if path, err := askOpenAPISpecLocation(group); err != nil {
		return err
	} else {
		openAPISpecPath = path
	}

	var stainlessConfigPath string
	if cmd.IsSet("stainless-config") {
		stainlessConfigPath = cmd.String("stainless-config")
		group.Property("stainless_config", stainlessConfigPath)
	} else if path, err := chooseStainlessConfigLocation(group); err != nil {
		return err
	} else {
		stainlessConfigPath = path
	}

	config, err := NewWorkspaceConfig(projectSlug, openAPISpecPath, stainlessConfigPath)
	if err != nil {
		group.Error("Failed to create workspace config: %v", err)
		return fmt.Errorf("project created but workspace initialization failed: %v", err)
	}

	if err = config.Save(); err != nil {
		group.Error("Failed to save workspace config: %v", err)
		return fmt.Errorf("project created but workspace initialization failed: %v", err)
	}

	group.Success("Workspace initialized at %s", config.ConfigPath)

	console.Spacer()

	if err := downloadConfigFiles(ctx, cc.client, config); err != nil {
		return fmt.Errorf("project created but downloading configuration files failed: %v", err)
	}

	console.Spacer()

	if len(targets) > 0 {
		if err := configureTargets(projectSlug, targets, &config); err != nil {
			return fmt.Errorf("target configuration failed: %v", err)
		}
	}

	console.Spacer()

	console.Info("Waiting for build to complete...")

	// Try to get the latest build for this project (which should have been created automatically)
	build, err := getLatestBuild(ctx, cc.client, projectSlug, "main")
	if err != nil {
		return fmt.Errorf("expected build to exist after project creation, but none found: %v", err)
	}

	downloadPaths := make(map[stainless.Target]string, len(config.Targets))
	for targetName, targetConfig := range config.Targets {
		downloadPaths[stainless.Target(targetName)] = targetConfig.OutputPath
	}

	model := buildCompletionModel{cbuild.NewModel(cc.client, ctx, *build, downloadPaths)}
	_, err = tea.NewProgram(model).Run()
	if err != nil {
		console.Warn(err.Error())
	}

	console.Spacer()

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

// askExistingOpenAPISpec provides the location of an _existing_ openapi spec. We first ask how the user would like
// provide the openapi spec, either 1. from computer, 2. from url, or 3. from an example. Then, we should fan
// out to the various options.
func askExistingOpenAPISpec(group console.Group) (content string, err error) {
	type Source string
	const (
		SourceComputer Source = "computer"
		SourceURL      Source = "url"
		SourceExample  Source = "example"
	)

	var source Source
	err = group.Field(huh.NewSelect[Source]().
		Title("openapi_spec").
		Description("How would you like to provide your OpenAPI spec?").
		Options(
			huh.NewOption("From a file on my computer", SourceComputer),
			huh.NewOption("From a URL", SourceURL),
			huh.NewOption("Use an example", SourceExample),
		).
		Value(&source))
	if err != nil {
		return "", err
	}

	switch source {
	case SourceComputer:
		// Use file picker to select file
		var filePath string
		err = group.Field(huh.NewFilePicker().
			Picking(true).
			CurrentDirectory(".").
			Title("openapi_spec (file)").
			Description("Select your OpenAPI spec file").
			ShowHidden(true).
			Height(10).
			Value(&filePath))
		if err != nil {
			return "", err
		}

		if filePath == "" {
			return "", fmt.Errorf("no file selected")
		}

		// Read the file
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read OpenAPI spec from %s: %w", filePath, err)
		}

		group.Property("openapi", filePath)

		return string(fileBytes), nil

	case SourceURL:
		// Ask for URL
		var urlStr string
		err = group.Field(huh.NewInput().
			Title("openapi_spec (url)").
			Description("Enter the URL to your OpenAPI spec file").
			Value(&urlStr).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("URL is required")
				}
				return nil
			}))
		if err != nil {
			return "", err
		}

		// Fetch content from URL
		resp, err := http.Get(urlStr)
		if err != nil {
			return "", fmt.Errorf("failed to fetch OpenAPI spec from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to fetch OpenAPI spec: HTTP %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}

		group.Property("openapi", urlStr)

		return string(bodyBytes), nil

	case SourceExample:
		group.Property("openapi", "petstore.yml")

		return exampleSpecJSON, nil

	default:
		group.Property("openapi_spec", "example.yml")
		return "", fmt.Errorf("unknown source: %s", source)
	}
}

// askOpenAPISpecLocation should ask where to store the OpenAPI spec on disk.
func askOpenAPISpecLocation(group console.Group) (string, error) {
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
	err := group.Field(huh.NewInput().
		Title("openapi_spec").
		Description("Path to where the OpenAPI spec file should be stored").
		Value(&path).
		Placeholder(suggestion))

	if err != nil {
		return "", err
	}

	group.Property("openapi_spec", path)
	return path, nil
}

func chooseStainlessConfigLocation(group console.Group) (string, error) {
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
	err := group.Field(huh.NewInput().
		Title("stainless_config").
		Description("Path where the Stainless config file should be stored").
		Value(&path).
		Placeholder(suggestion))

	if err != nil {
		return "", err
	}

	group.Property("stainless_config", path)
	return path, nil
}

// downloadConfigFiles downloads the OpenAPI spec and Stainless config from the API
func downloadConfigFiles(ctx context.Context, client stainless.Client, config WorkspaceConfig) error {
	if config.StainlessConfig == "" {
		return fmt.Errorf("No destination for the stainless configuration file")
	}
	if config.OpenAPISpec == "" {
		return fmt.Errorf("No destination for the OpenAPI spec file")
	}

	group := console.Info("Downloading configuration files")
	params := stainless.ProjectConfigGetParams{
		Project: stainless.String(config.Project),
		Include: stainless.String("openapi"),
	}

	configRes, err := client.Projects.Configs.Get(ctx, params)
	if err != nil {
		return fmt.Errorf("config download failed: %v", err)
	}

	// Helper function to write a file with confirmation if it exists
	writeFileWithConfirm := func(path string, content []byte, description string) error {
		// Create parent directories if they don't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for config file: %w", err)
		}

		// Check if the file exists and is nonempty
		if fileInfo, err := os.Stat(path); err == nil && fileInfo.Size() > 0 {
			// If contents are identical, this is a no-op
			existingContent, readErr := os.ReadFile(path)
			if readErr == nil && string(existingContent) == string(content) {
				return nil
			}

			shouldOverwrite, _, err := group.Confirm(nil, "", fmt.Sprintf("File %s already exists", path), "Do you want to overwrite it?", true)
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
		if err := writeFileWithConfirm(config.OpenAPISpec, []byte(openAPISpec), "OpenAPI spec"); err != nil {
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

	group := console.Info("Configuring targets...")

	// Initialize target configs with default absolute paths
	targetConfigs := make(map[stainless.Target]*TargetConfig, len(targets))
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	for _, target := range targets {
		defaultPath := filepath.Join("sdks", fmt.Sprintf("%s-%s", slug, target))
		targetConfigs[target] = &TargetConfig{OutputPath: Resolve(cwd, defaultPath)}
	}

	// Create form fields for each target
	pathVars := make(map[stainless.Target]*string, len(targets))
	fields := make([]huh.Field, 0, len(targets))

	for _, target := range targets {
		pathVar := targetConfigs[target].OutputPath
		pathVars[target] = &pathVar
		fields = append(fields, huh.NewInput().
			Title(fmt.Sprintf("%s output path", target)).
			Value(pathVars[target]))
	}

	// Run the form
	form := huh.NewForm(huh.NewGroup(fields...)).
		WithTheme(console.GetFormTheme(1)).
		WithKeyMap(console.GetFormKeyMap())
	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get target output paths: %v", err)
	}

	// Update config with user-provided paths (convert to absolute)
	for target, pathVar := range pathVars {
		if path := strings.TrimSpace(*pathVar); path != "" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s target: %w", target, err)
			}
			targetConfigs[target] = &TargetConfig{OutputPath: absPath}
		} else {
			delete(targetConfigs, target)
		}
	}

	// Save updated config
	config.Targets = targetConfigs
	if err := config.Save(); err != nil {
		group.Error("Failed to update workspace config with target paths: %v", err)
		return fmt.Errorf("workspace config update failed: %v", err)
	}

	for target, targetConfig := range targetConfigs {
		group.Property(string(target)+".output_path", targetConfig.OutputPath)
	}

	group.Success("Targets configured to output locally")
	return nil
}

// TargetInfo represents a target with its display name and default selection
type TargetInfo struct {
	DisplayName     string
	Name            stainless.Target
	DefaultSelected bool
}

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

func isValidTarget(targetInfos []TargetInfo, name stainless.Target) bool {
	for _, info := range targetInfos {
		if info.Name == name {
			return true
		}
	}
	return false
}

func targetInfoToOptions(targets []TargetInfo) []huh.Option[string] {
	options := make([]huh.Option[string], len(targets))
	for i, target := range targets {
		options[i] = huh.NewOption(target.DisplayName, string(target.Name)).Selected(target.DefaultSelected)
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

	buildObj := stainlessutils.NewBuild(*build)
	if len(buildObj.Languages()) == 0 {
		return targetInfo
	}

	return slices.DeleteFunc(targetInfo, func(item TargetInfo) bool {
		for name := range config.Targets {
			if name == item.Name {
				return false
			}
		}
		for _, target := range buildObj.Languages() {
			if target == item.Name {
				return false
			}
		}
		return true
	})
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

// findFile searches for a file by name in the current directory tree
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
