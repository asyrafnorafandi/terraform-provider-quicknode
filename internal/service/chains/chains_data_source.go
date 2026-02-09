// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package chains

import (
	"context"
	"fmt"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &chainsDataSource{}
	_ datasource.DataSourceWithConfigure = &chainsDataSource{}
)

// NewChainsDataSource is a helper function to simplify the provider implementation.
func NewChainsDataSource() datasource.DataSource {
	return &chainsDataSource{}
}

// chainsDataSource is the data source implementation.
type chainsDataSource struct {
	client *client.Client
}

// chainsDataSourceModel maps the data source schema data.
type chainsDataSourceModel struct {
	Chains []chainsModel `tfsdk:"chains"`
}

// chainsModel maps chains schema data.
type chainsModel struct {
	Slug     types.String    `tfsdk:"slug"`
	Networks []networksModel `tfsdk:"networks"`
}

// networksModel maps networks schema data.
type networksModel struct {
	Slug types.String `tfsdk:"slug"`
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *chainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chains"
}

// Schema defines the schema for the data source.
func (d *chainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of chains from the QuickNode API.",
		Attributes: map[string]schema.Attribute{
			"chains": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slug": schema.StringAttribute{
							Computed:    true,
							Description: "The slug of the chain.",
						},
						"networks": schema.ListNestedAttribute{
							Computed:    true,
							Description: "The list of networks for the chain.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"slug": schema.StringAttribute{
										Computed:    true,
										Description: "The slug of the network.",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "The name of the network.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *chainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state chainsDataSourceModel

	tflog.Debug(ctx, "Reading QuickNode chains")

	chains, err := d.client.GetChains(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Chains",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Received QuickNode chains", map[string]interface{}{
		"count": len(chains),
	})

	for _, chain := range chains {
		chainState := chainsModel{
			Slug: types.StringValue(chain.Slug),
		}

		for _, network := range chain.Networks {
			chainState.Networks = append(chainState.Networks, networksModel{
				Slug: types.StringValue(network.Slug),
				Name: types.StringValue(network.Name),
			})
		}

		state.Chains = append(state.Chains, chainState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *chainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
