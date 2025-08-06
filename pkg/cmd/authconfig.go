package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AuthConfig stores the OAuth credentials
type AuthConfig struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
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

// AuthConfigPath returns the path to the auth config file
func AuthConfigPath() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "auth.json"), nil
}

// LoadAuthConfig loads the auth config from disk
func LoadAuthConfig() (*AuthConfig, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "auth.json")
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var config AuthConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveAuthConfig saves the auth config to disk
func SaveAuthConfig(config *AuthConfig, configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// RemoveAuthConfig removes the auth config file
func RemoveAuthConfig() error {
	configPath, err := AuthConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // Already removed
	}

	return os.Remove(configPath)
}

// HasAuthConfig checks if an auth config file exists
func HasAuthConfig() bool {
	configPath, err := AuthConfigPath()
	if err != nil {
		return false
	}
	
	_, err = os.Stat(configPath)
	return !os.IsNotExist(err)
}