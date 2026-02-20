package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/logrusorgru/aurora/v4"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// NewProgram wraps tea.NewProgram with better handling for tty environments
func NewProgram(model tea.Model, opts ...tea.ProgramOption) *tea.Program {
	// Always output to stderr, in case we want to also output JSON so that the json is redirectable e.g. to jq.
	opts = append(opts, tea.WithOutput(os.Stderr))

	// If not a TTY, use stdin and disable renderer
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		opts = append(opts,
			tea.WithInput(os.Stdin),
			tea.WithoutRenderer(),
		)
	}

	return tea.NewProgram(model, opts...)
}

// Group represents a nested logging group
type Group struct {
	prefix string
	indent int
	silent bool
}

func NewGroup(silent bool) Group {
	return Group{silent: silent}
}

func Header(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SHeader(format, args...))
	return Group{indent: 1}
}

func Info(format string, args ...any) Group {
	return Group{}.Info(format, args...)
}

func Property(key, msg string) Group {
	return Group{}.Property(key, msg)
}

func Progress(format string, args ...any) Group {
	return Group{}.Progress(format, args...)
}

func Error(format string, args ...any) Group {
	return Group{}.Error(format, args...)
}

func Warn(format string, args ...any) Group {
	return Group{}.Warn(format, args...)
}

func Success(format string, args ...any) Group {
	return Group{}.Success(format, args...)
}

func Confirm(cmd *cli.Command, flagName, title, description string, defaultValue bool) (bool, error) {
	value, _, err := Group{}.Confirm(cmd, flagName, title, description, defaultValue)
	return value, err
}

func Field(field huh.Field) error {
	return Group{}.Field(field)
}

func Spacer() {
	fmt.Fprintf(os.Stderr, "\n")
}

func Hyperlink(url, text string) string {
	return ansi.SetHyperlink(url, "") + text + ansi.ResetHyperlink()
}

func SHeader(format string, args ...any) string {
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s\n", aurora.BgCyan(" "+msg+" ").Black().Bold())
}

func SInfo(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightBlue("✱"), msg)
}

func SProperty(indent int, key, msg string) string {
	indentStr := strings.Repeat("  ", indent)
	return fmt.Sprintf("%s%s %s %s\n", indentStr, aurora.Cyan("✱"), key, aurora.Gray(13, msg))
}

func SProgress(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.Cyan("✱"), msg)
}

func SError(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightRed("✗"), msg)
}

func SWarn(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightYellow("!"), msg)
}

func SSuccess(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightGreen("✓"), msg)
}

func (g Group) Info(format string, args ...any) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SInfo(g.indent, format, args...))
	}
	return Group{prefix: "i", indent: g.indent + 1}
}

func (g Group) Property(key, msg string) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SProperty(g.indent, key, msg))
	}
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Progress(format string, args ...any) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SProgress(g.indent, format, args...))
	}
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Error(format string, args ...any) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SError(g.indent, format, args...))
	}
	return Group{prefix: "✗", indent: g.indent + 1}
}

func (g Group) Warn(format string, args ...any) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SWarn(g.indent, format, args...))
	}
	return Group{prefix: "!", indent: g.indent + 1}
}

func (g Group) Success(format string, args ...any) Group {
	if !g.silent {
		fmt.Fprint(os.Stderr, SSuccess(g.indent, format, args...))
	}
	return Group{prefix: "✓", indent: g.indent + 1}
}

func (g Group) Confirm(cmd *cli.Command, flagName, title, description string, defaultValue bool) (bool, Group, error) {
	if cmd != nil && flagName != "" && cmd.IsSet(flagName) {
		return cmd.Bool(flagName), Group{prefix: "✱", indent: g.indent + 1}, nil
	}

	foreground := lipgloss.Color("15")
	cyanBright := lipgloss.Color("6")

	t := GetFormTheme(g.indent)
	t.Focused.Title = t.Focused.Title.Foreground(foreground)
	t.Focused.Base = t.Focused.Base.
		SetString("\b\b" + lipgloss.NewStyle().Foreground(cyanBright).Render("✱")).
		PaddingLeft(2)

	value := defaultValue
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Value(&value),
		),
	).WithTheme(t).WithKeyMap(GetFormKeyMap()).Run()

	return value, Group{prefix: "✱", indent: g.indent + 1}, err
}

func (g Group) Field(field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).WithTheme(GetFormTheme(g.indent)).Run()
}

// spinnerModel handles the spinner UI while executing an operation
type spinnerModel struct {
	spinner spinner.Model
	message string
	indent  int
	execute func() error
	err     error
	done    bool
}

type operationCompleteMsg struct {
	err error
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runOperation)
}

func (m spinnerModel) runOperation() tea.Msg {
	err := m.execute()
	return operationCompleteMsg{err: err}
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case operationCompleteMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		return ""
	}
	indentStr := strings.Repeat("  ", m.indent)
	return indentStr + m.spinner.View() + " " + m.message
}

// Spinner runs the given operation with a spinner and message
func Spinner(message string, operation func() error) error {
	return spinnerWithIndent(0, message, operation)
}

// Spinner runs the given operation with a spinner and message, respecting the group's indent
func (g Group) Spinner(message string, operation func() error) error {
	return spinnerWithIndent(g.indent, message, operation)
}

func spinnerWithIndent(indent int, message string, operation func() error) error {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	model := spinnerModel{
		spinner: s,
		message: message,
		indent:  indent,
		execute: operation,
	}

	finalModel, err := NewProgram(model).Run()
	if err != nil {
		return err
	}

	if m, ok := finalModel.(spinnerModel); ok && m.err != nil {
		return m.err
	}

	return nil
}
