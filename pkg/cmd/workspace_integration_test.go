package cmd

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stainless-api/stainless-api-cli/internal/mockstainless"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"slices"
)

type workspaceFixture struct {
	Dir                 string
	WorkspaceProject    string
	OverrideProject     string
	WorkspaceOASPath    string
	WorkspaceConfigPath string
	OverrideOASPath     string
	OverrideConfigPath  string
	WorkspaceOAS        string
	WorkspaceConfig     string
	OverrideOAS         string
	OverrideConfig      string
}

func TestWorkspaceProjectAutofillIntegration(t *testing.T) {
	t.Parallel()

	server := newMockServer(t)

	t.Run("workspace", func(t *testing.T) {
		fixture := newWorkspaceFixture(t)
		server.ResetRequests()
		runCLI(t, fixture.Dir, server.URL(), "projects", "retrieve", "--api-key", "string")

		request := findRequest(t, server.Requests(), "GET", "/v0/projects/"+fixture.WorkspaceProject)
		assertProjectPathSuffix(t, request, fixture.WorkspaceProject)
	})

	t.Run("flag override", func(t *testing.T) {
		fixture := newWorkspaceFixture(t)
		server.ResetRequests()
		runCLI(t, fixture.Dir, server.URL(), "projects", "retrieve", "--api-key", "string", "--project", fixture.OverrideProject)

		request := findRequest(t, server.Requests(), "GET", "/v0/projects/"+fixture.OverrideProject)
		assertProjectPathSuffix(t, request, fixture.OverrideProject)
	})

	t.Run("without workspace fails", func(t *testing.T) {
		server.ResetRequests()
		runCLIExpectError(t, t.TempDir(), server.URL(), "projects", "retrieve", "--api-key", "string")
	})
}

func TestWorkspaceFileAutofillIntegration(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		args         []string
		server       *mockServer
		requestIndex func(t *testing.T, requests []mockstainless.RecordedRequest) mockstainless.RecordedRequest
		assertFiles  func(t *testing.T, request mockstainless.RecordedRequest, fixture workspaceFixture)
		expectError  bool
	}

	cases := []testCase{
		{
			name:   "builds create",
			args:   []string{"builds", "create", "--api-key", "string", "--wait", "none"},
			server: newMockServer(t),
			requestIndex: func(t *testing.T, requests []mockstainless.RecordedRequest) mockstainless.RecordedRequest {
				return findRequest(t, requests, "POST", "/v0/builds")
			},
			assertFiles: assertBuildCreateFiles,
		},
		{
			name:   "lint",
			args:   []string{"lint", "--api-key", "string"},
			server: newMockServer(t),
			requestIndex: func(t *testing.T, requests []mockstainless.RecordedRequest) mockstainless.RecordedRequest {
				return findRequest(t, requests, "POST", "/api/generate/spec")
			},
			assertFiles: assertLintFiles,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run("workspace", func(t *testing.T) {
				fixture := newWorkspaceFixture(t)
				tc.server.ResetRequests()
				if tc.expectError {
					runCLIExpectError(t, fixture.Dir, tc.server.URL(), tc.args...)
				} else {
					runCLI(t, fixture.Dir, tc.server.URL(), tc.args...)
				}

				request := tc.requestIndex(t, tc.server.Requests())
				tc.assertFiles(t, request, fixture)
				assertProjectBody(t, request, fixture.WorkspaceProject)
			})

			t.Run("openapi override suppresses workspace config", func(t *testing.T) {
				fixture := newWorkspaceFixture(t)
				tc.server.ResetRequests()
				args := append(slices.Clone(tc.args), "--oas", fixture.OverrideOASPath)
				if tc.name == "lint" {
					runCLIExpectError(t, fixture.Dir, tc.server.URL(), args...)
					assertNoRequest(t, tc.server.Requests(), "POST", "/api/generate/spec")
					return
				}

				runCLI(t, fixture.Dir, tc.server.URL(), args...)
				request := tc.requestIndex(t, tc.server.Requests())
				assertProjectBody(t, request, fixture.WorkspaceProject)
				assertOpenAPIOnly(t, request, fixture.OverrideOAS)
			})

			t.Run("config override suppresses workspace openapi", func(t *testing.T) {
				fixture := newWorkspaceFixture(t)
				tc.server.ResetRequests()
				args := append(slices.Clone(tc.args), "--config", fixture.OverrideConfigPath)
				if tc.name == "lint" {
					runCLIExpectError(t, fixture.Dir, tc.server.URL(), args...)
					assertNoRequest(t, tc.server.Requests(), "POST", "/api/generate/spec")
					return
				}

				runCLI(t, fixture.Dir, tc.server.URL(), args...)
				request := tc.requestIndex(t, tc.server.Requests())
				assertProjectBody(t, request, fixture.WorkspaceProject)
				assertConfigOnly(t, request, fixture.OverrideConfig)
			})

			t.Run("alias flags override workspace", func(t *testing.T) {
				fixture := newWorkspaceFixture(t)
				tc.server.ResetRequests()
				args := append(slices.Clone(tc.args),
					"--project", fixture.OverrideProject,
					"--oas", fixture.OverrideOASPath,
					"--config", fixture.OverrideConfigPath,
				)
				if tc.expectError {
					runCLIExpectError(t, fixture.Dir, tc.server.URL(), args...)
				} else {
					runCLI(t, fixture.Dir, tc.server.URL(), args...)
				}

				request := tc.requestIndex(t, tc.server.Requests())
				tc.assertFiles(t, request, workspaceFixture{
					WorkspaceProject: fixture.OverrideProject,
					WorkspaceOAS:     fixture.OverrideOAS,
					WorkspaceConfig:  fixture.OverrideConfig,
				})
				assertProjectBody(t, request, fixture.OverrideProject)
			})

			t.Run("long flags override workspace", func(t *testing.T) {
				fixture := newWorkspaceFixture(t)
				tc.server.ResetRequests()
				args := append(slices.Clone(tc.args),
					"--project", fixture.OverrideProject,
					"--openapi-spec", fixture.OverrideOASPath,
					"--stainless-config", fixture.OverrideConfigPath,
				)
				if tc.expectError {
					runCLIExpectError(t, fixture.Dir, tc.server.URL(), args...)
				} else {
					runCLI(t, fixture.Dir, tc.server.URL(), args...)
				}

				request := tc.requestIndex(t, tc.server.Requests())
				tc.assertFiles(t, request, workspaceFixture{
					WorkspaceProject: fixture.OverrideProject,
					WorkspaceOAS:     fixture.OverrideOAS,
					WorkspaceConfig:  fixture.OverrideConfig,
				})
				assertProjectBody(t, request, fixture.OverrideProject)
			})

			t.Run("without workspace fails", func(t *testing.T) {
				tc.server.ResetRequests()
				runCLIExpectError(t, t.TempDir(), tc.server.URL(), tc.args...)
			})
		})
	}
}

type mockServer struct {
	server *httptest.Server
	mock   *mockstainless.Mock
}

type mockOption func(*mockstainless.Mock)

var (
	buildCLIOnce sync.Once
	buildCLIPath string
	buildCLIErr  error
)

var ()

func newMockServer(t *testing.T, opts ...mockOption) *mockServer {
	t.Helper()

	mock := mockstainless.NewMock(
		mockstainless.WithDefaultOrg(),
		mockstainless.WithDefaultProject(),
		mockstainless.WithDefaultCompareBuild(),
	)
	for _, opt := range opts {
		opt(mock)
	}

	server := httptest.NewServer(mock.Server())
	t.Cleanup(func() {
		server.Close()
		mock.Cleanup()
	})

	return &mockServer{server: server, mock: mock}
}

func (s *mockServer) URL() string {
	return s.server.URL
}

func (s *mockServer) ResetRequests() {
	s.mock.ResetRequests()
}

func (s *mockServer) Requests() []mockstainless.RecordedRequest {
	return s.mock.Requests()
}

func runCLI(t *testing.T, dir string, baseURL string, args ...string) string {
	t.Helper()
	return runCLIWithExpectation(t, dir, baseURL, false, args...)
}

func runCLIExpectError(t *testing.T, dir string, baseURL string, args ...string) string {
	t.Helper()
	return runCLIWithExpectation(t, dir, baseURL, true, args...)
}

func runCLIWithExpectation(t *testing.T, dir string, baseURL string, expectError bool, args ...string) string {
	t.Helper()

	binary := buildCLI(t)
	commandArgs := append([]string{"--base-url", baseURL}, args...)
	t.Logf("Testing command: %s %s", binary, strings.Join(commandArgs, " "))

	cmd := exec.Command(binary, commandArgs...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TERM=dumb", "CI=1")
	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)

	if expectError {
		assert.Error(t, err, "Expected command to fail\nOutput: %s", output)
	} else {
		assert.NoError(t, err, "Test failed\nError: %v\nOutput: %s", err, output)
	}

	return output
}

func buildCLI(t *testing.T) string {
	t.Helper()

	buildCLIOnce.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			buildCLIErr = assert.AnError
			return
		}

		repoRoot := filepath.Join(filepath.Dir(filename), "..", "..")
		buildCLIPath = filepath.Join(os.TempDir(), "stl-workspace-integration-bin")

		cmd := exec.Command("go", "build", "-o", buildCLIPath, "./cmd/stl")
		cmd.Dir = repoRoot
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("go build output:\n%s", output)
			buildCLIErr = err
		}
	})

	require.NoError(t, buildCLIErr)
	return buildCLIPath
}

func newWorkspaceFixture(t *testing.T) workspaceFixture {
	t.Helper()

	dir := t.TempDir()
	fixture := workspaceFixture{
		Dir:                 dir,
		WorkspaceProject:    "workspace-project",
		OverrideProject:     "flag-project",
		WorkspaceOASPath:    filepath.Join(dir, "workspace-openapi.yaml"),
		WorkspaceConfigPath: filepath.Join(dir, "workspace-stainless.yaml"),
		OverrideOASPath:     filepath.Join(dir, "override-openapi.yaml"),
		OverrideConfigPath:  filepath.Join(dir, "override-stainless.yaml"),
		WorkspaceOAS:        "workspace-openapi",
		WorkspaceConfig:     "workspace-config",
		OverrideOAS:         "override-openapi",
		OverrideConfig:      "override-config",
	}

	require.NoError(t, os.WriteFile(fixture.WorkspaceOASPath, []byte(fixture.WorkspaceOAS), 0644))
	require.NoError(t, os.WriteFile(fixture.WorkspaceConfigPath, []byte(fixture.WorkspaceConfig), 0644))
	require.NoError(t, os.WriteFile(fixture.OverrideOASPath, []byte(fixture.OverrideOAS), 0644))
	require.NoError(t, os.WriteFile(fixture.OverrideConfigPath, []byte(fixture.OverrideConfig), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".stainless"), 0755))

	workspaceConfig := map[string]any{
		"project":          fixture.WorkspaceProject,
		"openapi_spec":     "../" + filepath.Base(fixture.WorkspaceOASPath),
		"stainless_config": "../" + filepath.Base(fixture.WorkspaceConfigPath),
	}
	data, err := json.Marshal(workspaceConfig)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".stainless", "workspace.json"), data, 0644))

	return fixture
}

func findRequest(t *testing.T, requests []mockstainless.RecordedRequest, method string, pathSuffix string) mockstainless.RecordedRequest {
	t.Helper()
	for _, request := range requests {
		if request.Method == method && request.Path == pathSuffix {
			return request
		}
	}
	require.Failf(t, "request not found", "method=%s path=%s requests=%v", method, pathSuffix, requests)
	return mockstainless.RecordedRequest{}
}

func assertNoRequest(t *testing.T, requests []mockstainless.RecordedRequest, method string, path string) {
	t.Helper()
	for _, request := range requests {
		if request.Method == method && request.Path == path {
			require.Failf(t, "unexpected request", "method=%s path=%s requests=%v", method, path, requests)
		}
	}
}

func assertProjectBody(t *testing.T, request mockstainless.RecordedRequest, project string) {
	t.Helper()
	assert.Equal(t, project, gjson.Get(request.Body, "project").String())
}

func assertProjectPathSuffix(t *testing.T, request mockstainless.RecordedRequest, project string) {
	t.Helper()
	assert.Equal(t, "/v0/projects/"+project, request.Path)
}

func assertBuildCreateFiles(t *testing.T, request mockstainless.RecordedRequest, fixture workspaceFixture) {
	t.Helper()
	assert.Equal(t, fixture.WorkspaceOAS, gjson.Get(request.Body, "revision.openapi\\.yaml.content").String())
	assert.Equal(t, fixture.WorkspaceConfig, gjson.Get(request.Body, "revision.stainless\\.yaml.content").String())
}

func assertLintFiles(t *testing.T, request mockstainless.RecordedRequest, fixture workspaceFixture) {
	t.Helper()
	assert.Equal(t, fixture.WorkspaceOAS, gjson.Get(request.Body, "source.openapi_spec").String())
	assert.Equal(t, fixture.WorkspaceConfig, gjson.Get(request.Body, "source.stainless_config").String())
}

func assertOpenAPIOnly(t *testing.T, request mockstainless.RecordedRequest, openapi string) {
	t.Helper()
	if request.Path == "/v0/builds" {
		assert.Equal(t, openapi, gjson.Get(request.Body, "revision.openapi\\.yaml.content").String())
		assert.False(t, gjson.Get(request.Body, "revision.stainless\\.yaml").Exists())
		return
	}
	assert.Equal(t, openapi, gjson.Get(request.Body, "head.revision.openapi\\.yaml.content").String())
	assert.False(t, gjson.Get(request.Body, "head.revision.stainless\\.yaml").Exists())
}

func assertConfigOnly(t *testing.T, request mockstainless.RecordedRequest, config string) {
	t.Helper()
	if request.Path == "/v0/builds" {
		assert.Equal(t, config, gjson.Get(request.Body, "revision.stainless\\.yaml.content").String())
		assert.False(t, gjson.Get(request.Body, "revision.openapi\\.yaml").Exists())
		return
	}
	assert.Equal(t, config, gjson.Get(request.Body, "head.revision.stainless\\.yaml.content").String())
	assert.False(t, gjson.Get(request.Body, "head.revision.openapi\\.yaml").Exists())
}
