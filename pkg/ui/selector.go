package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type template struct {
	name, desc string
}

func (t template) Title() string       { return t.name }
func (t template) Description() string { return t.desc }
func (t template) FilterValue() string { return t.name }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.choice = m.list.SelectedItem().(template).name
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return fmt.Sprintf("Selected template: %s\n", m.choice)
	}
	if m.quitting {
		return "Operation cancelled\n"
	}
	return docStyle.Render(m.list.View())
}



func SelectTemplate(detectedTemplate string) (string, error) {
	templates := []list.Item{
		template{name: "🏗  Blank", desc: "A minimal template"},
		template{name: "🤖 Langchain.js", desc: "LangChain with Next.js"},
		template{name: "🚀 OpenAI Node.js", desc: "OpenAI API with Express + React"},
		template{name: "⚡ OpenAI Python", desc: "OpenAI API with FastAPI + Vue"},
		template{name: "🔥 Hugging Face", desc: "Hugging Face AI with FastAPI"},
		template{name: "🎯 Llama C++", desc: "Llama.cpp with Next.js"},
		template{name: "🌟 Vertex AI", desc: "Google Vertex AI with Flask"},
		template{name: "🔮 MERN Stack", desc: "Full-stack MongoDB, Express, React, and Node.js"},
		template{name: "⚛️ PERN Stack", desc: "Full-stack PostgreSQL, Express, React, and Node.js"},
	}

	// Set default template if one was detected
	defaultIndex := 0
	if detectedTemplate != "" {
		for i, t := range templates {
			if t.(template).name == detectedTemplate {
				defaultIndex = i
				break
			}
		}
	}

	l := list.New(templates, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🚀 Select a Template"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Select(defaultIndex)

	m := model{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if finalModel.(model).quitting {
		return "", fmt.Errorf("template selection cancelled")
	}

	return finalModel.(model).choice, nil
}
