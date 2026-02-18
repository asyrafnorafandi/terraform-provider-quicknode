locals {
  whitelisted_ips = [
    "10.20.10.0/24",
    "10.20.11.0/24",
    "10.20.12.0/24",
  ]
}

resource "quicknode_endpoint" "example" {
  chain   = "optimism"
  network = "optimism-sepolia"
  label   = "test-chain"
}

resource "quicknode_endpoint_whitelist_ip" "example" {
  for_each    = toset(local.whitelisted_ips)
  ip          = each.value
  endpoint_id = quicknode_endpoint.example.id
}
