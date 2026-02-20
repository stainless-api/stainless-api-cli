package diagnostics

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stainless-api/stainless-api-go"
)

var update = flag.Bool("update", false, "update snapshot files")

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.ANSI)
	os.Exit(m.Run())
}

func mustDiags(t *testing.T, jsonStr string) []stainless.BuildDiagnostic {
	t.Helper()
	var d []stainless.BuildDiagnostic
	if err := json.Unmarshal([]byte(jsonStr), &d); err != nil {
		t.Fatalf("failed to unmarshal diagnostics JSON: %v", err)
	}
	return d
}

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

func TestViewDiagnostics(t *testing.T) {
	var out strings.Builder

	// no diagnostics
	out.WriteString(ViewDiagnostics(nil, 10))
	out.WriteString("\n")

	// notes only (hidden, treated as empty)
	out.WriteString(ViewDiagnostics(mustDiags(t, `[
		{"code": "StyleSuggestion", "level": "note", "message": "Consider camelCase", "ignored": false, "more": null}
	]`), 10))
	out.WriteString("\n")

	// fetch error
	out.WriteString(ViewDiagnosticsError(errors.New("connection refused")))
	out.WriteString("\n")

	// mixed: errors, warnings, notes, refs, more content, truncation
	out.WriteString(ViewDiagnostics(mustDiags(t, `[
		{
			"code": "MissingField",
			"level": "error",
			"message": "The field 'name' is required but missing",
			"ignored": false,
			"more": null,
			"oas_ref": "/paths/~1users/post/requestBody",
			"config_ref": "/endpoints/~1users/post"
		},
		{
			"code": "FatalError",
			"level": "fatal",
			"message": "Build failed due to configuration error",
			"ignored": false,
			"more": {"type": "markdown", "markdown": "Check your stainless.yml for syntax errors.\nSee docs for details."},
			"oas_ref": "/paths/~1users"
		},
		{
			"code": "DeprecatedUsage",
			"level": "warning",
			"message": "The x-deprecated extension is deprecated",
			"ignored": false,
			"more": null,
			"oas_ref": "/paths/~1foo/get"
		},
		{
			"code": "StyleSuggestion",
			"level": "note",
			"message": "Consider using camelCase",
			"ignored": false,
			"more": null
		},
		{
			"code": "Err3",
			"level": "error",
			"message": "Truncated away",
			"ignored": false,
			"more": null
		}
	]`), 3))

	snapshot(t, "view_diagnostics", out.String())
}
