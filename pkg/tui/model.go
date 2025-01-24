package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 1)
)

// DeploymentConfig holds the configuration for deployment
type DeploymentConfig struct {
	AppName        string
	DeploymentName string
	DatabaseType   string
	BackendType    string
	FrontendType   string
	GithubUsername string
	GithubToken    string
}

// Model represents the TUI state
type Model struct {
	config      DeploymentConfig
	steps       []Step
	currentStep int
	inputs      []textinput.Model
	quitting    bool
}

// Step represents a wizard step
type Step struct {
	Title       string
	Description string
	InputLabel  string
	Completed   bool
	Error       error
}

// NewModel creates a new TUI model
func NewModel() Model {
	steps := []Step{
		{
			Title:       "Application Name",
			Description: "Enter a name for your application",
			InputLabel:  "App name:",
		},
		{
			Title:       "Deployment Name",
			Description: "Enter a name for this deployment",
			InputLabel:  "Deployment name:",
		},
		{
			Title:       "Database Type",
			Description: "What type of database are you using?",
			InputLabel:  "Database type:",
		},
		{
			Title:       "Backend Type",
			Description: "What backend framework are you using?",
			InputLabel:  "Backend type:",
		},
		{
			Title:       "Frontend Type",
			Description: "What frontend framework are you using?",
			InputLabel:  "Frontend type:",
		},
		{
			Title:       "GitHub Username",
			Description: "Enter your GitHub username for container registry access",
			InputLabel:  "Username:",
		},
		{
			Title:       "GitHub Token",
			Description: "Enter your GitHub Personal Access Token (with read packages permission)",
			InputLabel:  "Token:",
		},
	}

	m := Model{
		steps:       steps,
		currentStep: 0,
		inputs:      make([]textinput.Model, len(steps)),
	}

	// Initialize text inputs
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = steps[i].InputLabel
		t.Focus()
		t.CharLimit = 156
		t.Width = 50
		m.inputs[i] = t
	}

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles model updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.currentStep >= len(m.steps)-1 {
				// Save configuration
				m.config = DeploymentConfig{
					AppName:        m.inputs[0].Value(),
					DeploymentName: m.inputs[1].Value(),
					DatabaseType:   m.inputs[2].Value(),
					BackendType:    m.inputs[3].Value(),
					FrontendType:   m.inputs[4].Value(),
					GithubUsername: m.inputs[5].Value(),
					GithubToken:    m.inputs[6].Value(),
				}
				return m, tea.Quit
			}
			m.steps[m.currentStep].Completed = true
			m.currentStep++
			return m, nil
		}
	}

	// Handle input updates
	if m.currentStep < len(m.inputs) {
		m.inputs[m.currentStep], cmd = m.inputs[m.currentStep].Update(msg)
	}

	return m, cmd
}

// View renders the model
func (m Model) View() string {
	if m.quitting {
		return "Thanks for using Nexlayer CLI!\n"
	}

	s := fmt.Sprintf("\n%s\n", titleStyle.Render(" Nexlayer CLI - Describe Your Application "))

	// Show current step
	step := m.steps[m.currentStep]
	s += fmt.Sprintf("%s\n", titleStyle.Render(step.Title))
	s += fmt.Sprintf("%s\n\n", infoStyle.Render(step.Description))

	// Show input for current step
	s += m.inputs[m.currentStep].View() + "\n\n"

	// Show progress
	s += fmt.Sprintf("Step %d of %d\n\n", m.currentStep+1, len(m.steps))

	// Show help
	s += infoStyle.Render("Press Enter to continue, q to quit")

	return s
}
