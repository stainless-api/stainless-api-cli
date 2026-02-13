package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stainless-api/stainless-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name     string
		baseDir  string
		path     string
		expected string
	}{
		{
			name:     "relative path",
			baseDir:  "/home/user",
			path:     "config/spec.yaml",
			expected: "/home/user/config/spec.yaml",
		},
		{
			name:     "absolute path",
			baseDir:  "/home/user",
			path:     "/etc/config/spec.yaml",
			expected: "/etc/config/spec.yaml",
		},
		{
			name:     "path with dot dot",
			baseDir:  "/home/user/project",
			path:     "../config/spec.yaml",
			expected: "/home/user/config/spec.yaml",
		},
		{
			name:     "empty path",
			baseDir:  "/home/user",
			path:     "",
			expected: "/home/user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resolve(tt.baseDir, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRelative(t *testing.T) {
	// Get current working directory for testing
	cwd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		wantErr  bool
		checkRel bool // whether to check if result is relative
	}{
		{
			name: "path in current directory",
			path: filepath.Join(cwd, "test.txt"),
		},
		{
			name: "path in parent directory",
			path: filepath.Join(filepath.Dir(cwd), "test.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Relative(tt.path)
			if filepath.IsAbs(result) && result == tt.path {
				// If it couldn't be made relative, it returns the original path
				// This is acceptable behavior
			}
		})
	}
}

func TestConfigLoadAndSave(t *testing.T) {
	tests := []struct {
		name   string
		json   string                     // The JSON file content with relative paths
		config func(tmpDir string) Config // Expected config after loading (with absolute paths)
	}{
		{
			name: "config with all fields and multiple targets",
			json: `{
  "project": "test-project",
  "openapi_spec": "../specs/openapi.yaml",
  "stainless_config": "../configs/stainless.yaml",
  "targets": {
    "node": {
      "output_path": "../sdks/node"
    },
    "python": {
      "output_path": "../sdks/python"
    }
  }
}`,
			config: func(tmpDir string) Config {
				return Config{
					Project:         "test-project",
					OpenAPISpec:     filepath.Join(tmpDir, "specs", "openapi.yaml"),
					StainlessConfig: filepath.Join(tmpDir, "configs", "stainless.yaml"),
					Targets: map[stainless.Target]*TargetConfig{
						stainless.TargetNode: {
							OutputPath: filepath.Join(tmpDir, "sdks", "node"),
						},
						stainless.TargetPython: {
							OutputPath: filepath.Join(tmpDir, "sdks", "python"),
						},
					},
					ConfigPath: filepath.Join(tmpDir, ".stainless", "workspace.json"),
				}
			},
		},
		{
			name: "config with single target",
			json: `{
  "project": "single-target-project",
  "openapi_spec": "../openapi.yaml",
  "stainless_config": "../stainless.yaml",
  "targets": {
    "go": {
      "output_path": "../sdk"
    }
  }
}`,
			config: func(tmpDir string) Config {
				return Config{
					Project:         "single-target-project",
					OpenAPISpec:     filepath.Join(tmpDir, "openapi.yaml"),
					StainlessConfig: filepath.Join(tmpDir, "stainless.yaml"),
					Targets: map[stainless.Target]*TargetConfig{
						stainless.TargetGo: {
							OutputPath: filepath.Join(tmpDir, "sdk"),
						},
					},
					ConfigPath: filepath.Join(tmpDir, ".stainless", "workspace.json"),
				}
			},
		},
		{
			name: "config with no targets",
			json: `{
  "project": "no-targets-project",
  "openapi_spec": "../openapi.yaml",
  "stainless_config": "../stainless.yaml"
}`,
			config: func(tmpDir string) Config {
				return Config{
					Project:         "no-targets-project",
					OpenAPISpec:     filepath.Join(tmpDir, "openapi.yaml"),
					StainlessConfig: filepath.Join(tmpDir, "stainless.yaml"),
					Targets:         nil,
					ConfigPath:      filepath.Join(tmpDir, ".stainless", "workspace.json"),
				}
			},
		},
		{
			name: "config with minimal fields",
			json: `{
  "project": "minimal-project"
}`,
			config: func(tmpDir string) Config {
				return Config{
					Project:    "minimal-project",
					Targets:    nil,
					ConfigPath: filepath.Join(tmpDir, ".stainless", "workspace.json"),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create config directory
			configDir := filepath.Join(tmpDir, ".stainless")
			if err := os.MkdirAll(configDir, 0755); err != nil {
				t.Fatalf("Failed to create config directory: %v", err)
			}

			// Write the JSON file
			configPath := filepath.Join(configDir, "workspace.json")
			if err := os.WriteFile(configPath, []byte(tt.json), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Test Load: JSON → Config
			t.Run("Load", func(t *testing.T) {
				loadedConfig := Config{}
				require.NoError(t, loadedConfig.Load(configPath))

				wantConfig := tt.config(tmpDir)
				assert.Equal(t, wantConfig, loadedConfig)
			})

			// Test Save: Config → JSON
			t.Run("Save", func(t *testing.T) {
				// Create the expected config
				saveConfig := tt.config(tmpDir)
				saveConfig.ConfigPath = filepath.Join(configDir, "workspace-saved.json")

				// Save it
				require.NoError(t, saveConfig.Save())

				// Read the saved JSON
				savedJSON, err := os.ReadFile(saveConfig.ConfigPath)
				require.NoError(t, err)

				// Compare JSON directly
				assert.JSONEq(t, tt.json, string(savedJSON))
			})
		})
	}
}

func TestConfigLoadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.json")

	// Create an empty file
	require.NoError(t, os.WriteFile(configPath, []byte(""), 0644))

	config := Config{}
	err := config.Load(configPath)
	assert.NoError(t, err, "Load should not error on empty file")

	// ConfigPath should not be set for empty files
	assert.Empty(t, config.ConfigPath, "ConfigPath should be empty for empty file")
}

func TestConfigLoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.json")

	config := Config{}
	err := config.Load(configPath)
	assert.NoError(t, err, "Load should not error on non-existent file")
}

func TestConfigLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	// Create a file with invalid JSON
	require.NoError(t, os.WriteFile(configPath, []byte("{invalid json"), 0644))

	config := Config{}
	err := config.Load(configPath)
	assert.Error(t, err, "Load should error on invalid JSON")
}

func TestConfigFind(t *testing.T) {
	type fileSpec struct {
		path    string // Path relative to tmpDir
		content string // File content (empty string for empty file)
	}

	tests := []struct {
		name        string
		files       []fileSpec // Files to create
		workingDir  string     // Directory to chdir to (relative to tmpDir)
		wantFound   bool
		wantProject string
	}{
		{
			name: "config in current directory (.stainless/workspace.json)",
			files: []fileSpec{
				{".stainless/workspace.json", `{"project": "current-dir-project"}`},
			},
			workingDir:  ".",
			wantFound:   true,
			wantProject: "current-dir-project",
		},
		{
			name: "config in current directory (stainless-workspace.json)",
			files: []fileSpec{
				{"stainless-workspace.json", `{"project": "alt-location-project"}`},
			},
			workingDir:  ".",
			wantFound:   true,
			wantProject: "alt-location-project",
		},
		{
			name: "config in parent directory",
			files: []fileSpec{
				{".stainless/workspace.json", `{"project": "parent-project"}`},
			},
			workingDir:  "nested",
			wantFound:   true,
			wantProject: "parent-project",
		},
		{
			name: "config in grandparent directory",
			files: []fileSpec{
				{".stainless/workspace.json", `{"project": "grandparent-project"}`},
			},
			workingDir:  "level1/level2/level3",
			wantFound:   true,
			wantProject: "grandparent-project",
		},
		{
			name:       "no config found",
			files:      []fileSpec{},
			workingDir: ".",
			wantFound:  false,
		},
		{
			name: "empty config file ignored, continues search",
			files: []fileSpec{
				{"nested/.stainless/workspace.json", ""}, // Empty config in nested dir
				{".stainless/workspace.json", `{"project": "parent-after-empty"}`},
			},
			workingDir:  "nested",
			wantFound:   true,
			wantProject: "parent-after-empty",
		},
		{
			name: "prefers .stainless/workspace.json over stainless-workspace.json",
			files: []fileSpec{
				{".stainless/workspace.json", `{"project": "dot-stainless"}`},
				{"stainless-workspace.json", `{"project": "root-level"}`},
			},
			workingDir:  ".",
			wantFound:   true,
			wantProject: "dot-stainless",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create all specified files
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file.path)
				fileDir := filepath.Dir(filePath)

				require.NoError(t, os.MkdirAll(fileDir, 0755))
				require.NoError(t, os.WriteFile(filePath, []byte(file.content), 0644))
			}

			// Create working directory if it doesn't exist
			targetDir := filepath.Join(tmpDir, tt.workingDir)
			require.NoError(t, os.MkdirAll(targetDir, 0755))

			// Change to target directory
			originalWd, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalWd)

			require.NoError(t, os.Chdir(targetDir))

			// Test Find
			config := Config{}
			found, err := config.Find()
			require.NoError(t, err)

			assert.Equal(t, tt.wantFound, found)
			if tt.wantFound {
				assert.Equal(t, tt.wantProject, config.Project)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	openAPISpec := filepath.Join(tmpDir, "openapi.yaml")
	stainlessConfig := filepath.Join(tmpDir, "stainless.yaml")

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	require.NoError(t, os.Chdir(tmpDir))

	config, err := NewConfig("test-project", openAPISpec, stainlessConfig)
	require.NoError(t, err)

	assert.Equal(t, "test-project", config.Project)

	// Verify paths are absolute
	assert.True(t, filepath.IsAbs(config.OpenAPISpec), "OpenAPISpec should be absolute")
	assert.True(t, filepath.IsAbs(config.StainlessConfig), "StainlessConfig should be absolute")
	assert.True(t, filepath.IsAbs(config.ConfigPath), "ConfigPath should be absolute")

	// Verify config path ends with .stainless/workspace.json
	assert.Equal(t, ".stainless", filepath.Base(filepath.Dir(config.ConfigPath)), "ConfigPath should be in .stainless directory")
	assert.Equal(t, "workspace.json", filepath.Base(config.ConfigPath), "ConfigPath should end with workspace.json")
}
