data "quicknode_chains" "example" {}

output "chains" {
  value = data.quicknode_chains.example
}
