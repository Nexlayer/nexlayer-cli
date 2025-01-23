package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleHeading = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B61FF")).
			Bold(true).
			Padding(1, 0)

	styleSelected = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1).
			MarginTop(1)

	styleUnselected = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.NoColor{}).
			Padding(1).
			MarginTop(1)
)

// suggestion represents a single AI suggestion
type suggestion struct {
	title   string
	content string
}

func (s suggestion) Title() string       { return s.title }
func (s suggestion) Description() string { return "" }
func (s suggestion) FilterValue() string { return s.title }

type model struct {
	spinner    spinner.Model
	textInput  textinput.Model
	list       list.Model
	client     AIClient
	err        error
	loading    bool
	renderer   *glamour.TermRenderer
	showDetail bool
	showHelp   bool
}

func newModel(client AIClient) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7B61FF"))

	ti := textinput.New()
	ti.Placeholder = "What would you like help with?"
	ti.Focus()

	// Initialize glamour renderer for Markdown
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	// Initialize empty list
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Suggestions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styleHeading
	l.Styles.NoItems = lipgloss.NewStyle().Margin(1, 0)

	return model{
		spinner:   s,
		textInput: ti,
		list:      l,
		client:    client,
		renderer:  renderer,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.loading && !m.showDetail && m.textInput.Value() != "" {
				m.loading = true
				m.list.SetItems(nil)
				return m, m.getSuggestions
			}
		case "esc":
			if m.showDetail {
				m.showDetail = false
				return m, nil
			}
			return m, tea.Quit
		case "tab":
			if !m.loading && len(m.list.Items()) > 0 {
				m.showDetail = !m.showDetail
			}
			return m, nil
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := styleHeading.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case suggestionsMsg:
		m.loading = false
		items := make([]list.Item, len(msg.suggestions))
		for i, s := range msg.suggestions {
			// Extract title from the first line
			parts := strings.SplitN(s, "\n", 2)
			title := parts[0]
			content := s
			if len(parts) > 1 {
				content = parts[1]
			}
			items[i] = suggestion{title: title, content: content}
		}
		m.list.SetItems(items)
		m.err = msg.err
		return m, nil
	}

	if !m.loading && !m.showDetail {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	if !m.loading && len(m.list.Items()) > 0 {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(styleHeading.Render("🤖 Nexlayer AI Suggest"))
	s.WriteString("\n\n")

	if m.showHelp {
		s.WriteString("Keyboard Shortcuts:\n")
		s.WriteString("  Enter: Submit query\n")
		s.WriteString("  Tab: Toggle view\n")
		s.WriteString("  ?: Toggle this help screen\n")
		s.WriteString("  q: Quit\n")
		s.WriteString("  Esc: Go back or quit\n")
		return s.String()
	}

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Thinking...\n", m.spinner.View()))
	} else if m.showDetail && len(m.list.Items()) > 0 {
		// Show detailed view of selected suggestion
		item := m.list.SelectedItem().(suggestion)
		rendered, _ := m.renderer.Render(item.content)
		s.WriteString(styleSelected.Render(rendered))
		s.WriteString("\n\nPress ESC to go back, TAB to toggle view\n")
	} else {
		s.WriteString(m.textInput.View())
		s.WriteString("\n\n")

		if m.err != nil {
			s.WriteString(fmt.Sprintf("Error: %v\n", m.err))
		}

		if len(m.list.Items()) > 0 {
			s.WriteString(m.list.View())
			s.WriteString("\n\nPress TAB to view details, q to quit\n")
		}
	}

	return s.String()
}

type suggestionsMsg struct {
	suggestions []string
	err         error
}

func (m model) getSuggestions() tea.Msg {
	ctx := context.Background()
	suggestions, err := m.client.GetSuggestions(ctx, "general", map[string]interface{}{
		"query": m.textInput.Value(),
	})
	return suggestionsMsg{suggestions, err}
}

// RunModsUI starts the mods-powered UI
func RunModsUI(client AIClient) error {
	p := tea.NewProgram(newModel(client))
	_, err := p.Run()
	return err
}
