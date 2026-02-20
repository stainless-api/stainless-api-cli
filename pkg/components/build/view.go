package build

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	s := strings.Builder{}
	buildObj := stainlessutils.NewBuild(m.Build)
	languages := buildObj.Languages()
	for _, target := range languages {
		s.WriteString(ViewBuildPipeline(m.Build, target, m.Downloads, m.CommitOnly, m.Spinner))
		s.WriteString("\n")
	}

	return s.String()
}

// ViewBuildPipeline renders the build pipeline for a target as one or two lines.
// Line 1: codegen status text + optional download
// Line 2: post-commit steps (lint/build/test), only when !commitOnly
func ViewBuildPipeline(build stainless.Build, target stainless.Target, downloads map[stainless.Target]DownloadStatus, commitOnly bool, sp spinner.Model) string {
	langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	buildObj := stainlessutils.NewBuild(build)
	buildTarget := buildObj.BuildTarget(target)
	if buildTarget == nil {
		return ""
	}

	var line1 strings.Builder
	line1.WriteString(langStyle.Render(fmt.Sprintf("%-13s", string(target))) + " ")

	// Render commit step status as human-readable text
	commitStep := buildTarget.Commit
	switch commitStep.Status {
	case "", "not_started", "queued":
		line1.WriteString(grayStyle.Render("queued"))
	case "in_progress":
		line1.WriteString(grayStyle.Render("generating ") + sp.View())
	case "completed":
		conclusion := commitStep.Conclusion
		switch conclusion {
		case "merge_conflict", "upstream_merge_conflict":
			pr := commitStep.MergeConflictPr
			prURL := fmt.Sprintf("https://github.com/%s/%s/pull/%.0f", pr.Repo.Owner, pr.Repo.Name, pr.Number)
			line1.WriteString(yellowStyle.Render(console.Hyperlink(prURL, "merge conflict")))
		case "fatal":
			line1.WriteString(redStyle.Render("fatal error"))
		case "payment_required":
			line1.WriteString(redStyle.Render("payment required"))
		case "cancelled":
			line1.WriteString(grayStyle.Render("cancelled"))
		case "timed_out":
			line1.WriteString(redStyle.Render("timed out"))
		case "noop":
			line1.WriteString(grayStyle.Render("no-op"))
		case "success", "note", "warning", "error", "version_bump":
			// These conclusions all produce a commit
			commit := commitStep.Commit
			sha := commit.Sha
			if len(sha) > 7 {
				sha = sha[:7]
			}
			commitURL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", commit.Repo.Owner, commit.Repo.Name, commit.Sha)
			additions := commit.Stats.Additions
			deletions := commit.Stats.Deletions
			line1.WriteString(console.Hyperlink(commitURL, sha))
			if additions > 0 || deletions > 0 {
				line1.WriteString(" " + grayStyle.Render("(") +
					greenStyle.Render(fmt.Sprintf("+%d", additions)) +
					grayStyle.Render("/") +
					redStyle.Render(fmt.Sprintf("-%d", deletions)) +
					grayStyle.Render(")"))
			} else {
				line1.WriteString(" " + grayStyle.Render("(unchanged)"))
			}
			switch conclusion {
			case "error":
				line1.WriteString(" with " + redStyle.Render("error") + " diagnostic(s)")
			case "warning":
				line1.WriteString(" with " + yellowStyle.Render("warning") + " diagnostic(s)")
			}
		default:
			line1.WriteString(grayStyle.Render(conclusion))
		}
	}

	// Line 2: post-commit steps (lint/build/test) + download
	// Always emit a second line for consistent layout.
	var line2Parts []string
	// 14 chars padding = "%-13s" label + 1 space separator
	padding := strings.Repeat(" ", 14)

	if !commitOnly {
		for _, step := range buildTarget.Steps() {
			if step == "commit" {
				continue
			}
			stepStatus, _, stepConclusion := buildTarget.StepInfo(step)
			if stepStatus == "" {
				continue
			}
			line2Parts = append(line2Parts, ViewStepSymbol(stepStatus, stepConclusion)+" "+step)
		}
	}

	if download, ok := downloads[target]; ok {
		line2Parts = append(line2Parts, ViewStepSymbol(download.Status, download.Conclusion)+" "+"download")
		if download.Conclusion == "failure" && download.Error != "" {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
			line1.WriteString("\n" + errorStyle.Render("  Error: "+download.Error))
		}
	}

	if len(line2Parts) > 0 {
		line1.WriteString("\n" + padding + strings.Join(line2Parts, "  "))
	} else {
		line1.WriteString("\n")
	}

	return line1.String()
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
		case "success", "note":
			return greenStyle.Render("✓")
		case "warning":
			return yellowStyle.Render("⚠")
		case "failure", "error":
			return redStyle.Render("⚠")
		case "fatal":
			return redStyle.Render("✗")
		case "merge_conflict", "upstream_merge_conflict":
			return yellowStyle.Render("m")
		default:
			return grayStyle.Render(conclusion)
		}
	default:
		return grayStyle.Render("○")
	}
}
