package console

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// GetFormTheme returns the standard huh theme used across all forms
func GetFormTheme(indent int) *huh.Theme {
	t := huh.ThemeBase()

	grayBright := lipgloss.Color("251")
	gray := lipgloss.Color("243")
	primary := lipgloss.Color("6")
	primaryBright := lipgloss.Color("14")
	error := lipgloss.Color("1")

	t.Form.Base = t.Form.Base.PaddingLeft(indent * 2)
	t.Group.Title = t.Group.Title.Foreground(primary).PaddingBottom(1)
	t.Group.Description = t.Group.Description.Foreground(gray)

	t.Focused.Title = t.Focused.Title.Foreground(primary).Bold(true)
	t.Focused.Base = t.Focused.Base.
		BorderLeft(false).
		SetString("\b\b" + lipgloss.NewStyle().Foreground(primaryBright).Render("✱")).
		PaddingLeft(2)
	t.Focused.Description = t.Focused.Description.Foreground(gray).Width(70)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(gray)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))

	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(error)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(error)

	t.Blurred.Title = t.Blurred.Title.Foreground(grayBright).Bold(true)
	t.Blurred.Base = t.Blurred.Base.
		Foreground(grayBright).
		BorderLeft(false).
		SetString("\b\b" + lipgloss.NewStyle().Foreground(primary).Render("✱")).
		PaddingLeft(2)
	t.Blurred.Description = t.Blurred.Description.Foreground(gray).Width(70)
	t.Blurred.TextInput.Placeholder = t.Blurred.TextInput.Placeholder.Foreground(gray)
	t.Blurred.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")

	return t
}

// GetFormKeyMap returns the standard huh keymap used across all forms
func GetFormKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Input.AcceptSuggestion = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "complete"),
	)
	keyMap.Input.Next = key.NewBinding(
		key.WithKeys("tab", "down", "enter"),
		key.WithHelp("enter", "next"),
	)
	keyMap.Input.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab", "back"),
	)
	return keyMap
}
