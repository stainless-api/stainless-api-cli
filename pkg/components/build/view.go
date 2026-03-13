package build

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

var (
	headerLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")).Bold(true)
	headerIDStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// ViewHeader renders a styled header with a label badge, build ID, config commit, and relative timestamp.
func ViewHeader(label string, b stainless.Build) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(headerLabelStyle.Render(" " + label + " "))
	if b.ID != "" {
		s.WriteString("  ")
		s.WriteString(headerIDStyle.Render(b.ID))
	}
	configCommit := b.ConfigCommit
	if len(configCommit) > 7 {
		configCommit = configCommit[:7]
	}
	if configCommit != "" {
		s.WriteString("  ")
		s.WriteString(headerIDStyle.Render(configCommit))
	}
	if !b.CreatedAt.IsZero() {
		s.WriteString("  ")
		s.WriteString(headerIDStyle.Render(relativeTime(b.CreatedAt)))
	}
	s.WriteString("\n\n")
	return s.String()
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	s := strings.Builder{}
	buildObj := stainlessutils.NewBuild(m.Build)
	languages := buildObj.Languages()
	for _, target := range languages {
		s.WriteString(ViewBuildPipeline(m.Build, target, m.Downloads, m.CommitOnly, m.Spinner))
	}

	return s.String()
}

// commitStatusWidth is the fixed visible width for the commit status column,
// based on the longest expected content: "71d249c (unchanged) with error diagnostic(s)"
const commitStatusWidth = 44

// ViewBuildPipeline renders the build pipeline for a target on a single line.
// Format: <language> <commit-status (padded)> <step symbols>
func ViewBuildPipeline(build stainless.Build, target stainless.Target, downloads map[stainless.Target]DownloadStatus, commitOnly bool, sp spinner.Model) string {
	langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	buildObj := stainlessutils.NewBuild(build)
	buildTarget := buildObj.BuildTarget(target)
	if buildTarget == nil {
		return ""
	}

	// Build commit status text
	var commitStatus strings.Builder
	commitStep := buildTarget.Commit
	switch commitStep.Status {
	case "", "not_started", "queued":
		commitStatus.WriteString(grayStyle.Render("queued"))
	case "in_progress":
		commitStatus.WriteString(grayStyle.Render("generating ") + sp.View())
	case "completed":
		conclusion := commitStep.Conclusion
		switch conclusion {
		case "merge_conflict", "upstream_merge_conflict":
			pr := commitStep.MergeConflictPr
			prURL := fmt.Sprintf("https://github.com/%s/%s/pull/%.0f", pr.Repo.Owner, pr.Repo.Name, pr.Number)
			commitStatus.WriteString(yellowStyle.Render(console.Hyperlink(prURL, fmt.Sprintf("merge conflict #%.0f", pr.Number))))
		case "fatal":
			commitStatus.WriteString(redStyle.Render("fatal error"))
		case "payment_required":
			commitStatus.WriteString(redStyle.Render("payment required"))
		case "cancelled":
			commitStatus.WriteString(grayStyle.Render("cancelled"))
		case "timed_out":
			commitStatus.WriteString(redStyle.Render("timed out"))
		case "noop":
			commitStatus.WriteString(grayStyle.Render("no-op"))
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
			commitStatus.WriteString(console.Hyperlink(commitURL, sha))
			if additions > 0 || deletions > 0 {
				commitStatus.WriteString(" " + grayStyle.Render("(") +
					greenStyle.Render(fmt.Sprintf("+%d", additions)) +
					grayStyle.Render("/") +
					redStyle.Render(fmt.Sprintf("-%d", deletions)) +
					grayStyle.Render(")"))
			} else {
				commitStatus.WriteString(" " + grayStyle.Render("(unchanged)"))
			}
			switch conclusion {
			case "error":
				commitStatus.WriteString(" with " + redStyle.Render("error") + " diagnostic(s)")
			case "warning":
				commitStatus.WriteString(" with " + yellowStyle.Render("warning") + " diagnostic(s)")
			}
		default:
			commitStatus.WriteString(grayStyle.Render(conclusion))
		}
	}

	// Pad commit status to fixed width so step symbols align vertically
	statusStr := commitStatus.String()
	if pad := commitStatusWidth - lipgloss.Width(statusStr); pad > 0 {
		statusStr += strings.Repeat(" ", pad)
	}

	// Build the line
	var line strings.Builder
	line.WriteString(langStyle.Render(fmt.Sprintf("%-13s", string(target))) + " ")
	line.WriteString(statusStr)

	// Collect post-commit steps + download (only when commit step is completed)
	var stepParts []string
	if !commitOnly && commitStep.Status == "completed" {
		for _, step := range buildTarget.Steps() {
			if step == "commit" {
				continue
			}
			stepStatus, stepURL, stepConclusion := buildTarget.StepInfo(step)
			if stepStatus == "" {
				continue
			}
			stepLabel := step
			if stepURL != "" {
				stepLabel = console.Hyperlink(stepURL, step)
			}
			stepParts = append(stepParts, ViewStepSymbol(stepStatus, stepConclusion)+" "+stepLabel)
		}
	}

	if download, ok := downloads[target]; ok {
		stepParts = append(stepParts, ViewStepSymbol(download.Status, download.Conclusion)+" "+"download")
		if download.Conclusion == "failure" && download.Error != "" {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
			line.WriteString("  " + strings.Join(stepParts, "  "))
			line.WriteString("\n" + errorStyle.Render("  Error: "+download.Error))
			line.WriteString("\n")
			return line.String()
		}
	}

	if len(stepParts) > 0 {
		line.WriteString("  " + strings.Join(stepParts, "  "))
	}
	line.WriteString("\n")

	return line.String()
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
		case "cancelled", "skipped":
			return grayStyle.Render("⊘")
		case "merge_conflict", "upstream_merge_conflict":
			return yellowStyle.Render("m")
		default:
			return grayStyle.Render(conclusion)
		}
	default:
		return grayStyle.Render("○")
	}
}
