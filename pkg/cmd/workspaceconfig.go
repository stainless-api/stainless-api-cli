package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
