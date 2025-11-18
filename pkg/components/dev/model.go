package dev

import (
	"context"
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/components/diagnostics"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

var ErrUserCancelled = errors.New("user cancelled")

type Model struct {
	Err error

	Client      stainless.Client
	Ctx         context.Context
	Watch       bool
	IsCompleted bool
	start       func() (*stainless.Build, error)
	Branch      string
	view        string

	// models

	Help        help.Model
	Build       build.Model
	Diagnostics diagnostics.Model
}

type TickMsg time.Time
type ErrorMsg error
type FileChangeMsg struct{}

func NewModel(client stainless.Client, ctx context.Context, branch string, fn func() (*stainless.Build, error), downloadPaths map[stainless.Target]string, watch bool) Model {
	return Model{
		start:       fn,
		Client:      client,
		Ctx:         ctx,
		Branch:      branch,
		Help:        help.New(),
		Build:       build.NewModel(client, ctx, stainless.Build{}, downloadPaths),
		Diagnostics: diagnostics.NewModel(client, ctx, nil),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		}),
		func() tea.Msg {
			res, err := m.start()
			if err != nil {
				return ErrorMsg(err)
			}
			return build.FetchBuildMsg(*res)
		},
		m.Build.Init(),
		m.Diagnostics.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Help.Width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Err = ErrUserCancelled
			cmds = append(cmds, tea.Quit)
		case "enter":
			if m.Watch {
				cmds = append(cmds, tea.Quit)
			}
		}

	case build.TickMsg, build.DownloadMsg, build.ErrorMsg:
		m.Build, cmd = m.Build.Update(msg)
		cmds = append(cmds, cmd)

	case diagnostics.FetchDiagnosticsMsg:
		m.Diagnostics, cmd = m.Diagnostics.Update(msg)
		cmds = append(cmds, tea.Sequence(
			m.updateView("build diagnostics"),
			cmd,
		))

	case diagnostics.ErrorMsg:
		m.Diagnostics, cmd = m.Diagnostics.Update(msg)
		cmds = append(cmds, cmd)

	case build.FetchBuildMsg:
		// Check if all commit steps are completed, and if so trigger the diagnostics fetch process.
		if !m.IsCompleted {
			allCommitsCompleted := true
			buildObj := stainlessutils.NewBuild(stainless.Build(msg))
			for _, target := range buildObj.Languages() {
				buildTarget := buildObj.BuildTarget(target)
				if buildTarget != nil && !buildTarget.IsCommitCompleted() {
					allCommitsCompleted = false
					break
				}
			}
			if allCommitsCompleted {
				m.IsCompleted = true
				cmds = append(cmds, m.Diagnostics.FetchDiagnostics(buildObj.Build.ID))
			}
		}

		m.Build, cmd = m.Build.Update(msg)
		cmds = append(cmds, tea.Sequence(
			m.updateView("header"),
			cmd,
		))

	case ErrorMsg:
		m.Err = msg
		cmds = append(cmds, tea.Quit)

	case FileChangeMsg:
		// File change detected, exit with success
		cmds = append(cmds, tea.Quit)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) ShortHelp() []key.Binding {
	if m.Watch {
		return []key.Binding{
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rebuild")),
		}
	} else {
		return []key.Binding{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}
	}
}

func (m Model) FullHelp() [][]key.Binding {
	if m.Watch {
		return [][]key.Binding{{
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rebuild")),
		}}
	} else {
		return [][]key.Binding{{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}}
	}
}
