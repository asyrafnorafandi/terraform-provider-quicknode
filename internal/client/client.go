// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"time"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/api"
)

// HostURL - Default QuickNode URL.
const HostURL string = "https://api.quicknode.com"

// Client wraps the generated QuickNode API client.
type Client struct {
	API *api.ClientWithResponses
}

// NewClient creates a new QuickNode API client.
func NewClient(endpoint, apiKey *string) (*Client, error) {
	host := HostURL
	if endpoint != nil {
		host = *endpoint
	}

	key := ""
	if apiKey != nil {
		key = *apiKey
	}

	c, err := api.NewClientWithResponses(host,
		api.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
		api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("User-Agent", "terraform-provider-quicknode")
			req.Header.Set("Accept", "application/json")
			if key != "" {
				req.Header.Set("x-api-key", key)
			}
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Client{API: c}, nil
}
