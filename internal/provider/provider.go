// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/service/chains"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/service/endpoints"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &quicknodeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &quicknodeProvider{
			version: version,
		}
	}
}

// quicknodeProvider is the provider implementation.
type quicknodeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// quicknodeProviderModel maps provider schema data to a Go type.
type quicknodeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

// Metadata returns the provider type name.
func (p *quicknodeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "quicknode"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *quicknodeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "QuickNode is a leading blockchain infrastructure provider founded in 2017, offering " +
			"high-performance, 99.99% uptime RPC nodes, APIs, and developer tools for 80+ networks " +
			"(including Ethereum, Solana, and Bitcoin). Trusted by industry leaders, it enables developers " +
			"to build, scale, and deploy dApps, NFTs, and analytics tools without managing backend infrastructure.\n\n" +
			"Configure the provider by providing your API key and optionally the endpoint to use for the QuickNode API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The endpoint to use for the QuickNode API. Can also be set with the `QUICKNODE_ENDPOINT` environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The API key to use for the QuickNode API. Can also be set with the `QUICKNODE_API_KEY` environment variable.",
			},
		},
	}
}

// Configure prepares a QuickNode API client for data sources and resources.
func (p *quicknodeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config quicknodeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown QuickNode API Endpoint",
			"The provider cannot create the QuickNode API client as there is an unknown configuration value for the QuickNode API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the QUICKNODE_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown QuickNode API Key",
			"The provider cannot create the QuickNode API client as there is an unknown configuration value for the QuickNode API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the QUICKNODE_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("QUICKNODE_ENDPOINT")
	apiKey := os.Getenv("QUICKNODE_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing QuickNode API Endpoint",
			"The provider cannot create the QuickNode API client as there is a missing or empty value for the QuickNode API endpoint. "+
				"Set the endpoint value in the configuration or use the QUICKNODE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing QuickNode API Key",
			"The provider cannot create the QuickNode API client as there is a missing or empty value for the QuickNode API key. "+
				"Set the api_key value in the configuration or use the QUICKNODE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new QuickNode client using the configuration values
	client, err := client.NewClient(&endpoint, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create QuickNode API Client",
			"An unexpected error occurred when creating the QuickNode API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"QuickNode Client Error: "+err.Error(),
		)
		return
	}

	// Make the QuickNode client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *quicknodeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		chains.NewChainsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *quicknodeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		endpoints.NewEndpointResource,
	}
}
