locals {
  whitelisted_domains = [
    "rpc.example.com",
    "rpc.op-sepolia.example.com",
  ]
}

resource "quicknode_endpoint" "example" {
  chain   = "optimism"
  network = "optimism-sepolia"
  label   = "test-chain"

  security_options = {
    tokens       = true
    referrers    = false
    jwts         = false
    ips          = false
    domain_masks = true # Must be set to true to use the whitelist_domain_mask resource
    hsts         = false
    cors         = true
  }
}

resource "quicknode_endpoint_whitelist_domain_mask" "example" {
  for_each    = toset(local.whitelisted_domains)
  domain_mask = each.value
  endpoint_id = quicknode_endpoint.example.id
}
