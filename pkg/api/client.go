// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"context"
)

// Client defines the interface for interacting with the Nexlayer API
type Client interface {
	// Deployment Operations
	StartDeployment(ctx context.Context, appID string, configPath string) (*APIResponse[DeploymentResponse], error)
	GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*APIResponse[Deployment], error)
	ListDeployments(ctx context.Context) (*APIResponse[[]Deployment], error)
	GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error)

	// Domain Operations
	SaveCustomDomain(ctx context.Context, appID string, domain string) (*APIResponse[DomainResponse], error)
	ListCustomDomains(ctx context.Context, appID string) (*APIResponse[[]Domain], error)
	RemoveCustomDomain(ctx context.Context, appID string, domain string) (*APIResponse[struct{}], error)
}
