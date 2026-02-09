# Contributing to terraform-provider-quicknode

Thank you for your interest in contributing to the QuickNode Terraform provider!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/asyrafnorafandi/terraform-provider-quicknode.git
   cd terraform-provider-quicknode
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the provider:
   ```bash
   make build
   ```

## Making Changes

1. Fork the repository and create a feature branch from `main`.
2. Make your changes, following the conventions below.
3. Add or update tests as needed.
4. Run `make fmt` and `make lint` to ensure code formatting and linting pass.
5. Run `make test` to verify all tests pass.
6. If you changed provider schema or examples, run `make generate` and commit the generated docs.
7. Open a pull request against `main`.

## Conventions

- **Plugin Framework only** — this provider uses `terraform-plugin-framework`. The linter blocks imports from the legacy `terraform-plugin-sdk/v2`.
- New data sources go in `internal/service/<name>/` implementing `datasource.DataSource`.
- New resources follow the same pattern implementing `resource.Resource`.
- Register new data sources/resources in `internal/provider/provider.go`.
- The client layer (`internal/client/`) handles raw HTTP calls; the service layer maps between API models and Terraform schema types.
- All exported client methods should accept a `context.Context` as the first parameter.

## Running Acceptance Tests

Acceptance tests create real resources against the QuickNode API. You will need valid credentials:

```bash
export QUICKNODE_ENDPOINT="https://api.quicknode.com/v0"
export QUICKNODE_API_KEY="your-api-key"
make testacc
```

## Reporting Issues

Please use [GitHub Issues](https://github.com/asyrafnorafandi/terraform-provider-quicknode/issues) to report bugs or request features.

## License

By contributing, you agree that your contributions will be licensed under the [MPL-2.0 License](LICENSE).
