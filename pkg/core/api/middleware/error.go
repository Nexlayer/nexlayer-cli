// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// WithErrorHandling wraps an APIClient with error handling middleware
func WithErrorHandling(next api.APIClient) api.APIClient {
	return &errorHandler{next: next}
}

type errorHandler struct {
	next api.APIClient
}

func (h *errorHandler) GetDeploymentInfo(ctx context.Context, namespace, appID string) (*schema.APIResponse[schema.Deployment], error) {
	resp, err := h.next.GetDeploymentInfo(ctx, namespace, appID)
	if err != nil {
		return nil, h.handleError(err)
	}
	return resp, nil
}

func (h *errorHandler) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	resp, err := h.next.ListDeployments(ctx)
	if err != nil {
		return nil, h.handleError(err)
	}
	return resp, nil
}

func (h *errorHandler) SaveCustomDomain(ctx context.Context, appID, domain string) (*schema.APIResponse[struct{}], error) {
	resp, err := h.next.SaveCustomDomain(ctx, appID, domain)
	if err != nil {
		return nil, h.handleError(err)
	}
	return resp, nil
}

func (h *errorHandler) StartDeployment(ctx context.Context, appID, yamlFile string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	resp, err := h.next.StartDeployment(ctx, appID, yamlFile)
	if err != nil {
		return nil, h.handleError(err)
	}
	return resp, nil
}

func (h *errorHandler) GetLogs(ctx context.Context, namespace, appID string, follow bool, tail int) ([]string, error) {
	logs, err := h.next.GetLogs(ctx, namespace, appID, follow, tail)
	if err != nil {
		return nil, h.handleError(err)
	}
	return logs, nil
}

func (h *errorHandler) SendFeedback(ctx context.Context, text string) error {
	err := h.next.SendFeedback(ctx, text)
	if err != nil {
		return h.handleError(err)
	}
	return nil
}

func (h *errorHandler) handleError(err error) error {
	if apiErr, ok := err.(*schema.APIError); ok {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("authentication failed: %w", err)
		case http.StatusNotFound:
			return fmt.Errorf("resource not found: %w", err)
		case http.StatusBadRequest:
			return fmt.Errorf("invalid request: %w", err)
		default:
			return fmt.Errorf("API error: %w", err)
		}
	}
	return fmt.Errorf("unexpected error: %w", err)
}
