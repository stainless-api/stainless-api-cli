package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AuthConfig stores the OAuth credentials
type AuthConfig struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`

	ConfigPath string `json:"-"`
}

// ConfigDir returns the directory where config files are stored
func ConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".config", "stainless")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

// Find searches for and loads the auth config from the standard location.
// Returns (true, nil) if config file exists and was successfully loaded.
// Returns (false, nil) if config file doesn't exist or is empty (not an error).
// Returns (false, error) if config file exists but failed to load due to an error.
func (config *AuthConfig) Find() (bool, error) {
	if config.ConfigPath != "" {
		return true, nil
	}

	configDir, err := ConfigDir()
	if err != nil {
		return false, fmt.Errorf("failed to get config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "auth.json")
	if _, err := os.Stat(configPath); err == nil {
		// Config file exists, attempt to load it
		err := config.Load(configPath)
		if err != nil {
			return false, err
		}
		// Check if the config was actually loaded (ConfigPath is only set for valid configs)
		if config.ConfigPath != "" {
			return true, nil
		}
		// File exists but is empty or invalid - treat as not found (not an error)
	}

	// Config file doesn't exist - this is not an error
	return false, nil
}

// Load loads the auth config from a specific path.
// Returns nil if the file doesn't exist (not treated as an error).
// Returns nil if the file exists but is empty (not treated as an error).
// Returns error only if the file exists but fails to parse or read.
// Only sets ConfigPath if a valid config with AccessToken is successfully loaded.
func (config *AuthConfig) Load(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist - this is not an error, just means no auth config
			return nil
		}
		return fmt.Errorf("failed to open auth config file %s: %w", configPath, err)
	}
	defer file.Close()

	// Check if file is empty
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", configPath, err)
	}
	if info.Size() == 0 {
		// File exists but is empty - this is not an error, treat as no auth config
		return nil
	}

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return fmt.Errorf("failed to parse auth config file %s: %w", configPath, err)
	}

	// Only set ConfigPath if we successfully loaded a config with an access token
	if config.AccessToken != "" {
		config.ConfigPath = configPath
	}
	return nil
}

// Save saves the auth config to disk
func (config *AuthConfig) Save() error {
	if config.ConfigPath == "" {
		return fmt.Errorf("no config path set")
	}

	file, err := os.Create(config.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create auth config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// Remove removes the auth config file
func (config *AuthConfig) Remove() error {
	if config.ConfigPath == "" {
		return fmt.Errorf("config has not been loaded yet")
	}

	// File doesn't actually exist
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(config.ConfigPath)
}

// Exists checks if the auth config file exists
func (config *AuthConfig) Exists() bool {
	if config.ConfigPath == "" {
		return false
	}
	_, err := os.Stat(config.ConfigPath)
	return !os.IsNotExist(err)
}

// NewAuthConfig creates a new AuthConfig with ConfigPath populated.
// Use this when creating a new config that you plan to save.
// For loading existing configs, use &AuthConfig{} and call Find() or Load().
func NewAuthConfig() (*AuthConfig, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	return &AuthConfig{
		ConfigPath: filepath.Join(configDir, "auth.json"),
	}, nil
}
