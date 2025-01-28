package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type progressBar struct {
	mu      sync.Mutex
	message string
	width   int
	started time.Time
}

func newProgressBar(msg string) ProgressTracker {
	return &progressBar{
		message: msg,
		width:   40,
		started: time.Now(),
	}
}

func (p *progressBar) Update(progress float64, msg string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Calculate the filled width
	filled := int(progress / 100 * float64(p.width))
	bar := strings.Repeat("=", filled) + strings.Repeat("-", p.width-filled)

	// Calculate elapsed time
	elapsed := time.Since(p.started).Round(time.Second)

	// Update message if provided
	if msg != "" {
		p.message = msg
	}

	// Clear the current line and print the progress bar
	fmt.Printf("\r\033[K[%s] %.1f%% %s (%s)", 
		bar,
		progress,
		p.message,
		elapsed,
	)
}

func (p *progressBar) Complete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Clear the current line
	fmt.Print("\r\033[K")
	
	// Print completion message
	color.Green("✓ %s (%.2fs)", p.message, time.Since(p.started).Seconds())
}
