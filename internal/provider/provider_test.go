// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"quicknode": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("QUICKNODE_ENDPOINT"); v == "" {
		t.Fatal("QUICKNODE_ENDPOINT must be set for acceptance tests")
	}

	if v := os.Getenv("QUICKNODE_API_KEY"); v == "" {
		t.Fatal("QUICKNODE_API_KEY must be set for acceptance tests")
	}
}
