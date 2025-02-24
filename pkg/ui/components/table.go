// Package components provides reusable UI components for the CLI
package components

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/ui/styles"
)

// Table represents a formatted table for CLI output
type Table struct {
	header []string
	rows   [][]string
	style  TableStyle
}

// TableStyle represents the styling options for a table
type TableStyle struct {
	HeaderStyle   func(string) string
	CellStyle     func(string) string
	BorderStyle   func(string) string
	ShowBorders   bool
	ColumnPadding int
}

// DefaultTableStyle returns the default table styling
func DefaultTableStyle() TableStyle {
	return TableStyle{
		HeaderStyle: func(s string) string {
			return styles.TableHeader.Render(s)
		},
		CellStyle: func(s string) string {
			return styles.TableCell.Render(s)
		},
		BorderStyle: func(s string) string {
			return styles.TableBorder.Render(s)
		},
		ShowBorders:   true,
		ColumnPadding: 2,
	}
}

// NewTable creates a new CLI table
func NewTable() *Table {
	return &Table{
		style: DefaultTableStyle(),
	}
}

// WithStyle sets custom styling for the table
func (t *Table) WithStyle(style TableStyle) *Table {
	t.style = style
	return t
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
		fmt.Printf(t.style.HeaderStyle("%-*s"), widths[i]+t.style.ColumnPadding, h)
	}
	fmt.Println()

	// Print separator if borders are enabled
	if t.style.ShowBorders {
		for _, w := range widths {
			fmt.Print(t.style.BorderStyle(strings.Repeat("-", w+t.style.ColumnPadding)))
		}
		fmt.Println()
	}

	// Print rows
	for _, row := range t.rows {
		for i, col := range row {
			fmt.Printf(t.style.CellStyle("%-*s"), widths[i]+t.style.ColumnPadding, col)
		}
		fmt.Println()
	}
	fmt.Println()

	return nil
}

// RenderToString returns the table as a string
func (t *Table) RenderToString() string {
	var sb strings.Builder

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

	// Add header
	sb.WriteString("\n")
	for i, h := range t.header {
		fmt.Fprintf(&sb, t.style.HeaderStyle("%-*s"), widths[i]+t.style.ColumnPadding, h)
	}
	sb.WriteString("\n")

	// Add separator if borders are enabled
	if t.style.ShowBorders {
		for _, w := range widths {
			sb.WriteString(t.style.BorderStyle(strings.Repeat("-", w+t.style.ColumnPadding)))
		}
		sb.WriteString("\n")
	}

	// Add rows
	for _, row := range t.rows {
		for i, col := range row {
			fmt.Fprintf(&sb, t.style.CellStyle("%-*s"), widths[i]+t.style.ColumnPadding, col)
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	return sb.String()
}
