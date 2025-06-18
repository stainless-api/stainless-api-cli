package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WorkspaceConfig stores workspace-level configuration
type WorkspaceConfig struct {
	Project         string `json:"project"`
	OpenAPISpec     string `json:"openapi_spec,omitempty"`
	StainlessConfig string `json:"stainless_config,omitempty"`

	ConfigPath string `json:"-"`
}

// Find searches for a stainless-workspace.json file starting from the current directory
// and moving up to parent directories until found or root is reached
func (config *WorkspaceConfig) Find() (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	for {
		configPath := filepath.Join(dir, "stainless-workspace.json")
		if _, err := os.Stat(configPath); err == nil {
			// Found config file
			err := config.Load(configPath)
			if err != nil {
				return false, err
			}
			return true, nil
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
		return err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return err
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

func NewWorkspaceConfig(projectName, openAPISpec, stainlessConfig string) (*WorkspaceConfig, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &WorkspaceConfig{
		Project:         projectName,
		OpenAPISpec:     openAPISpec,
		StainlessConfig: stainlessConfig,
		ConfigPath:      filepath.Join(dir, "stainless-workspace.json"),
	}, nil
}
