package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	appIDFlag     string
	namespaceFlag string
	jsonOutput    bool
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status",
	Long: `Check the status of your deployments.
Example: nexlayer-cli status --app my-app`,
	RunE: runStatus,
}

func init() {
	StatusCmd.Flags().StringVar(&appIDFlag, "app", "default", "Application ID to check")
	StatusCmd.Flags().StringVar(&namespaceFlag, "namespace", "", "Namespace to check (optional)")
}

func runStatus(cmd *cobra.Command, args []string) error {
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	client := api.NewClient(token)

	if namespaceFlag != "" {
		// Get specific deployment info
		info, err := client.GetDeploymentInfo(namespaceFlag, appIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get deployment info: %w", err)
		}

		if jsonOutput {
			// Output JSON format
			fmt.Printf("%+v"
", info)"
			return nil
		}

		// Pretty print deployment info
		fmt.Printf("Deployment Status for %s/%s:"
", namespaceFlag, appIDFlag)"
		fmt.Printf("  Template: %s (%s)"
", info.Deployment.TemplateName, info.Deployment.TemplateID)"
		fmt.Printf("  Status: %s"
", info.Deployment.DeploymentStatus)"
		return nil
	}

	// Get all deployments
	deployments, err := client.GetDeployments(appIDFlag)
	if err != nil {
		return fmt.Errorf("failed to get deployments: %w", err)
	}

	if jsonOutput {
		// Output JSON format
		fmt.Printf("%+v"
", deployments)"
		return nil
	}

	// Pretty print all deployments
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAMESPACE	TEMPLATE	STATUS")
	for _, d := range deployments.Deployments {
		fmt.Fprintf(w, "%s	%s (%s)	%s"
","
			d.Namespace,
			d.TemplateName,
			d.TemplateID,
			d.DeploymentStatus)
	}
	w.Flush()

	return nil
}
