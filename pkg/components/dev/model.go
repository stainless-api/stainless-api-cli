package dev

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/components/diagnostics"
	"github.com/stainless-api/stainless-api-cli/pkg/stainlessutils"
	"github.com/stainless-api/stainless-api-go"
)

var ErrUserCancelled = errors.New("user cancelled")

// WaitMode represents the level of waiting for build completion.
type WaitMode int

const (
	WaitNone   WaitMode = iota // Don't wait
	WaitCommit                 // Wait for commit only
	WaitAll                    // Wait for everything including workflows
)

type Model struct {
	Err error

	Client      stainless.Client
	Ctx         context.Context
	Watch       bool
	IsCompleted bool
	start       func() (*stainless.Build, error)
	Branch      string
	view        string
	label    string
	waitMode WaitMode
	Indent   string

	// models

	Help        help.Model
	Build       build.Model
	Diagnostics diagnostics.Model
}

type ErrorMsg error
type FileChangeMsg struct{}

type ModelConfig struct {
	Client        stainless.Client
	Ctx           context.Context
	Branch        string
	Start         func() (*stainless.Build, error)
	DownloadPaths map[stainless.Target]string
	Watch         bool
	Label         string   // Header label, defaults to "PREVIEW"
	WaitMode      WaitMode // When non-zero, auto-quits after diagnostics are fetched and build targets reach completion
	Indent        string   // Prefix for every non-empty output line (e.g. "  ")
}

func NewModel(cfg ModelConfig) Model {
	label := cfg.Label
	if label == "" {
		label = "PREVIEW"
	}
	return Model{
		start:       cfg.Start,
		Client:      cfg.Client,
		Ctx:         cfg.Ctx,
		Branch:      cfg.Branch,
		Watch:       cfg.Watch,
		label:    label,
		waitMode: cfg.WaitMode,
		Indent:   cfg.Indent,
		Help:        help.New(),
		Build:       build.NewModel(cfg.Client, cfg.Ctx, stainless.Build{}, cfg.Branch, cfg.DownloadPaths),
		Diagnostics: diagnostics.NewModel(cfg.Client, cfg.Ctx, nil),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Build.Init(),
		m.Diagnostics.Init(),
		func() tea.Msg {
			res, err := m.start()
			if err != nil {
				return ErrorMsg(err)
			}
			return build.FetchBuildMsg(*res)
		},
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

	case build.TickMsg, build.DownloadMsg, build.ErrorMsg, spinner.TickMsg:
		m.Build, cmd = m.Build.Update(msg)
		cmds = append(cmds, cmd)
		if m.Build.Err != nil {
			m.Err = m.Build.Err
		}

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

	// Auto-quit when WaitMode is set and build targets have reached completion
	if m.waitMode > WaitNone && m.diagnosticsFetched() && m.isComplete() {
		return m, tea.Sequence(tea.Batch(cmds...), tea.Quit)
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

func (m Model) diagnosticsFetched() bool {
	return m.Diagnostics.Diagnostics != nil || m.Diagnostics.Err != nil
}

func (m Model) isComplete() bool {
	buildObj := stainlessutils.NewBuild(m.Build.Build)
	for _, target := range buildObj.Languages() {
		buildTarget := buildObj.BuildTarget(target)
		if buildTarget == nil {
			return false
		}

		// Check if download is completed (if applicable)
		if buildTarget.IsCommitCompleted() && buildTarget.IsGoodCommitConclusion() {
			if download, ok := m.Build.Downloads[target]; ok {
				if download.Status != "completed" {
					return false
				}
			}
		}

		// Check if target is done based on wait mode
		done := buildTarget.IsCommitCompleted()
		if m.waitMode >= WaitAll {
			done = buildTarget.IsCompleted()
		}

		if !done {
			return false
		}
	}
	return true
}
