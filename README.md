<div align="center">
  <br />
    <a href="https://astar.network/" target="_blank">
      <img src="https://mma.prnewswire.com/media/1988572/QuickNode_Logo.jpg?p=twitter" alt="Project Banner">
    </a>
  <br />

  <div>
    <img alt="Terraform" src="https://img.shields.io/badge/-Terraform-844FBA?style=for-the-badge&logo=terraform&logoColor=white" />
    <img alt="Go" src="https://img.shields.io/badge/-Go-2496ED?style=for-the-badge&logo=go&logoColor=white" />
  </div>

</div>

# Terraform Provider for QuickNode

![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/asyrafnorafandi/terraform-provider-quicknode?style=for-the-badge&logo=Github&logoColor=white)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/asyrafnorafandi/terraform-provider-quicknode/test.yml?style=for-the-badge&logo=go)
![GitHub Release](https://img.shields.io/github/v/release/asyrafnorafandi/terraform-provider-quicknode?style=for-the-badge&logo=rocket&color=green)

The QuickNode Terraform provider allows you to manage [QuickNode](https://www.quicknode.com/) blockchain infrastructure resources using Terraform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.13
- [Go](https://golang.org/doc/install) >= 1.25 (to build the provider plugin)

## Usage

```hcl
terraform {
  required_providers {
    quicknode = {
      source = "asyrafnorafandi/quicknode"
    }
  }
}

provider "quicknode" {
  # Can also be set with the QUICKNODE_API_KEY environment variable
  api_key = var.quicknode_api_key
}

resource "quicknode_endpoint" "example" {
  chain   = "optimism"
  network = "optimism-sepolia"
  label   = "optimisim-sepolia-test-chain"

  security_options = {
    tokens          = true
    referrers       = false
    jwts            = false
    ips             = true
    domain_masks    = true
    hsts            = false
    cors            = true
    request_filters = true
  }

  tags = ["env:staging", "chain:optimism"]
}

# Whitelist specific IPs for the endpoint
locals {
  whitelisted_ips = [
    "10.20.10.0/24",
    "10.20.11.0/24",
  ]
}

resource "quicknode_endpoint_whitelist_ip" "example" {
  for_each    = toset(local.whitelisted_ips)
  ip          = each.value
  endpoint_id = quicknode_endpoint.example.id
}

# Whitelist specific domain masks for the endpoint
locals {
  whitelisted_domains = [
    "rpc.example.com",
    "rpc.op-sepolia.example.com",
  ]
}

resource "quicknode_endpoint_whitelist_domain_mask" "example" {
  for_each    = toset(local.whitelisted_domains)
  domain_mask = each.value
  endpoint_id = quicknode_endpoint.example.id
}

# Whitelist specific RPC methods for the endpoint
resource "quicknode_endpoint_whitelist_methods" "example" {
  method      = ["eth_blockNumber", "eth_getBalance", "eth_chainId"]
  endpoint_id = quicknode_endpoint.example.id
}

output "optimism_chain" {
  value = quicknode_endpoint.example
}
```

## Authentication

The provider requires a QuickNode API key. You can configure it in two ways:

1. **Environment variables** (recommended):
   ```bash
   export QUICKNODE_API_KEY="your-api-key"
   ```

2. **Provider configuration** (not recommended for production):
   ```hcl
   provider "quicknode" {
     api_key = "your-api-key"
   }
   ```

You can optionally override the API base URL (defaults to `https://api.quicknode.com`):
```bash
export QUICKNODE_ENDPOINT="https://api.quicknode.com"
```

## Resources

- `quicknode_endpoint` - Creates and manages a QuickNode RPC endpoint.
- `quicknode_endpoint_whitelist_ip` - Manages IP whitelist entries for an endpoint.
- `quicknode_endpoint_whitelist_domain_mask` - Manages domain mask whitelist entries for an endpoint.
- `quicknode_endpoint_whitelist_methods` - Manages RPC method whitelist (request filters) for an endpoint.

## Data Sources

- `quicknode_chains` - Fetches the list of supported blockchain chains and their networks.
- `quicknode_endpoint` - Returns info for a specific endpoint.
- `quicknode_endpoints` - Lists info for all available endpoints.

## Developing the Provider

### Building

```bash
make build
```

### Running Tests

Unit tests:
```bash
make test
```

Acceptance tests (creates real resources):
```bash
export QUICKNODE_API_KEY="your-api-key"
make testacc
```

### Code Generation

The API client is generated from the [QuickNode OpenAPI spec](https://www.quicknode.com/api-docs/v0/swagger.json) using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen). Generated code lives in `internal/api/`.

To regenerate after updating the spec:
```bash
cd internal/api
go generate ./...
```

### Generating Documentation

```bash
make generate
```

Documentation is generated from the provider schema and example configurations in `examples/`. Do not edit files in `docs/` directly.

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
