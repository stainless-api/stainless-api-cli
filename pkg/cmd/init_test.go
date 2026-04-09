package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestInitNonInteractive(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		// Flag configuration
		project    string
		targets    string // comma-separated, empty means omit flag
		oasFlag    string // empty means omit flag
		configFlag string // empty means omit flag

		isNewProject   bool
		expectError    bool
		expectErrorMsg string

		// Assertions on the created workspace
		wantTargets []string
	}

	cases := []testCase{
		// ── New project ──────────────────────────────────────────────
		{
			name:         "new project with targets and config",
			project:      "brand-new",
			targets:      "python,typescript",
			oasFlag:      "openapi.json",
			configFlag:   "stainless.yml",
			isNewProject: true,
			wantTargets:  []string{"python", "typescript"},
		},
		{
			name:         "new project with targets, no config",
			project:      "brand-new-no-cfg",
			targets:      "go",
			oasFlag:      "openapi.json",
			isNewProject: true,
			wantTargets:  []string{"go"},
		},
		{
			name:           "new project without targets fails",
			project:        "brand-new-no-tgt",
			oasFlag:        "openapi.json",
			isNewProject:   true,
			expectError:    true,
			expectErrorMsg: "--targets",
		},
		{
			name:           "new project without openapi-spec fails",
			project:        "brand-new-no-oas",
			targets:        "python",
			isNewProject:   true,
			expectError:    true,
			expectErrorMsg: "--openapi-spec",
		},
		// ── Existing project ─────────────────────────────────────────
		{
			name:        "existing project with all flags",
			project:     "acme-api",
			targets:     "python,typescript",
			oasFlag:     "openapi.json",
			configFlag:  "stainless.yml",
			wantTargets: []string{"python", "typescript"},
		},
		{
			name:        "existing project without targets uses server targets",
			project:     "acme-api",
			oasFlag:     "openapi.json",
			configFlag:  "stainless.yml",
			wantTargets: []string{"typescript", "python", "go"},
		},
		{
			name:        "existing project without config",
			project:     "acme-api",
			targets:     "typescript",
			oasFlag:     "openapi.json",
			wantTargets: []string{"typescript"},
		},
		{
			name:        "existing project with config, no explicit targets",
			project:     "acme-api",
			oasFlag:     "openapi.json",
			configFlag:  "stainless.yml",
			wantTargets: []string{"typescript", "python", "go"},
		},
		// ── Missing required flags ───────────────────────────────────
		{
			name:           "no project flag fails",
			oasFlag:        "openapi.json",
			expectError:    true,
			expectErrorMsg: "--project",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Each subtest gets its own mock server to avoid shared state.
			server := newMockServer(t)
			dir := t.TempDir()

			// Create dummy spec/config files in the working directory.
			oasContent := `{"openapi":"3.1.0","info":{"title":"Test","version":"1.0.0"},"paths":{}}`
			require.NoError(t, os.WriteFile(filepath.Join(dir, "openapi.json"), []byte(oasContent), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(dir, "stainless.yml"), []byte("client:\n  name: Test\n"), 0644))

			args := []string{"init", "--api-key", "test-key"}
			if tc.project != "" {
				args = append(args, "--project", tc.project)
			}
			if tc.targets != "" {
				args = append(args, "--targets", tc.targets)
			}
			if tc.oasFlag != "" {
				args = append(args, "--openapi-spec", tc.oasFlag)
			}
			if tc.configFlag != "" {
				args = append(args, "--stainless-config", tc.configFlag)
			}

			output := runCLIWithExpectation(t, dir, server.URL(), tc.expectError, args...)

			if tc.expectError {
				if tc.expectErrorMsg != "" {
					assert.Contains(t, output, tc.expectErrorMsg)
				}
				return
			}

			// Verify .stainless/workspace.json was created.
			wsPath := filepath.Join(dir, ".stainless", "workspace.json")
			data, err := os.ReadFile(wsPath)
			require.NoError(t, err, "workspace.json should exist")

			ws := gjson.ParseBytes(data)
			assert.Equal(t, tc.project, ws.Get("project").String(), "workspace project")

			// Paths in workspace.json are relative to .stainless/, so
			// a file at dir/openapi.json becomes ../openapi.json.
			if tc.oasFlag != "" {
				assert.Equal(t, "../"+tc.oasFlag, ws.Get("openapi_spec").String(), "workspace openapi_spec")
			}
			if tc.configFlag != "" {
				assert.Equal(t, "../"+tc.configFlag, ws.Get("stainless_config").String(), "workspace stainless_config")
			}

			// Verify targets were configured.
			targets := ws.Get("targets")
			for _, want := range tc.wantTargets {
				assert.True(t, targets.Get(want).Exists(), "target %q should be configured", want)
			}

			// Verify correct API calls were made.
			if tc.isNewProject {
				req := findRequest(t, server.Requests(), "POST", "/v0/projects")
				assert.Equal(t, tc.project, gjson.Get(req.Body, "slug").String())
			} else {
				findRequest(t, server.Requests(), "GET", "/v0/projects")
			}
		})
	}
}

func TestInitNonInteractiveRequests(t *testing.T) {
	t.Parallel()

	t.Run("new project sends openapi spec content in revision", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		oasContent := `{"openapi":"3.1.0","info":{"title":"ReqTest","version":"1.0.0"}}`
		require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.json"), []byte(oasContent), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "req-test-project",
			"--targets", "python",
			"--openapi-spec", "spec.json",
		)

		req := findRequest(t, server.Requests(), "POST", "/v0/projects")
		assert.Equal(t, "req-test-project", gjson.Get(req.Body, "slug").String())
		assert.Equal(t, oasContent, gjson.Get(req.Body, "revision.openapi\\.json.content").String())
		assert.False(t, gjson.Get(req.Body, "revision.stainless\\.yml").Exists())
		assert.False(t, gjson.Get(req.Body, "revision.stainless\\.json").Exists())
	})

	t.Run("new project sends stainless config when flag provided", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		oasContent := `{"openapi":"3.1.0","info":{"title":"CfgTest","version":"1.0.0"}}`
		cfgContent := "client:\n  name: CfgTest\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.json"), []byte(oasContent), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "cfg.yml"), []byte(cfgContent), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "cfg-test-project",
			"--targets", "typescript",
			"--openapi-spec", "spec.json",
			"--stainless-config", "cfg.yml",
		)

		req := findRequest(t, server.Requests(), "POST", "/v0/projects")
		assert.Equal(t, oasContent, gjson.Get(req.Body, "revision.openapi\\.json.content").String())
		assert.Equal(t, cfgContent, gjson.Get(req.Body, "revision.stainless\\.yml.content").String())
	})

	t.Run("new project sends targets in create request", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.json"), []byte(`{}`), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "tgt-test-project",
			"--targets", "python,go",
			"--openapi-spec", "spec.json",
		)

		req := findRequest(t, server.Requests(), "POST", "/v0/projects")
		targetsResult := gjson.Get(req.Body, "targets")
		require.True(t, targetsResult.IsArray())
		var got []string
		for _, v := range targetsResult.Array() {
			got = append(got, v.String())
		}
		assert.ElementsMatch(t, []string{"python", "go"}, got)
	})

	t.Run("existing project does not POST to create", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.json"), []byte(`{}`), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "cfg.yml"), []byte("client:\n  name: E\n"), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "acme-api",
			"--openapi-spec", "spec.json",
			"--stainless-config", "cfg.yml",
		)

		assertNoRequest(t, server.Requests(), "POST", "/v0/projects")
	})
}

func TestInitNonInteractiveWorkspaceContents(t *testing.T) {
	t.Parallel()

	t.Run("workspace.json has correct structure", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "my-spec.json"), []byte(`{}`), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "my-config.yml"), []byte("x: 1\n"), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "acme-api",
			"--targets", "python,go",
			"--openapi-spec", "my-spec.json",
			"--stainless-config", "my-config.yml",
		)

		data, err := os.ReadFile(filepath.Join(dir, ".stainless", "workspace.json"))
		require.NoError(t, err)

		var ws map[string]any
		require.NoError(t, json.Unmarshal(data, &ws))
		assert.Equal(t, "acme-api", ws["project"])
		assert.Equal(t, "../my-spec.json", ws["openapi_spec"])
		assert.Equal(t, "../my-config.yml", ws["stainless_config"])

		targets, ok := ws["targets"].(map[string]any)
		require.True(t, ok, "targets should be a map")
		assert.Contains(t, targets, "python")
		assert.Contains(t, targets, "go")
	})

	t.Run("default stainless_config path when flag omitted", func(t *testing.T) {
		t.Parallel()
		server := newMockServer(t)
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.json"), []byte(`{}`), 0644))

		runCLI(t, dir, server.URL(), "init",
			"--api-key", "test-key",
			"--project", "acme-api",
			"--targets", "typescript",
			"--openapi-spec", "spec.json",
		)

		data, err := os.ReadFile(filepath.Join(dir, ".stainless", "workspace.json"))
		require.NoError(t, err)

		ws := gjson.ParseBytes(data)
		configPath := ws.Get("stainless_config").String()
		assert.NotEmpty(t, configPath, "stainless_config should have a default value")
	})
}
