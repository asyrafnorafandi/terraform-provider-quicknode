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
	_ datasource.DataSource              = &endpointsDataSource{}
	_ datasource.DataSourceWithConfigure = &endpointsDataSource{}
)

// NewEndpointsDataSource is a helper function to simplify the provider implementation.
func NewEndpointsDataSource() datasource.DataSource {
	return &endpointsDataSource{}
}

// endpointsDataSource is the data source implementation.
type endpointsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *endpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoints"
}

// Schema defines the schema for the data source.
func (d *endpointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists info for all available endpoints.",
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "The number of endpoints to return.",
				Optional:    true,
			},
			"offset": schema.Int64Attribute{
				Description: "The offset to start from.",
				Optional:    true,
			},
			// TODO: Add tag support
			// "tag_labels": schema.ListAttribute{
			// 	Description: "The total number of endpoints.",
			// 	Optional:    true,
			// },
			"endpoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "A unique identifier for the created endpoint.",
							Computed:    true,
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
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *endpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.EndpointsDataSourceModel

	// Read the user's config (the values they set in the .tf file)
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode
	endpoints, err := d.client.ListEndpoints(ctx, config.Limit.ValueInt64(), config.Offset.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing QuickNode Endpoints",
			"Could not list QuickNode endpoints: "+err.Error(),
		)
		return
	}

	var state models.EndpointsDataSourceModel
	state.Limit = config.Limit
	state.Offset = config.Offset
	for _, endpoint := range *endpoints {
		state.Endpoints = append(state.Endpoints, struct {
			ID      types.String `tfsdk:"id"`
			Label   types.String `tfsdk:"label"`
			Chain   types.String `tfsdk:"chain"`
			Network types.String `tfsdk:"network"`
			HTTPURL types.String `tfsdk:"http_url"`
			WSSURL  types.String `tfsdk:"wss_url"`
		}{
			ID:      types.StringValue(endpoint.ID),
			Label:   types.StringValue(endpoint.Label),
			Chain:   types.StringValue(endpoint.Chain),
			Network: types.StringValue(endpoint.Network),
			HTTPURL: types.StringValue(endpoint.HTTPURL),
			WSSURL:  types.StringValue(endpoint.WSSURL),
		})
	}
	// // Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *endpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
