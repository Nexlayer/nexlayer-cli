package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

const (
	// UI Constants - used for consistent sizing across components
	defaultWidth    = 50 // Default width for UI components
	defaultBoxWidth = 19 // Default width for boxes and containers

	// Colors - using standard ANSI color codes for consistent appearance
	colorSuccess = "#00FF00" // Green for success messages
	colorInfo    = "#87CEEB" // Light blue for informational messages
	colorError   = "196"     // Red for error messages
	colorHeading = "205"     // Pink for headings

	// Error messages
	errNilError     = "an error occurred but no details were provided"
	errEmptyTitle   = "title cannot be empty"
	errInvalidRange = "invalid range: current value cannot be greater than total"
)

var (
	// Styles - predefined styles for consistent UI appearance
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorSuccess)).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorInfo)).
			Italic(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSuccess))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSuccess)).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorInfo))

	headingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorHeading)).
			Bold(true).
			Margin(1, 0, 1, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorError)).
			Bold(true)
)

// RenderTitle renders a styled title box with consistent width.
// Returns an error message if title is empty.
func RenderTitle(title string) string {
	if title == "" {
		return RenderErrorMessage(fmt.Errorf(errEmptyTitle))
	}
	return titleStyle.Width(defaultWidth).Render(title)
}

// RenderProgressBar renders a progress bar with percentage.
// Returns an error message if current is greater than total.
func RenderProgressBar(current, total int) string {
	if current > total {
		return RenderErrorMessage(fmt.Errorf(errInvalidRange))
	}

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	percent := float64(current) / float64(total)
	return progressStyle.Render(prog.ViewAs(percent))
}

// RenderArchitectureDiagram renders an ASCII art diagram of the stack
func RenderArchitectureDiagram(database, backend, frontend string) string {
	var sb strings.Builder

	emptyBox := "+-------------------+"
	line := "         |         "
	arrow := "         v         "

	// Helper to center text in box
	centerText := func(text string) string {
		padding := (defaultBoxWidth - len(text)) / 2
		return fmt.Sprintf("|%s%s%s|",
			strings.Repeat(" ", padding),
			text,
			strings.Repeat(" ", defaultBoxWidth-len(text)-padding))
	}

	// Add components from top to bottom
	if database != "" {
		sb.WriteString(emptyBox + "\n")
		sb.WriteString(centerText(database) + "\n")
		sb.WriteString(emptyBox + "\n")
		if backend != "" || frontend != "" {
			sb.WriteString(line + "\n")
			sb.WriteString(arrow + "\n")
		}
	}

	if backend != "" {
		sb.WriteString(emptyBox + "\n")
		sb.WriteString(centerText(backend) + "\n")
		sb.WriteString(emptyBox + "\n")
		if frontend != "" {
			sb.WriteString(line + "\n")
			sb.WriteString(arrow + "\n")
		}
	}

	if frontend != "" {
		sb.WriteString(emptyBox + "\n")
		sb.WriteString(centerText(frontend) + "\n")
		sb.WriteString(emptyBox + "\n")
	}

	return sb.String()
}

// RenderYAMLPreview renders a preview of the YAML configuration
func RenderYAMLPreview(database, backend, frontend string) string {
	var sb strings.Builder
	sb.WriteString("📜 YAML Preview:\n")
	sb.WriteString("pods:\n")

	if database != "" {
		sb.WriteString(fmt.Sprintf("  - type: database\n    name: %s\n", database))
	}
	if backend != "" {
		sb.WriteString(fmt.Sprintf("  - type: backend\n    name: %s\n", backend))
	}
	if frontend != "" {
		sb.WriteString(fmt.Sprintf("  - type: frontend\n    name: %s\n", frontend))
	}

	return sb.String()
}

// RenderSuccessMessage renders a success message with optional next steps
func RenderSuccessMessage(text string, nextSteps ...string) string {
	var sb strings.Builder
	sb.WriteString(successStyle.Render(fmt.Sprintf("✅ %s\n", text)))

	if len(nextSteps) > 0 {
		sb.WriteString("\nNext steps:\n")
		for i, step := range nextSteps {
			sb.WriteString(fmt.Sprintf("%d️⃣ %s\n", i+1, step))
		}
	}

	return sb.String()
}

// RenderHeading renders a heading with consistent styling
func RenderHeading(text string) string {
	return headingStyle.Render(text)
}

// RenderErrorMessage renders an error message with consistent styling.
// If err is nil, returns a generic error message.
func RenderErrorMessage(err error) string {
	if err == nil {
		return errorStyle.Render(fmt.Sprintf("❌ Error: %s", errNilError))
	}
	return errorStyle.Render(fmt.Sprintf("❌ Error: %s", err.Error()))
}

// RenderArchitecturePreview renders a preview of the stack architecture.
// If components is empty, returns an informational message.
func RenderArchitecturePreview(components []string) string {
	if len(components) == 0 {
		return RenderInfoMessage("No components selected")
	}

	var sb strings.Builder
	sb.WriteString(RenderHeading("Stack Architecture"))

	for _, component := range components {
		sb.WriteString(infoStyle.Render(fmt.Sprintf("• %s\n", component)))
	}

	return sb.String()
}

// RenderInfoMessage renders an info message with consistent styling.
// If text is empty, returns an empty string.
func RenderInfoMessage(text string) string {
	if text == "" {
		return ""
	}
	return infoStyle.Render(text)
}
