package mockstainless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// Page wraps data in a paginated response envelope.
func Page(data []M) M {
	return M{
		"data":        data,
		"next_cursor": "",
	}
}

// newServeMux creates an http.Handler with all mock endpoints registered.
func newServeMux(m *Mock) http.Handler {
	authCounter := &CallCounter{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/oauth/device", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, M{
			"device_code":               "demo_device_code_abc123",
			"user_code":                 "DEMO-CODE",
			"verification_uri":          "https://app.stainless.com/activate",
			"verification_uri_complete": "https://app.stainless.com/activate?code=DEMO-CODE",
			"expires_in":                300,
			"interval":                  1,
		})
	})

	mux.HandleFunc("POST /v0/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		count := authCounter.Increment()
		if count <= m.AuthPendingCount {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "authorization_pending",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"access_token":  "demo_access_token_xyz789",
			"refresh_token": "demo_refresh_token_abc456",
			"token_type":    "bearer",
		})
	})

	mux.HandleFunc("GET /v0/orgs", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, M{
			"data":     m.Orgs,
			"has_more": false,
		})
	})

	mux.HandleFunc("GET /v0/projects", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, Page(m.Projects))
	})

	mux.HandleFunc("POST /v0/projects", func(w http.ResponseWriter, r *http.Request) {
		body := mustReadBody(r)
		slug := gjson.GetBytes(body, "slug").String()
		displayName := gjson.GetBytes(body, "display_name").String()
		org := gjson.GetBytes(body, "org").String()

		if slug == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "slug is required"})
			return
		}
		if displayName == "" {
			displayName = slug
		}

		var targets []any
		gjson.GetBytes(body, "targets").ForEach(func(_, v gjson.Result) bool {
			targets = append(targets, v.String())
			return true
		})
		if len(targets) == 0 {
			targets = []any{"typescript", "python", "go"}
		}

		project := M{
			"slug":         slug,
			"display_name": displayName,
			"object":       "project",
			"org":          org,
			"config_repo":  fmt.Sprintf("https://github.com/%s/%s", org, slug),
			"targets":      targets,
		}

		m.mu.Lock()
		m.Projects = append(m.Projects, project)
		m.mu.Unlock()

		// Create a build so the post-creation build-wait step succeeds.
		if len(m.Builds) > 0 {
			m.CreateBuildFromTemplate(m.Builds[0])
		}

		writeJSON(w, http.StatusOK, project)
	})

	mux.HandleFunc("GET /v0/projects/{project}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("project")
		for _, p := range m.Projects {
			if p["slug"] == slug {
				writeJSON(w, http.StatusOK, p)
				return
			}
		}
		writeJSON(w, http.StatusNotFound, M{"error": "project not found"})
	})

	mux.HandleFunc("PATCH /v0/projects/{project}", func(w http.ResponseWriter, r *http.Request) {
		project := r.PathValue("project")
		if project == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		body := mustReadBody(r)
		writeJSON(w, http.StatusOK, M{
			"slug":         project,
			"display_name": gjson.GetBytes(body, "display_name").String(),
			"object":       "project",
		})
	})

	mux.HandleFunc("POST /v0/projects/{project}/generate_commit_message", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"message": "mock commit message",
		})
	})

	mux.HandleFunc("GET /v0/projects/{project}/configs", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, m.ProjectConfigs)
	})

	mux.HandleFunc("POST /v0/projects/{project}/configs/guess", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"stainless.yml": M{
				"content": "# guessed",
			},
		})
	})

	mux.HandleFunc("POST /v0/projects/{project}/branches", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		body := mustReadBody(r)
		writeJSON(w, http.StatusOK, M{
			"branch": gjson.GetBytes(body, "branch").String(),
			"object": "project_branch",
		})
	})

	mux.HandleFunc("GET /v0/projects/{project}/branches", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, Page([]M{
			{"branch": "main", "object": "project_branch"},
		}))
	})

	mux.HandleFunc("GET /v0/projects/{project}/branches/{branch}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"branch": r.PathValue("branch"),
			"object": "project_branch",
		})
	})

	mux.HandleFunc("DELETE /v0/projects/{project}/branches/{branch}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{"deleted": true})
	})

	mux.HandleFunc("PUT /v0/projects/{project}/branches/{branch}/rebase", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"branch": r.PathValue("branch"),
			"object": "project_branch",
		})
	})

	mux.HandleFunc("PUT /v0/projects/{project}/branches/{branch}/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"branch": r.PathValue("branch"),
			"object": "project_branch",
		})
	})

	mux.HandleFunc("POST /v0/builds", func(w http.ResponseWriter, r *http.Request) {
		body := mustReadBody(r)
		if gjson.GetBytes(body, "project").String() == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		if len(m.Builds) == 0 {
			writeJSON(w, http.StatusNotFound, M{"error": "missing build"})
			return
		}
		build := m.CreateBuildFromTemplate(m.Builds[0])
		if build == nil {
			writeJSON(w, http.StatusNotFound, M{"error": "missing build"})
			return
		}
		writeJSON(w, http.StatusOK, build.Snapshot())
	})

	mux.HandleFunc("GET /v0/builds", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("project") == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
			return
		}
		builds := make([]M, len(m.Builds))
		for i, b := range m.Builds {
			builds[i] = b.Snapshot()
		}
		writeJSON(w, http.StatusOK, Page(builds))
	})

	mux.HandleFunc("GET /v0/builds/{id}", func(w http.ResponseWriter, r *http.Request) {
		if pb := m.GetBuild(r.PathValue("id")); pb != nil {
			writeJSON(w, http.StatusOK, pb.Snapshot())
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	mux.HandleFunc("GET /v0/builds/{id}/diagnostics", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, Page(m.Diagnostics(r.PathValue("id"))))
	})

	mux.HandleFunc("GET /v0/build_target_outputs", func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("build_id")
		target := r.URL.Query().Get("target")
		outputType := r.URL.Query().Get("output")
		sourceType := r.URL.Query().Get("type")

		pb := m.GetBuild(buildID)
		if pb == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		targetData, ok := pb.CompletedData[target]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Extract commit info: target["commit"]["commit"]["sha"] and repo info
		commitStep, _ := targetData["commit"].(M)
		commitObj, _ := commitStep["commit"].(M)
		repo, _ := commitObj["repo"].(M)
		sha, _ := commitObj["sha"].(string)
		owner, _ := repo["owner"].(string)
		name, _ := repo["name"].(string)

		if sha == "" || owner == "" || name == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gitURL := fmt.Sprintf("https://github.com/%s/%s", owner, name)
		ref := sha

		// Use local git repo if available
		if repoPath, localRef, ok := m.GetGitRepo(owner, name); ok {
			gitURL = "file://" + repoPath
			ref = localRef
		}

		switch outputType {
		case "git":
			writeJSON(w, http.StatusOK, M{
				"output": "git",
				"target": target,
				"type":   sourceType,
				"url":    gitURL,
				"ref":    ref,
				"token":  "mock_token_123",
			})
		default:
			writeJSON(w, http.StatusOK, M{
				"output": "url",
				"target": target,
				"type":   sourceType,
				"url":    gitURL + "/archive/" + ref + ".tar.gz",
			})
		}
	})

	if m.CompareBuild != nil {
		mux.HandleFunc("POST /v0/builds/compare", func(w http.ResponseWriter, r *http.Request) {
			body := mustReadBody(r)
			if gjson.GetBytes(body, "project").String() == "" {
				writeJSON(w, http.StatusBadRequest, M{"error": "project is required"})
				return
			}
			headBuild := m.CreateBuildFromTemplate(m.CompareBuild.PreviewBuild)
			if headBuild == nil {
				writeJSON(w, http.StatusNotFound, M{"error": "missing preview build"})
				return
			}
			head := cloneMap(m.CompareBuild.Head)
			head["id"] = headBuild.ID
			head["created_at"] = time.Now().Format(time.RFC3339)
			writeJSON(w, http.StatusOK, M{
				"base": m.CompareBuild.Base,
				"head": head,
			})
		})
	}

	mux.HandleFunc("POST /api/generate/spec", func(w http.ResponseWriter, r *http.Request) {
		body := mustReadBody(r)
		if gjson.GetBytes(body, "project").String() == "" ||
			gjson.GetBytes(body, "source.openapi_spec").String() == "" ||
			gjson.GetBytes(body, "source.stainless_config").String() == "" {
			writeJSON(w, http.StatusBadRequest, M{"error": "project, openapi_spec, and stainless_config are required"})
			return
		}
		writeJSON(w, http.StatusOK, M{
			"spec": M{
				"diagnostics": M{},
			},
		})
	})

	// Add simulated latency to all requests (except health checks).
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := mustReadBody(r)
		m.RecordRequest(RecordedRequest{
			Method:   r.Method,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
			Body:     string(body),
		})
		r.Body = io.NopCloser(bytes.NewReader(body))

		if r.URL.Path != "/health" {
			time.Sleep(150 * time.Millisecond)
		}
		mux.ServeHTTP(w, r)
	})
}

func mustReadBody(r *http.Request) []byte {
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body
}
