// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const defaultClientID = "stl_client_0001u04Vo1IWoSe0Mwinw2SVuuO3hTkvL"

var authLogin = cli.Command{
	Name:  "login",
	Usage: "Authenticate with Stainless API",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "client-id",
			Value: defaultClientID,
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

func authenticate(ctx context.Context, cmd *cli.Command, forceAuthentication bool) error {
	if apiKey := os.Getenv("STAINLESS_API_KEY"); apiKey != "" && !forceAuthentication {
		return nil
	}

	config, err := NewAuthConfig()
	if err != nil {
		Error("Failed to create config: %v", err)
		return fmt.Errorf("authentication failed")
	}

	if !forceAuthentication {
		if found, err := config.Find(); err == nil && found && config.AccessToken != "" {
			return nil
		}
	}

	cc := getAPICommandContext(cmd)
	clientID := cmd.String("client-id")
	scope := "*"
	authResult, err := startDeviceFlow(ctx, cmd, cc.client, clientID, scope)
	if err != nil {
		return err
	}

	config.AccessToken = authResult.AccessToken
	config.RefreshToken = authResult.RefreshToken
	config.TokenType = authResult.TokenType

	if err := config.Save(); err != nil {
		Error("Failed to save authentication: %v", err)
		return fmt.Errorf("authentication failed")
	}
	Success("Authentication successful! Your credentials have been saved to %s", config.ConfigPath)
	return nil
}

func handleAuthLogin(ctx context.Context, cmd *cli.Command) error {
	return authenticate(ctx, cmd, true)
}

func handleAuthLogout(ctx context.Context, cmd *cli.Command) error {
	config := &AuthConfig{}
	found, err := config.Find()
	if err != nil {
		return fmt.Errorf("failed to find auth config: %v", err)
	}

	if !found {
		Warn("No active session found.")
		return nil
	}

	if err := config.Remove(); err != nil {
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
	config := &AuthConfig{}
	found, err := config.Find()
	if err != nil {
		return fmt.Errorf("failed to find auth config: %v", err)
	}

	if !found || config.AccessToken == "" {
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
func startDeviceFlow(ctx context.Context, cmd *cli.Command, client stainless.Client, clientID, scope string) (*AuthConfig, error) {
	var deviceResponse struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
	}

	err := client.Post(ctx, "api/oauth/device", map[string]string{
		"client_id":   clientID,
		"scope":       scope,
		"device_name": getDeviceName(),
	}, &deviceResponse)

	if err != nil {
		return nil, err
	}

	group := Info("To authenticate, visit the verification URL")
	group.Property("url", deviceResponse.VerificationURIComplete)
	group.Property("code", deviceResponse.UserCode)

	ok, _, err := group.Confirm(cmd, "browser", "Open browser?", "", true)
	if err != nil {
		return nil, err
	} else if ok {
		if err := browser.OpenURL(deviceResponse.VerificationURIComplete); err == nil {
			group.Info("Opening browser...")
		} else {
			group.Warn("Could not open browser")
		}
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
				errorStr := gjson.Get(apierr.RawJSON(), "error").String()
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

// getDeviceName generates a descriptive device name for OAuth device flow
func getDeviceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown Device"
	}

	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows fallback
	}

	// Get OS name in a more human-readable format
	osName := runtime.GOOS
	switch osName {
	case "darwin":
		osName = "macOS"
	case "windows":
		osName = "Windows"
	case "linux":
		osName = "Linux"
	default:
		osName = cases.Title(language.English).String(osName)
	}

	if username != "" {
		return fmt.Sprintf("Stainless CLI on %s (%s, %s)", hostname, username, osName)
	}

	return fmt.Sprintf("Stainless CLI on %s (%s)", hostname, osName)
}
