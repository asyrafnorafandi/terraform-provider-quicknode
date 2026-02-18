// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default QuickNode URL.
const HostURL string = "https://api.quicknode.com/v0"

// Client holds the configuration for the QuickNode API client.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	APIKey     string //nolint:gosec // G117: not a hardcoded secret, this holds the user-provided API key
	UserAgent  string
}

// NewClient creates a new QuickNode API client.
func NewClient(endpoint, apiKey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
		UserAgent:  "terraform-provider-quicknode",
	}

	if endpoint != nil {
		c.HostURL = *endpoint
	}

	if apiKey == nil {
		return &c, nil
	}

	c.APIKey = *apiKey

	return &c, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.APIKey != "" {
		req.Header.Set("x-api-key", c.APIKey)
	}

	res, err := c.HTTPClient.Do(req) //nolint:gosec // G704: URL is constructed from provider config, not arbitrary user input
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}
