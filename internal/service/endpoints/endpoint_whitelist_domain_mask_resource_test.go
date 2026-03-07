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

func TestAccEndpointWhitelistDomainMaskResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccEndpointWhitelistDomainMaskResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_domain_mask.test", "id"),
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_domain_mask.test", "endpoint_id"),
					resource.TestCheckResourceAttrSet("quicknode_endpoint_whitelist_domain_mask.test", "domain_mask"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "quicknode_endpoint_whitelist_domain_mask.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["quicknode_endpoint_whitelist_domain_mask.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return rs.Primary.Attributes["endpoint_id"] + "/" + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

const testAccEndpointWhitelistDomainMaskResourceConfig = `
resource "quicknode_endpoint" "test" {
  chain   = "optimism"
  network = "optimism-sepolia"

  security_options = {
    tokens          = true
    referrers       = false
    jwts            = false
    ips             = false
    domain_masks    = true
    hsts            = false
    cors            = true
    request_filters = true
  }
}

resource "quicknode_endpoint_whitelist_domain_mask" "test" {
  domain_mask = "rpc.example.com"
  endpoint_id = quicknode_endpoint.test.id
}
`
