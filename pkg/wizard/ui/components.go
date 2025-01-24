package ui

import (
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

const headerWidth = 48

var (
	// Colors
	hotPink = lipgloss.Color("#FF1493")
	neonBlue = lipgloss.Color("#00FFFF")
	neonGreen = lipgloss.Color("#39FF14")
	neonPurple = lipgloss.Color("#9D00FF")
	white = lipgloss.Color("#FFFFFF")

	// Box characters
	boxChar  = "═"  // Double horizontal line
	lineChar = "═"  // Double horizontal line
	vertChar = "║"  // Double vertical line
	topLeft = "╔"   // Top left corner
	topRight = "╗"  // Top right corner
	bottomLeft = "╚" // Bottom left corner
	bottomRight = "╝" // Bottom right corner

	// Styles
	BoxStyle = lipgloss.NewStyle().
		BorderForeground(neonBlue)

	TitleStyle = lipgloss.NewStyle().
		Foreground(white).
		Background(hotPink).
		Bold(true).
		Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(neonGreen).
		Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
		Foreground(neonPurple)

	CallToActionStyle = lipgloss.NewStyle().
		Foreground(neonBlue).
		Bold(true)

	YamlKeyStyle = lipgloss.NewStyle().
		Foreground(hotPink).
		Bold(true)

	YamlValueStyle = lipgloss.NewStyle().
		Foreground(neonGreen)
)

// RenderHeader creates a boxed header
func RenderHeader(title string) string {
	border := strings.Repeat(lineChar, headerWidth-2)
	header := fmt.Sprintf("%s%s%s", topLeft, border, topRight)
	
	// Ensure title fits within the box
	maxTitleLen := headerWidth - 4 // -4 for the vertical bars and spaces
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}
	
	// Center the title
	totalPadding := headerWidth - 2 - len(title) // -2 for the vertical bars
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding
	
	titleLine := fmt.Sprintf("%s%s%s%s%s",
		vertChar,
		strings.Repeat(" ", leftPadding),
		TitleStyle.Render(title),
		strings.Repeat(" ", rightPadding),
		vertChar,
	)

	footer := fmt.Sprintf("%s%s%s", bottomLeft, strings.Repeat(lineChar, headerWidth-2), bottomRight)
	
	return BoxStyle.Render(fmt.Sprintf("%s\n%s\n%s", header, titleLine, footer))
}

// RenderWelcome returns the welcome screen content
func RenderWelcome() string {
	return fmt.Sprintf(`%s
%s

%s
%s
%s
%s

%s
`,
		RenderHeader(" Welcome to the Nexlayer Deployment Wizard!"),
		CallToActionStyle.Render("Let's launch your app in minutes!"),
		TitleStyle.Render("Here's how it works:"),
		CallToActionStyle.Render("1 Name your app and pick your stack"),
		CallToActionStyle.Render("2 Deploy instantly with pre-configured settings"),
		CallToActionStyle.Render("3 Come back anytime to customize and scale"),
		CallToActionStyle.Render("Press Enter to get started!"),
	)
}

// RenderNamePrompt returns the app name prompt screen
func RenderNamePrompt(name string, isValid bool) string {
	var status string
	if name != "" {
		if isValid {
			status = fmt.Sprintf("\n%s\n\n%s",
				SelectedStyle.Render(fmt.Sprintf(" Great! Your application will be deployed as: `%s`", name)),
				CallToActionStyle.Render("Press Enter to continue."),
			)
		} else {
			status = "\n" + UnselectedStyle.Render(" Invalid name! Please use only letters, numbers, dashes (-), or underscores (_).")
		}
	}

	return fmt.Sprintf(`%s
%s

> %s%s
`,
		RenderHeader("Step 1: Name Your Application"),
		CallToActionStyle.Render(" Let's give your app a name. What will it be?"),
		name,
		status,
	)
}

// RenderProgress returns a progress bar
func RenderProgress(percent int) string {
	width := 20
	filled := width * percent / 100
	bar := BoxStyle.Render(strings.Repeat("=", filled) + ">" + strings.Repeat(" ", width-filled-1))
	return fmt.Sprintf("%s [%s] %d%%", 
		TitleStyle.Render("Progress:"),
		bar,
		percent,
	)
}

// RenderComponentSelection returns the component selection screen
func RenderComponentSelection(title string, options []string, selected int, progress int, stack map[string]string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(RenderHeader(title) + "\n")
	sb.WriteString(CallToActionStyle.Render("Use ↑ and ↓ to navigate. Press Enter to select.") + "\n\n")

	// Options
	for i, opt := range options {
		if i == selected {
			sb.WriteString(SelectedStyle.Render("[✅] " + opt) + "\n")
		} else {
			sb.WriteString(UnselectedStyle.Render("[ ] " + opt) + "\n")
		}
	}
	sb.WriteString("\n")

	// Progress bar
	sb.WriteString(RenderProgress(progress) + "\n\n")

	// YAML Preview
	if stack["database"] != "" || stack["backend"] != "" || stack["frontend"] != "" {
		sb.WriteString("\n" + YamlKeyStyle.Render(" YAML Preview:") + "\n")
		sb.WriteString(YamlKeyStyle.Render("pods:") + "\n")
		
		if stack["database"] != "" && stack["database"] != "Skip (no database)" {
			sb.WriteString(renderPodYAML("database", stack["database"], "mongo:6.0", map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "mongo",
				"MONGO_INITDB_ROOT_PASSWORD": "passw0rd",
			}))
		}
		
		if stack["backend"] != "" && stack["backend"] != "Skip (no backend)" {
			sb.WriteString(renderPodYAML("backend", stack["backend"], "python:3.9", map[string]string{
				"DATABASE_URL": "mongodb://mongo:passw0rd@mongodb:27017",
			}))
		}
		
		if stack["frontend"] != "" && stack["frontend"] != "Skip (no frontend)" {
			sb.WriteString(renderPodYAML("frontend", stack["frontend"], "nginx:latest", map[string]string{
				"BACKEND_URL": "http://backend:8000",
			}))
		}
	}

	return sb.String()
}

// renderPodYAML renders a pod's YAML configuration
func renderPodYAML(podType, name, tag string, vars map[string]string) string {
	var sb strings.Builder
	
	sb.WriteString("  " + YamlKeyStyle.Render("- type: ") + YamlValueStyle.Render(podType) + "\n")
	sb.WriteString("    " + YamlKeyStyle.Render("name: ") + YamlValueStyle.Render(name) + "\n")
	sb.WriteString("    " + YamlKeyStyle.Render("tag: ") + YamlValueStyle.Render(tag) + "\n")
	
	if len(vars) > 0 {
		sb.WriteString("    " + YamlKeyStyle.Render("vars:") + "\n")
		for k, v := range vars {
			sb.WriteString("      " + YamlKeyStyle.Render("- key: ") + YamlValueStyle.Render(k) + "\n")
			sb.WriteString("        " + YamlKeyStyle.Render("value: ") + YamlValueStyle.Render(v) + "\n")
		}
	}
	
	return sb.String()
}
