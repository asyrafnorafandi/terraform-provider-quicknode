resource "quicknode_endpoint" "example" {
  chain   = "optimism"
  network = "optimism-sepolia"
  label   = "test-chain"

  security_options = {
    tokens          = true
    referrers       = false
    jwts            = false
    ips             = false
    domain_masks    = false
    hsts            = false
    cors            = true
    request_filters = true # Must be set to true to use the whitelist_methods resource
  }
}

resource "quicknode_endpoint_whitelist_methods" "example" {
  method      = ["eth_blockNumber", "eth_getBalance"]
  endpoint_id = quicknode_endpoint.example.id
}
