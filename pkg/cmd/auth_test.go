package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mockstainless"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthLogin(t *testing.T) {
	mock := mockstainless.NewMock(mockstainless.WithDeviceAuth(1))
	server := httptest.NewServer(mock.Server())
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
	assert.Equal(t, "demo_access_token_xyz789", saved.AccessToken)
	assert.Equal(t, "demo_refresh_token_abc456", saved.RefreshToken)
	assert.Equal(t, "bearer", saved.TokenType)
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
