// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/api"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &endpointWhitelistMethodsResource{}
	_ resource.ResourceWithConfigure   = &endpointWhitelistMethodsResource{}
	_ resource.ResourceWithImportState = &endpointWhitelistMethodsResource{}
)

// NewEndpointWhitelistMethodsResource is a helper function to simplify the provider implementation.
func NewEndpointWhitelistMethodsResource() resource.Resource {
	return &endpointWhitelistMethodsResource{}
}

// endpointWhitelistMethodsResource is the resource implementation.
type endpointWhitelistMethodsResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *endpointWhitelistMethodsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_whitelist_methods"
}

// Schema defines the schema for the resource.
func (r *endpointWhitelistMethodsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a new endpoint whitelist method (request filter) in the QuickNode API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A unique identifier for the created request filter.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"method": schema.SetAttribute{
				Description: "The set of RPC method names to whitelist.",
				Required:    true,
				ElementType: types.StringType,
			},
			"endpoint_id": schema.StringAttribute{
				Description: "The ID of the endpoint to create the request filter for.",
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
func (r *endpointWhitelistMethodsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan models.EndpointWhitelistMethodsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	methods := expandStringSet(ctx, plan.Method)

	// Create new request filter.
	createResp, err := r.client.API.CreateRequestFilterWithResponse(ctx, plan.EndpointID.ValueString(), api.CreateRequestFilterJSONRequestBody{
		Method: &methods,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist method",
			"Could not create endpoint whitelist method, unexpected error: "+err.Error(),
		)
		return
	}
	if createResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist method",
			fmt.Sprintf("API returned status %d: %s", createResp.StatusCode(), string(createResp.Body)),
		)
		return
	}

	if createResp.JSON200 == nil || createResp.JSON200.Data == nil || createResp.JSON200.Data.Id == nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist method",
			"API returned an empty response",
		)
		return
	}

	plan.ID = types.StringValue(*createResp.JSON200.Data.Id)

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *endpointWhitelistMethodsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state models.EndpointWhitelistMethodsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode.
	showResp, err := r.client.API.ShowEndpointWithResponse(ctx, state.EndpointID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+state.EndpointID.ValueString()+": "+err.Error(),
		)
		return
	}
	if showResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			fmt.Sprintf("API returned status %d: %s", showResp.StatusCode(), string(showResp.Body)),
		)
		return
	}

	endpoint := showResp.JSON200.Data

	// Find specific request filter by ID.
	found := false
	if endpoint.Security.RequestFilters != nil {
		for _, rf := range *endpoint.Security.RequestFilters {
			if rf.Id != nil && *rf.Id == state.ID.ValueString() {
				found = true
				if rf.Method != nil {
					state.Method = flattenStringSet(*rf.Method)
				}
				break
			}
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointWhitelistMethodsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan models.EndpointWhitelistMethodsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	methods := expandStringSet(ctx, plan.Method)

	// Update request filter.
	updateResp, err := r.client.API.UpdateRequestFilterWithResponse(ctx, plan.EndpointID.ValueString(), plan.ID.ValueString(), api.UpdateRequestFilterJSONRequestBody{
		Method: &methods,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Endpoint Whitelist Method",
			"Could not update endpoint whitelist method, unexpected error: "+err.Error(),
		)
		return
	}
	if updateResp.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Endpoint Whitelist Method",
			fmt.Sprintf("API returned status %d: %s", updateResp.StatusCode(), string(updateResp.Body)),
		)
		return
	}

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointWhitelistMethodsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state models.EndpointWhitelistMethodsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing request filter.
	deleteResp, err := r.client.API.DeleteRequestFilterWithResponse(ctx, state.EndpointID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint Whitelist Method",
			"Could not delete endpoint whitelist method, unexpected error: "+err.Error(),
		)
		return
	}
	if deleteResp.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint Whitelist Method",
			fmt.Sprintf("API returned status %d: %s", deleteResp.StatusCode(), string(deleteResp.Body)),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *endpointWhitelistMethodsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *endpointWhitelistMethodsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: endpoint_id/request_filter_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("endpoint_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// expandStringSet converts a Terraform set of strings to a Go string slice.
func expandStringSet(ctx context.Context, set types.Set) []string {
	var result []string
	set.ElementsAs(ctx, &result, false)
	return result
}

// flattenStringSet converts a Go string slice to a Terraform set of strings.
func flattenStringSet(values []string) types.Set {
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	s, _ := types.SetValue(types.StringType, elems)
	return s
}
