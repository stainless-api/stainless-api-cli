package stainlessviews

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

// ViewData contains all the data needed to render the build view
type ViewData struct {
	Build       *stainless.Build
	Diagnostics []stainless.BuildDiagnostic
	Duration    time.Duration
	Branch      string
	Downloads   map[stainless.Target]DownloadStatus
	HelpText    string
}

// ViewPart represents a single part of the build view
type ViewPart struct {
	Name string
	View func(ViewData, *strings.Builder)
}

// ViewBuild renders the build view from a given state onwards
func ViewBuild(data ViewData, currentView string) string {
	s := strings.Builder{}

	startIndex := 0
	if currentView != "" {
		for i, part := range parts {
			if part.Name == currentView {
				startIndex = i + 1
				break
			}
		}
	}

	for i := startIndex; i < len(parts); i++ {
		parts[i].View(data, &s)
	}

	return s.String()
}

// ViewBuildRange renders the build view from startView to endView (inclusive)
// Returns empty string if views are not found
func ViewBuildRange(data ViewData, startView, endView string) string {
	var output strings.Builder

	startIndex := 0
	if startView != "" {
		found := false
		for i, part := range parts {
			if part.Name == startView {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			startIndex = -1
		}
	} else {
		startIndex = 0
	}

	endIndex := -1
	for i, part := range parts {
		if part.Name == endView {
			endIndex = i
			break
		}
	}

	if endIndex == -1 || startIndex == -1 {
		return ""
	}

	for i := startIndex; i <= endIndex; i++ {
		parts[i].View(data, &output)
	}

	return output.String()
}

// GetViewParts returns the ordered list of view parts
var parts = []ViewPart{
	{
		Name: "header",
		View: func(d ViewData, s *strings.Builder) {
			buildIDStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")).Bold(true)
			if d.Build != nil {
				fmt.Fprintf(s, "\n\n%s %s\n\n", buildIDStyle.Render(" BUILD "), d.Build.ID)
			} else {
				fmt.Fprintf(s, "\n\n%s\n\n", buildIDStyle.Render(" BUILD "))
			}
		},
	},
	{
		Name: "build diagnostics",
		View: func(d ViewData, s *strings.Builder) {
			if d.Diagnostics == nil {
				s.WriteString(console.SProperty(0, "build diagnostics", "(waiting for build to finish)"))
			} else {
				s.WriteString(ViewDiagnosticsPrint(d.Diagnostics, 10))
			}
		},
	},
	{
		Name: "duration",
		View: func(d ViewData, s *strings.Builder) {
			s.WriteString(console.SProperty(0, "duration", d.Duration.Round(time.Second).String()))
		},
	},
	{
		Name: "studio",
		View: func(d ViewData, s *strings.Builder) {
			if d.Build != nil {
				url := fmt.Sprintf("https://app.stainless.com/%s/%s/studio?branch=%s", d.Build.Org, d.Build.Project, d.Branch)
				s.WriteString(console.SProperty(0, "studio", console.Hyperlink(url, url)))
			}
		},
	},
	{
		Name: "build_status",
		View: func(d ViewData, s *strings.Builder) {
			s.WriteString("\n")
			if d.Build != nil {
				buildObj := stainlessutils.NewBuild(d.Build)
				languages := buildObj.Languages()
				// Target rows with colors
				for _, target := range languages {
					pipeline := ViewBuildPipeline(d.Build, target, d.Downloads)
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
		Name: "help",
		View: func(d ViewData, s *strings.Builder) {
			s.WriteString(d.HelpText)
		},
	},
}

// DownloadStatus represents the download status and path for a target
type DownloadStatus struct {
	Status string
	Path   string
}

// ViewBuildPipeline renders the build pipeline for a target
func ViewBuildPipeline(build *stainless.Build, target stainless.Target, downloads map[stainless.Target]DownloadStatus) string {
	buildObj := stainlessutils.NewBuild(build)
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
		pipeline.WriteString(symbol + " " + console.Hyperlink(url, step))
	}

	if download, ok := downloads[target]; ok {
		if download.Status == "not started" {
			// do nothing
		} else if download.Status == "started" {
			pipeline.WriteString(" → " + "downloading")
		} else if download.Status == "completed" {
			pipeline.WriteString(" → " + "downloaded")
		} else {
			pipeline.WriteString(" → " + download.Status)
		}
	}

	return pipeline.String()
}

// ViewStepSymbol returns a colored symbol for a build step status
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

// renderMarkdown renders markdown content using glamour
func renderMarkdown(content string) string {
	width, _, err := term.GetSize(uintptr(os.Stdout.Fd()))
	if err != nil || width <= 0 || width > 120 {
		width = 120
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
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

// ViewDiagnosticsPrint renders build diagnostics with formatting
func ViewDiagnosticsPrint(diagnostics []stainless.BuildDiagnostic, maxDiagnostics int) string {
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
