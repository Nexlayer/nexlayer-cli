// Formatted with gofmt -s
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Editor represents the YAML editor component
type Editor struct {
	textarea textarea.Model
	err      error
}

// NewEditor creates a new YAML editor
func NewEditor(content string) Editor {
	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()

	return Editor{
		textarea: ta,
	}
}

// Init initializes the editor
func (e Editor) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles editor updates
func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return e, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			// TODO: Implement save functionality
			return e, nil
		}
	}

	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

// View renders the editor
func (e Editor) View() string {
	if e.err != nil {
		return fmt.Sprintf("Error: %v", e.err)
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1)

	return style.Render(e.textarea.View()) + "\n\n" +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("Press Ctrl+S to save, ESC to exit")
}

// GetContent returns the editor content
func (e Editor) GetContent() string {
	return e.textarea.Value()
}

// SetContent sets the editor content
func (e *Editor) SetContent(content string) {
	e.textarea.SetValue(content)
}
