# List all chains from QuickNode
data "quicknode_chains" "all" {}

output "chains" {
  value = data.quicknode_chains.all
}
