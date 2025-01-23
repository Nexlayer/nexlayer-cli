package wizard

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// WizardCmd represents the wizard command
var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive setup wizard",
	Long: `Interactive setup wizard for Nexlayer.
This will guide you through:
1. Setting up your Nexlayer account
2. Creating your first application
3. Deploying your first service`,
	RunE: runWizard,
}

type step struct {
	title       string
	description string
	action      func() error
}

func runWizard(cmd *cobra.Command, args []string) error {
	// Create spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Running setup wizard..."
	s.Start()
	defer s.Stop()

	steps := []step{
		{
			title:       "Account Setup",
			description: "Set up your Nexlayer account",
			action: func() error {
				s.Stop()
				fmt.Println("🔑 Setting up your Nexlayer account...")
				if err := setupAccount(); err != nil {
					return err
				}
				s.Start()
				return nil
			},
		},
		{
			title:       "Application Creation",
			description: "Create your first application",
			action: func() error {
				s.Stop()
				fmt.Println("📦 Creating your first application...")
				if err := createApplication(); err != nil {
					return err
				}
				s.Start()
				return nil
			},
		},
		{
			title:       "Service Deployment",
			description: "Deploy your first service",
			action: func() error {
				s.Stop()
				fmt.Println("🚀 Deploying your first service...")
				if err := deployService(); err != nil {
					return err
				}
				s.Start()
				return nil
			},
		},
	}

	// Run each step
	for _, step := range steps {
		if err := step.action(); err != nil {
			return err
		}
	}

	s.Stop()
	fmt.Println("✨ Setup wizard completed successfully!")
	return nil
}

func setupAccount() error {
	// Create list model
	items := []list.Item{
		item{title: "GitHub", desc: "Sign in with GitHub"},
		item{title: "GitLab", desc: "Sign in with GitLab"},
		item{title: "Email", desc: "Sign in with Email"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Sign-in Method"

	p := tea.NewProgram(model{list: l})
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	selected := m.(model).list.SelectedItem().(item).title
	switch selected {
	case "GitHub", "GitLab":
		return signInWithOAuth(selected)
	case "Email":
		return signInWithEmail()
	default:
		return fmt.Errorf("invalid sign-in method: %s", selected)
	}
}

func createApplication() error {
	var appName string
	fmt.Print("Enter application name: ")
	fmt.Scanln(&appName)

	if appName == "" {
		return fmt.Errorf("application name cannot be empty")
	}

	// Create app directory
	if err := os.MkdirAll(appName, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = appName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

func deployService() error {
	// Get list of services
	services := []string{"frontend", "backend", "database"}

	// Create list model
	var items []list.Item
	for _, service := range services {
		items = append(items, item{title: service, desc: fmt.Sprintf("Deploy %s service", service)})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Service to Deploy"

	p := tea.NewProgram(model{list: l})
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	selected := m.(model).list.SelectedItem().(item).title
	fmt.Printf("Selected service: %s\n", selected)

	// Deploy selected service
	return nil
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
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
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func signInWithOAuth(provider string) error {
	// TO DO: implement OAuth sign-in
	return fmt.Errorf("OAuth sign-in with %s provider not implemented yet", provider)
}

func signInWithEmail() error {
	// TO DO: implement email sign-in
	return nil
}
