data "quicknode_endpoints" "example" {
  limit  = 10
  offset = 0
}

output "endpoints" {
  value = data.quicknode_endpoints.example
}
