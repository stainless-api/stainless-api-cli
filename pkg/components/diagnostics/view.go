package diagnostics

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/stainless-api/stainless-api-go"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	noteStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	codeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	refStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type sourceResolver struct {
	parsed map[string]parsedSource
}

type parsedSource struct {
	file *ast.File
	err  error
}

// levelLabel returns the colored level prefix and bracket-wrapped code for a diagnostic.
func levelLabel(level stainless.BuildDiagnosticLevel, code string) string {
	var levelStr string
	switch level {
	case stainless.BuildDiagnosticLevelFatal:
		levelStr = errorStyle.Render("fatal")
		code = errorStyle.UnsetBold().Render("[" + code + "]")
	case stainless.BuildDiagnosticLevelError:
		levelStr = errorStyle.Render("error")
		code = errorStyle.UnsetBold().Render("[" + code + "]")
	case stainless.BuildDiagnosticLevelWarning:
		levelStr = warningStyle.Render("warning")
		code = warningStyle.UnsetBold().Render("[" + code + "]")
	case stainless.BuildDiagnosticLevelNote:
		levelStr = noteStyle.Render("note")
		code = noteStyle.Render("[" + code + "]")
	default:
		levelStr = code
		code = ""
	}
	if code != "" {
		return levelStr + code
	}
	return levelStr
}

// ViewDiagnosticsError renders an error when fetching diagnostics fails
func ViewDiagnosticsError(err error) string {
	return errorStyle.Render("error") + ": failed to fetch diagnostics: " + err.Error() + "\n"
}

// ViewDiagnostics renders build diagnostics in Rust-style formatting.
// Notes are hidden by default. oasPath and configPath should be display paths,
// typically relative to the current working directory.
func ViewDiagnostics(diagnostics []stainless.BuildDiagnostic, maxDiagnostics int, oasPath, configPath string) string {
	// Filter out notes
	var visible []stainless.BuildDiagnostic
	for _, d := range diagnostics {
		if d.Level != stainless.BuildDiagnosticLevelNote {
			visible = append(visible, d)
		}
	}

	if len(visible) == 0 {
		grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		return grayStyle.Render("(no diagnostics)") + "\n"
	}

	var s strings.Builder
	resolver := sourceResolver{parsed: map[string]parsedSource{}}

	truncated := false
	shown := len(visible)
	if maxDiagnostics >= 0 && len(visible) > maxDiagnostics {
		truncated = true
		shown = maxDiagnostics
	}

	rendered := 0
	for _, diag := range visible {
		if maxDiagnostics >= 0 && rendered >= maxDiagnostics {
			break
		}

		if rendered > 0 {
			s.WriteString("\n")
		}
		rendered++

		// Header: error[Code]: message
		s.WriteString(levelLabel(diag.Level, diag.Code))
		s.WriteString(": ")
		s.WriteString(diag.Message)
		s.WriteString("\n")

		// Source references
		if diag.OasRef != "" {
			s.WriteString(refStyle.Render("  --> " + resolver.resolveRef(oasPath, "openapi.yml", diag.OasRef)))
			s.WriteString("\n")
		}
		if diag.ConfigRef != "" {
			s.WriteString(refStyle.Render("  --> " + resolver.resolveRef(configPath, "stainless.yml", diag.ConfigRef)))
			s.WriteString("\n")
		}

		// Additional content from More field
		if diag.More.AsAny() != nil {
			switch more := diag.More.AsAny().(type) {
			case stainless.BuildDiagnosticMoreMarkdown:
				text := strings.TrimSpace(more.Markdown)
				if text != "" {
					for _, line := range strings.Split(text, "\n") {
						s.WriteString("  ")
						s.WriteString(line)
						s.WriteString("\n")
					}
				}
			case stainless.BuildDiagnosticMoreRaw:
				text := strings.TrimSpace(more.Raw)
				if text != "" {
					for _, line := range strings.Split(text, "\n") {
						s.WriteString("  ")
						s.WriteString(line)
						s.WriteString("\n")
					}
				}
			}
		}
	}

	if truncated {
		s.WriteString(fmt.Sprintf("\n... and %d more diagnostics\n", len(visible)-shown))
	}

	return s.String()
}

func (r *sourceResolver) resolveRef(path, fallbackLabel, pointer string) string {
	label := sourceLabel(path, fallbackLabel)
	if line, column, ok := r.resolvePointer(path, pointer); ok {
		return fmt.Sprintf("%s:%d:%d: %s", label, line, column, pointer)
	}
	return label + ": " + pointer
}

func sourceLabel(path, fallbackLabel string) string {
	if path == "" {
		return fallbackLabel
	}
	return path
}

func (r *sourceResolver) resolvePointer(displayPath, pointer string) (int, int, bool) {
	path, ok := resolveSourcePath(displayPath)
	if !ok {
		return 0, 0, false
	}

	parsed, ok := r.parsed[path]
	if !ok {
		content, err := os.ReadFile(path)
		if err != nil {
			parsed = parsedSource{err: err}
		} else {
			file, err := parser.ParseBytes(content, 0)
			parsed = parsedSource{file: file, err: err}
		}
		r.parsed[path] = parsed
	}
	if parsed.err != nil || parsed.file == nil {
		return 0, 0, false
	}

	node, ok := resolveJSONPointer(parsed.file, pointer)
	if !ok {
		return 0, 0, false
	}

	token := node.GetToken()
	if token == nil || token.Position == nil {
		return 0, 0, false
	}
	return token.Position.Line, token.Position.Column, true
}

func resolveSourcePath(displayPath string) (string, bool) {
	if displayPath == "" {
		return "", false
	}

	path, err := filepath.Abs(displayPath)
	if err != nil {
		return "", false
	}
	return path, true
}

func resolveJSONPointer(file *ast.File, pointer string) (ast.Node, bool) {
	node := firstDocumentNode(file)
	if node == nil {
		return nil, false
	}

	segments, ok := parseJSONPointer(pointer)
	if !ok {
		return nil, false
	}

	for _, segment := range segments {
		var found bool
		node, found = descendNode(node, segment)
		if !found {
			return nil, false
		}
	}

	return node, true
}

func firstDocumentNode(file *ast.File) ast.Node {
	for _, doc := range file.Docs {
		if doc.Body == nil || doc.Body.Type() == ast.DirectiveType {
			continue
		}
		return doc.Body
	}
	return nil
}

func parseJSONPointer(pointer string) ([]string, bool) {
	if pointer == "" || pointer == "#" {
		return nil, true
	}

	switch {
	case strings.HasPrefix(pointer, "#/"):
		pointer = pointer[1:]
	case !strings.HasPrefix(pointer, "/"):
		return nil, false
	}

	parts := strings.Split(pointer[1:], "/")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		unescaped, err := url.PathUnescape(part)
		if err != nil {
			return nil, false
		}
		part = unescaped
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")
		segments = append(segments, part)
	}
	return segments, true
}

func descendNode(node ast.Node, segment string) (ast.Node, bool) {
	switch node := node.(type) {
	case *ast.MappingNode:
		for _, value := range node.Values {
			if mapKeyString(value.Key) == segment {
				return value.Value, true
			}
		}
	case *ast.SequenceNode:
		idx, err := strconv.Atoi(segment)
		if err != nil || idx < 0 || idx >= len(node.Values) {
			return nil, false
		}
		return node.Values[idx], true
	}
	return nil, false
}

func mapKeyString(key ast.MapKeyNode) string {
	if key == nil || key.GetToken() == nil {
		return ""
	}

	value := key.GetToken().Value
	if len(value) == 0 {
		return value
	}

	switch value[0] {
	case '"':
		unquoted, err := strconv.Unquote(value)
		if err == nil {
			return unquoted
		}
	case '\'':
		if len(value) > 1 && value[len(value)-1] == '\'' {
			return value[1 : len(value)-1]
		}
	}

	return value
}
