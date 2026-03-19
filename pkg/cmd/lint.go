package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-cli/pkg/components/build"
	"github.com/stainless-api/stainless-api-cli/pkg/components/diagnostics"
	"github.com/stainless-api/stainless-api-cli/pkg/console"
	"github.com/stainless-api/stainless-api-cli/pkg/workspace"
	"github.com/stainless-api/stainless-api-go"
	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var lintCommand = cli.Command{
	Name:  "lint",
	Usage: "Lint your stainless configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "Project name to use for the build",
		},
		&cli.StringFlag{
			Name:    "openapi-spec",
			Aliases: []string{"oas"},
			Usage:   "Path to OpenAPI spec file",
		},
		&cli.StringFlag{
			Name:    "stainless-config",
			Aliases: []string{"config"},
			Usage:   "Path to Stainless config file",
		},
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Watch for files to change and re-run linting",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Bool("watch") {
			// Clear the screen and move the cursor to the top
			fmt.Print("\033[2J\033[H")
			os.Stdout.Sync()
		}
		return runLinter(ctx, cmd)
	},
}

type lintModel struct {
	spinner     spinner.Model
	diagnostics []stainless.BuildDiagnostic
	error       error
	watching    bool
	skipped     bool
	ctx         context.Context
	cmd         *cli.Command
	client      stainless.Client
	wc          workspace.Config
	stopPolling chan struct{}
	help        help.Model
}

func (m lintModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m lintModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.watching = false
			m.error = build.ErrUserCancelled
			return m, tea.Quit
		} else if msg.String() == "enter" {
			m.watching = false
			m.skipped = true
			return m, tea.Quit
		}

	case diagnosticsMsg:
		m.diagnostics = msg.diagnostics
		m.error = msg.err
		m.ctx = msg.ctx
		m.cmd = msg.cmd
		m.client = msg.client

		if m.watching {
			return m, func() tea.Msg {
				if err := waitTillConfigChanges(m.ctx, m.cmd, m.wc); err != nil {
					log.Fatal(err)
				}
				return configChangedEvent{}
			}
		}
		return m, tea.Quit

	case configChangedEvent:
		m.diagnostics = nil // Clear diagnostics while linting
		return m, getDiagnosticsCmd(m.ctx, m.cmd, m.client)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		return m, nil
	}

	return m, nil
}

func (m lintModel) View() string {
	var content string
	if m.error != nil {
		content = "Linting failed!"
	} else if m.diagnostics == nil {
		if m.skipped {
			content = "Skipped!"
		} else {
			content = m.spinner.View() + " Linting"
		}
	} else {
		content = diagnostics.ViewDiagnostics(m.diagnostics, -1, workspace.Relative(m.wc.OpenAPISpec), workspace.Relative(m.wc.StainlessConfig))
		if m.skipped {
			content += "\nContinuing..."
		} else if m.watching {
			content += "\n" + m.spinner.View() + " Waiting for configuration changes"
		}
	}

	content += "\n" + m.help.View(m)
	return content
}

type diagnosticsMsg struct {
	diagnostics []stainless.BuildDiagnostic
	err         error
	ctx         context.Context
	cmd         *cli.Command
	client      stainless.Client
}

func getDiagnosticsCmd(ctx context.Context, cmd *cli.Command, client stainless.Client) tea.Cmd {
	return func() tea.Msg {
		diagnostics, err := getDiagnostics(ctx, cmd, client)
		return diagnosticsMsg{
			diagnostics: diagnostics,
			err:         err,
			ctx:         ctx,
			cmd:         cmd,
			client:      client,
		}
	}
}

type GenerateSpecParams struct {
	Project string `json:"project"`
	Source  struct {
		Type            string `json:"type"`
		OpenAPISpec     string `json:"openapi_spec"`
		StainlessConfig string `json:"stainless_config"`
	} `json:"source"`
}

func getDiagnostics(ctx context.Context, cmd *cli.Command, client stainless.Client) ([]stainless.BuildDiagnostic, error) {
	var specParams GenerateSpecParams
	specParams.Project = cmd.String("project")

	specParams.Source.Type = "upload"

	var (
		oas    []byte
		config []byte
		err    error
	)

	if _, config, err = convertFileFlag(cmd, "stainless-config"); err != nil {
		return nil, err
	}
	if _, oas, err = convertFileFlag(cmd, "openapi-spec"); err != nil {
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("You must provide a stainless configuration file with `--config /path/to/stainless.yml` or run this command from an initialized workspace.")
	}
	if oas == nil {
		return nil, fmt.Errorf("You must provide an OpenAPI specification with `--oas /path/to/openapi.json` or run this command from an initialized workspace.")
	}

	specParams.Source.StainlessConfig = string(config)
	specParams.Source.OpenAPISpec = string(oas)

	options := []option.RequestOption{}
	if cmd.Bool("debug") {
		options = append(options, debugMiddlewareOption)
	}
	var result []byte
	err = client.Post(
		ctx,
		"api/generate/spec",
		specParams,
		&result,
		options...,
	)
	if err != nil {
		return nil, err
	}

	transform := "spec.diagnostics.@values.@flatten.#(ignored==false)#"
	jsonObj := gjson.Parse(string(result)).Get(transform)
	var diagnostics []stainless.BuildDiagnostic
	json.Unmarshal([]byte(jsonObj.Raw), &diagnostics)
	return diagnostics, nil
}

func (m lintModel) ShortHelp() []key.Binding {
	return []key.Binding{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}
}

func (m lintModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{{key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl-c", "quit"))}}
}

func runLinter(ctx context.Context, cmd *cli.Command) error {
	client := stainless.NewClient(getDefaultRequestOptions(cmd)...)

	wc := getWorkspace(ctx)

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

	m := lintModel{
		spinner:     s,
		watching:    cmd.Bool("watch"),
		ctx:         ctx,
		cmd:         cmd,
		client:      client,
		wc:          wc,
		stopPolling: make(chan struct{}),
		help:        help.New(),
	}

	p := console.NewProgram(m, tea.WithContext(ctx))

	// Start the diagnostics process
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to let the UI initialize
		p.Send(getDiagnosticsCmd(ctx, cmd, client)())
	}()

	model, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := model.(lintModel)
	if finalModel.stopPolling != nil {
		close(finalModel.stopPolling)
	}

	if finalModel.error != nil {
		return finalModel.error
	}

	if !cmd.Bool("watch") {
		for _, d := range finalModel.diagnostics {
			if d.Ignored {
				continue
			}
			switch d.Level {
			case stainless.BuildDiagnosticLevelFatal, stainless.BuildDiagnosticLevelError, stainless.BuildDiagnosticLevelWarning:
				return cli.Exit("", 1)
			}
		}
	}

	return nil
}
