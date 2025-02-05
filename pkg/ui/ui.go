package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	// titleStyle defines the style for primary titles.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ff00"))

	// subtitleStyle defines the style for subtitles.
	subtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#888888"))

	// errorStyle defines the style for error messages.
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff0000"))

	// successStyle defines the style for success messages.
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ff00"))

	// tableStyle defines the style for tables.
	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#888888"))
)

// RenderTitle renders a title with an optional subtitle.
func RenderTitle(title string, subtitle ...string) string {
	result := titleStyle.Render(title)
	if len(subtitle) > 0 {
		result += "\n" + subtitleStyle.Render(subtitle[0])
	}
	return result
}

// RenderTitleWithBorder renders a title enclosed in a decorative border.
func RenderTitleWithBorder(title string) string {
	return titleStyle.Copy().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Render(title)
}

// RenderError formats an error message in a bold red style.
func RenderError(msg string) string {
	return errorStyle.Render(fmt.Sprintf("Error: %s", msg))
}

// RenderSuccess formats a success message in a bold green style.
func RenderSuccess(msg string) string {
	return successStyle.Render(msg)
}

// RenderWarning formats a warning message in yellow.
func RenderWarning(msg string) string {
	return color.YellowString(fmt.Sprintf("Warning: %s", msg))
}

// RenderInfo formats an informational message in blue.
func RenderInfo(msg string) string {
	return color.BlueString(msg)
}

// RenderTable creates a textual table from headers and rows.
func RenderTable(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}

	// Calculate column widths.
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

	var sb strings.Builder

	// Render header row.
	for i, h := range headers {
		fmt.Fprintf(&sb, "%-*s", widths[i]+2, h)
	}
	sb.WriteString("\n")

	// Render separator.
	for _, w := range widths {
		sb.WriteString(strings.Repeat("-", w+2))
	}
	sb.WriteString("\n")

	// Render each row.
	for _, row := range rows {
		for i, cell := range row {
			fmt.Fprintf(&sb, "%-*s", widths[i]+2, cell)
		}
		sb.WriteString("\n")
	}

	return tableStyle.Render(sb.String())
}
