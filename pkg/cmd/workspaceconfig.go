package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stainless-api/stainless-api-go"
)

// Resolve converts a path to absolute, resolving it relative to baseDir if it's not already absolute
func Resolve(baseDir, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(baseDir, path))
}

func Relative(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}

	rel, err := filepath.Rel(cwd, path)
	if err != nil {
		return path
	}

	return rel
}

// TargetConfig stores configuration for a specific SDK target
type TargetConfig struct {
	OutputPath string `json:"output_path"`
}

// WorkspaceConfigExport represents the on-disk format with relative paths
type WorkspaceConfigExport struct {
	Project         string                             `json:"project"`
	OpenAPISpec     string                             `json:"openapi_spec,omitempty"`
	StainlessConfig string                             `json:"stainless_config,omitempty"`
	Targets         map[stainless.Target]*TargetConfig `json:"targets,omitempty"`
}

// WorkspaceConfig stores workspace-level configuration with absolute paths
type WorkspaceConfig struct {
	Project         string
	OpenAPISpec     string
	StainlessConfig string
	Targets         map[stainless.Target]*TargetConfig

	ConfigPath string
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
		for _, configPath := range []string{filepath.Join(dir, ".stainless", "workspace.json"), filepath.Join(dir, "stainless-workspace.json")} {
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

	// Load into export format (with relative paths)
	var export WorkspaceConfigExport
	if err := json.NewDecoder(file).Decode(&export); err != nil {
		return fmt.Errorf("failed to parse workspace config file %s: %w", configPath, err)
	}

	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)

	// Convert relative paths to absolute paths
	config.Project = export.Project
	config.ConfigPath = configPath

	if export.OpenAPISpec != "" {
		config.OpenAPISpec = Resolve(configDir, export.OpenAPISpec)
	}

	if export.StainlessConfig != "" {
		config.StainlessConfig = Resolve(configDir, export.StainlessConfig)
	}

	// Convert target paths to absolute
	if export.Targets != nil {
		config.Targets = make(map[stainless.Target]*TargetConfig, len(export.Targets))
		for target, targetConfig := range export.Targets {
			config.Targets[target] = &TargetConfig{
				OutputPath: Resolve(configDir, targetConfig.OutputPath),
			}
		}
	}

	return nil
}

func (config *WorkspaceConfig) Save() error {
	// Create parent directories if they don't exist
	configDir := filepath.Dir(config.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for config file: %w", err)
	}

	// Convert absolute paths to relative paths for export
	export := WorkspaceConfigExport{
		Project: config.Project,
	}

	// Convert paths to relative (fallback to absolute if conversion fails)
	if config.OpenAPISpec != "" {
		if relPath, err := filepath.Rel(configDir, config.OpenAPISpec); err == nil {
			export.OpenAPISpec = relPath
		} else {
			println(err.Error())
			export.OpenAPISpec = config.OpenAPISpec
		}
	}

	if config.StainlessConfig != "" {
		if relPath, err := filepath.Rel(configDir, config.StainlessConfig); err == nil {
			export.StainlessConfig = relPath
		} else {
			println(err.Error())
			export.StainlessConfig = config.StainlessConfig
		}
	}

	if config.Targets != nil {
		export.Targets = make(map[stainless.Target]*TargetConfig, len(config.Targets))
		for target, targetConfig := range config.Targets {
			outputPath := targetConfig.OutputPath
			if relPath, err := filepath.Rel(configDir, outputPath); err == nil {
				outputPath = relPath
			}
			export.Targets[target] = &TargetConfig{
				OutputPath: outputPath,
			}
		}
	}

	file, err := os.Create(config.ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(export)
}

func NewWorkspaceConfig(projectName, openAPISpecPath, stainlessConfigPath string) (WorkspaceConfig, error) {
	dir, err := os.Getwd()
	if err != nil {
		return WorkspaceConfig{}, err
	}

	// Convert paths to absolute
	absOpenAPISpec, err := filepath.Abs(openAPISpecPath)
	if err != nil {
		return WorkspaceConfig{}, fmt.Errorf("failed to get absolute path for OpenAPI spec: %w", err)
	}

	absStainlessConfig, err := filepath.Abs(stainlessConfigPath)
	if err != nil {
		return WorkspaceConfig{}, fmt.Errorf("failed to get absolute path for Stainless config: %w", err)
	}

	return WorkspaceConfig{
		Project:         projectName,
		OpenAPISpec:     absOpenAPISpec,
		StainlessConfig: absStainlessConfig,
		Targets:         nil,
		ConfigPath:      filepath.Join(dir, ".stainless", "workspace.json"),
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
