package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessviews"
)

func (m BuildModel) View() string {
	data := stainlessviews.ViewData{
		Build:       m.build,
		Diagnostics: m.diagnostics,
		Duration:    m.getBuildDuration(),
		Branch:      m.branch,
		Downloads:   m.downloads,
		HelpText:    m.help.View(m),
	}
	return stainlessviews.ViewBuild(data, m.view)
}

// updateView updates the view state and prints everything from current state to target state to scrollback
func (m *BuildModel) updateView(targetState string) tea.Cmd {
	data := stainlessviews.ViewData{
		Build:       m.build,
		Diagnostics: m.diagnostics,
		Duration:    m.getBuildDuration(),
		Branch:      m.branch,
		Downloads:   m.downloads,
		HelpText:    m.help.View(*m),
	}

	output := stainlessviews.ViewBuildRange(data, m.view, targetState)

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
