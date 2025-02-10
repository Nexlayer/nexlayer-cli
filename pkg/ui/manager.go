package ui

import (
	"fmt"
	"io"
	"sync"
)

// progressState tracks the state of a progress operation
type progressState struct {
	Message string
	Status  string
}

// Manager handles UI interactions and progress tracking
type Manager struct {
	out    io.Writer
	mu     sync.Mutex
	tracks map[string]*progressState
}

// NewManager creates a new UI manager
func NewManager(out io.Writer) *Manager {
	return &Manager{
		out:    out,
		tracks: make(map[string]*progressState),
	}
}

// StartProgress starts tracking progress for a given ID
func (m *Manager) StartProgress(id string, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tracks[id] = &progressState{
		Message: msg,
		Status:  "running",
	}
	fmt.Fprintf(m.out, "%s...\n", RenderHighlight(msg))
}

// CompleteProgress marks progress as complete for a given ID
func (m *Manager) CompleteProgress(id string, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if track, ok := m.tracks[id]; ok {
		track.Status = "complete"
		track.Message = msg
	}
	fmt.Fprintf(m.out, "%s\n", RenderSuccess(msg))
}

// FailProgress marks progress as failed for a given ID
func (m *Manager) FailProgress(id string, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if track, ok := m.tracks[id]; ok {
		track.Status = "failed"
		track.Message = msg
	}
	fmt.Fprintf(m.out, "%s\n", RenderError(msg))
}
