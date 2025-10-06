package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stainless-api/stainless-api-go"
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
	Action: runLinter,
}

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	Margin(1, 0, 0, 0)

type lintModel struct {
	spinner      spinner.Model
	diagnostics  []stainless.BuildDiagnostic
	error        error
	watching     bool
	watchedFiles map[string]time.Time
	ctx          context.Context
	cmd          *cli.Command
	cc           *apiCommandContext
	stopPolling  chan struct{}
	help         helpModel
}

type helpModel struct {
	width  int
	height int
}

type fileChangedEvent struct{}

func waitForFileChanges(m lintModel) tea.Cmd {
	return func() tea.Msg {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for file, lastModTime := range m.watchedFiles {
					stat, err := os.Stat(file)
					if err != nil {
						continue
					}

					if stat.ModTime().After(lastModTime) {
						m.watchedFiles[file] = stat.ModTime()
						return fileChangedEvent{}
					}
				}
			case <-m.stopPolling:
				return nil
			}
		}
	}
}

func (m lintModel) Init() tea.Cmd {
	if m.watching {
		// Clear the screen and move the cursor to the top
		fmt.Print("\033[2J\033[H")
		os.Stdout.Sync()
	}
	return m.spinner.Tick
}

func (m lintModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.watching = false
			return m, tea.Quit
		}

	case diagnosticsMsg:
		m.diagnostics = msg.diagnostics
		m.error = msg.err
		m.ctx = msg.ctx
		m.cmd = msg.cmd
		m.cc = msg.cc

		if m.watching {
			return m, waitForFileChanges(m)
		}
		return m, tea.Quit

	case fileChangedEvent:
		m.diagnostics = nil // Clear diagnostics while linting
		return m, getDiagnosticsCmd(m.ctx, m.cmd, m.cc)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.help.width = msg.Width
		m.help.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m lintModel) View() string {
	if m.error != nil {
		return fmt.Sprintf("Error: %s\n", m.error)
	}

	var content string
	if m.diagnostics == nil {
		content = m.spinner.View() + " Linting"
	} else {
		content = ViewDiagnosticsPrint(m.diagnostics, -1)
		if m.watching {
			content += "\n" + m.spinner.View() + " Waiting for configuration changes"
		}
	}

	helpText := helpStyle.Render("Press Ctrl+C to exit")
	return content + helpText
}

type diagnosticsMsg struct {
	diagnostics []stainless.BuildDiagnostic
	err         error
	ctx         context.Context
	cmd         *cli.Command
	cc          *apiCommandContext
}

func getDiagnosticsCmd(ctx context.Context, cmd *cli.Command, cc *apiCommandContext) tea.Cmd {
	return func() tea.Msg {
		diagnostics, err := getDiagnostics(ctx, cmd, cc)
		return diagnosticsMsg{
			diagnostics: diagnostics,
			err:         err,
			ctx:         ctx,
			cmd:         cmd,
			cc:          cc,
		}
	}
}

func runLinter(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

	m := lintModel{
		spinner:      s,
		watching:     cmd.Bool("watch"),
		ctx:          ctx,
		cmd:          cmd,
		cc:           cc,
		watchedFiles: make(map[string]time.Time),
		stopPolling:  make(chan struct{}),
		help:         helpModel{},
	}

	if m.watching {
		if err := setupFileWatcher(&m, cmd, cc); err != nil {
			return err
		}
	}

	p := tea.NewProgram(m, tea.WithContext(ctx))

	// Start the diagnostics process
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to let the UI initialize
		p.Send(getDiagnosticsCmd(ctx, cmd, cc)())
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

	// If not in watch mode and we have blocking diagnostics, exit with error code
	if !cmd.Bool("watch") && hasBlockingDiagnostic(finalModel.diagnostics) {
		os.Exit(1)
	}

	return nil
}

func setupFileWatcher(m *lintModel, cmd *cli.Command, cc *apiCommandContext) error {
	// Watch OpenAPI spec file
	openapiSpecPath := cc.workspaceConfig.OpenAPISpec
	if cmd.IsSet("openapi-spec") {
		openapiSpecPath = cmd.String("openapi-spec")
	}

	if err := addFileToWatch(m, openapiSpecPath); err != nil {
		return err
	}

	// Watch Stainless config file
	stainlessConfigPath := cc.workspaceConfig.StainlessConfig
	if cmd.IsSet("stainless-config") {
		stainlessConfigPath = cmd.String("stainless-config")
	}

	return addFileToWatch(m, stainlessConfigPath)
}

func addFileToWatch(m *lintModel, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", path, err)
	}
	m.watchedFiles[path] = stat.ModTime()
	return nil
}
