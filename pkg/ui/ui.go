// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"strings"
)

// Table represents a formatted table for CLI output
type Table struct {
	header []string
	rows   [][]string
}

// NewTable creates a new CLI table
func NewTable() *Table {
	return &Table{}
}

// AddHeader adds column headers to the table
func (t *Table) AddHeader(headers ...string) {
	t.header = headers
}

// AddRow adds a row to the table
func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
}

// Render prints the table to stdout
func (t *Table) Render() error {
	// Calculate column widths
	widths := make([]int, len(t.header))
	for i, h := range t.header {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, col := range row {
			if len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	// Print header
	fmt.Println()
	for i, h := range t.header {
		fmt.Printf("%-*s  ", widths[i], h)
	}
	fmt.Println()

	// Print separator
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w) + "  ")
	}
	fmt.Println()

	// Print rows
	for _, row := range t.rows {
		for i, col := range row {
			fmt.Printf("%-*s  ", widths[i], col)
		}
		fmt.Println()
	}
	fmt.Println()
	return nil
}

// RenderTitleWithBorder renders a title with a border
func RenderTitleWithBorder(title string) {
	width := len(title) + 4
	border := strings.Repeat("=", width)
	fmt.Printf("\n%s\n  %s  \n%s\n\n", border, title, border)
}

// RenderBox renders text in a box
func RenderBox(text string) {
	lines := strings.Split(text, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	border := strings.Repeat("-", maxWidth+4)
	fmt.Printf("\n%s\n", border)
	for _, line := range lines {
		fmt.Printf("| %-*s |\n", maxWidth, line)
	}
	fmt.Printf("%s\n\n", border)
}

// RenderWelcome renders a welcome message
func RenderWelcome(text string) {
	RenderBox(text)
}

// RenderHighlight renders highlighted text
func RenderHighlight(text string) {
	fmt.Printf("\n✨ %s\n", text)
}

// RenderSuccess renders a success message
func RenderSuccess(text string) {
	fmt.Printf("\n✅ %s\n", text)
}

// RenderError renders an error message
func RenderError(text string) {
	fmt.Printf("\n❌ %s\n", text)
}

// Spinner represents a CLI progress spinner
type Spinner struct {
	message string
	active  bool
}

// NewSpinner creates a new CLI spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{message: message}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.active = true
	fmt.Printf("%s... ", s.message)
}

// Stop ends the spinner animation
func (s *Spinner) Stop() {
	if s.active {
		fmt.Println("done")
		s.active = false
	}
}
