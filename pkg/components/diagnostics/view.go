package diagnostics

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-go"
	"golang.org/x/term"
)

// ViewDiagnosticsError renders an error when fetching diagnostics fails
func ViewDiagnosticsError(err error) string {
	var s strings.Builder
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	s.WriteString(console.SProperty(0, "build diagnostics", errorStyle.Render("(error: "+err.Error()+")")))
	return s.String()
}

// ViewDiagnosticIcon returns a colored icon for a diagnostic level
func ViewDiagnosticIcon(level stainless.BuildDiagnosticLevel) string {
	switch level {
	case stainless.BuildDiagnosticLevelFatal:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("(F)")
	case stainless.BuildDiagnosticLevelError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("(E)")
	case stainless.BuildDiagnosticLevelWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("(W)")
	case stainless.BuildDiagnosticLevelNote:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("(i)")
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("•")
	}
}

var renderer *glamour.TermRenderer

func init() {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 || width > 120 {
		width = 120
	}
	renderer, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
}

// renderMarkdown renders markdown content using glamour
func renderMarkdown(content string) string {
	if renderer == nil {
		return content
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return content
	}

	return strings.Trim(rendered, "\n ")
}

// countDiagnosticsBySeverity counts diagnostics by severity level
func countDiagnosticsBySeverity(diagnostics []stainless.BuildDiagnostic) (fatal, errors, warnings, notes int) {
	for _, diag := range diagnostics {
		switch diag.Level {
		case stainless.BuildDiagnosticLevelFatal:
			fatal++
		case stainless.BuildDiagnosticLevelError:
			errors++
		case stainless.BuildDiagnosticLevelWarning:
			warnings++
		case stainless.BuildDiagnosticLevelNote:
			notes++
		}
	}
	return
}

// ViewDiagnostics renders build diagnostics with formatting
func ViewDiagnostics(diagnostics []stainless.BuildDiagnostic, maxDiagnostics int) string {
	var s strings.Builder

	if len(diagnostics) > 0 {
		// Count diagnostics by severity
		fatal, errors, warnings, notes := countDiagnosticsBySeverity(diagnostics)

		// Create summary string
		var summaryParts []string
		if fatal > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d fatal", fatal))
		}
		if errors > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d errors", errors))
		}
		if warnings > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d warnings", warnings))
		}
		if notes > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d notes", notes))
		}

		summary := strings.Join(summaryParts, ", ")
		if summary != "" {
			summary = fmt.Sprintf(" (%s)", summary)
		}

		var sub strings.Builder

		if maxDiagnostics >= 0 && len(diagnostics) > maxDiagnostics {
			sub.WriteString(fmt.Sprintf("Showing first %d of %d diagnostics:\n", maxDiagnostics, len(diagnostics)))
		}

		for i, diag := range diagnostics {
			if maxDiagnostics >= 0 && i >= maxDiagnostics {
				break
			}

			levelIcon := ViewDiagnosticIcon(diag.Level)
			codeStyle := lipgloss.NewStyle().Bold(true)

			if i > 0 {
				sub.WriteString("\n")
			}
			sub.WriteString(fmt.Sprintf("%s %s\n", levelIcon, codeStyle.Render(diag.Code)))
			sub.WriteString(fmt.Sprintf("%s\n", renderMarkdown(diag.Message)))

			if diag.Code == "FatalError" {
				switch more := diag.More.AsAny().(type) {
				case stainless.BuildDiagnosticMoreMarkdown:
					sub.WriteString(fmt.Sprintf("%s\n", renderMarkdown(more.Markdown)))
				case stainless.BuildDiagnosticMoreRaw:
					sub.WriteString(fmt.Sprintf("%s\n", more.Raw))
				}
			}

			// Show source references if available
			if diag.OasRef != "" {
				refStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
				sub.WriteString(fmt.Sprintf("    %s\n", refStyle.Render("OpenAPI: "+diag.OasRef)))
			}
			if diag.ConfigRef != "" {
				refStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
				sub.WriteString(fmt.Sprintf("    %s\n", refStyle.Render("Config: "+diag.ConfigRef)))
			}
		}

		s.WriteString(console.SProperty(0, "build diagnostics", summary))
		s.WriteString(lipgloss.NewStyle().
			Padding(0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("208")).
			Render(strings.TrimRight(sub.String(), "\n")),
		)
	} else {
		s.WriteString(console.SProperty(0, "build diagnostics", "(no errors or warnings)"))
	}

	return s.String()
}
