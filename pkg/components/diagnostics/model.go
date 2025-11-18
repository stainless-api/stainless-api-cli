package diagnostics

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stainless-api/stainless-api-go"
)

var ErrUserCancelled = errors.New("user cancelled")

type Model struct {
	Diagnostics []stainless.BuildDiagnostic
	Client      stainless.Client
	Ctx         context.Context
	Err         error
}

type FetchDiagnosticsMsg []stainless.BuildDiagnostic
type ErrorMsg error

func NewModel(client stainless.Client, ctx context.Context, diagnostics []stainless.BuildDiagnostic) Model {
	return Model{
		Client:      client,
		Ctx:         ctx,
		Diagnostics: diagnostics,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Err = ErrUserCancelled
			return m, tea.Quit
		}

	case FetchDiagnosticsMsg:
		m.Diagnostics = msg

	case ErrorMsg:
		m.Err = msg
	}

	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		return ViewDiagnosticsError(m.Err)
	}
	if m.Diagnostics == nil {
		return ""
	}
	return ViewDiagnostics(m.Diagnostics, 10)
}

func (m Model) FetchDiagnostics(buildID string) tea.Cmd {
	return func() tea.Msg {
		if buildID == "" {
			return ErrorMsg(errors.New("no build ID provided"))
		}

		diags := []stainless.BuildDiagnostic{}
		diagnostics := m.Client.Builds.Diagnostics.ListAutoPaging(m.Ctx, buildID, stainless.BuildDiagnosticListParams{
			Limit: stainless.Float(100),
		})

		for diagnostics.Next() {
			diag := diagnostics.Current()
			if !diag.Ignored {
				diags = append(diags, diag)
			}
		}

		if err := diagnostics.Err(); err != nil {
			return ErrorMsg(err)
		}

		return FetchDiagnosticsMsg(diags)
	}
}
