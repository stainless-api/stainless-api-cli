package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-go"
)

func (m BuildModel) View() string {
	s := strings.Builder{}

	startIndex := 0
	if m.view != "" {
		for i, part := range parts {
			if part.name == string(m.view) {
				startIndex = i + 1
				break
			}
		}
	}

	for i := startIndex; i < len(parts); i++ {
		parts[i].view(m, &s)
	}

	return s.String()
}

// updateView updates the view state and prints everything from current state to target state to scrollback
func (m *BuildModel) updateView(targetState string) tea.Cmd {
	// Find current state index
	currentIndex := -1
	if m.view != "" {
		for i, part := range parts {
			if part.name == m.view {
				currentIndex = i
				break
			}
		}
	} else {
		currentIndex = -1 // Start from beginning if no current view
	}

	// Find target state index
	targetIndex := -1
	for i, part := range parts {
		if part.name == targetState {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return nil // Target state not found
	}

	// Build output from current state to target state
	var output strings.Builder
	startIndex := currentIndex + 1
	if currentIndex == -1 {
		startIndex = 0
	}

	for i := startIndex; i <= targetIndex; i++ {
		parts[i].view(*m, &output)
	}

	// Update model state
	m.view = targetState

	// Print to scrollback
	if output.Len() > 0 {
		out := output.String()
		if out[len(out)-1] == '\n' {
			return tea.Println(out[:len(out)-1])
		} else {
			return tea.Println(out)
		}
	}
	return nil
}

type buildViewState string

var parts = []struct {
	name string
	view func(BuildModel, *strings.Builder)
}{
	{
		name: "header",
		view: func(m BuildModel, s *strings.Builder) {
			buildIDStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")).Bold(true)
			if m.build != nil {
				s.WriteString(fmt.Sprintf("\n\n%s %s\n\n", buildIDStyle.Render(" BUILD "), m.build.ID))
			} else {
				s.WriteString(fmt.Sprintf("\n\n%s\n\n", buildIDStyle.Render(" BUILD ")))
			}
		},
	},
	{
		name: "diagnostics",
		view: func(m BuildModel, s *strings.Builder) {
			if m.diagnostics == nil {
				s.WriteString(SProperty(0, "diagnostics", "waiting for build to finish"))
			} else {
				s.WriteString(ViewDiagnosticsPrint(m.diagnostics))
			}
		},
	},
	{
		name: "duration",
		view: func(m BuildModel, s *strings.Builder) {
			duration := m.getBuildDuration()
			s.WriteString(SProperty(0, "duration", duration.Round(time.Second).String()))
		},
	},
	{
		name: "studio",
		view: func(m BuildModel, s *strings.Builder) {
			if m.build != nil {
				url := fmt.Sprintf("https://app.stainless.com/%s/%s/studio?branch=%s", m.build.Org, m.build.Project, m.branch)
				s.WriteString(SProperty(0, "studio", Hyperlink(url, url)))
			}
		},
	},
	{
		name: "build_status",
		view: func(m BuildModel, s *strings.Builder) {
			s.WriteString("\n")
			if m.build != nil {
				buildObj := NewBuildObject(m.build)
				languages := buildObj.Languages()
				// Target rows with colors
				for _, target := range languages {
					pipeline := ViewBuildPipeline(m.build, target, m.downloads)
					langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
					s.WriteString(fmt.Sprintf("%s %s\n", langStyle.Render(fmt.Sprintf("%-13s", string(target))), pipeline))
				}

				s.WriteString("\n")

				completed := 0
				building := 0
				for _, target := range languages {
					buildTarget := buildObj.BuildTarget(target)
					if buildTarget != nil {
						if buildTarget.IsCompleted() {
							completed++
						} else if buildTarget.IsInProgress() {
							building++
						}
					}
				}

				statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
				statusText := fmt.Sprintf("%d completed, %d building, %d pending\n",
					completed, building, len(languages)-completed-building)
				s.WriteString(statusStyle.Render(statusText))
			}
		},
	},
	{
		name: "help",
		view: func(m BuildModel, s *strings.Builder) {
			s.WriteString("\n")
			s.WriteString(ViewHelpMenu())
		},
	},
}

func ViewBuildPipeline(build *stainless.BuildObject, target stainless.Target, downloads map[stainless.Target]struct {
	status string
	path   string
}) string {
	buildObj := NewBuildObject(build)
	buildTarget := buildObj.BuildTarget(target)
	if buildTarget == nil {
		return ""
	}

	stepOrder := buildTarget.Steps()
	var pipeline strings.Builder

	for _, step := range stepOrder {
		status, url, conclusion := buildTarget.StepInfo(step)
		if status == "" {
			continue // Skip steps that don't exist for this target
		}
		symbol := ViewStepSymbol(status, conclusion)
		if pipeline.Len() > 0 {
			pipeline.WriteString(" → ")
		}
		pipeline.WriteString(symbol + " " + Hyperlink(url, step))
	}

	if download, ok := downloads[target]; ok {
		if download.status == "not started" {
			// do nothing
		} else if download.status == "started" {
			pipeline.WriteString(" → " + "downloading")
		} else if download.status == "completed" {
			pipeline.WriteString(" → " + "downloaded")
		} else {
			pipeline.WriteString(" → " + download.status)
		}
	}

	return pipeline.String()
}

func ViewStepSymbol(status, conclusion string) string {
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	switch status {
	case "not_started", "queued":
		return grayStyle.Render("○")
	case "in_progress":
		return yellowStyle.Render("●")
	case "completed":
		switch conclusion {
		case "success":
			return greenStyle.Render("✓")
		case "failure":
			return redStyle.Render("✗")
		case "warning":
			return yellowStyle.Render("⚠")
		default:
			return greenStyle.Render("✓")
		}
	default:
		return grayStyle.Render("○")
	}
}

func ViewDiagnosticIcon(level stainless.BuildDiagnosticListResponseLevel) string {
	switch level {
	case stainless.BuildDiagnosticListResponseLevelFatal:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("💀")
	case stainless.BuildDiagnosticListResponseLevelError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("❌")
	case stainless.BuildDiagnosticListResponseLevelWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("⚠️")
	case stainless.BuildDiagnosticListResponseLevelNote:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("ℹ️")
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("•")
	}
}

// ViewHelpMenu creates a styled help menu inspired by huh help component
func ViewHelpMenu() string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	helpItems := []struct {
		key  string
		desc string
	}{
		{"enter", "rebuild"},
		{"ctrl+c", "exit"},
	}

	var parts []string
	for _, item := range helpItems {
		parts = append(parts,
			keyStyle.Render(item.key)+
				sepStyle.Render(" ")+
				descStyle.Render(item.desc))
	}

	return strings.Join(parts, sepStyle.Render(" • "))
}

// renderMarkdown renders markdown content using glamour
func renderMarkdown(content string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		return content
	}

	rendered, err := r.Render(content)
	if err != nil {
		return content
	}

	return strings.Trim(rendered, "\n ")
}

func countDiagnosticsBySeverity(diagnostics []stainless.BuildDiagnosticListResponse) (fatal, errors, warnings, notes int) {
	for _, diag := range diagnostics {
		switch diag.Level {
		case stainless.BuildDiagnosticListResponseLevelFatal:
			fatal++
		case stainless.BuildDiagnosticListResponseLevelError:
			errors++
		case stainless.BuildDiagnosticListResponseLevelWarning:
			warnings++
		case stainless.BuildDiagnosticListResponseLevelNote:
			notes++
		}
	}
	return
}

func ViewDiagnosticsPrint(diagnostics []stainless.BuildDiagnosticListResponse) string {
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
		maxDiagnostics := 10

		if len(diagnostics) > maxDiagnostics {
			sub.WriteString(fmt.Sprintf("Showing first %d of %d diagnostics:\n", maxDiagnostics, len(diagnostics)))
		}

		for i, diag := range diagnostics {
			if i >= maxDiagnostics {
				break
			}

			levelIcon := ViewDiagnosticIcon(diag.Level)
			codeStyle := lipgloss.NewStyle().Bold(true)

			if i > 0 {
				sub.WriteString("\n")
			}
			sub.WriteString(fmt.Sprintf("%s %s\n", levelIcon, codeStyle.Render(diag.Code)))
			sub.WriteString(fmt.Sprintf("%s\n", renderMarkdown(diag.Message)))

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

		s.WriteString(SProperty(0, "diagnostics", summary))
		s.WriteString(lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			Render(strings.TrimRight(sub.String(), "\n")),
		)
		s.WriteString("\n\n")
	} else {
		s.WriteString(SProperty(0, "diagnostics", "(no errors or warnings)"))
	}

	return s.String()
}
