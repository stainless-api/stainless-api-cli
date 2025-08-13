package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stainless-api/stainless-api-go"
)

// TargetConfig stores configuration for a specific SDK target
type TargetConfig struct {
	OutputPath string `json:"output_path"`
}

// WorkspaceConfig stores workspace-level configuration
type WorkspaceConfig struct {
	Project         string                   `json:"project"`
	OpenAPISpec     string                   `json:"openapi_spec,omitempty"`
	StainlessConfig string                   `json:"stainless_config,omitempty"`
	Targets         map[string]*TargetConfig `json:"targets,omitempty"`

	ConfigPath string `json:"-"`
}

// Find searches for a stainless-workspace.json file starting from the current directory
// and moving up to parent directories until found or root is reached
func (config *WorkspaceConfig) Find() (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	if config.ConfigPath != "" {
		return true, nil
	}

	for {
		configPath := filepath.Join(dir, "stainless-workspace.json")
		if _, err := os.Stat(configPath); err == nil {
			// Found config file
			err := config.Load(configPath)
			if err != nil {
				return false, err
			}
			// Check if the config was actually loaded (not empty)
			if config.ConfigPath != "" {
				return true, nil
			}
			// File exists but is empty, continue searching
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// At root directory
			return false, nil
		}
		dir = parent
	}
}

func (config *WorkspaceConfig) Load(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open workspace config file %s: %w", configPath, err)
	}
	defer file.Close()

	// Check if file is empty
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", configPath, err)
	}
	if info.Size() == 0 {
		// File is empty, treat as if no config exists
		return nil
	}

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return fmt.Errorf("failed to parse workspace config file %s: %w", configPath, err)
	}
	config.ConfigPath = configPath
	return nil
}

func (config *WorkspaceConfig) Save() error {
	file, err := os.Create(config.ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func NewWorkspaceConfig(projectName, openAPISpec, stainlessConfig string) (WorkspaceConfig, error) {
	return NewWorkspaceConfigWithTargets(projectName, openAPISpec, stainlessConfig, nil)
}

func NewWorkspaceConfigWithTargets(projectName, openAPISpec, stainlessConfig string, targets map[string]*TargetConfig) (WorkspaceConfig, error) {
	dir, err := os.Getwd()
	if err != nil {
		return WorkspaceConfig{}, err
	}

	return WorkspaceConfig{
		Project:         projectName,
		OpenAPISpec:     openAPISpec,
		StainlessConfig: stainlessConfig,
		Targets:         targets,
		ConfigPath:      filepath.Join(dir, "stainless-workspace.json"),
	}, nil
}

type projectInfo struct {
	Name string
	Org  string
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
			return "./" + filename
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
			return "./" + filename
		}
	}
	return ""
}
