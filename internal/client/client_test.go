// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	c, err := NewClient(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if c.API == nil {
		t.Fatal("expected non-nil API client")
	}
}

func TestNewClient_WithEndpointAndAPIKey(t *testing.T) {
	endpoint := "https://custom.api.example.com"
	apiKey := "test-key-123"

	c, err := NewClient(&endpoint, &apiKey)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if c.API == nil {
		t.Fatal("expected non-nil API client")
	}
}

func TestNewClient_SetsHeaders(t *testing.T) {
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	apiKey := "test-api-key"
	c, err := NewClient(&server.URL, &apiKey)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Make a real request through the generated client to verify headers.
	_, _ = c.API.ChainsWithResponse(context.Background())

	if got := receivedHeaders.Get("Accept"); got != "application/json" {
		t.Errorf("expected Accept %q, got %q", "application/json", got)
	}

	if got := receivedHeaders.Get("User-Agent"); got != "terraform-provider-quicknode" {
		t.Errorf("expected User-Agent %q, got %q", "terraform-provider-quicknode", got)
	}

	if got := receivedHeaders.Get("x-api-key"); got != apiKey {
		t.Errorf("expected x-api-key %q, got %q", apiKey, got)
	}
}

func TestNewClient_NoAPIKey_OmitsHeader(t *testing.T) {
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	c, err := NewClient(&server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	_, _ = c.API.ChainsWithResponse(context.Background())

	if got := receivedHeaders.Get("x-api-key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}
