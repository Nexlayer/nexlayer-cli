// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package ui

// Manager defines the interface for UI operations.
type Manager interface {
	// Rendering methods.
	RenderTitle(title string, subtitle ...string) string
	RenderTitleWithBorder(title string) string
	RenderError(msg string) string
	RenderSuccess(msg string) string
	RenderWarning(msg string) string
	RenderInfo(msg string) string
	RenderTable(headers []string, rows [][]string) string

	// Progress tracking.
	StartProgress(msg string) ProgressTracker
	UpdateProgress(progress float64, msg string)
	CompleteProgress()
}

// ProgressTracker tracks progress for long-running operations.
type ProgressTracker interface {
	Update(progress float64, msg string)
	Complete()
}

type manager struct {
	currentProgress ProgressTracker
}

// NewManager creates and returns a new UI manager instance.
func NewManager() Manager {
	return &manager{}
}

func (m *manager) RenderTitle(title string, subtitle ...string) string {
	return RenderTitle(title, subtitle...)
}

func (m *manager) RenderTitleWithBorder(title string) string {
	return RenderTitleWithBorder(title)
}

func (m *manager) RenderError(msg string) string {
	return RenderError(msg)
}

func (m *manager) RenderSuccess(msg string) string {
	return RenderSuccess(msg)
}

func (m *manager) RenderWarning(msg string) string {
	return RenderWarning(msg)
}

func (m *manager) RenderInfo(msg string) string {
	return RenderInfo(msg)
}

func (m *manager) RenderTable(headers []string, rows [][]string) string {
	return RenderTable(headers, rows)
}

func (m *manager) StartProgress(msg string) ProgressTracker {
	m.currentProgress = newProgressBar(msg)
	return m.currentProgress
}

func (m *manager) UpdateProgress(progress float64, msg string) {
	if m.currentProgress != nil {
		m.currentProgress.Update(progress, msg)
	}
}

func (m *manager) CompleteProgress() {
	if m.currentProgress != nil {
		m.currentProgress.Complete()
		m.currentProgress = nil
	}
}
