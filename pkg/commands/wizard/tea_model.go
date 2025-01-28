package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7571F9")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A5A5A")).
			MarginBottom(1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7571F9")).
			Bold(true)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A5A5A"))

	cursorStyle = focusedStyle.Copy()

	noStyle = lipgloss.NewStyle()
)

type step int

const (
	projectName step = iota
	templateSelect
	envVars
	done
)

type model struct {
	step        step
	projectName textinput.Model
	template    string
	templates   []string
	cursor      int
	spinner     spinner.Model
	quitting    bool
	err         error
}

func initialModel() model {
	// Initialize text input for project name
	ti := textinput.New()
	ti.Placeholder = "awesome-app"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 40

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		step:        projectName,
		projectName: ti,
		templates: []string{
			"langchain-nextjs",
			"langchain-fastapi",
			"mern",
			"pern",
			"mean",
		},
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.step {
			case projectName:
				if m.projectName.Value() != "" {
					m.step = templateSelect
				}
			case templateSelect:
				m.template = m.templates[m.cursor]
				m.step = done
				return m, tea.Quit
			}

		case "up", "k":
			if m.step == templateSelect {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.templates) - 1
				}
			}

		case "down", "j":
			if m.step == templateSelect {
				m.cursor++
				if m.cursor >= len(m.templates) {
					m.cursor = 0
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case projectName:
		m.projectName, cmd = m.projectName.Update(msg)
		return m, cmd
	case templateSelect:
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var s strings.Builder

	switch m.step {
	case projectName:
		s.WriteString(titleStyle.Render("🚀 Welcome to Nexlayer!"))
		s.WriteString("\n")
		s.WriteString(subtitleStyle.Render("Let's create something amazing together."))
		s.WriteString("\n\n")
		s.WriteString("Project name: ")
		s.WriteString(m.projectName.View())
		s.WriteString("\n")
		s.WriteString(subtitleStyle.Render("Press Enter to continue"))

	case templateSelect:
		s.WriteString(titleStyle.Render("📦 Choose a Template"))
		s.WriteString("\n")
		s.WriteString(subtitleStyle.Render("Select the perfect starting point for your project."))
		s.WriteString("\n\n")

		for i, choice := range m.templates {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				choice = focusedStyle.Render(choice)
			} else {
				choice = blurredStyle.Render(choice)
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
		}

	case done:
		s.WriteString(titleStyle.Render("🎉 All Set!"))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Project: %s\n", focusedStyle.Render(m.projectName.Value())))
		s.WriteString(fmt.Sprintf("Template: %s\n", focusedStyle.Render(m.template)))
		s.WriteString("\n")
		s.WriteString("Run the following command to deploy:\n")
		s.WriteString(focusedStyle.Render("  nexlayer deploy\n"))
	}

	return s.String()
}
