package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/errors"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
)

// Client handles template storage and retrieval
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new template registry client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// ListTemplates lists all available templates
func (c *Client) ListTemplates(ctx context.Context) ([]types.NexlayerTemplate, error) {
	if ctx == nil {
		return nil, errors.NewError(errors.ErrConfigInvalid, "context is required", nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/templates", nil)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error creating request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error listing templates", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	var templates []types.NexlayerTemplate
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error decoding response", err)
	}

	return templates, nil
}

// GetTemplate retrieves a template from the registry
func (c *Client) GetTemplate(ctx context.Context, name string) (*types.NexlayerTemplate, error) {
	if ctx == nil {
		return nil, errors.NewError(errors.ErrConfigInvalid, "context is required", nil)
	}

	if name == "" {
		return nil, errors.NewError(errors.ErrTemplateInvalid, "template name is required", nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/templates/"+name, nil)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error creating request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error getting template", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewError(errors.ErrTemplateNotFound, fmt.Sprintf("template %s not found", name), nil)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	var template types.NexlayerTemplate
	if err := json.NewDecoder(resp.Body).Decode(&template); err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error decoding response", err)
	}

	return &template, nil
}

// GetTemplateVersions retrieves all versions of a template
func (c *Client) GetTemplateVersions(ctx context.Context, name string) ([]string, error) {
	if ctx == nil {
		return nil, errors.NewError(errors.ErrConfigInvalid, "context is required", nil)
	}

	if name == "" {
		return nil, errors.NewError(errors.ErrTemplateInvalid, "template name is required", nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/templates/"+name+"/versions", nil)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error creating request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error getting template versions", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewError(errors.ErrTemplateNotFound, fmt.Sprintf("template %s not found", name), nil)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, errors.NewError(errors.ErrRegistryUnavailable, "error decoding response", err)
	}

	return versions, nil
}

// PublishTemplate publishes a template to the registry
func (c *Client) PublishTemplate(ctx context.Context, template *types.NexlayerTemplate) error {
	if ctx == nil {
		return errors.NewError(errors.ErrConfigInvalid, "context is required", nil)
	}

	if template == nil {
		return errors.NewError(errors.ErrTemplateInvalid, "template is required", nil)
	}

	data, err := json.Marshal(template)
	if err != nil {
		return errors.NewError(errors.ErrTemplateInvalid, "error marshaling template", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/templates/"+template.Name, bytes.NewReader(data))
	if err != nil {
		return errors.NewError(errors.ErrRegistryUnavailable, "error creating request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.NewError(errors.ErrRegistryUnavailable, "error publishing template", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewError(errors.ErrRegistryUnavailable, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	return nil
}
