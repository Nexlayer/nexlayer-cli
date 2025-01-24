package types

// CostEstimate represents cost estimation for resources
type CostEstimate struct {
	TotalCost     float64       `json:"total_cost"`
	ResourceCosts []ResourceCost `json:"resource_costs"`
	Currency      string        `json:"currency"`
}

// ResourceCost represents cost for a specific resource
type ResourceCost struct {
	ResourceType string  `json:"type"`
	MonthlyCost  float64 `json:"monthly_cost"`
	Description  string  `json:"description"`
}
