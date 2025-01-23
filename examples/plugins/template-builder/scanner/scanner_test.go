package scanner

import (
	"testing"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityScanner(t *testing.T) {
	scanner := NewSecurityScanner()

	t.Run("Scan Node.js API Template", func(t *testing.T) {
		template := &types.NexlayerTemplate{
			Name: "node-api",
			Services: []types.Service{
				{
					Name: "api",
					Ports: []types.PortConfig{
						{Port: 3000, Protocol: "http"},
					},
					Environment: map[string]string{
						"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
						"API_KEY":      "secret-key",
					},
				},
			},
		}

		results, err := scanner.ScanTemplate(template)
		require.NoError(t, err)
		assert.NotEmpty(t, results)

		// Check for exposed secrets
		var foundSecretIssue bool
		for _, issue := range results {
			if issue.Type == "exposed_secret" {
				foundSecretIssue = true
				break
			}
		}
		assert.True(t, foundSecretIssue, "Should detect exposed secrets in environment variables")
	})

	t.Run("Scan Python Web Template", func(t *testing.T) {
		template := &types.NexlayerTemplate{
			Name: "python-web",
			Services: []types.Service{
				{
					Name: "web",
					Ports: []types.PortConfig{
						{Port: 8080, Protocol: "http"},
						{Port: 8443, Protocol: "https"},
					},
				},
			},
		}

		results, err := scanner.ScanTemplate(template)
		require.NoError(t, err)
		assert.NotEmpty(t, results)

		// Check for insecure ports
		var foundInsecurePortIssue bool
		for _, issue := range results {
			if issue.Type == "insecure_port" {
				foundInsecurePortIssue = true
				break
			}
		}
		assert.True(t, foundInsecurePortIssue, "Should detect insecure HTTP ports")
	})

	t.Run("Scan Empty Template", func(t *testing.T) {
		template := &types.NexlayerTemplate{
			Name: "empty",
		}

		results, err := scanner.ScanTemplate(template)
		require.NoError(t, err)
		assert.Empty(t, results, "Empty template should have no security issues")
	})

	t.Run("Scan Nil Template", func(t *testing.T) {
		_, err := scanner.ScanTemplate(nil)
		assert.Error(t, err, "Should return error for nil template")
	})
}

func TestCostScanner(t *testing.T) {
	scanner := NewCostScanner()

	t.Run("Resource Cost Estimation", func(t *testing.T) {
		template := &types.NexlayerTemplate{
			Name: "web-app",
			Services: []types.Service{
				{
					Name: "web",
					Resources: types.ResourceRequests{
						CPU:    "1",
						Memory: "2Gi",
					},
				},
				{
					Name: "db",
					Resources: types.ResourceRequests{
						CPU:    "2",
						Memory: "4Gi",
					},
				},
			},
			Resources: map[string]types.Resource{
				"web": {
					Type: "compute",
					Storage: []types.Storage{
						{Size: "20Gi", Type: "ssd"},
					},
				},
				"db": {
					Type: "compute",
					Storage: []types.Storage{
						{Size: "100Gi", Type: "ssd"},
					},
				},
			},
		}

		estimate, err := scanner.EstimateCosts(template)
		require.NoError(t, err)
		assert.NotZero(t, estimate.TotalCost)
		assert.NotEmpty(t, estimate.ResourceCosts)

		// Verify cost breakdown
		var foundStorageCost bool
		for _, cost := range estimate.ResourceCosts {
			if cost.Type == "storage" {
				foundStorageCost = true
				assert.NotZero(t, cost.MonthlyCost)
				break
			}
		}
		assert.True(t, foundStorageCost, "Should include storage costs in estimation")
	})

	t.Run("Network Cost Estimation", func(t *testing.T) {
		template := &types.NexlayerTemplate{
			Name: "cdn-app",
			Services: []types.Service{
				{
					Name: "cdn",
					Resources: types.ResourceRequests{
						CPU:    "1",
						Memory: "2Gi",
					},
				},
			},
			Resources: map[string]types.Resource{
				"cdn": {
					Type: "compute",
					Network: types.Network{
						Ingress:  "100Gi",
						Egress:   "1Ti",
						Requests: "1M",
					},
				},
			},
		}

		estimate, err := scanner.EstimateCosts(template)
		require.NoError(t, err)
		assert.NotZero(t, estimate.TotalCost)

		// Verify network cost breakdown
		var foundNetworkCost bool
		for _, cost := range estimate.ResourceCosts {
			if cost.Type == "network" {
				foundNetworkCost = true
				assert.NotZero(t, cost.MonthlyCost)
				break
			}
		}
		assert.True(t, foundNetworkCost, "Should include network costs in estimation")
	})
}

func TestScannerIntegration(t *testing.T) {
	securityScanner := NewSecurityScanner()
	costScanner := NewCostScanner()

	template := &types.NexlayerTemplate{
		Name: "full-stack-app",
		Services: []types.Service{
			{
				Name: "api",
				Ports: []types.PortConfig{
					{Port: 3000, Protocol: "http"},
				},
				Environment: map[string]string{
					"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
				},
				Resources: types.ResourceRequests{
					CPU:    "1",
					Memory: "2Gi",
				},
			},
		},
		Resources: map[string]types.Resource{
			"api": {
				Type: "compute",
				Storage: []types.Storage{
					{Size: "20Gi", Type: "ssd"},
				},
			},
		},
	}

	t.Run("Full Scan", func(t *testing.T) {
		// Security scan
		securityResults, err := securityScanner.ScanTemplate(template)
		require.NoError(t, err)
		assert.NotEmpty(t, securityResults, "Should find security issues due to HTTP port and exposed database URL")

		// Verify specific security issues
		var foundInsecurePort, foundExposedSecret bool
		for _, issue := range securityResults {
			switch issue.Type {
			case "insecure_port":
				foundInsecurePort = true
			case "exposed_secret":
				foundExposedSecret = true
			}
		}
		assert.True(t, foundInsecurePort, "Should detect insecure HTTP port")
		assert.True(t, foundExposedSecret, "Should detect exposed database URL")

		// Cost estimation
		costEstimate, err := costScanner.EstimateCosts(template)
		require.NoError(t, err)
		assert.NotZero(t, costEstimate.TotalCost)
	})
}
