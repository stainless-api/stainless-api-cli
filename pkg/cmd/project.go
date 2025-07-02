// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/stainless-api/stainless-api-cli/pkg/jsonflag"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var projectsCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new project",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "display-name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "display_name",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "org",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "org",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "slug",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "slug",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "targets",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "targets.#",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "+target",
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
		&cli.StringFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Usage:   "Path to Stainless config file",
		},
		&cli.BoolFlag{
			Name:  "init-workspace",
			Usage: "Initialize workspace configuration after creating project",
			Value: true, // Default to true
		},
		&cli.BoolFlag{
			Name:  "download-config",
			Usage: "Download the project configuration after creation",
			Value: true,
		},
	},
	Action:          handleProjectsCreate,
	HideHelpCommand: true,
}

var projectsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve a project by name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
	},
	Action:          handleProjectsRetrieve,
	HideHelpCommand: true,
}

var projectsUpdate = cli.Command{
	Name:  "update",
	Usage: "Update a project's properties",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "project",
		},
		&jsonflag.JSONStringFlag{
			Name: "display-name",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Body,
				Path: "display_name",
			},
		},
	},
	Action:          handleProjectsUpdate,
	HideHelpCommand: true,
}

var projectsList = cli.Command{
	Name:  "list",
	Usage: "List projects in an organization, from oldest to newest",
	Flags: []cli.Flag{
		&jsonflag.JSONStringFlag{
			Name: "cursor",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "cursor",
			},
		},
		&jsonflag.JSONFloatFlag{
			Name: "limit",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "limit",
			},
		},
		&jsonflag.JSONStringFlag{
			Name: "org",
			Config: jsonflag.JSONConfig{
				Kind: jsonflag.Query,
				Path: "org",
			},
		},
	},
	Action:          handleProjectsList,
	HideHelpCommand: true,
}

func handleProjectsCreate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	// Define available target languages
	availableTargets := []huh.Option[string]{
		huh.NewOption("TypeScript", "typescript").Selected(true),
		huh.NewOption("Python", "python").Selected(true),
		huh.NewOption("Go", "go"),
		huh.NewOption("Java", "java"),
		huh.NewOption("Kotlin", "kotlin"),
		huh.NewOption("Ruby", "ruby"),
		huh.NewOption("Terraform", "terraform"),
		huh.NewOption("C#", "csharp"),
		huh.NewOption("PHP", "php"),
	}

	// Get values from flags
	org := cmd.String("org")
	projectName := cmd.String("display-name") // Keep display-name flag for compatibility
	if projectName == "" {
		projectName = cmd.String("slug") // Also check slug flag for compatibility
	}
	targetsFlag := cmd.String("targets")
	openAPISpec := cmd.String("openapi-spec")
	stainlessConfig := cmd.String("stainless-config")
	initWorkspace := cmd.Bool("init-workspace")
	downloadConfig := cmd.Bool("download-config")

	// Convert comma-separated targets flag to slice for multi-select
	var selectedTargets []string
	if targetsFlag != "" {
		for _, target := range strings.Split(targetsFlag, ",") {
			selectedTargets = append(selectedTargets, strings.TrimSpace(target))
		}
	}

	// Pre-fill OpenAPI spec and Stainless config if found and not provided via flags
	if openAPISpec == "" {
		openAPISpec = findOpenAPISpec()
	}
	if stainlessConfig == "" {
		stainlessConfig = findStainlessConfig()
	}

	// Check if all required values are provided via flags
	// Stainless config is optional, but OpenAPI spec is required
	allValuesProvided := org != "" && projectName != "" && openAPISpec != ""

	if !allValuesProvided {
		Info("Creating a new project...")

		// Fetch available organizations for suggestions
		orgs := fetchUserOrgs(cc.client, ctx)

		// Auto-fill with first organization if org is empty and orgs are available
		if org == "" && len(orgs) > 0 {
			org = orgs[0]
		}

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
			).Title("Page (1/2)"),
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
				huh.NewInput().
					Title("stainless_config (optional)").
					Description("Path where the Stainless config will be saved. If file exists, it will be used for project creation; otherwise, the auto-generated config will be saved here").
					Placeholder("stainless.yml").
					Value(&stainlessConfig),
				huh.NewConfirm().
					Title("download_config").
					Description("Download project configuration after creation").
					Value(&downloadConfig),
			).Title("Page (2/2)"),
		).WithTheme(GetFormTheme()).WithKeyMap(GetFormKeyMap())

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get project configuration: %v", err)
		}

		// Generate slug from project name
		slug := nameToSlug(projectName)

		Property("organization", org)
		Property("project_name", projectName)
		Property("slug", slug)
		if len(selectedTargets) > 0 {
			Property("targets", strings.Join(selectedTargets, ", "))
		}
		if openAPISpec != "" {
			Property("openapi_spec", openAPISpec)
		}
		if stainlessConfig != "" {
			Property("stainless_config", stainlessConfig)
		}

		// Set the flag values so the JSONFlag middleware can pick them up
		cmd.Set("org", org)
		cmd.Set("display-name", projectName)
		cmd.Set("slug", slug)
		for _, target := range selectedTargets {
			cmd.Set("+target", target)
		}
		if openAPISpec != "" {
			cmd.Set("openapi-spec", openAPISpec)
		}
		if stainlessConfig != "" {
			cmd.Set("stainless-config", stainlessConfig)
		}
	} else {
		// Generate slug from project name for non-interactive mode too
		slug := nameToSlug(projectName)
		cmd.Set("slug", slug)
	}

	// Inject file contents into the API payload if files are provided or found
	if openAPISpec != "" {
		content, err := os.ReadFile(openAPISpec)
		if err == nil {
			// Inject the actual file content into the project creation payload
			jsonflag.Mutate(jsonflag.Body, "revision.openapi\\.yml.content", string(content))
		}
	}

	if stainlessConfig != "" {
		content, err := os.ReadFile(stainlessConfig)
		if err == nil {
			// Inject the actual file content into the project creation payload
			jsonflag.Mutate(jsonflag.Body, "revision.openapi\\.stainless\\.yml.content", string(content))
		}
	}

	params := stainless.ProjectNewParams{}
	res, err := cc.client.Projects.New(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	if !allValuesProvided {
		Success("Project created successfully")
	}
	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))

	// Initialize workspace if requested
	if initWorkspace {
		Info("Initializing workspace configuration...")

		// Use the same project name (slug) for workspace initialization
		slug := nameToSlug(projectName)
		config, err := NewWorkspaceConfig(slug, openAPISpec, stainlessConfig)
		if err != nil {
			Error("Failed to create workspace config: %v", err)
			return fmt.Errorf("project created but workspace initialization failed: %v", err)
		}

		err = config.Save()
		if err != nil {
			Error("Failed to save workspace config: %v", err)
			return fmt.Errorf("project created but workspace initialization failed: %v", err)
		}

		Success("Workspace initialized")
	}

	// Download project configuration if requested
	if downloadConfig {
		Info("Downloading project configuration...")

		// Use the same project name (slug) for config download
		slug := nameToSlug(projectName)

		params := stainless.ProjectConfigGetParams{
			Project: stainless.String(slug),
		}

		configData := []byte{}
		_, err := cc.client.Projects.Configs.Get(
			ctx,
			params,
			option.WithMiddleware(cc.AsMiddleware()),
			option.WithResponseBodyInto(&configData),
		)
		if err != nil {
			Error("Failed to download project config: %v", err)
			return fmt.Errorf("project created but config download failed: %v", err)
		}

		// Determine where to save the config file
		var configPath string

		if stainlessConfig != "" {
			configPath = stainlessConfig
		}
		if configPath == "" {
			configPath = "stainless.yml"
		}

		// Write the config to file
		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			Error("Failed to save project config to %s: %v", configPath, err)
			return fmt.Errorf("project created but config save failed: %v", err)
		}

		Success("Project configuration downloaded to %s", configPath)
	}

	return nil
}

func handleProjectsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectGetParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Get(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleProjectsUpdate(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectUpdateParams{}
	if cmd.IsSet("project") {
		params.Project = stainless.String(cmd.Value("project").(string))
	}
	res, err := cc.client.Projects.Update(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleProjectsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	params := stainless.ProjectListParams{}
	res, err := cc.client.Projects.List(
		context.TODO(),
		params,
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ColorizeJSON(res.RawJSON(), os.Stdout))
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
