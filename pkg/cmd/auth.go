// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var authLogin = cli.Command{
	Name:  "login",
	Usage: "Authenticate with Stainless API",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "client-id",
			Value: "stl_client_0001u04Vo1IWoSe0Mwinw2SVuuO3hTkvL",
			Usage: "OAuth client ID",
		},
	},
	Action:          handleAuthLogin,
	HideHelpCommand: true,
}

var authLogout = cli.Command{
	Name:            "logout",
	Usage:           "Log out and remove saved credentials",
	Action:          handleAuthLogout,
	HideHelpCommand: true,
}

var authStatus = cli.Command{
	Name:            "status",
	Usage:           "Check authentication status",
	Action:          handleAuthStatus,
	HideHelpCommand: true,
}

func handleAuthLogin(ctx context.Context, cmd *cli.Command) error {
	clientID := cmd.String("client-id")
	scope := "openapi:read project:write project:read"
	config, err := StartDeviceFlow(clientID, scope)
	if err != nil {
		return err
	}
	if err := SaveAuthConfig(config); err != nil {
		return fmt.Errorf("%s", au.Red(fmt.Sprintf("Failed to save authentication: %v", err)))
	}
	fmt.Printf("%s %s\n", au.BrightGreen("✱"), "Authentication successful! Your credentials have been saved.")
	return nil
}

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
func SaveAuthConfig(config *AuthConfig) error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "auth.json")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func handleAuthLogout(ctx context.Context, cmd *cli.Command) error {
	configDir, err := ConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "auth.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println(au.BrightYellow("No active session found."))
		return nil
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("failed to remove auth file: %v", err)
	}

	fmt.Printf("%s %s\n", au.BrightGreen("✱"), "Successfully logged out.")
	return nil
}

func handleAuthStatus(ctx context.Context, cmd *cli.Command) error {
	// Check for API key in environment variables first
	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey != "" {
		fmt.Printf("%s %s\n", au.BrightGreen("✱"), "Authenticated via STAINLESS_API_KEY environment variable")
		return nil
	}

	// Check for token in config file
	config, err := LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to load auth config: %v", err)
	}

	if config == nil {
		fmt.Printf("%s %s\n", au.BrightYellow("✱"), "Not logged in.")
		return nil
	}

	// If we have a config file with a token
	fmt.Printf("%s %s\n", au.BrightGreen("✱"), "Authenticated via saved credentials")

	// Show a truncated version of the token for verification
	if len(config.AccessToken) > 10 {
		truncatedToken := config.AccessToken[:5] + "..." + config.AccessToken[len(config.AccessToken)-5:]
		fmt.Printf("Token: %s\n", truncatedToken)
	}

	return nil
}

// StartDeviceFlow initiates the OAuth 2.0 device flow
func StartDeviceFlow(clientID, scope string) (*AuthConfig, error) {
	deviceEndpoint := "https://api.stainless.com/api/oauth/device"
	payload := map[string]string{
		"client_id": clientID,
		"scope":     scope,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create JSON payload: %v", err)
	}
	req, err := http.NewRequest("POST", deviceEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to initiate device flow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth: device flow initiation failed with status %d: %s", resp.StatusCode, body)
	}

	var deviceResponse struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deviceResponse); err != nil {
		return nil, fmt.Errorf("failed to parse device response: %v", err)
	}

	if err := browser.OpenURL(deviceResponse.VerificationURIComplete); err != nil {
		fmt.Println()
		fmt.Printf("To authenticate, visit %s\n", au.Hyperlink(deviceResponse.VerificationURI, deviceResponse.VerificationURI))
		fmt.Printf("and enter this code %s\n", au.Bold(deviceResponse.UserCode))
		fmt.Println()
		fmt.Printf("Or navigate to this URL %s\n", au.Hyperlink(deviceResponse.VerificationURIComplete, deviceResponse.VerificationURIComplete))
	} else {
		fmt.Printf("Browser opened to %s\n", au.Hyperlink(deviceResponse.VerificationURIComplete, deviceResponse.VerificationURIComplete))
	}

	return pollForToken(
		clientID,
		deviceResponse.DeviceCode,
		// Hard-code to 1 second for now, instead of using deviceResponse.Interval
		1,
		deviceResponse.ExpiresIn,
	)
}

// pollForToken polls the token endpoint until the user completes authentication
func pollForToken(clientID, deviceCode string, interval, expiresIn int) (*AuthConfig, error) {
	tokenEndpoint := "https://api.stainless.com/v0/oauth/token"

	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second

	fmt.Println(au.BrightBlack("Waiting for authentication to complete..."))

	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)

		data := url.Values{}
		data.Set("client_id", clientID)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("auth: failed to poll for token: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// If we got a successful response, parse the token
		if resp.StatusCode == http.StatusOK {
			var tokenResponse struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				TokenType    string `json:"token_type"`
			}

			if err := json.Unmarshal(body, &tokenResponse); err != nil {
				return nil, fmt.Errorf("auth: failed to parse token response: %v", err)
			}

			return &AuthConfig{
				AccessToken:  tokenResponse.AccessToken,
				RefreshToken: tokenResponse.RefreshToken,
				TokenType:    tokenResponse.TokenType,
			}, nil
		}

		// If we got an error, check if it's "authorization_pending" and continue polling
		var errorResponse struct {
			Error struct {
				Error string `json:"error"`
			} `json:"error"`
		}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("%s", fmt.Sprintf("could not parse authentication error %d: %s", resp.StatusCode, err.Error()))
		}
		if errorResponse.Error.Error == "authorization_pending" {
			// This is expected, continue polling
			continue
		}
		return nil, fmt.Errorf("auth: %s", errorResponse.Error.Error)
	}
	return nil, fmt.Errorf("auth: timed out")
}

// GetClientOptions returns the request options for API calls
func getClientOptions() []option.RequestOption {
	options := []option.RequestOption{}

	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey != "" {
		return options
	}

	// Add authentication if available
	config, err := LoadAuthConfig()
	if err == nil && config != nil {
		options = append(options, option.WithAPIKey(config.AccessToken))
	}

	// Add default project from workspace config if available
	var workspaceConfig WorkspaceConfig
	found, err := workspaceConfig.Find()
	if err == nil && found && workspaceConfig.Project != "" {
		options = append(options, option.WithProject(workspaceConfig.Project))
	}

	return options
}
