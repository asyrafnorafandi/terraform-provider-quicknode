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

data "quicknode_endpoint" "example" {
  id = "111111"
}

output "endpoint" {
  value = data.quicknode_endpoint.example
}
