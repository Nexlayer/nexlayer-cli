package wizard

import (
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/wizard/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/wizard/ui"
)

type wizardState int

const (
	wizardStateWelcome wizardState = iota
	wizardStateAppName
	wizardStateDatabase
	wizardStateBackend
	wizardStateFrontend
	wizardStateReview
)

type stackSelection struct {
	appName  string
	database string
	backend  string
	frontend string
}

// wizardModel represents the wizard's model
type wizardModel struct {
	state          wizardState
	stack          stackSelection
	err            error
	quitting       bool
	nameInput      textinput.Model
	dbSelected     int
	backSelected   int
	frontSelected  int
	aiEnabled      bool
	aiSuggestions  []string
	aiClient       *ai.Client
}

var (
	validNameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	databaseOptions = []string{
		"MongoDB",
		"PostgreSQL",
		"MySQL",
		"Redis",
		"Skip (no database)",
	}

	backendOptions = []string{
		"Express.js",
		"FastAPI",
		"Spring Boot",
		"Django",
		"Skip (no backend)",
	}

	frontendOptions = []string{
		"React",
		"Vue.js",
		"Angular",
		"Svelte",
		"Skip (no frontend)",
	}
)

// NewWizardCmd creates a new wizard command
func NewWizardCmd() *cobra.Command {
	var useAI bool
	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive wizard to create a new deployment",
		Long:  "Start an interactive wizard that helps you create a new deployment by selecting components and configuring settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(newWizardModel(useAI))
			model, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to run wizard: %w", err)
			}

			m := model.(wizardModel)
			if m.quitting {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&useAI, "ai", false, "Enable AI-powered suggestions")
	return cmd
}

func newWizardModel(useAI bool) wizardModel {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-app"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	m := wizardModel{
		state:     wizardStateWelcome,
		aiEnabled: useAI,
		nameInput: ti,
	}

	if useAI {
		client, err := ai.NewClient()
		if err != nil {
			m.err = fmt.Errorf("failed to initialize AI client: %w", err)
			return m
		}
		m.aiClient = client
	}

	return m
}

func (m wizardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyUp:
			switch m.state {
			case wizardStateDatabase:
				if m.dbSelected > 0 {
					m.dbSelected--
				}
			case wizardStateBackend:
				if m.backSelected > 0 {
					m.backSelected--
				}
			case wizardStateFrontend:
				if m.frontSelected > 0 {
					m.frontSelected--
				}
			}
		case tea.KeyDown:
			switch m.state {
			case wizardStateDatabase:
				if m.dbSelected < len(databaseOptions)-1 {
					m.dbSelected++
				}
			case wizardStateBackend:
				if m.backSelected < len(backendOptions)-1 {
					m.backSelected++
				}
			case wizardStateFrontend:
				if m.frontSelected < len(frontendOptions)-1 {
					m.frontSelected++
				}
			}
		case tea.KeyEnter:
			switch m.state {
			case wizardStateWelcome:
				m.state = wizardStateAppName
				m.nameInput.Focus()
			case wizardStateAppName:
				if validNameRegex.MatchString(m.nameInput.Value()) {
					m.stack.appName = m.nameInput.Value()
					m.state = wizardStateDatabase
				}
			case wizardStateDatabase:
				m.stack.database = databaseOptions[m.dbSelected]
				m.state = wizardStateBackend
			case wizardStateBackend:
				m.stack.backend = backendOptions[m.backSelected]
				m.state = wizardStateFrontend
			case wizardStateFrontend:
				m.stack.frontend = frontendOptions[m.frontSelected]
				m.state = wizardStateReview
			case wizardStateReview:
				return m, tea.Quit
			}
		}
	}

	// Handle text input
	if m.state == wizardStateAppName {
		var inputCmd tea.Cmd
		m.nameInput, inputCmd = m.nameInput.Update(msg)
		cmd = tea.Batch(cmd, inputCmd)
	}

	return m, cmd
}

func (m wizardModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	switch m.state {
	case wizardStateWelcome:
		return ui.RenderWelcome()
	case wizardStateAppName:
		return ui.RenderNamePrompt(
			m.nameInput.Value(),
			validNameRegex.MatchString(m.nameInput.Value()),
		)
	case wizardStateDatabase:
		return ui.RenderComponentSelection(
			"Step 2: Choose Your Database",
			databaseOptions,
			m.dbSelected,
			map[string]string{
				"database": m.stack.database,
				"backend":  m.stack.backend,
				"frontend": m.stack.frontend,
			},
		)
	case wizardStateBackend:
		return ui.RenderComponentSelection(
			"Step 3: Choose Your Backend",
			backendOptions,
			m.backSelected,
			map[string]string{
				"database": m.stack.database,
				"backend":  m.stack.backend,
				"frontend": m.stack.frontend,
			},
		)
	case wizardStateFrontend:
		return ui.RenderComponentSelection(
			"Step 4: Choose Your Frontend",
			frontendOptions,
			m.frontSelected,
			map[string]string{
				"database": m.stack.database,
				"backend":  m.stack.backend,
				"frontend": m.stack.frontend,
			},
		)
	case wizardStateReview:
		return ui.RenderComponentSelection(
			"Step 5: Review Your Stack",
			[]string{"Your stack is ready to deploy!"},
			0,
			map[string]string{
				"database": m.stack.database,
				"backend":  m.stack.backend,
				"frontend": m.stack.frontend,
			},
		)
	default:
		return "Unknown state"
	}
}
