package components

import (
	"fmt"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/ui/styles"
	"github.com/briandowns/spinner"
)

// Spinner represents a loading spinner with message
type Spinner struct {
	spinner *spinner.Spinner
	message string
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + styles.SpinnerText.Render(message)
	return &Spinner{
		spinner: s,
		message: message,
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.spinner.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.message = message
	s.spinner.Suffix = " " + styles.SpinnerText.Render(message)
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", styles.SuccessIcon.Render("✓"), styles.SuccessText.Render(message))
}

// Error stops the spinner and shows an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", styles.ErrorIcon.Render("✗"), styles.ErrorText.Render(message))
}

// Info stops the spinner and shows an info message
func (s *Spinner) Info(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", styles.InfoIcon.Render("i"), styles.InfoText.Render(message))
}

// Warning stops the spinner and shows a warning message
func (s *Spinner) Warning(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", styles.WarningIcon.Render("!"), styles.WarningText.Render(message))
}
