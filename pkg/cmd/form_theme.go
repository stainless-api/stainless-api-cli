package cmd

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// GetFormTheme returns the standard huh theme used across all forms
func GetFormTheme() *huh.Theme {
	theme := huh.ThemeBase()
	theme.Group.Title = theme.Group.Title.Foreground(lipgloss.Color("8")).PaddingBottom(1)

	theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color("6")).Bold(true)
	theme.Focused.Base = theme.Focused.Base.BorderForeground(lipgloss.Color("240"))
	theme.Focused.Description = theme.Focused.Description.Foreground(lipgloss.Color("8"))
	theme.Focused.TextInput.Placeholder = theme.Focused.TextInput.Placeholder.Foreground(lipgloss.Color("8"))

	theme.Blurred.Title = theme.Blurred.Title.Foreground(lipgloss.Color("6")).Bold(true)
	theme.Blurred.Base = theme.Blurred.Base.Foreground(lipgloss.Color("251"))
	theme.Blurred.Description = theme.Blurred.Description.Foreground(lipgloss.Color("8"))
	theme.Blurred.TextInput.Placeholder = theme.Blurred.TextInput.Placeholder.Foreground(lipgloss.Color("8"))

	return theme
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
		key.WithHelp("tab/↓/enter", "next"),
	)
	keyMap.Input.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "previous"),
	)
	return keyMap
}
