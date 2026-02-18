// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"fmt"
	"strings"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &endpointWhitelistIPResource{}
	_ resource.ResourceWithConfigure   = &endpointWhitelistIPResource{}
	_ resource.ResourceWithImportState = &endpointWhitelistIPResource{}
)

// NewEndpointResource is a helper function to simplify the provider implementation.
func NewEndpointWhitelistIPResource() resource.Resource {
	return &endpointWhitelistIPResource{}
}

// endpointResource is the resource implementation.
type endpointWhitelistIPResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *endpointWhitelistIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_whitelist_ip"
}

// Schema defines the schema for the resource.
func (r *endpointWhitelistIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a new endpoint in the QuickNode API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A unique identifier for the created endpoint.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Description: "The IP address to whitelist.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_id": schema.StringAttribute{
				Description: "The ID of the endpoint to whitelist the IP address for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create a new resource.
func (r *endpointWhitelistIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.EndpointWhitelistIPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new order
	whitelistIP, err := r.client.CreateEndpointWhitelistIP(ctx, plan.EndpointID.ValueString(), plan.IP.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist IP",
			"Could not create endpoint whitelist IP, unexpected error: "+err.Error(),
		)
		return
	}
	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(whitelistIP.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *endpointWhitelistIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.EndpointWhitelistIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode
	endpoint, err := r.client.GetEndpoint(ctx, state.EndpointID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Find specific IP in IP array
	for _, ip := range endpoint.Security.IPs {
		if ip.ID == state.ID.ValueString() {
			state.IP = types.StringValue(ip.IP)
			break
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointWhitelistIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Whitelisted IPs cannot be updated in-place. This is a bug in the provider.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointWhitelistIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.EndpointWhitelistIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteEndpointWhitelistIP(ctx, state.EndpointID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint Whitelist IP",
			"Could not delete endpoint whitelist IP, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *endpointWhitelistIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// ImportState imports the state of the resource into the Terraform state.
func (r *endpointWhitelistIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: endpoint_id/ip_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("endpoint_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
