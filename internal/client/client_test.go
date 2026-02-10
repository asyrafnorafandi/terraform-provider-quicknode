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

	if c.HostURL != HostURL {
		t.Errorf("expected HostURL %q, got %q", HostURL, c.HostURL)
	}

	if c.APIKey != "" {
		t.Errorf("expected empty APIKey, got %q", c.APIKey)
	}

	if c.UserAgent != "terraform-provider-quicknode" {
		t.Errorf("expected UserAgent %q, got %q", "terraform-provider-quicknode", c.UserAgent)
	}
}

func TestNewClient_WithEndpointAndAPIKey(t *testing.T) {
	endpoint := "https://custom.api.example.com"
	apiKey := "test-key-123"

	c, err := NewClient(&endpoint, &apiKey)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if c.HostURL != endpoint {
		t.Errorf("expected HostURL %q, got %q", endpoint, c.HostURL)
	}

	if c.APIKey != apiKey {
		t.Errorf("expected APIKey %q, got %q", apiKey, c.APIKey)
	}
}

func TestDoRequest_SetsHeaders(t *testing.T) {
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	apiKey := "test-api-key"
	c, _ := NewClient(&server.URL, &apiKey)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := c.doRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := receivedHeaders.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", got)
	}

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

func TestDoRequest_NoAPIKey_OmitsHeader(t *testing.T) {
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, _ := NewClient(&server.URL, nil)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := c.doRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := receivedHeaders.Get("x-api-key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}

func TestDoRequest_Non200_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	c, _ := NewClient(&server.URL, nil)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := c.doRequest(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}

	expected := `status: 401, body: {"error":"unauthorized"}`
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestDoRequest_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	c, _ := NewClient(&server.URL, nil)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := c.doRequest(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestDoRequest_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, _ := NewClient(&server.URL, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := c.doRequest(ctx, req)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
