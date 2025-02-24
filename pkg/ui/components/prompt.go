package components

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
)

// Prompt represents a user input prompt
type Prompt struct {
	label     string
	validator func(string) error
}

// NewPrompt creates a new prompt with the given label
func NewPrompt(label string) *Prompt {
	return &Prompt{
		label: label,
	}
}

// WithValidator adds a validator function to the prompt
func (p *Prompt) WithValidator(validator func(string) error) *Prompt {
	p.validator = validator
	return p
}

// Run runs the prompt and returns the user's input
func (p *Prompt) Run() (string, error) {
	prompt := promptui.Prompt{
		Label:    p.label,
		Validate: p.validator,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return strings.TrimSpace(result), nil
}

// Select represents a selection prompt
type Select struct {
	label    string
	items    []string
	selected int
}

// NewSelect creates a new selection prompt with the given label and items
func NewSelect(label string, items []string) *Select {
	return &Select{
		label: label,
		items: items,
	}
}

// WithSelected sets the default selected item
func (s *Select) WithSelected(index int) *Select {
	s.selected = index
	return s
}

// Run runs the selection prompt and returns the selected item and index
func (s *Select) Run() (string, int, error) {
	prompt := promptui.Select{
		Label: s.label,
		Items: s.items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("%s {{ . | underline }}", lipgloss.NewStyle().Foreground(styles.Primary).Render(">")),
			Inactive: "  {{ . }}",
			Selected: fmt.Sprintf("%s {{ . }}", styles.SuccessIcon.Render("✓")),
		},
		CursorPos: s.selected,
	}

	index, result, err := prompt.Run()
	if err != nil {
		return "", -1, fmt.Errorf("selection failed: %w", err)
	}

	return result, index, nil
}

// Confirm represents a yes/no confirmation prompt
type Confirm struct {
	message      string
	defaultValue bool
}

// NewConfirm creates a new confirmation prompt with the given message
func NewConfirm(message string) *Confirm {
	return &Confirm{
		message:      message,
		defaultValue: false,
	}
}

// WithDefault sets the default value for the confirmation
func (c *Confirm) WithDefault(value bool) *Confirm {
	c.defaultValue = value
	return c
}

// Run runs the confirmation prompt and returns the user's choice
func (c *Confirm) Run() (bool, error) {
	options := []string{"Yes", "No"}
	defaultIndex := 1
	if c.defaultValue {
		defaultIndex = 0
	}

	prompt := promptui.Select{
		Label: c.message,
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("%s {{ . | underline }}", lipgloss.NewStyle().Foreground(styles.Primary).Render(">")),
			Inactive: "  {{ . }}",
			Selected: fmt.Sprintf("%s {{ . }}", styles.SuccessIcon.Render("✓")),
		},
		CursorPos: defaultIndex,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("confirmation failed: %w", err)
	}

	return index == 0, nil
}
