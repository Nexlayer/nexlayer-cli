package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	outputFile   string
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize service connections",
	Long: `Generate a visual diagram of service connections in your application.
Supports multiple output formats:
- ascii (default, prints to terminal)
- dot (Graphviz DOT file)
- png (PNG image)
- svg (SVG image)`,
	RunE: runVisualize,
}

func init() {
	visualizeCmd.Flags().StringVar(&appName, "app", "", "Application name")
	visualizeCmd.Flags().StringVar(&outputFormat, "format", "ascii", "Output format (ascii, dot, png, svg)")
	visualizeCmd.Flags().StringVar(&outputFile, "output", "", "Output file (required for non-ascii formats)")

	visualizeCmd.MarkFlagRequired("app")
}

type ServiceConnection struct {
	From        string
	To          string
	Description string
}

func runVisualize(cmd *cobra.Command, args []string) error {
	// Get auth token
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create API client
	client := api.NewClient("https://app.nexlayer.io")

	// Get service connections
	connections, err := client.GetServiceConnections(appName, token)
	if err != nil {
		return fmt.Errorf("failed to get service connections: %w", err)
	}

	switch outputFormat {
	case "ascii":
		return visualizeAscii(connections)
	case "dot", "png", "svg":
		if outputFile == "" {
			return fmt.Errorf("--output flag is required for %s format", outputFormat)
		}
		return visualizeGraphviz(connections, outputFormat, outputFile)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func visualizeAscii(connections []ServiceConnection) error {
	if len(connections) == 0 {
		fmt.Println("No service connections found")
		return nil
	}

	fmt.Println("Service Connections:")
	fmt.Println("-------------------")
	for _, conn := range connections {
		fmt.Printf("%s --> %s: %s\n", conn.From, conn.To, conn.Description)
	}
	return nil
}

func visualizeGraphviz(connections []ServiceConnection, format, outputFile string) error {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}
	defer graph.Close()

	// Create nodes and edges
	nodes := make(map[string]*graphviz.Node)
	for _, conn := range connections {
		// Create nodes if they don't exist
		if _, ok := nodes[conn.From]; !ok {
			node, err := graph.CreateNode(conn.From)
			if err != nil {
				return fmt.Errorf("failed to create node: %w", err)
			}
			nodes[conn.From] = node
		}
		if _, ok := nodes[conn.To]; !ok {
			node, err := graph.CreateNode(conn.To)
			if err != nil {
				return fmt.Errorf("failed to create node: %w", err)
			}
			nodes[conn.To] = node
		}

		// Create edge
		edge, err := graph.CreateEdge(conn.Description, nodes[conn.From], nodes[conn.To])
		if err != nil {
			return fmt.Errorf("failed to create edge: %w", err)
		}
		edge.SetLabel(conn.Description)
	}

	// Render graph
	if format == "dot" {
		return g.RenderDOT(graph, outputFile)
	}
	return g.RenderFilename(graph, graphviz.Format(format), outputFile)
}
