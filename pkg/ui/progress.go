package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// progressBar implements the ProgressTracker interface.
type progressBar struct {
	mu      sync.Mutex
	message string
	width   int
	started time.Time
}

// newProgressBar creates a new progress bar instance with a given message.
func newProgressBar(msg string) ProgressTracker {
	return &progressBar{
		message: msg,
		width:   40,
		started: time.Now(),
	}
}

// Update refreshes the progress bar display with the given percentage and message.
func (p *progressBar) Update(progress float64, msg string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Calculate filled width based on percentage.
	filled := int((progress / 100.0) * float64(p.width))
	bar := strings.Repeat("=", filled) + strings.Repeat("-", p.width-filled)

	// Calculate elapsed time.
	elapsed := time.Since(p.started).Round(time.Second)

	// Update message if provided.
	if msg != "" {
		p.message = msg
	}

	// Clear the current line and print the progress bar.
	fmt.Printf("\r\033[K[%s] %.1f%% %s (%s)", bar, progress, p.message, elapsed)
}

// Complete finalizes the progress bar, clears the line, and prints a completion message.
func (p *progressBar) Complete() {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Print("\r\033[K")
	color.Green("✓ %s (%.2fs)", p.message, time.Since(p.started).Seconds())
}
