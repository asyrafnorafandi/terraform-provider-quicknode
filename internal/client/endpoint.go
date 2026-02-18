// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var endpointsURL = "/endpoints"

// GetEndpoint returns a sepcific endpoint based on ID from the QuickNode API.
func (c *Client) GetEndpoint(ctx context.Context, id string) (*models.EndpointModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s/%s", c.HostURL, endpointsURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data  models.EndpointModel `json:"data"`
		Error string               `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return &response.Data, nil
}

// CreateEndpoint creates a new endpoint in the QuickNode API.
func (c *Client) CreateEndpoint(ctx context.Context, endpoint models.EndpointResourceModel) (*models.EndpointModel, error) {
	jsonBody, err := json.Marshal(map[string]interface{}{
		"chain":   endpoint.Chain.ValueString(),
		"network": endpoint.Network.ValueString(),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.HostURL, endpointsURL), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data  models.EndpointModel `json:"data"`
		Error string               `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return &response.Data, nil
}

// PatchEndpoint updates an existing endpoint in the QuickNode API.
func (c *Client) PatchEndpoint(ctx context.Context, endpoint models.EndpointResourceModel) error {
	jsonBody, err := json.Marshal(map[string]interface{}{
		"label": endpoint.Label.ValueString(),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s%s/%s", c.HostURL, endpointsURL, endpoint.ID.ValueString()), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var response struct {
		Data  bool   `json:"data"`
		Error string `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return nil
}

// PatchEndpointSecurity updates the security options for an existing endpoint.
func (c *Client) PatchEndpointSecurity(ctx context.Context, id string, options map[string]string) error {
	jsonBody, err := json.Marshal(map[string]interface{}{
		"options": options,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s%s/%s/security_options", c.HostURL, endpointsURL, id), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var response struct {
		Data  bool   `json:"data"`
		Error string `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return nil
}

// DeleteEndpoint deletes an existing endpoint in the QuickNode API.
func (c *Client) DeleteEndpoint(ctx context.Context, endpoint models.EndpointResourceModel) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s/%s", c.HostURL, endpointsURL, endpoint.ID.ValueString()), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var response struct {
		Result bool `json:"result"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	tflog.Info(ctx, "Deleted endpoint:", map[string]interface{}{"Result": response.Result})

	return nil
}

// ListEndpoints returns a list of endpoints from the QuickNode API.
func (c *Client) ListEndpoints(ctx context.Context, limit int64, offset int64) (*[]models.EndpointModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?limit=%d&offset=%d", c.HostURL, endpointsURL, limit, offset), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data  []models.EndpointModel `json:"data"`
		Error string                 `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return &response.Data, nil
}

// Create a new whitelist IP for an existing endpoint.
func (c *Client) CreateEndpointWhitelistIP(ctx context.Context, endpointID string, ip string) (*models.EndpointSecurityWhitelistIPsModel, error) {
	jsonBody, err := json.Marshal(map[string]interface{}{
		"ip": ip,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s/%s/security/ips", c.HostURL, endpointsURL, endpointID), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data  models.EndpointSecurityWhitelistIPsModel `json:"data"`
		Error string                                   `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return &response.Data, nil
}

// Delete a whitelist IP for an existing endpoint.
func (c *Client) DeleteEndpointWhitelistIP(ctx context.Context, endpointID string, whitelistIPID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s/%s/security/ips/%s", c.HostURL, endpointsURL, endpointID, whitelistIPID), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var response struct {
		Data  bool   `json:"data"`
		Error string `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return fmt.Errorf("QuickNode API Error: %s", response.Error)
	}

	return nil
}
