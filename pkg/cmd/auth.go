// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
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
	cc := getAPICommandContext(cmd)
	clientID := cmd.String("client-id")
	scope := "openapi:read project:write project:read"
	config, err := startDeviceFlow(ctx, cc.client, clientID, scope)
	if err != nil {
		return err
	}
	if err := SaveAuthConfig(config); err != nil {
		Error("Failed to save authentication: %v", err)
		return fmt.Errorf("authentication failed")
	}
	Success("Authentication successful! Your credentials have been saved.")
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
		Warn("No active session found.")
		return nil
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("failed to remove auth file: %v", err)
	}

	Success("Successfully logged out.")
	return nil
}

func handleAuthStatus(ctx context.Context, cmd *cli.Command) error {
	// Check for API key in environment variables first
	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey != "" {
		Success("Authenticated via STAINLESS_API_KEY environment variable")
		return nil
	}

	// Check for token in config file
	config, err := LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to load auth config: %v", err)
	}

	if config == nil {
		Warn("Not logged in.")
		return nil
	}

	// If we have a config file with a token
	group := Success("Authenticated via saved credentials")

	// Show a truncated version of the token for verification
	if len(config.AccessToken) > 10 {
		truncatedToken := config.AccessToken[:5] + "..." + config.AccessToken[len(config.AccessToken)-5:]
		group.Property("token", truncatedToken)
	}

	return nil
}

// startDeviceFlow initiates the OAuth 2.0 device flow
func startDeviceFlow(ctx context.Context, client stainless.Client, clientID, scope string) (*AuthConfig, error) {
	var deviceResponse struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
	}

	err := client.Post(ctx, "api/oauth/device", map[string]string{
		"client_id": clientID,
		"scope":     scope,
	}, &deviceResponse)

	if err != nil {
		return nil, err
	}

	if err := browser.OpenURL(deviceResponse.VerificationURIComplete); err != nil {
		group := Info("To authenticate, visit the verification URL")
		group.Property("url", deviceResponse.VerificationURI)
		group.Property("code", deviceResponse.UserCode)
		group.Property("direct_url", deviceResponse.VerificationURIComplete)
	} else {
		group := Info("Browser opened")
		group.Property("url", deviceResponse.VerificationURIComplete)
	}

	return pollForToken(
		ctx,
		client,
		clientID,
		deviceResponse.DeviceCode,
		// Hard-code to 1 second for now, instead of using deviceResponse.Interval
		1,
		deviceResponse.ExpiresIn,
	)
}

// pollForToken polls the token endpoint until the user completes authentication
func pollForToken(ctx context.Context, client stainless.Client, clientID, deviceCode string, interval, expiresIn int) (*AuthConfig, error) {
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second

	Progress("Waiting for authentication to complete...")

	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)

		data := url.Values{}
		data.Set("client_id", clientID)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		var tokenResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			TokenType    string `json:"token_type"`
		}
		err := client.Post(ctx, "v0/oauth/token",
			strings.NewReader(data.Encode()),
			&tokenResponse,
			option.WithHeader("Content-Type", "application/x-www-form-urlencoded"),
		)

		if err != nil {
			var apierr *stainless.Error
			if errors.As(err, &apierr) {
				// If we got an error, check if it's "authorization_pending" and continue polling
				errorStr := gjson.Get(apierr.RawJSON(), "error.error").String()
				// This is expected, continue polling
				if errorStr == "authorization_pending" {
					continue
				}
			}

			return nil, fmt.Errorf("auth: %w", err)
		}

		return &AuthConfig{
			AccessToken:  tokenResponse.AccessToken,
			RefreshToken: tokenResponse.RefreshToken,
			TokenType:    tokenResponse.TokenType,
		}, nil
	}
	return nil, fmt.Errorf("auth: timed out")
}
