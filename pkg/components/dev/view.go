package dev

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
)

func (m Model) View() string {
	s := strings.Builder{}

	idx := slices.IndexFunc(parts, func(part ViewPart) bool {
		return part.Name == m.view
	}) + 1

	for i := idx; i < len(parts); i++ {
		parts[i].View(&m, &s)
	}

	return s.String()
}

// ViewPart represents a single part of the build view
type ViewPart struct {
	Name string
	View func(*Model, *strings.Builder)
}

var parts = []ViewPart{
	{
		Name: "header",
		View: func(m *Model, s *strings.Builder) {
			buildIDStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")).Bold(true)
			if m.Build.ID != "" {
				fmt.Fprintf(s, "\n\n%s %s\n\n", buildIDStyle.Render(" BUILD "), m.Build.ID)
			} else {
				fmt.Fprintf(s, "\n\n%s\n\n", buildIDStyle.Render(" BUILD "))
			}
		},
	},
	{
		Name: "build diagnostics",
		View: func(m *Model, s *strings.Builder) {
			if m.Diagnostics.Diagnostics == nil {
				s.WriteString(console.SProperty(0, "build diagnostics", "(waiting for build to finish)"))
			} else {
				s.WriteString(m.Diagnostics.View())
			}
		},
	},
	{
		Name: "studio",
		View: func(m *Model, s *strings.Builder) {
			if m.Build.ID != "" {
				url := fmt.Sprintf("https://app.stainless.com/%s/%s/studio?branch=%s", m.Build.Org, m.Build.Project, m.Branch)
				s.WriteString(console.SProperty(0, "studio", console.Hyperlink(url, url)))
			}
		},
	},
	{
		Name: "build_status",
		View: func(m *Model, s *strings.Builder) {
			s.WriteString("\n")
			s.WriteString(m.Build.View())
		},
	},
	{
		Name: "help",
		View: func(m *Model, s *strings.Builder) {
			s.WriteString("\n")
			s.WriteString(m.Help.View(m))
		},
	},
}

// updateView updates the view state and prints everything from current state to target state to scrollback
func (m *Model) updateView(targetState string) tea.Cmd {
	// Don't update if targetState is behind current index
	currentIndex := slices.IndexFunc(parts, func(part ViewPart) bool {
		return part.Name == m.view
	})
	targetIndex := slices.IndexFunc(parts, func(part ViewPart) bool {
		return part.Name == targetState
	})

	if targetIndex < currentIndex {
		return nil
	}

	output := ViewBuildRange(m, m.view, targetState)

	// Update model state
	m.view = targetState

	// Print to scrollback
	if len(output) > 0 {
		if output[len(output)-1] == '\n' {
			return tea.Println(output[:len(output)-1])
		} else {
			return tea.Println(output)
		}
	}
	return nil
}

// ViewBuildRange renders the build view from startView to endView (inclusive)
// Returns empty string if views are not found
func ViewBuildRange(m *Model, startView, endView string) string {
	var output strings.Builder

	startIndex := slices.IndexFunc(parts, func(part ViewPart) bool {
		return part.Name == startView
	}) + 1

	endIndex := slices.IndexFunc(parts, func(part ViewPart) bool {
		return part.Name == endView
	})

	for i := startIndex; i <= endIndex; i++ {
		parts[i].View(m, &output)
	}

	return output.String()
}
