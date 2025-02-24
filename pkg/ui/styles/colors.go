// Package styles provides consistent styling for the CLI UI
package styles

import "github.com/charmbracelet/lipgloss"

// Colors used throughout the CLI
var (
	// Brand colors
	Primary   = lipgloss.Color("#00B4D8")
	Secondary = lipgloss.Color("#0077B6")
	Accent    = lipgloss.Color("#90E0EF")

	// Status colors
	Success = lipgloss.Color("#2ECC71")
	Error   = lipgloss.Color("#E74C3C")
	Warning = lipgloss.Color("#F1C40F")
	Info    = lipgloss.Color("#3498DB")

	// Text colors
	TextPrimary   = lipgloss.Color("#2C3E50")
	TextSecondary = lipgloss.Color("#7F8C8D")
	TextMuted     = lipgloss.Color("#BDC3C7")

	// Background colors
	BgPrimary   = lipgloss.Color("#FFFFFF")
	BgSecondary = lipgloss.Color("#F8F9FA")
	BgDark      = lipgloss.Color("#343A40")
)

// Common styles used throughout the CLI
var (
	// Text styles
	Title = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	Text = lipgloss.NewStyle().
		Foreground(TextPrimary)

	Muted = lipgloss.NewStyle().
		Foreground(TextMuted)

	// Status styles
	SuccessText = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	WarningText = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	ErrorText = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	InfoText = lipgloss.NewStyle().
			Foreground(Info).
			Bold(true)

	// Box styles
	Box = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(TextMuted).
		Padding(1)

	BoxPrimary = Box.Copy().
			BorderForeground(Primary)

	BoxSuccess = Box.Copy().
			BorderForeground(Success)

	BoxError = Box.Copy().
			BorderForeground(Error)

	BoxWarning = Box.Copy().
			BorderForeground(Warning)

	BoxInfo = Box.Copy().
		BorderForeground(Info)

	// Table styles
	TableHeader = Text.Copy().
			Padding(0, 1)

	TableCell = Text.Copy().
			Padding(0, 1)

	TableBorder = Text.Copy().
			Foreground(TextMuted)

	// Spinner styles
	SpinnerText = Text.Copy().
			Foreground(TextSecondary)

	// Status icons
	SuccessIcon = Text.Copy().
			Foreground(Success).
			Bold(true)

	ErrorIcon = Text.Copy().
			Foreground(Error).
			Bold(true)

	WarningIcon = Text.Copy().
			Foreground(Warning).
			Bold(true)

	InfoIcon = Text.Copy().
			Foreground(Info).
			Bold(true)
)
