package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRenderTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "Valid title",
			title:   "Hello World",
			wantErr: false,
		},
		{
			name:    "Empty title",
			title:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderTitle(tt.title)
			if tt.wantErr {
				assert.Contains(t, result, "❌")
				assert.Contains(t, result, "Error")
			} else {
				assert.Contains(t, result, tt.title)
			}
		})
	}
}

func TestRenderErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "Nil error",
			err:  nil,
			want: "",
		},
		{
			name: "Error message",
			err:  fmt.Errorf("test error"),
			want: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderErrorMessage(tt.err)
			assert.Contains(t, result, tt.want)
			if tt.err != nil {
				assert.Contains(t, result, "❌")
			}
		})
	}
}

func TestRenderArchitecturePreview(t *testing.T) {
	tests := []struct {
		name      string
		stack     []string
		wantLines int
	}{
		{
			name:      "Empty stack",
			stack:     []string{},
			wantLines: 1, // Should at least show a header or empty message
		},
		{
			name:      "Single component",
			stack:     []string{"frontend"},
			wantLines: 2,
		},
		{
			name:      "Full stack",
			stack:     []string{"frontend", "backend", "database"},
			wantLines: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preview := RenderArchitecturePreview(tt.stack)
			lines := strings.Count(preview, "\n") + 1
			if lines < tt.wantLines {
				t.Errorf("RenderArchitecturePreview() got %v lines, want at least %v", lines, tt.wantLines)
			}
		})
	}
}

func TestStyleConsistency(t *testing.T) {
	// Test that styles are consistent across different renders
	title := RenderHeading("Test")
	errorMsg := RenderErrorMessage(fmt.Errorf("test"))
	preview := RenderArchitecturePreview([]string{"test"})

	// All rendered content should use lipgloss styling
	if !strings.Contains(title, lipgloss.NewStyle().String()) ||
		!strings.Contains(errorMsg, lipgloss.NewStyle().String()) ||
		!strings.Contains(preview, lipgloss.NewStyle().String()) {
		t.Error("Not all components use lipgloss styling")
	}
}

func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name    string
		current int
		total   int
		wantErr bool
	}{
		{
			name:    "Valid progress",
			current: 50,
			total:   100,
			wantErr: false,
		},
		{
			name:    "Zero progress",
			current: 0,
			total:   100,
			wantErr: false,
		},
		{
			name:    "Complete progress",
			current: 100,
			total:   100,
			wantErr: false,
		},
		{
			name:    "Invalid progress",
			current: 150,
			total:   100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressBar(tt.current, tt.total)
			if tt.wantErr {
				assert.Contains(t, result, "❌")
			} else {
				assert.NotContains(t, result, "❌")
			}
		})
	}
}

func TestRenderInfoMessage(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Valid message",
			text: "Info message",
			want: "Info message",
		},
		{
			name: "Empty message",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderInfoMessage(tt.text)
			assert.Equal(t, tt.want == "", result == "")
			if tt.want != "" {
				assert.Contains(t, result, tt.want)
			}
		})
	}
}
