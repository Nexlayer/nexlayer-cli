// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// NewCommand creates a new status command
func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [namespace] [application-id]",
		Short: "Get deployment status",
		Long: `Get detailed status information about your deployments.
If no namespace and application ID are provided, lists all deployments.
If namespace and application ID are provided, shows detailed information about that specific deployment.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 {
				// Get detailed info for specific deployment
				namespace := args[0]
				appID := args[1]
				return getDeploymentInfo(cmd.Context(), client, namespace, appID)
			}
			// List all deployments
			return listDeployments(cmd.Context(), client)
		},
	}

	return cmd
}

func getDeploymentInfo(ctx context.Context, client *api.Client, namespace, appID string) error {
	info, err := client.GetDeploymentInfo(ctx, namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	// Print deployment info
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("\n%s\n\n", bold("Deployment Information"))
	fmt.Printf("Namespace:     %s\n", info.Namespace)
	fmt.Printf("Template:      %s (%s)\n", info.TemplateName, info.TemplateID)
	fmt.Printf("Status:        %s\n", formatStatus(info.DeploymentStatus))
	fmt.Println()

	return nil
}

func listDeployments(ctx context.Context, client *api.Client) error {
	// Get deployments with pagination
	deployments, err := client.GetDeployments(ctx, "")
	if err != nil {
		return err
	}

	if len(deployments) == 0 {
		fmt.Println("No deployments found")
		return nil
	}

	// Define table columns
	columns := []table.Column{
		{Title: "Namespace", Width: 20},
		{Title: "Template ID", Width: 36},
		{Title: "Template Name", Width: 30},
		{Title: "Status", Width: 15},
	}

	// Create rows
	var rows []table.Row
	for _, d := range deployments {
		rows = append(rows, table.Row{
			d.Namespace,
			d.TemplateID,
			d.TemplateName,
			formatStatus(d.DeploymentStatus),
		})
	}

	// Initialize table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Create and run the Bubble Tea program
	m := model{t}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func formatStatus(status string) string {
	switch status {
	case "running":
		return color.GreenString("● Running")
	case "pending":
		return color.YellowString("○ Pending")
	case "failed":
		return color.RedString("✕ Failed")
	case "stopped":
		return color.BlueString("■ Stopped")
	default:
		return status
	}
}
