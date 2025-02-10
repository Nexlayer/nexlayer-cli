// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package ui provides utilities for rendering styled text and UI elements in the CLI.
// It uses lipgloss and color packages to create a consistent and visually appealing
// command-line interface with support for colors, borders, and formatted tables.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

func init() {
	// Force color output
	color.NoColor = false
}

var (
	// Colors
	primaryColor = lipgloss.AdaptiveColor{Light: "#FF69B4", Dark: "#FF69B4"} // Hot pink for selected/active items
	backgroundColor = lipgloss.AdaptiveColor{Light: "#1B1B1B", Dark: "#1B1B1B"} // Dark background
	textColor = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"} // White text
	accentColor = lipgloss.AdaptiveColor{Light: "#36B2B2", Dark: "#36B2B2"} // Teal for progress/success
	errorColor = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"} // Red for errors

	// Common styles
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	// Text styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)

	subtitleStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(accentColor)

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(errorColor)

	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)

	highlightStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)

	// Container styles
	tableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(accentColor)

	boxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1).
		Width(60)
)

// RenderTitle renders a title with an optional subtitle.
func RenderTitle(title string) string {
	return titleStyle.Render(title)
}

// RenderHighlight renders text in the highlight style
func RenderHighlight(text string) string {
	return highlightStyle.Render(text)
}

// RenderBox renders text in a box
func RenderBox(text string) string {
	return boxStyle.Render(text)
}

// RenderTitleWithBorder renders a title enclosed in a decorative border using titleStyle.
// Returns the formatted title string with a border.
func RenderTitleWithBorder(title string) string {
	return titleStyle.Copy().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Render(title)
}

// RenderError formats an error message in a bold red style using errorStyle.
// Returns the formatted error message string.
func RenderError(msg string) string {
	return errorStyle.Render(fmt.Sprintf("Error: %s", msg))
}

// RenderSuccess formats a success message in a bold green style using successStyle.
// Returns the formatted success message string.
func RenderSuccess(msg string) string {
	return successStyle.Render(msg)
}

// RenderWarning formats a warning message in yellow using color.YellowString.
// Returns the formatted warning message string.
func RenderWarning(msg string) string {
	return color.YellowString(fmt.Sprintf("Warning: %s", msg))
}

// RenderInfo formats an informational message in blue using color.BlueString.
// Returns the formatted informational message string.
func RenderInfo(msg string) string {
	return color.BlueString(msg)
}

// RenderTable creates a textual table from headers and rows.
func RenderTable(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Build the table
	var sb strings.Builder

	// Header
	for i, h := range headers {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], h))
	}
	sb.WriteString("\n")

	// Separator
	for i, width := range colWidths {
		if i > 0 {
			sb.WriteString("-+-")
		}
		sb.WriteString(strings.Repeat("-", width))
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				sb.WriteString(" | ")
			}
			if i < len(colWidths) {
				sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
			}
		}
		sb.WriteString("\n")
	}

	return tableStyle.Render(sb.String())
}
