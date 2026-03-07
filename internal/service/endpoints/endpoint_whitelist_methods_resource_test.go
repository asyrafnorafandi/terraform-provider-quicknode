// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints_test

import (
	"fmt"
	"testing"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEndpointWhitelistMethodResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccEndpointWhitelistMethodResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_methods.test", "id"),
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_methods.test", "endpoint_id"),
					resource.TestCheckResourceAttr("quicknode_endpoint_whitelist_methods.test", "method.#", "2"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "quicknode_endpoint_whitelist_methods.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["quicknode_endpoint_whitelist_methods.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return rs.Primary.Attributes["endpoint_id"] + "/" + rs.Primary.Attributes["id"], nil
				},
			},
			// Update testing.
			{
				Config: testAccEndpointWhitelistMethodResourceConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_methods.test", "id"),
					resource.TestCheckResourceAttr("quicknode_endpoint_whitelist_methods.test", "method.#", "3"),
				),
			},
		},
	})
}

const testAccEndpointWhitelistMethodResourceConfig = `
resource "quicknode_endpoint" "test" {
  chain   = "optimism"
  network = "optimism-sepolia"

  security_options = {
    tokens          = true
    referrers       = false
    jwts            = false
    ips             = false
    domain_masks    = false
    hsts            = false
    cors            = true
    request_filters = true
  }
}

resource "quicknode_endpoint_whitelist_methods" "test" {
  method      = ["eth_blockNumber", "eth_getBalance"]
  endpoint_id = quicknode_endpoint.test.id
}
`

const testAccEndpointWhitelistMethodResourceConfigUpdated = `
resource "quicknode_endpoint" "test" {
  chain   = "optimism"
  network = "optimism-sepolia"

  security_options = {
    tokens          = true
    referrers       = false
    jwts            = false
    ips             = false
    domain_masks    = false
    hsts            = false
    cors            = true
    request_filters = true
  }
}

resource "quicknode_endpoint_whitelist_methods" "test" {
  method      = ["eth_blockNumber", "eth_getBalance", "eth_chainId"]
  endpoint_id = quicknode_endpoint.test.id
}
`
