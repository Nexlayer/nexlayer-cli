package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	titleStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
)

// RenderTitleWithBorder renders a title with a border
func RenderTitleWithBorder(title string) string {
	width := len(title) + 4
	border := strings.Repeat("=", width)
	return fmt.Sprintf("\n%s\n  %s\n%s\n", border, title, border)
}

// RenderErrorMessage renders an error message in red
func RenderErrorMessage(err error) string {
	return color.RedString("Error: %v", err)
}

// RenderSuccessMessage renders a success message in green
func RenderSuccessMessage(msg string) string {
	return color.GreenString("Success: %s", msg)
}

// RenderProgressBar renders a progress bar
func RenderProgressBar(progress float64) string {
	width := 40
	filled := int(progress / 100 * float64(width))
	bar := strings.Repeat("=", filled) + strings.Repeat("-", width-filled)
	return fmt.Sprintf("[%s] %.0f%%", bar, progress)
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

	return sb.String()
}
