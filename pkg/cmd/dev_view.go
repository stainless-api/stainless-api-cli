package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
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
				s.WriteString(SProperty(0, "Diagnostics", "waiting"))
			} else {
				s.WriteString(ViewDiagnosticsPrint(m.diagnostics))
			}
		},
	},
	{
		name: "duration",
		view: func(m BuildModel, s *strings.Builder) {
			duration := m.getBuildDuration()
			s.WriteString(SProperty(0, "Duration", duration.Round(time.Second).String()))
			s.WriteString("\n")
		},
	},
	{
		name: "build_status",
		view: func(m BuildModel, s *strings.Builder) {
			if m.build != nil {
				languages := getBuildLanguages(m.build)
				// Target rows with colors
				for _, target := range languages {
					pipeline := ViewPipeline(m.build, target)
					langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
					s.WriteString(fmt.Sprintf("%-13s %s\n", langStyle.Render(string(target)), pipeline))
				}

				s.WriteString("\n")

				completed := 0
				building := 0
				for _, target := range languages {
					if isBuildTargetCompleted(m.build, target) {
						completed++
					} else if isBuildTargetInProgress(m.build, target) {
						building++
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

func ViewPipeline(build *stainless.BuildObject, target stainless.Target) string {
	buildTarget := getBuildTarget(build, target)
	if buildTarget == nil {
		return ""
	}

	stepOrder := getBuildSteps(buildTarget)
	var pipeline strings.Builder

	for _, step := range stepOrder {
		stepUnion := getStepUnion(buildTarget, step)
		if stepUnion == nil {
			continue // Skip steps that don't exist for this target
		}
		symbol := ViewStepSymbol(stepUnion)
		if pipeline.Len() > 0 {
			pipeline.WriteString(" → ")
		}
		pipeline.WriteString(symbol + " " + step)
	}

	return pipeline.String()
}

func ViewStepSymbol(stepUnion any) string {
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	status, conclusion := extractStepInfo(stepUnion)

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

// buildTargetOptions creates huh options from the list of configured targets
func buildTargetOptions(configuredTargets []stainless.Target) []huh.Option[string] {
	var options []huh.Option[string]

	targetDisplayNames := map[stainless.Target]string{
		"typescript": "TypeScript",
		"python":     "Python",
		"go":         "Go",
		"java":       "Java",
		"kotlin":     "Kotlin",
		"ruby":       "Ruby",
		"terraform":  "Terraform",
		"cli":        "CLI",
		"php":        "PHP",
		"csharp":     "C#",
		"node":       "Node.js",
	}

	for _, target := range configuredTargets {
		displayName := targetDisplayNames[target]
		if displayName == "" {
			displayName = string(target)
		}
		options = append(options, huh.NewOption(displayName, string(target)))
	}

	return options
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

func ViewDiagnosticsPrint(diagnostics []stainless.BuildDiagnosticListResponse) string {
	var s strings.Builder

	if len(diagnostics) > 0 {
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

		s.WriteString(SProperty(0, "Diagnostics", ""))
		s.WriteString(lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			Render(strings.TrimRight(sub.String(), "\n")),
		)
		s.WriteString("\n\n")
	} else {
		s.WriteString(SProperty(0, "Diagnostics", "(no diagnostics)"))
	}

	return s.String()
}
