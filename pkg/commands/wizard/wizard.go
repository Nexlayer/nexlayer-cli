package wizard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	githubAuthURL = "https://app.staging.nexlayer.io/auth/github"
)

// ComponentOption represents a stack component option
type ComponentOption struct {
	Name       string
	Desc       string
	PodType    string
	Tag        string
	ExposeHttp bool
}

// Component configurations
var (
	frontendOptions = []ComponentOption{
		{
			Name:       "react",
			Desc:       "React.js",
			PodType:    "nginx",
			Tag:        "katieharris/mern-react-todo:latest",
			ExposeHttp: true,
		},
		{
			Name:       "angular",
			Desc:       "Angular",
			PodType:    "nginx",
			Tag:        "katieharris/mean-angular-todo:latest",
			ExposeHttp: true,
		},
		{
			Name:       "vue",
			Desc:       "Vue.js",
			PodType:    "nginx",
			Tag:        "katieharris/vue-todo:latest",
			ExposeHttp: true,
		},
	}

	backendOptions = []ComponentOption{
		{
			Name:       "express",
			Desc:       "Express.js",
			PodType:    "node",
			Tag:        "katieharris/mern-express-todo:latest",
			ExposeHttp: true,
		},
		{
			Name:       "flask",
			Desc:       "Flask",
			PodType:    "python",
			Tag:        "katieharris/flask-todo:latest",
			ExposeHttp: true,
		},
	}

	databaseOptions = []ComponentOption{
		{
			Name:       "mongodb",
			Desc:       "MongoDB",
			PodType:    "mongodb",
			Tag:        "mongo:6.0",
			ExposeHttp: false,
		},
		{
			Name:       "mysql",
			Desc:       "MySQL",
			PodType:    "mysql",
			Tag:        "mysql:8.0",
			ExposeHttp: false,
		},
	}

	registryOptions = []ComponentOption{
		{
			Name: "ghcr",
			Desc: "GitHub Container Registry (recommended)",
			Tag:  "ghcr.io",
		},
		{
			Name: "dockerhub",
			Desc: "Docker Hub",
			Tag:  "docker.io",
		},
	}
)

// NewWizardCmd creates a new wizard command
func NewWizardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive deployment wizard",
		Long: `An enhanced interactive wizard to help you deploy your application.
Features a beautiful UI with real-time architecture visualization and YAML preview.`,
		RunE: runWizard,
	}

	return cmd
}

type wizardState int

const (
	wizardStateInit wizardState = iota
	wizardStateAppName
	wizardStateRegistry
	wizardStateStack
	wizardStateConfirm
)

type stackSelection struct {
	appName  string
	registry string
	frontend string
	backend  string
	database string
}

type wizardModel struct {
	currentState wizardState
	ready        bool
	stack        stackSelection
	err          error
	quitting     bool
	list         list.Model
}

func newWizardModel() wizardModel {
	m := wizardModel{
		currentState: wizardStateInit,
		ready:        true,
		stack:        stackSelection{},
	}
	return m
}

func (m wizardModel) Init() tea.Cmd {
	return nil
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m wizardModel) View() string {
	if m.quitting {
		return "Thanks for using Nexlayer!\n"
	}

	var s strings.Builder

	s.WriteString(ui.RenderHeading("🚀 Welcome to Nexlayer!\n"))
	s.WriteString("Let's deploy your application.\n\n")

	switch m.currentState {
	case wizardStateInit:
		s.WriteString("Press any key to start...\n")
	case wizardStateAppName:
		s.WriteString("What's your application name?\n")
	case wizardStateRegistry:
		s.WriteString("Choose your container registry:\n")
		s.WriteString(m.renderRegistryOptions())
	case wizardStateStack:
		s.WriteString("Choose your stack components:\n")
		s.WriteString(m.renderStackOptions())
	case wizardStateConfirm:
		s.WriteString("Review your configuration:\n")
		s.WriteString(m.generateStackYAML())
	}

	if m.err != nil {
		s.WriteString("\n" + ui.RenderErrorMessage(m.err))
	}

	return s.String()
}

func (m wizardModel) renderRegistryOptions() string {
	var s strings.Builder
	s.WriteString("Available registries:\n")
	for _, opt := range registryOptions {
		s.WriteString(fmt.Sprintf("  %s (%s)\n", opt.Name, opt.Desc))
	}
	return s.String()
}

func (m wizardModel) renderStackOptions() string {
	var s strings.Builder
	s.WriteString("Frontend:\n")
	for _, opt := range frontendOptions {
		s.WriteString(fmt.Sprintf("  %s (%s)\n", opt.Name, opt.Desc))
	}
	s.WriteString("\nBackend:\n")
	for _, opt := range backendOptions {
		s.WriteString(fmt.Sprintf("  %s (%s)\n", opt.Name, opt.Desc))
	}
	s.WriteString("\nDatabase:\n")
	for _, opt := range databaseOptions {
		s.WriteString(fmt.Sprintf("  %s (%s)\n", opt.Name, opt.Desc))
	}
	return s.String()
}

func (m wizardModel) generateStackYAML() string {
	if m.stack.appName == "" {
		return ""
	}

	var yaml strings.Builder
	yaml.WriteString(fmt.Sprintf("name: %s\n", m.stack.appName))
	yaml.WriteString("\n")

	// Add registry configuration
	if m.stack.registry != "" {
		if opt := findComponentOption(registryOptions, m.stack.registry); opt != nil {
			yaml.WriteString("registry:\n")
			yaml.WriteString(fmt.Sprintf("  type: %s\n", opt.Name))
			yaml.WriteString(fmt.Sprintf("  url: %s\n", opt.Tag))
			yaml.WriteString("\n")
		}
	}

	yaml.WriteString("pods:\n")

	if m.stack.database != "" {
		if opt := findComponentOption(databaseOptions, m.stack.database); opt != nil {
			yaml.WriteString(fmt.Sprintf("  - type: %s\n", opt.PodType))
			yaml.WriteString(fmt.Sprintf("    name: %s\n", m.stack.database))
			yaml.WriteString(fmt.Sprintf("    tag: %s\n", opt.Tag))
			yaml.WriteString("    vars:\n")
			yaml.WriteString("      - key: MONGO_INITDB_ROOT_USERNAME\n")
			yaml.WriteString("        value: mongo\n")
			yaml.WriteString("      - key: MONGO_INITDB_ROOT_PASSWORD\n")
			yaml.WriteString("        value: passw0rd\n")
		}
	}

	if m.stack.backend != "" {
		if opt := findComponentOption(backendOptions, m.stack.backend); opt != nil {
			yaml.WriteString(fmt.Sprintf("  - type: %s\n", opt.PodType))
			yaml.WriteString(fmt.Sprintf("    name: %s\n", m.stack.backend))
			yaml.WriteString(fmt.Sprintf("    tag: %s\n", opt.Tag))
			yaml.WriteString("    vars:\n")
			yaml.WriteString("      - key: DATABASE_URL\n")
			yaml.WriteString("        value: mongo://mongodb-service\n")
		}
	}

	if m.stack.frontend != "" {
		if opt := findComponentOption(frontendOptions, m.stack.frontend); opt != nil {
			yaml.WriteString(fmt.Sprintf("  - type: %s\n", opt.PodType))
			yaml.WriteString(fmt.Sprintf("    name: %s\n", m.stack.frontend))
			yaml.WriteString(fmt.Sprintf("    tag: %s\n", opt.Tag))
			yaml.WriteString("    vars:\n")
			yaml.WriteString("      - key: BACKEND_URL\n")
			yaml.WriteString("        value: backend-service\n")
		}
	}

	return yaml.String()
}

// findComponentOption finds a component option by name
func findComponentOption(options []ComponentOption, name string) *ComponentOption {
	for _, opt := range options {
		if opt.Name == name {
			return &opt
		}
	}
	return nil
}

func setupAccount() error {
	fmt.Println("Opening GitHub authentication in your browser...")
	fmt.Printf("Please visit: %s\n", githubAuthURL)

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", githubAuthURL).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", githubAuthURL).Start()
	case "darwin":
		err = exec.Command("open", githubAuthURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("Failed to open browser automatically. Please visit %s manually.\n", githubAuthURL)
	}

	fmt.Println("\nPress Enter once you've completed the GitHub authentication...")
	fmt.Scanln() // Wait for user to press enter

	return nil
}

// runWizard is the main entry point for the wizard
func runWizard(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(newWizardModel())
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running wizard: %w", err)
	}

	finalModel := model.(wizardModel)
	if finalModel.quitting {
		return fmt.Errorf("wizard cancelled")
	}

	return nil
}

// item implements list.Item interface
type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
