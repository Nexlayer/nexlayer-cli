package wizard

import (
	"fmt"
	"os"
	"regexp"
	"time"

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
	progress       int
	targetProgress int
	lastTick       time.Time
	aiEnabled      bool
	aiSuggestions  []string
	aiClient       *ai.Client
}

var (
	validNameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	databaseOptions = []string{
		"MongoDB",
		"MySQL",
		"PostgreSQL",
		"Skip (no database)",
	}

	backendOptions = []string{
		"Express.js",
		"Flask",
		"Django",
		"Skip (no backend)",
	}

	frontendOptions = []string{
		"React",
		"Angular",
		"Vue.js",
		"Skip (no frontend)",
	}
)

// tickMsg represents a progress bar tick
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NewWizardCmd creates a new wizard command
func NewWizardCmd() *cobra.Command {
	var useAI bool

	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive wizard to set up your application",
		Long: `Interactive wizard to help you set up your application.
Use --ai flag to enable AI-powered recommendations (requires OPENAI_API_KEY).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if useAI {
				openAIKey := os.Getenv("OPENAI_API_KEY")
				if openAIKey == "" {
					return fmt.Errorf("OPENAI_API_KEY environment variable is required when using --ai flag")
				}
			}
			
			p := tea.NewProgram(newWizardModel(useAI))
			_, err := p.Run()
			return err
		},
	}

	cmd.Flags().BoolVar(&useAI, "ai", false, "Enable AI-powered recommendations (requires OPENAI_API_KEY)")
	return cmd
}

func newWizardModel(useAI bool) wizardModel {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-app"
	ti.Focus()

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
	return tea.Batch(textinput.Blink, tickCmd())
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up":
			switch m.state {
			case wizardStateDatabase:
				m.dbSelected = (m.dbSelected - 1 + len(databaseOptions)) % len(databaseOptions)
			case wizardStateBackend:
				m.backSelected = (m.backSelected - 1 + len(backendOptions)) % len(backendOptions)
			case wizardStateFrontend:
				m.frontSelected = (m.frontSelected - 1 + len(frontendOptions)) % len(frontendOptions)
			}
		case "down":
			switch m.state {
			case wizardStateDatabase:
				m.dbSelected = (m.dbSelected + 1) % len(databaseOptions)
			case wizardStateBackend:
				m.backSelected = (m.backSelected + 1) % len(backendOptions)
			case wizardStateFrontend:
				m.frontSelected = (m.frontSelected + 1) % len(frontendOptions)
			}
		case "enter":
			switch m.state {
			case wizardStateWelcome:
				m.state = wizardStateAppName
			case wizardStateAppName:
				if validNameRegex.MatchString(m.nameInput.Value()) {
					m.stack.appName = m.nameInput.Value()
					m.state = wizardStateDatabase
					m.targetProgress = 25
				}
			case wizardStateDatabase:
				m.stack.database = databaseOptions[m.dbSelected]
				m.state = wizardStateBackend
				m.targetProgress = 50
			case wizardStateBackend:
				m.stack.backend = backendOptions[m.backSelected]
				m.state = wizardStateFrontend
				m.targetProgress = 75
			case wizardStateFrontend:
				m.stack.frontend = frontendOptions[m.frontSelected]
				m.state = wizardStateReview
				m.targetProgress = 100
			case wizardStateReview:
				return m, tea.Quit
			}
		}
	case tickMsg:
		if m.progress < m.targetProgress {
			m.progress++
			cmds = append(cmds, tickCmd())
		} else if m.progress > m.targetProgress {
			m.progress--
			cmds = append(cmds, tickCmd())
		}
	}

	switch m.state {
	case wizardStateAppName:
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Always keep ticking for smooth animation
	if len(cmds) == 0 {
		cmds = append(cmds, tickCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m wizardModel) View() string {
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
			m.progress,
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
			m.progress,
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
			m.progress,
			map[string]string{
				"database": m.stack.database,
				"backend":  m.stack.backend,
				"frontend": m.stack.frontend,
			},
		)
	case wizardStateReview:
		stack := map[string]string{
			"database": m.stack.database,
			"backend":  m.stack.backend,
			"frontend": m.stack.frontend,
		}
		view := ui.RenderHeader("Review Your Stack") + "\n" +
			"\n"
		
		// Display the selected components
		if stack["database"] != "" && stack["database"] != "Skip (no database)" {
			view += ui.SelectedStyle.Render("✅ Database: ") + stack["database"] + "\n"
		}
		if stack["backend"] != "" && stack["backend"] != "Skip (no backend)" {
			view += ui.SelectedStyle.Render("✅ Backend: ") + stack["backend"] + "\n"
		}
		if stack["frontend"] != "" && stack["frontend"] != "Skip (no frontend)" {
			view += ui.SelectedStyle.Render("✅ Frontend: ") + stack["frontend"] + "\n"
		}
		view += "\n" +
			" Ready to deploy? Press Enter to confirm."
		return view
	default:
		return ""
	}
}
