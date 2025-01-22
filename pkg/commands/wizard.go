package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	advanced bool
)

type item struct {
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

type frameworkSelector struct {
	list     list.Model
	choice   string
	quitting bool
}

func newFrameworkSelector() *frameworkSelector {
	frameworks := []list.Item{
		item{"MERN Stack (MongoDB, Express, React, Node.js)"},
		item{"MEAN Stack (MongoDB, Express, Angular, Node.js)"},
		item{"MEVN Stack (MongoDB, Express, Vue.js, Node.js)"},
		item{"MNFA Stack (MongoDB, NestJS, Flutter, AWS)"},
		item{"PDN Stack (PostgreSQL, Django, Next.js)"},
		item{"PERN Stack (PostgreSQL, Express, React, Node.js)"},
		item{"Custom Project"},
	}

	// Get terminal dimensions
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	listHeight := min(height-4, len(frameworks)+3) // Account for title and borders

	l := list.New(frameworks, list.NewDefaultDelegate(), width-4, listHeight)
	l.Title = "Choose your framework:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = l.Styles.Title.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	return &frameworkSelector{
		list: l,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *frameworkSelector) Init() tea.Cmd {
	return nil
}

func (m *frameworkSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.choice = m.list.SelectedItem().(item).Title()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *frameworkSelector) View() string {
	return m.list.View()
}

func init() {
	WizardCmd.Flags().BoolVarP(&advanced, "advanced", "a", false, "Show advanced options")
}

// WizardCmd represents the wizard command
var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive setup wizard",
	Long: `An interactive wizard to help you get started with Nexlayer.
Example: nexlayer wizard`,
	RunE: runWizard,
}

func runWizard(cmd *cobra.Command, args []string) error {
	// Check for authentication
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		fmt.Println("⚠️ You need to login first!")
		fmt.Println("🔑 Running 'nexlayer login'...")
		
		if err := runLogin(cmd, args); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
	}

	// Create a new spinner
	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = " "
	s.FinalMSG = "✨ Setup complete!\n"

	// ASCII Art welcome
	color.Blue(`
 _   _           _                       
| \ | |         | |                      
|  \| | _____  _| | __ _ _   _  ___ _ __ 
| . ' |/ _ \ \/ / |/ _' | | | |/ _ \ '__|
| |\  |  __/>  <| | (_| | |_| |  __/ |   
\_| \_/\___/_/\_\_|\__,_|\__, |\___|_|   
                          __/ |           
                         |___/            
`)

	color.Green("\n🚀 Welcome to Nexlayer CLI!\n")
	fmt.Println("Let's get you from zero to deployment in minutes.")

	// Framework selection
	selector := newFrameworkSelector()
	p := tea.NewProgram(selector)
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to select framework: %w", err)
	}

	framework := m.(*frameworkSelector).choice
	if framework == "" {
		return fmt.Errorf("no framework selected")
	}

	// Project name input
	projectName := "hello-world"
	fmt.Printf("\nEnter your project name [%s]: ", projectName)
	fmt.Scanln(&projectName)
	if projectName == "" {
		projectName = "hello-world"
	}

	// Create project directory
	s.Suffix = " Creating project structure..."
	s.Start()
	
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		s.Stop()
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate sample app based on framework
	s.Suffix = " Generating sample application..."
	if err := generateSampleApp(projectDir, framework); err != nil {
		s.Stop()
		return fmt.Errorf("failed to generate sample app: %w", err)
	}

	// Create nexlayer.yaml
	s.Suffix = " Creating configuration..."
	if err := createConfig(projectDir, framework); err != nil {
		s.Stop()
		return fmt.Errorf("failed to create config: %w", err)
	}

	s.Stop()

	// Success message and next steps
	color.Green("\n🎉 Project created successfully!")
	fmt.Printf("\nYour new project is ready in ./%s\n", projectName)
	fmt.Println("\nNext steps:")
	color.Blue("  1. cd %s", projectName)
	color.Blue("  2. nexlayer deploy")
	
	if advanced {
		fmt.Println("\nAdvanced options:")
		color.Blue("  • Scale your app:     nexlayer scale --replicas 3")
		color.Blue("  • View logs:          nexlayer logs")
		color.Blue("  • Custom domain:      nexlayer domain add")
		color.Blue("  • Environment vars:   nexlayer env set")
	} else {
		fmt.Println("\nTip: Run 'nexlayer wizard --advanced' to see more options")
	}

	return nil
}

func generateSampleApp(dir, framework string) error {
	// Create a spinner for visual feedback
	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = " "
	s.Suffix = " Fetching template..."
	s.Start()
	defer s.Stop()

	// Create a temporary directory for cloning
	tmpDir, err := os.MkdirTemp("", "nexlayer-template-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Map framework to template repository
	var repoURL string
	switch {
	case strings.Contains(strings.ToLower(framework), "mern"):
		repoURL = "https://github.com/Nexlayer/MERN-Todo-List.git"
	case strings.Contains(strings.ToLower(framework), "mean"):
		repoURL = "https://github.com/Nexlayer/MEAN-Todo-List.git"
	case strings.Contains(strings.ToLower(framework), "mevn"):
		repoURL = "https://github.com/Nexlayer/MEVN-Todo-App.git"
	case strings.Contains(strings.ToLower(framework), "mnfa"):
		repoURL = "https://github.com/Nexlayer/MNFA-User-Store-App.git"
	case strings.Contains(strings.ToLower(framework), "pdn"):
		repoURL = "https://github.com/Nexlayer/PDN-Todo-List.git"
	case strings.Contains(strings.ToLower(framework), "pern"):
		repoURL = "https://github.com/Nexlayer/PERN-Todo-List.git"
	default:
		// For custom projects, create a basic structure
		return createBasicStructure(dir)
	}

	// Clone the template repository
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	// Copy template files to project directory
	if err := copyDir(tmpDir, dir); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Update package.json name if it exists
	if err := updatePackageJSON(dir); err != nil {
		return fmt.Errorf("failed to update package.json: %w", err)
	}

	return nil
}

func createBasicStructure(dir string) error {
	// Create basic directories
	dirs := []string{
		"src",
		"docs",
		"tests",
		"config",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(dir, d), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	// Create README.md
	readme := filepath.Join(dir, "README.md")
	content := `# Custom Nexlayer Project

This is a custom project created with Nexlayer CLI.

## Getting Started

1. Add your application code to the src directory
2. Configure your application in the config directory
3. Add tests to the tests directory
4. Deploy with: nexlayer deploy

Need help? Run 'nexlayer wizard --advanced'`

	if err := os.WriteFile(readme, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	return nil
}

func updatePackageJSON(dir string) error {
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Update the name to match the directory name
	packageJSON["name"] = filepath.Base(dir)

	updatedData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(packageJSONPath, updatedData, 0644)
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Skip .git directory and related files
			if strings.Contains(srcPath, ".git") {
				continue
			}

			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func createConfig(dir, framework string) error {
	configFile := filepath.Join(dir, "nexlayer.yaml")
	content := fmt.Sprintf(`name: sample-app
version: 1.0.0
type: %s

environment:
  PORT: "8080"

resources:
  - name: web
    type: service
    properties:
      port: 8080
      replicas: 1
`, strings.ToLower(strings.Split(framework, " ")[0]))

	return os.WriteFile(configFile, []byte(content), 0644)
}
