// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"fmt"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &endpointDataSource{}
	_ datasource.DataSourceWithConfigure = &endpointDataSource{}
)

// NewEndpointDataSource is a helper function to simplify the provider implementation.
func NewEndpointDataSource() datasource.DataSource {
	return &endpointDataSource{}
}

// endpointDataSource is the data source implementation.
type endpointDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *endpointDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

// Schema defines the schema for the data source.
func (d *endpointDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns info for a specific endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A unique identifier for the created endpoint.",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "A descriptive label for the endpoint.",
				Computed:    true,
			},
			"chain": schema.StringAttribute{
				Description: "The blockchain the endpoint is associated with.",
				Computed:    true,
			},
			"network": schema.StringAttribute{
				Description: "The specific network of the blockchain.",
				Computed:    true,
			},
			"http_url": schema.StringAttribute{
				Description: "The HTTP URL to access the newly created endpoint.",
				Computed:    true,
			},
			"wss_url": schema.StringAttribute{
				Description: "The WebSocket URL to access the newly created endpoint.",
				Computed:    true,
			},
			"security_options": schema.SingleNestedAttribute{
				Description: "Security options for the endpoint.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"tokens": schema.BoolAttribute{
						Description: "Token-based authentication for the endpoint.",
						Computed:    true,
					},
					"referrers": schema.BoolAttribute{
						Description: "Referrer-based access control for the endpoint.",
						Computed:    true,
					},
					"jwts": schema.BoolAttribute{
						Description: "JWT-based authentication for the endpoint.",
						Computed:    true,
					},
					"ips": schema.BoolAttribute{
						Description: "IP-based access control for the endpoint.",
						Computed:    true,
					},
					"domain_masks": schema.BoolAttribute{
						Description: "Domain mask-based access control for the endpoint.",
						Computed:    true,
					},
					"hsts": schema.BoolAttribute{
						Description: "HTTP Strict Transport Security for the endpoint.",
						Computed:    true,
					},
					"cors": schema.BoolAttribute{
						Description: "Cross-Origin Resource Sharing for the endpoint.",
						Computed:    true,
					},
				},
			},
			"status": schema.StringAttribute{
				Description: "The status of the endpoint.",
				Computed:    true,
			},
			"multichain": schema.BoolAttribute{
				Description: "Whether the endpoint is multichain.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *endpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.EndpointResourceModel

	// Read the user's config (the values they set in the .tf file)
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode
	endpoint, err := d.client.GetEndpoint(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+config.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite endpoint with refreshed state
	var state models.EndpointResourceModel
	state.ID = types.StringValue(endpoint.ID)
	state.Label = types.StringValue(endpoint.Label)
	state.Chain = types.StringValue(endpoint.Chain)
	state.Network = types.StringValue(endpoint.Network)
	state.HTTPURL = types.StringValue(endpoint.HTTPURL)
	state.WSSURL = types.StringValue(endpoint.WSSURL)
	state.SecurityOptions = mapSecurityOptions(&endpoint.Security.Options)
	state.Status = types.StringValue(endpoint.Status)
	state.Multichain = types.BoolValue(endpoint.Multichain)

	// // Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *endpointDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
