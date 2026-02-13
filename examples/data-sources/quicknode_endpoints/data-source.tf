terraform {
  required_providers {
    quicknode = {
      source = "registry.terraform.io/asyrafnorafandi/quicknode"
    }
  }
}

provider "quicknode" {
  # Set via QUICKNODE_ENDPOINT environment variable, or override here:
  # endpoint = "https://api.quicknode.com/v0"

  # Set via QUICKNODE_API_KEY environment variable, or override here:
  # api_key = "QN_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

data "quicknode_endpoints" "example" {
  limit  = 10
  offset = 0
}

output "endpoints" {
  value = data.quicknode_endpoints.example
}
