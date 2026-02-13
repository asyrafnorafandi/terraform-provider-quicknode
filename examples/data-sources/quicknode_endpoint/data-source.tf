data "quicknode_endpoint" "example" {
  id = "111111"
}

output "endpoint" {
  value = data.quicknode_endpoint.example
}
