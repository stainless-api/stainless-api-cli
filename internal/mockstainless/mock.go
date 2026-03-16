package mockstainless

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// M is a shorthand for map[string]any, used throughout for JSON-serializable data.
type M = map[string]any

const (
	DefaultOrg     = "acme-corp"
	DefaultProject = "acme-api"
)

// Mock holds all data for a mock Stainless API server.
type Mock struct {
	Builds           []*ProgressiveBuild
	Orgs             []M
	Projects         []M
	ProjectConfigs   M
	CompareBuild     *CompareBuildConfig
	AuthPendingCount int

	mu             sync.Mutex
	buildIndex     map[string]*ProgressiveBuild
	enableGitRepos bool
	gitRepos       map[string]gitRepo // key: "owner/name"
	tempDir        string
}

type gitRepo struct {
	Path string // path to the git repo
	Ref  string // commit SHA on main branch
}

// CompareBuildConfig configures the POST /v0/builds/compare endpoint.
type CompareBuildConfig struct {
	Base         M
	Head         M
	PreviewBuild *ProgressiveBuild
}

func (m *Mock) init() {
	m.buildIndex = make(map[string]*ProgressiveBuild, len(m.Builds))
	for _, b := range m.Builds {
		m.buildIndex[b.ID] = b
	}
	if m.CompareBuild != nil && m.CompareBuild.PreviewBuild != nil {
		m.buildIndex[m.CompareBuild.PreviewBuild.ID] = m.CompareBuild.PreviewBuild
	}
	if m.enableGitRepos {
		m.initGitRepos()
	}
}

// Cleanup removes temporary resources (git repos).
func (m *Mock) Cleanup() {
	if m.tempDir != "" {
		os.RemoveAll(m.tempDir)
	}
}

// initGitRepos creates local git repos for each unique repo found in build CompletedData.
func (m *Mock) initGitRepos() {
	m.gitRepos = make(map[string]gitRepo)
	tempDir, err := os.MkdirTemp("", "mock-git-repos-*")
	if err != nil {
		return
	}
	m.tempDir = tempDir

	// Collect unique repos from all builds, including compare preview builds.
	type repoKey struct{ owner, name string }
	seen := map[repoKey]bool{}

	collectRepos := func(b *ProgressiveBuild) {
		if b == nil {
			return
		}
		for _, targetData := range b.CompletedData {
			commitStep, _ := targetData["commit"].(M)
			commitObj, _ := commitStep["commit"].(M)
			repo, _ := commitObj["repo"].(M)
			owner, _ := repo["owner"].(string)
			name, _ := repo["name"].(string)
			if owner == "" || name == "" {
				continue
			}
			key := repoKey{owner, name}
			if seen[key] {
				continue
			}
			seen[key] = true

			repoPath := filepath.Join(tempDir, name)
			ref, err := createMockGitRepo(repoPath, name)
			if err != nil {
				continue
			}
			m.gitRepos[owner+"/"+name] = gitRepo{Path: repoPath, Ref: ref}
		}
	}
	for _, b := range m.Builds {
		collectRepos(b)
	}
	if m.CompareBuild != nil {
		collectRepos(m.CompareBuild.PreviewBuild)
	}
}

// createMockGitRepo creates a git repo with a single commit and returns the commit SHA.
func createMockGitRepo(dir, name string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	cmds := [][]string{
		{"git", "-C", dir, "init", "-b", "main"},
		{"git", "-C", dir, "config", "user.email", "mock@example.com"},
		{"git", "-C", dir, "config", "user.name", "Mock"},
	}
	for _, args := range cmds {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return "", err
		}
	}

	// Create a README
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# "+name+"\n"), 0644); err != nil {
		return "", err
	}

	cmds = [][]string{
		{"git", "-C", dir, "add", "."},
		{"git", "-C", dir, "commit", "-m", "Initial commit"},
	}
	for _, args := range cmds {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return "", err
		}
	}

	// Get the commit SHA
	var out bytes.Buffer
	cmd := exec.Command("git", "-C", dir, "rev-parse", "HEAD")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out.Bytes())), nil
}

// GetGitRepo returns the local git repo info for a given owner/name, if available.
func (m *Mock) GetGitRepo(owner, name string) (path, ref string, ok bool) {
	repo, ok := m.gitRepos[owner+"/"+name]
	if !ok {
		return "", "", false
	}
	return repo.Path, repo.Ref, true
}

func (m *Mock) GetBuild(id string) *ProgressiveBuild {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buildIndex[id]
}

func (m *Mock) Diagnostics(id string) []M {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pb, ok := m.buildIndex[id]; ok && pb.Diagnostics != nil {
		return pb.Diagnostics
	}
	return []M{}
}

// MockOption configures a Mock via NewMock.
type MockOption func(*Mock)

// NewMock creates a new mock with the given options.
func NewMock(opts ...MockOption) *Mock {
	m := &Mock{}
	for _, opt := range opts {
		opt(m)
	}
	m.init()
	return m
}

// Server returns an http.Handler serving the mock's endpoints.
func (m *Mock) Server() http.Handler {
	return newServeMux(m)
}

// WithGitRepos creates local git repos for each target in the mock's builds.
// This enables the build_target_outputs endpoint to return file:// URLs that
// work with git fetch. Call Mock.Cleanup() when done to remove temp directories.
func WithGitRepos() MockOption {
	return func(m *Mock) {
		m.enableGitRepos = true
	}
}

// WithDeviceAuth sets how many "authorization_pending" responses the token
// endpoint returns before succeeding.
func WithDeviceAuth(pendingCount int) MockOption {
	return func(m *Mock) {
		m.AuthPendingCount = pendingCount
	}
}

// WithAutomaticDeviceAuth configures instant auth success (zero pending responses).
func WithAutomaticDeviceAuth() MockOption {
	return WithDeviceAuth(0)
}

// MockOrg describes an organization to register in the mock.
type MockOrg struct {
	Name        string // slug, required
	DisplayName string // defaults to Name
}

func (o MockOrg) toM() M {
	displayName := o.DisplayName
	if displayName == "" {
		displayName = o.Name
	}
	return M{
		"slug":                      o.Name,
		"display_name":              displayName,
		"object":                    "org",
		"enable_ai_commit_messages": false,
	}
}

// WithOrg adds an organization to the mock.
func WithOrg(org MockOrg) MockOption {
	return func(m *Mock) {
		m.Orgs = append(m.Orgs, org.toM())
	}
}

// MockProject describes a project to register in the mock.
type MockProject struct {
	Name        string              // slug, required
	DisplayName string              // defaults to Name
	Org         string              // defaults to first configured org's slug
	Targets     []string            // defaults to ["typescript", "python", "go"]
	Builds      []*ProgressiveBuild // added to the mock's build list
	Configs     M                   // project config files
}

func (p MockProject) toM(org string) M {
	displayName := p.DisplayName
	if displayName == "" {
		displayName = p.Name
	}
	targets := p.Targets
	if len(targets) == 0 {
		targets = []string{"typescript", "python", "go"}
	}
	return M{
		"slug":         p.Name,
		"display_name": displayName,
		"object":       "project",
		"org":          org,
		"config_repo":  fmt.Sprintf("https://github.com/%s/%s", org, p.Name),
		"targets":      targets,
	}
}

// WithProject adds a project (and its builds) to the mock.
func WithProject(project MockProject) MockOption {
	return func(m *Mock) {
		org := project.Org
		if org == "" && len(m.Orgs) > 0 {
			org = m.Orgs[0]["slug"].(string)
		}
		m.Projects = append(m.Projects, project.toM(org))
		m.Builds = append(m.Builds, project.Builds...)
		if project.Configs != nil {
			m.ProjectConfigs = project.Configs
		}
	}
}

// WithCompareBuild enables the POST /v0/builds/compare endpoint.
func WithCompareBuild(cfg CompareBuildConfig) MockOption {
	return func(m *Mock) {
		m.CompareBuild = &cfg
	}
}

// WithDefaultOrg adds the default acme-corp organization.
func WithDefaultOrg() MockOption {
	return WithOrg(MockOrg{Name: "acme-corp", DisplayName: "Acme Corp"})
}

// WithDefaultProject adds the default acme-api project with 4 demo builds and config files.
func WithDefaultProject() MockOption {
	return func(m *Mock) {
		now := time.Now()
		WithProject(MockProject{
			Name:        "acme-api",
			DisplayName: "Acme API",
			Org:         "acme-corp",
			Targets:     []string{"typescript", "python", "go"},
			Configs: M{
				"stainless.yml": M{
					"content": "# Stainless configuration\norganization:\n  name: acme-corp\n  docs_url: https://docs.acme.com\n\nclient:\n  name: Acme\n\nendpoints:\n  list_pets:\n    path: /pets\n    method: get\n  create_pet:\n    path: /pets\n    method: post\n  get_pet:\n    path: /pets/{id}\n    method: get\n",
				},
				"openapi.json": M{
					"content": "{\"openapi\":\"3.1.0\",\"info\":{\"title\":\"Acme Pet Store API\",\"version\":\"1.0.0\"},\"paths\":{\"/pets\":{\"get\":{\"summary\":\"List pets\"},\"post\":{\"summary\":\"Create pet\"}},\"/pets/{id}\":{\"get\":{\"summary\":\"Get pet\"}}}}",
				},
			},
			Builds: []*ProgressiveBuild{
				{
					ID:           "bui_0cmmtv8r2j000425s640dp4kwn",
					ConfigCommit: "e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4",
					Targets:      []string{"typescript", "python", "go"},
					CompletedData: map[string]M{
						"typescript": CompletedTarget("acme-corp", "acme-typescript", "f4e8a2c91d3b7056ef12c489a37d6b0e51f8c2a4", 247, 83),
						"python":     CompletedTarget("acme-corp", "acme-python", "b3a9d7e21c5f8046ea31d589c47b6a0f52e8d3b5", 189, 42),
						"go":         CompletedTarget("acme-corp", "acme-go", "7b3d9e1f25a8c460d2f7b91e3c5a8d0f64e2b7c1", 156, 61),
					},
					StartTime: now,
				},
				{
					ID:           "bui_0cmmtsksxj000425s640c55yf1",
					ConfigCommit: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0",
					Targets:      []string{"typescript", "python", "go", "java", "kotlin", "ruby"},
					CompletedData: map[string]M{
						"typescript": CompletedTarget("acme-corp", "acme-typescript", "f4e8a2c91d3b7056ef12c489a37d6b0e51f8c2a4", 247, 83),
						"python":     WarningTarget("acme-corp", "acme-python", "b3a9d7e21c5f8046ea31d589c47b6a0f52e8d3b5", 189, 42),
						"go":         ErrorTarget("acme-corp", "acme-go", "7b3d9e1f25a8c460d2f7b91e3c5a8d0f64e2b7c1", 156, 61),
						"java":       MergeConflictTarget("acme-corp", "acme-java", 42),
						"kotlin":     FatalTarget(),
						"ruby":       NotStartedTarget(),
					},
					Diagnostics: []M{
						Diagnostic("Schema/TypeMismatch", "error", "Expected `string` type but got `integer` in response schema for `get /pets/{pet_id}`.",
							WithOASRef("#/paths/%2Fpets%2F%7Bpet_id%7D/get/responses/200/content/application%2Fjson/schema/properties/age"),
							WithMore("The `age` property is declared as `string` in the schema but the example value is an integer.\n\nTo fix this, change the type to `integer` or update the example value."),
						),
						Diagnostic("Schema/CannotInferUnionVariantName", "warning", "Placeholder name generated for union variant.",
							WithOASRef("#/paths/%2Fpets%2F%7Bpet_id%7D/get/responses/200/content/application%2Fjson/schema/anyOf/1"),
							WithMore("We were unable to infer a good name for this union variant, so we gave it an arbitrary placeholder name.\n\nTo resolve this issue, do one of the following:\n\n- Define a [model](https://www.stainless.com/docs/guides/configure#models)\n- Set a `title` property on the schema\n- Extract the schema to `#/components/schemas`\n- Provide a name by adding an `x-stainless-variantName` property to the schema containing the name you want to use"),
						),
						Diagnostic("Schema/IsAmbiguous", "warning", "This schema does not have at least one of `type`,\n`oneOf`, `anyOf`, or `allOf`, so its type has been interpreted as `unknown`.",
							WithOASRef("#/components/schemas/PetMetadata"),
							WithMore("If the schema should have a specific type, then add `type` to it.\n\nIf the schema should accept anything, then add [`x-stainless-any: true`](https://www.stainless.com/docs/reference/openapi-support#unknown-and-any)\nto suppress this note."),
						),
						Diagnostic("Endpoint/IsIgnored", "note", "`get /internal/health` is in `unspecified_endpoints`, so code will not be\ngenerated for it.",
							WithOASRef("#/paths/%2Finternal%2Fhealth/get"),
							WithConfigRef("#/unspecified_endpoints/0"),
							WithMore("If this is intentional, then ignore this note. Otherwise, remove the endpoint from\n`unspecified_endpoints` and add it to `resources`."),
						),
						Diagnostic("Schema/ObjectHasNoProperties", "note", "This schema has neither `properties` nor `additionalProperties` so\nits type has been interpreted as `unknown`.",
							WithOASRef("#/paths/%2Fpets%2F%7Bpet_id%7D%2Fvaccinations/get/responses/200/content/application%2Fjson/schema"),
							WithMore("If the schema should be a map, then add [`additionalProperties`](https://json-schema.org/understanding-json-schema/reference/object#additionalproperties)\nto it.\n\nIf the schema should be an empty object type, then add `x-stainless-empty-object: true` to it."),
						),
					},
					StartTime: now.Add(-10 * time.Minute),
				},
				{
					ID:           "bui_0cmmtrmq4z000425s640hpf9gx",
					ConfigCommit: "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1",
					Targets:      []string{"typescript", "python", "go"},
					CompletedData: map[string]M{
						"typescript": CompletedTarget("acme-corp", "acme-typescript", "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0", 52, 18),
						"python": Target("completed",
							CommitCompleted("warning", WithCommitData("acme-corp", "acme-python", "e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2", 67, 23)),
							WithLint(CheckStepCompleted("success")),
							WithBuild(CheckStepCompleted("failure")),
							WithTest(CheckStepCompleted("skipped")),
						),
						"go": CompletedTarget("acme-corp", "acme-go", "c2a5f8e31b6d9074a3e5c8f12d7b4a69e0f3c5b8", 98, 34),
					},
					Diagnostics: []M{
						Diagnostic("Schema/CannotInferUnionVariantName", "warning", "Placeholder name generated for union variant.",
							WithOASRef("#/paths/%2Fpets/get/responses/200/content/application%2Fjson/schema/anyOf/0"),
						),
					},
					StartTime: now.Add(-2 * time.Hour),
				},
				{
					ID:           "bui_0cmmtg5e8n000425s6403bkywd",
					ConfigCommit: "c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2",
					Targets:      []string{"typescript", "python", "go"},
					CompletedData: map[string]M{
						"typescript": CompletedTarget("acme-corp", "acme-typescript", "a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8", 312, 98),
						"python":     CompletedTarget("acme-corp", "acme-python", "f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9", 276, 114),
						"go":         CompletedTarget("acme-corp", "acme-go", "b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0", 198, 77),
					},
					StartTime: now.Add(-24 * time.Hour),
				},
			},
		})(m)
	}
}

// WithDefaultCompareBuild adds a compare endpoint with a preview build.
func WithDefaultCompareBuild() MockOption {
	return func(m *Mock) {
		now := time.Now()
		previewTargets := []string{"typescript", "python", "go"}

		base := Build("build_preview_base_01",
			WithConfigCommit("c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2"),
			WithCreatedAt(now),
		)
		head := Build("build_preview_head_01",
			WithConfigCommit("d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3"),
			WithCreatedAt(now),
		)
		for _, t := range previewTargets {
			base["targets"].(M)[t] = CompletedTarget("acme-corp", "acme-"+t, "a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8", 0, 0)
			head["targets"].(M)[t] = NotStartedTarget()
		}

		previewBuild := &ProgressiveBuild{
			ID:           "build_preview_head_01",
			ConfigCommit: "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3",
			Targets:      previewTargets,
			CompletedData: map[string]M{
				"typescript": CompletedTarget("acme-corp", "acme-typescript", "e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4", 134, 47),
				"python":     WarningTarget("acme-corp", "acme-python", "a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6", 89, 31),
				"go":         ErrorTarget("acme-corp", "acme-go", "c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7", 156, 61),
			},
			Diagnostics: []M{
				Diagnostic("Schema/TypeMismatch", "error", "Expected `string` type but got `integer` in response schema for `get /pets/{pet_id}`.",
					WithOASRef("#/paths/%2Fpets%2F%7Bpet_id%7D/get/responses/200/content/application%2Fjson/schema/properties/age"),
				),
				Diagnostic("Schema/CannotInferUnionVariantName", "warning", "Placeholder name generated for union variant.",
					WithOASRef("#/paths/%2Fpets%2F%7Bpet_id%7D/get/responses/200/content/application%2Fjson/schema/anyOf/1"),
				),
			},
		}

		WithCompareBuild(CompareBuildConfig{
			Base:         base,
			Head:         head,
			PreviewBuild: previewBuild,
		})(m)
	}
}
