resource "quicknode_endpoint" "example" {
  chain   = "hedera"
  network = "hedera-testnet"
  label   = "test-chain"
}

output "endpoint" {
  value = quicknode_endpoint.example
}
