package tui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner represents a loading spinner
type Spinner struct {
	spinner *spinner.Spinner
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = " "
	return &Spinner{spinner: s}
}

// Start starts the spinner with a message
func (s *Spinner) Start(message string) {
	s.spinner.Suffix = fmt.Sprintf(" %s", message)
	s.spinner.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// Update updates the spinner message
func (s *Spinner) Update(message string) {
	s.spinner.Suffix = fmt.Sprintf(" %s", message)
}
