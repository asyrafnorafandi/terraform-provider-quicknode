resource "quicknode_endpoint" "example" {
  chain   = "optimism"
  network = "optimism-sepolia"
  label   = "test-chain"

  security_options = {
    tokens       = true
    referrers    = false
    jwts         = false
    ips          = false
    domain_masks = false
    hsts         = false
    cors         = true
  }
}

output "endpoint" {
  value = quicknode_endpoint.example
}
