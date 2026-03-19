package build

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stainless-api/stainless-api-go"
)

var update = flag.Bool("update", false, "update snapshot files")

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.ANSI)
	os.Exit(m.Run())
}

func mustBuild(t *testing.T, jsonStr string) stainless.Build {
	t.Helper()
	var b stainless.Build
	if err := json.Unmarshal([]byte(jsonStr), &b); err != nil {
		t.Fatalf("failed to unmarshal build JSON: %v", err)
	}
	return b
}

func newSpinner() spinner.Model {
	return spinner.New()
}

// snapshot compares got against the snapshot file testdata/<name>.snapshot.
// When -update is passed, it writes/overwrites the snapshot file instead.
func snapshot(t *testing.T, name string, got string) {
	t.Helper()
	path := filepath.Join("testdata", name+".snapshot")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("snapshot file %s not found; run with -update to create it: %v", path, err)
	}
	if string(want) != got {
		t.Errorf("snapshot mismatch for %s\nwant: %q\ngot:  %q\nrun with -update to update", name, string(want), got)
	}
}

const checkSteps = `"lint": {"status": "not_started"}, "build": {"status": "not_started"}, "test": {"status": "not_started"}`

func TestViewBuildPipeline(t *testing.T) {
	sp := newSpinner()
	var out strings.Builder
	dl := map[stainless.Target]DownloadStatus{"typescript": {Status: "not_started"}}

	// queued (not_started)
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "not_started"}, `+checkSteps+`, "status": "not_started", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// queued (queued status)
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "queued"}, `+checkSteps+`, "status": "not_started", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// in_progress
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "in_progress"}, `+checkSteps+`, "status": "codegen", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// success with changes
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 100, "deletions": 30, "total": 130}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// success unchanged
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 0, "deletions": 0, "total": 0}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// warning conclusion
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "warning", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 50, "deletions": 10, "total": 60}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// error conclusion
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "error", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 50, "deletions": 10, "total": 60}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// fatal
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "fatal"}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// merge_conflict
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "merge_conflict", "merge_conflict_pr": {"number": 42, "repo": {"owner": "org", "name": "repo"}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// payment_required
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "payment_required"}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// cancelled
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "cancelled"}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// timed_out
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "timed_out"}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// noop
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "noop"}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// version_bump
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "version_bump", "commit": {"sha": "def5678901234", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 3, "deletions": 3, "total": 6}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// post-commit steps (lint/build/test in various states)
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 10, "deletions": 5, "total": 15}}}, "lint": {"status": "not_started"}, "build": {"status": "in_progress"}, "test": {"status": "completed", "conclusion": "success", "url": ""}, "status": "postgen", "object": "build_target", "install_url": ""}}
	}`), "typescript", dl, false, sp))
	out.WriteString("\n\n")

	// commitOnly hides post-commit steps
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 10, "deletions": 5, "total": 15}}}, `+checkSteps+`, "status": "postgen", "object": "build_target", "install_url": ""}}
	}`), "typescript", nil, true, sp))
	out.WriteString("\n\n")

	// download success
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 10, "deletions": 5, "total": 15}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", map[stainless.Target]DownloadStatus{"typescript": {Status: "completed", Conclusion: "success"}}, true, sp))
	out.WriteString("\n\n")

	// download failure
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{
		"id": "build_1",
		"targets": {"typescript": {"commit": {"status": "completed", "conclusion": "success", "commit": {"sha": "abc1234567890", "tree_oid": "tree", "repo": {"owner": "org", "name": "repo"}, "stats": {"additions": 10, "deletions": 5, "total": 15}}}, `+checkSteps+`, "status": "completed", "object": "build_target", "install_url": ""}}
	}`), "typescript", map[stainless.Target]DownloadStatus{"typescript": {Status: "completed", Conclusion: "failure", Error: "connection refused"}}, true, sp))
	out.WriteString("\n\n")

	// nil target
	out.WriteString(ViewBuildPipeline(mustBuild(t, `{"id": "build_1", "targets": {}}`), "typescript", nil, false, sp))

	snapshot(t, "view_build_pipeline", out.String())
}
