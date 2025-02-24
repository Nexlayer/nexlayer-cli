package components

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Progress represents a progress bar with message
type Progress struct {
	total     int
	current   int
	width     int
	message   string
	completed bool
}

// NewProgress creates a new progress bar with the given total and width
func NewProgress(total, width int) *Progress {
	return &Progress{
		total:   total,
		width:   width,
		message: "",
	}
}

// Update updates the progress bar with the current value and message
func (p *Progress) Update(current int, message string) {
	p.current = current
	p.message = message
	p.render()
}

// Complete marks the progress bar as completed
func (p *Progress) Complete(message string) {
	p.current = p.total
	p.message = message
	p.completed = true
	p.render()
	fmt.Println()
}

// Error shows an error message and stops the progress bar
func (p *Progress) Error(message string) {
	p.completed = true
	fmt.Printf("%s %s\n", styles.ErrorIcon.Render("✗"), styles.ErrorText.Render(message))
}

func (p *Progress) render() {
	percentage := float64(p.current) / float64(p.total)
	filled := int(float64(p.width) * percentage)
	empty := p.width - filled

	// Create the progress bar
	bar := strings.Builder{}
	bar.WriteString("[")
	bar.WriteString(lipgloss.NewStyle().Foreground(styles.Primary).Render(strings.Repeat("=", filled)))
	bar.WriteString(strings.Repeat(" ", empty))
	bar.WriteString("]")

	// Format the percentage
	percentStr := fmt.Sprintf(" %3.0f%%", percentage*100)

	// Format the message
	messageStyle := styles.Text
	if p.completed {
		messageStyle = styles.SuccessText
	}

	// Print the progress bar
	fmt.Printf("\r%s %s %s",
		bar.String(),
		lipgloss.NewStyle().Foreground(styles.TextSecondary).Render(percentStr),
		messageStyle.Render(p.message))
}
