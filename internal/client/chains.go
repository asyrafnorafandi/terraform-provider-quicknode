// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"
)

var chainsURL = "/chains"

// GetChains returns the list of chains from the QuickNode API.
func (c *Client) GetChains(ctx context.Context) ([]models.ChainModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.HostURL, chainsURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []models.ChainModel `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
