package scanner

import (
	"fmt"
	"strings"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
)

// SecurityScanner scans templates for security issues
type SecurityScanner struct{}

// NewSecurityScanner creates a new SecurityScanner
func NewSecurityScanner() *SecurityScanner {
	return &SecurityScanner{}
}

// ScanTemplate scans a template for security issues
func (s *SecurityScanner) ScanTemplate(template *types.NexlayerTemplate) ([]types.SecurityIssue, error) {
	if template == nil {
		return nil, fmt.Errorf("template cannot be nil")
	}

	var issues []types.SecurityIssue

	// Scan services for security issues
	for _, service := range template.Services {
		// Check for insecure ports
		for _, port := range service.Ports {
			if strings.ToLower(port.Protocol) == "http" {
				issues = append(issues, types.SecurityIssue{
					Type:        "insecure_port",
					Severity:    "high",
					Description: fmt.Sprintf("Service %s is using insecure HTTP port %d", service.Name, port.Port),
					Context: map[string]string{
						"service": service.Name,
						"port":    fmt.Sprintf("%d", port.Port),
					},
				})
			}
		}

		// Check for sensitive environment variables
		for key, value := range service.Environment {
			if s.isSensitiveEnvVar(key) {
				issues = append(issues, types.SecurityIssue{
					Type:        "exposed_secret",
					Severity:    "critical",
					Description: fmt.Sprintf("Service %s has exposed sensitive environment variable: %s", service.Name, key),
					Context: map[string]string{
						"service": service.Name,
						"key":     key,
					},
				})
			} else if strings.Contains(value, "://") && strings.Contains(value, ":") {
				// Check for exposed credentials in URLs
				issues = append(issues, types.SecurityIssue{
					Type:        "exposed_secret",
					Severity:    "critical",
					Description: fmt.Sprintf("Service %s has exposed credentials in URL: %s", service.Name, key),
					Context: map[string]string{
						"service": service.Name,
						"key":     key,
					},
				})
			}
		}
	}

	return issues, nil
}

// isSensitiveEnvVar checks if an environment variable key is sensitive
func (s *SecurityScanner) isSensitiveEnvVar(key string) bool {
	sensitiveKeys := []string{
		"password",
		"secret",
		"token",
		"key",
		"cert",
		"api_key",
		"auth",
		"credential",
		"database_url",
		"db_url",
		"mongo_url",
		"redis_url",
		"postgres_url",
		"mysql_url",
	}

	key = strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(key, sensitive) {
			return true
		}
	}
	return false
}

// CostScanner estimates resource costs for templates
type CostScanner struct{}

// NewCostScanner creates a new CostScanner
func NewCostScanner() *CostScanner {
	return &CostScanner{}
}

// EstimateCosts estimates the costs for a template
func (c *CostScanner) EstimateCosts(template *types.NexlayerTemplate) (*types.CostEstimate, error) {
	if template == nil {
		return nil, fmt.Errorf("template cannot be nil")
	}

	var totalCost float64
	var resourceCosts []types.ResourceCost

	// Calculate compute costs
	for _, service := range template.Services {
		cpu := service.Resources.CPU
		memory := service.Resources.Memory

		if cpu != "" && memory != "" {
			computeCost := c.estimateComputeCost(cpu, memory)
			resourceCosts = append(resourceCosts, types.ResourceCost{
				Type:         "compute",
				MonthlyCost:  computeCost,
				Description:  fmt.Sprintf("Compute resources for service %s", service.Name),
			})
			totalCost += computeCost
		}
	}

	// Calculate storage costs
	for name, resource := range template.Resources {
		for _, storage := range resource.Storage {
			storageCost := c.estimateStorageCost(storage.Size, storage.Type)
			resourceCosts = append(resourceCosts, types.ResourceCost{
				Type:         "storage",
				MonthlyCost:  storageCost,
				Description:  fmt.Sprintf("Storage resources for %s", name),
			})
			totalCost += storageCost
		}

		// Calculate network costs if present
		if resource.Network.Ingress != "" || resource.Network.Egress != "" {
			networkCost := c.estimateNetworkCost(resource.Network)
			resourceCosts = append(resourceCosts, types.ResourceCost{
				Type:         "network",
				MonthlyCost:  networkCost,
				Description:  fmt.Sprintf("Network resources for %s", name),
			})
			totalCost += networkCost
		}
	}

	return &types.CostEstimate{
		TotalCost:     totalCost,
		ResourceCosts: resourceCosts,
		Currency:      "USD",
	}, nil
}

func (c *CostScanner) estimateComputeCost(cpu, memory string) float64 {
	// Simple cost estimation based on CPU and memory
	// In a real implementation, this would use cloud provider pricing
	cpuCost := 50.0  // $50 per CPU per month
	memCost := 10.0  // $10 per GB per month

	// Parse CPU value (assuming format like "1" or "0.5")
	var cpuValue float64
	fmt.Sscanf(cpu, "%f", &cpuValue)

	// Parse memory value (assuming format like "2Gi" or "512Mi")
	var memValue float64
	if strings.HasSuffix(memory, "Gi") {
		fmt.Sscanf(memory, "%fGi", &memValue)
	} else if strings.HasSuffix(memory, "Mi") {
		fmt.Sscanf(memory, "%fMi", &memValue)
		memValue = memValue / 1024 // Convert to GB
	}

	return (cpuValue * cpuCost) + (memValue * memCost)
}

func (c *CostScanner) estimateStorageCost(size, storageType string) float64 {
	// Simple storage cost estimation
	// In a real implementation, this would use cloud provider pricing
	var sizeGB float64
	if strings.HasSuffix(size, "Gi") {
		fmt.Sscanf(size, "%fGi", &sizeGB)
	} else if strings.HasSuffix(size, "Ti") {
		var sizeTB float64
		fmt.Sscanf(size, "%fTi", &sizeTB)
		sizeGB = sizeTB * 1024
	}

	var costPerGB float64
	switch strings.ToLower(storageType) {
	case "ssd":
		costPerGB = 0.10 // $0.10 per GB per month for SSD
	default:
		costPerGB = 0.05 // $0.05 per GB per month for standard storage
	}

	return sizeGB * costPerGB
}

func (c *CostScanner) estimateNetworkCost(network types.Network) float64 {
	// Simple network cost estimation
	// In a real implementation, this would use cloud provider pricing
	var cost float64

	// Ingress cost (usually free but included for completeness)
	if network.Ingress != "" {
		var ingressGB float64
		if strings.HasSuffix(network.Ingress, "Gi") {
			fmt.Sscanf(network.Ingress, "%fGi", &ingressGB)
		} else if strings.HasSuffix(network.Ingress, "Ti") {
			var ingressTB float64
			fmt.Sscanf(network.Ingress, "%fTi", &ingressTB)
			ingressGB = ingressTB * 1024
		}
		cost += ingressGB * 0.00 // Ingress is typically free
	}

	// Egress cost
	if network.Egress != "" {
		var egressGB float64
		if strings.HasSuffix(network.Egress, "Gi") {
			fmt.Sscanf(network.Egress, "%fGi", &egressGB)
		} else if strings.HasSuffix(network.Egress, "Ti") {
			var egressTB float64
			fmt.Sscanf(network.Egress, "%fTi", &egressTB)
			egressGB = egressTB * 1024
		}
		cost += egressGB * 0.09 // $0.09 per GB for egress
	}

	// Request cost
	if network.Requests != "" {
		var requests float64
		if strings.HasSuffix(network.Requests, "M") {
			fmt.Sscanf(network.Requests, "%fM", &requests)
			requests *= 1000000
		}
		cost += requests * 0.0000004 // $0.40 per million requests
	}

	return cost
}
