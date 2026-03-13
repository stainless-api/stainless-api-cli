package diagnostics

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-go"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	noteStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	codeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	refStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

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
// Notes are hidden by default. oasLabel and configLabel are the filenames
// shown in source references (e.g. "openapi.json", "stainless.yaml").
func ViewDiagnostics(diagnostics []stainless.BuildDiagnostic, maxDiagnostics int, oasLabel, configLabel string) string {
	if oasLabel == "" {
		oasLabel = "openapi.yml"
	}
	if configLabel == "" {
		configLabel = "stainless.yml"
	}
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
			s.WriteString(refStyle.Render("  --> " + oasLabel + ": " + diag.OasRef))
			s.WriteString("\n")
		}
		if diag.ConfigRef != "" {
			s.WriteString(refStyle.Render("  --> " + configLabel + ": " + diag.ConfigRef))
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
