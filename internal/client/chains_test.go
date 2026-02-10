// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetChains_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains" {
			t.Errorf("expected path /chains, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"slug": "ethereum",
					"networks": [
						{"slug": "mainnet", "name": "Mainnet"},
						{"slug": "goerli", "name": "Goerli"}
					]
				},
				{
					"slug": "solana",
					"networks": [
						{"slug": "mainnet-beta", "name": "Mainnet Beta"}
					]
				}
			]
		}`))
	}))
	defer server.Close()

	apiKey := "test-key"
	c, _ := NewClient(&server.URL, &apiKey)

	chains, err := c.GetChains(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(chains) != 2 {
		t.Fatalf("expected 2 chains, got %d", len(chains))
	}

	if chains[0].Slug != "ethereum" {
		t.Errorf("expected first chain slug %q, got %q", "ethereum", chains[0].Slug)
	}

	if len(chains[0].Networks) != 2 {
		t.Fatalf("expected 2 networks for ethereum, got %d", len(chains[0].Networks))
	}

	if chains[0].Networks[0].Slug != "mainnet" {
		t.Errorf("expected first network slug %q, got %q", "mainnet", chains[0].Networks[0].Slug)
	}

	if chains[0].Networks[0].Name != "Mainnet" {
		t.Errorf("expected first network name %q, got %q", "Mainnet", chains[0].Networks[0].Name)
	}

	if chains[1].Slug != "solana" {
		t.Errorf("expected second chain slug %q, got %q", "solana", chains[1].Slug)
	}

	if len(chains[1].Networks) != 1 {
		t.Fatalf("expected 1 network for solana, got %d", len(chains[1].Networks))
	}
}

func TestGetChains_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	apiKey := "test-key"
	c, _ := NewClient(&server.URL, &apiKey)

	chains, err := c.GetChains(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(chains) != 0 {
		t.Errorf("expected 0 chains, got %d", len(chains))
	}
}

func TestGetChains_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "invalid api key"}`))
	}))
	defer server.Close()

	apiKey := "bad-key"
	c, _ := NewClient(&server.URL, &apiKey)

	_, err := c.GetChains(context.Background())
	if err == nil {
		t.Fatal("expected error for unauthorized response, got nil")
	}
}

func TestGetChains_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	apiKey := "test-key"
	c, _ := NewClient(&server.URL, &apiKey)

	_, err := c.GetChains(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGetChains_ServerDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	c, _ := NewClient(&server.URL, nil)

	_, err := c.GetChains(context.Background())
	if err == nil {
		t.Fatal("expected error for closed server, got nil")
	}
}
