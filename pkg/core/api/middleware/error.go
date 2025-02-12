// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
)

// APIClient is an alias for the api.APIClient interface
type APIClient = api.APIClient

// WithErrorHandling wraps an APIClient with error handling middleware
func WithErrorHandling(next APIClient) APIClient {
	return &errorHandler{next: next}
}

type errorHandler struct {
	next APIClient
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
	logger := observability.NewLogger(observability.ERROR)

	// Check if it's already our error type
	if e, ok := err.(*errors.Error); ok {
		logger.Error(context.Background(), "Error occurred: %s", e.Message)
		return e
	}

	// Handle HTTP errors
	if httpErr, ok := err.(*schema.APIError); ok {
		switch httpErr.StatusCode {
		case http.StatusBadRequest:
			return errors.UserError(
				"Invalid request: " + h.cleanErrorMessage(httpErr.Message),
				httpErr,
			)
		case http.StatusUnauthorized:
			return errors.UserError(
				"Authentication failed. Please check your credentials.",
				httpErr,
			)
		case http.StatusForbidden:
			return errors.UserError(
				"You don't have permission to perform this action.",
				httpErr,
			)
		case http.StatusNotFound:
			return errors.UserError(
				"The requested resource was not found.",
				httpErr,
			)
		case http.StatusTooManyRequests:
			return errors.NetworkError(
				"Rate limit exceeded. Please try again later.",
				httpErr,
			)
		case http.StatusBadGateway, http.StatusServiceUnavailable:
			return errors.NetworkError(
				"The Nexlayer service is temporarily unavailable. Please try again later.",
				httpErr,
			)
		default:
			if httpErr.StatusCode >= 500 {
				return errors.SystemError(
					"An unexpected error occurred. Our team has been notified.",
					httpErr,
				)
			}
			return errors.InternalError(
				"An unexpected error occurred.",
				httpErr,
			)
		}
	}

	// Handle context cancellation
	if err == context.Canceled {
		return errors.UserError(
			"Operation cancelled by user.",
			err,
		)
	}

	// Handle context deadline
	if err == context.DeadlineExceeded {
		return errors.NetworkError(
			"Operation timed out. Please check your network connection and try again.",
			err,
		)
	}

	// Default to internal error
	logger.Error(context.Background(), "Unhandled error type: %v", err)
	return errors.InternalError(
		"An unexpected error occurred.",
		err,
	)
}

// cleanErrorMessage removes sensitive information from error messages
func (h *errorHandler) cleanErrorMessage(msg string) string {
	// Remove any potential token or key information
	msg = strings.ReplaceAll(msg, "token", "[REDACTED]")
	msg = strings.ReplaceAll(msg, "key", "[REDACTED]")

	// Remove any potential file paths
	if strings.Contains(msg, "/") {
		parts := strings.Split(msg, "/")
		msg = parts[len(parts)-1]
	}

	return msg
}
