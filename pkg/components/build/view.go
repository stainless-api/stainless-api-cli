package build

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

func (m Model) View() string {
	return View(m.Build, m.Downloads)
}

func View(build stainless.Build, downloads map[stainless.Target]DownloadStatus) string {
	s := strings.Builder{}
	buildObj := stainlessutils.NewBuild(build)
	languages := buildObj.Languages()
	// Target rows with colors
	for _, target := range languages {
		pipeline := ViewBuildPipeline(build, target, downloads)
		langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)

		s.WriteString(fmt.Sprintf("%s %s\n", langStyle.Render(fmt.Sprintf("%-13s", string(target))), pipeline))
	}

	// s.WriteString("\n")

	// completed := 0
	// building := 0
	// for _, target := range languages {
	// 	buildTarget := buildObj.BuildTarget(target)
	// 	if buildTarget != nil {
	// 		if buildTarget.IsCompleted() {
	// 			completed++
	// 		} else if buildTarget.IsInProgress() {
	// 			building++
	// 		}
	// 	}
	// }

	// 		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	// 		statusText := fmt.Sprintf("%d completed, %d building, %d pending\n",
	// 			completed, building, len(languages)-completed-building)
	// 		s.WriteString(statusStyle.Render(statusText))

	return s.String()
}

// View renders the build pipeline for a target
func ViewBuildPipeline(build stainless.Build, target stainless.Target, downloads map[stainless.Target]DownloadStatus) string {
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
			pipeline.WriteString("  ")
		}
		// align our naming of the commit step with the 
		if step == "commit" {
			step = "codegen"
		}
		pipeline.WriteString(symbol + " " + console.Hyperlink(url, step))
	}

	if download, ok := downloads[target]; ok {
		if download.Status == "not started" {
			// do nothing
		} else if download.Status == "started" {
			pipeline.WriteString("  " + "downloading")
		} else if download.Status == "completed" {
			pipeline.WriteString("  " + "downloaded")
		} else {
			pipeline.WriteString("  " + download.Status)
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
