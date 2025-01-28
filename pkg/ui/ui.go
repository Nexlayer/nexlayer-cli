package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ff00"))

	subtitleStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#888888"))

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff0000"))

	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ff00"))

	tableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#888888"))
)

// RenderTitle renders a title with optional subtitle
func RenderTitle(title string, subtitle ...string) string {
	result := titleStyle.Render(title)
	if len(subtitle) > 0 {
		result += "\n" + subtitleStyle.Render(subtitle[0])
	}
	return result
}

// RenderTitleWithBorder renders a title with a border
func RenderTitleWithBorder(title string) string {
	return titleStyle.Copy().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Render(title)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return errorStyle.Render(fmt.Sprintf("Error: %s", msg))
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return successStyle.Render(msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return color.YellowString(fmt.Sprintf("Warning: %s", msg))
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return color.BlueString(msg)
}

// RenderTable renders a table with headers and rows
func RenderTable(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build table
	var sb strings.Builder

	// Headers
	for i, h := range headers {
		fmt.Fprintf(&sb, "%-*s", widths[i]+2, h)
	}
	sb.WriteString("\n")

	// Separator
	for _, w := range widths {
		sb.WriteString(strings.Repeat("-", w+2))
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, cell := range row {
			fmt.Fprintf(&sb, "%-*s", widths[i]+2, cell)
		}
		sb.WriteString("\n")
	}

	return tableStyle.Render(sb.String())
}
