// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/api"
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
	_ resource.Resource                = &endpointWhitelistDomainMaskResource{}
	_ resource.ResourceWithConfigure   = &endpointWhitelistDomainMaskResource{}
	_ resource.ResourceWithImportState = &endpointWhitelistDomainMaskResource{}
)

// NewEndpointWhitelistDomainMaskResource is a helper function to simplify the provider implementation.
func NewEndpointWhitelistDomainMaskResource() resource.Resource {
	return &endpointWhitelistDomainMaskResource{}
}

// endpointWhitelistDomainMaskResource is the resource implementation.
type endpointWhitelistDomainMaskResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *endpointWhitelistDomainMaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_whitelist_domain_mask"
}

// Schema defines the schema for the resource.
func (r *endpointWhitelistDomainMaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a new endpoint whitelist domain mask in the QuickNode API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A unique identifier for the created endpoint whitelist domain mask.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_mask": schema.StringAttribute{
				Description: "The domain mask to whitelist.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_id": schema.StringAttribute{
				Description: "The ID of the endpoint to whitelist the domain mask for.",
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
func (r *endpointWhitelistDomainMaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan models.EndpointWhitelistDomainMaskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainMask := plan.DomainMask.ValueString()
	// Create new whitelist domain mask.
	createResp, err := r.client.API.CreateDomainMaskWithResponse(ctx, plan.EndpointID.ValueString(), api.CreateDomainMaskJSONRequestBody{
		DomainMask: &domainMask,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist domain mask",
			"Could not create endpoint whitelist domain mask, unexpected error: "+err.Error(),
		)
		return
	}
	if createResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating endpoint whitelist domain mask",
			fmt.Sprintf("API returned status %d: %s", createResp.StatusCode(), string(createResp.Body)),
		)
		return
	}

	// The CreateDomainMask response doesn't have a typed JSON200 in the spec, so we parse the raw body.
	var dmResp struct {
		Data struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"data"`
	}
	if err := json.Unmarshal(createResp.Body, &dmResp); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing whitelist domain mask response",
			"Could not parse response: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(dmResp.Data.ID)

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *endpointWhitelistDomainMaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state models.EndpointWhitelistDomainMaskResourceModel
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
			"Could not read QuickNode endpoint ID "+state.ID.ValueString()+": "+err.Error(),
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

	// Find specific domain mask in domain masks array.
	if endpoint.Security.DomainMasks != nil {
		for _, dm := range *endpoint.Security.DomainMasks {
			if dm.Id != nil && *dm.Id == state.ID.ValueString() {
				if dm.Domain != nil {
					state.DomainMask = types.StringValue(*dm.Domain)
				}
				break
			}
		}
	}

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointWhitelistDomainMaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Whitelisted domain masks cannot be updated in-place. This is a bug in the provider.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointWhitelistDomainMaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state models.EndpointWhitelistDomainMaskResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing whitelist domain mask.
	deleteResp, err := r.client.API.DeleteDomainMaskWithResponse(ctx, state.EndpointID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint Whitelist Domain Mask",
			"Could not delete endpoint whitelist domain mask, unexpected error: "+err.Error(),
		)
		return
	}
	if deleteResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint Whitelist Domain Mask",
			fmt.Sprintf("API returned status %d: %s", deleteResp.StatusCode(), string(deleteResp.Body)),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *endpointWhitelistDomainMaskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *endpointWhitelistDomainMaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: endpoint_id/domain_mask_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("endpoint_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
