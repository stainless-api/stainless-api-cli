package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthLogin(t *testing.T) {
	tokenCallCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/oauth/device":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"device_code":               "test_device_code",
				"user_code":                 "TEST-CODE",
				"verification_uri":          "https://example.com/activate",
				"verification_uri_complete": "https://example.com/activate?code=TEST-CODE",
				"expires_in":               300,
				"interval":                 1,
			})
		case "/v0/oauth/token":
			tokenCallCount++
			if tokenCallCount == 1 {
				// First poll returns authorization_pending
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "authorization_pending",
				})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"access_token":  "test_access_token",
				"refresh_token": "test_refresh_token",
				"token_type":    "bearer",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Use a temp dir as HOME so auth config doesn't touch real config
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("XDG_CONFIG_HOME", "")
	// Ensure no env API key so the device flow is triggered
	t.Setenv("STAINLESS_API_KEY", "")

	err := Command.Run(context.Background(), []string{
		"stl", "auth", "login",
		"--base-url", server.URL,
		"--browser=false",
	})
	require.NoError(t, err)

	// Verify the auth config was saved
	configPath := filepath.Join(tmpHome, ".config", "stainless", "auth.json")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err, "auth config file should exist")

	var saved AuthConfig
	require.NoError(t, json.Unmarshal(data, &saved))
	assert.Equal(t, "test_access_token", saved.AccessToken)
	assert.Equal(t, "test_refresh_token", saved.RefreshToken)
	assert.Equal(t, "bearer", saved.TokenType)

	// Verify the token endpoint was polled at least twice (once pending, once success)
	assert.GreaterOrEqual(t, tokenCallCount, 2)
}

func TestAuthLoad(t *testing.T) {
	var gotAuthHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"email":"test@example.com"}`))
	}))
	defer server.Close()

	// Set up a temp HOME with a pre-existing auth config
	tmpHome := t.TempDir()
	configDir := filepath.Join(tmpHome, ".config", "stainless")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configData, _ := json.Marshal(AuthConfig{
		AccessToken:  "saved_test_token",
		RefreshToken: "saved_refresh_token",
		TokenType:    "bearer",
	})
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "auth.json"), configData, 0644))

	t.Setenv("HOME", tmpHome)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("STAINLESS_API_KEY", "")

	err := Command.Run(context.Background(), []string{
		"stl", "user", "retrieve",
		"--base-url", server.URL,
	})
	require.NoError(t, err)

	assert.Equal(t, "Bearer saved_test_token", gotAuthHeader)
}
