// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"fmt"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/client"
	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &endpointResource{}
	_ resource.ResourceWithConfigure   = &endpointResource{}
	_ resource.ResourceWithImportState = &endpointResource{}
)

// NewEndpointResource is a helper function to simplify the provider implementation.
func NewEndpointResource() resource.Resource {
	return &endpointResource{}
}

// endpointResource is the resource implementation.
type endpointResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *endpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

// Schema defines the schema for the resource.
func (r *endpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"label": schema.StringAttribute{
				Description: "A descriptive label for the endpoint.",
				Computed:    true,
				Optional:    true,
			},
			"chain": schema.StringAttribute{
				Description: "The blockchain the endpoint is associated with.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network": schema.StringAttribute{
				Description: "The specific network of the blockchain.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_url": schema.StringAttribute{
				Description: "The HTTP URL to access the newly created endpoint.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"wss_url": schema.StringAttribute{
				Description: "The WebSocket URL to access the newly created endpoint.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_options": schema.SingleNestedAttribute{
				Description: "Security options for the endpoint.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"tokens": schema.BoolAttribute{
						Description: "Token-based authentication for the endpoint. (default: true)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"referrers": schema.BoolAttribute{
						Description: "Referrer-based access control for the endpoint. (default: false)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"jwts": schema.BoolAttribute{
						Description: "JWT-based authentication for the endpoint. (default: false)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"ips": schema.BoolAttribute{
						Description: "IP-based access control for the endpoint. (default: false)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"domain_masks": schema.BoolAttribute{
						Description: "Domain mask-based access control for the endpoint. (default: false)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"hsts": schema.BoolAttribute{
						Description: "HTTP Strict Transport Security for the endpoint. (default: false)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"cors": schema.BoolAttribute{
						Description: "Cross-Origin Resource Sharing for the endpoint. (default: true)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
			},
			"status": schema.StringAttribute{
				Description: "The status of the endpoint.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multichain": schema.BoolAttribute{
				Description: "Whether the endpoint is multichain.",
				Computed:    true,
			},
		},
	}
}

// securityString converts a Terraform bool to an API "enabled"/"disabled" string.
func securityString(val bool) string {
	if val {
		return "enabled"
	}
	return "disabled"
}

// mapSecurityOptions maps the API security options model to the Terraform resource model.
func mapSecurityOptions(api *models.EndpointSecurityOptionsModel) *models.SecurityOptionsResourceModel {
	if api == nil {
		return &models.SecurityOptionsResourceModel{
			Tokens:      types.BoolValue(true),
			Referrers:   types.BoolValue(false),
			JWTs:        types.BoolValue(false),
			IPs:         types.BoolValue(false),
			DomainMasks: types.BoolValue(false),
			HSTS:        types.BoolValue(false),
			CORS:        types.BoolValue(true),
		}
	}
	return &models.SecurityOptionsResourceModel{
		Tokens:      types.BoolValue(api.Tokens),
		Referrers:   types.BoolValue(api.Referrers),
		JWTs:        types.BoolValue(api.JWTs),
		IPs:         types.BoolValue(api.IPs),
		DomainMasks: types.BoolValue(api.DomainMasks),
		HSTS:        types.BoolValue(api.HSTS),
		CORS:        types.BoolValue(api.CORS),
	}
}

// buildSecurityOptionsMap builds a map from the Terraform resource model for the API call.
func buildSecurityOptionsMap(tf *models.SecurityOptionsResourceModel) map[string]string {
	if tf == nil {
		return map[string]string{
			"tokens":      "enabled",
			"referrers":   "disabled",
			"jwts":        "disabled",
			"ips":         "disabled",
			"domainMasks": "disabled",
			"hsts":        "disabled",
			"cors":        "enabled",
		}
	}
	return map[string]string{
		"tokens":      securityString(tf.Tokens.ValueBool()),
		"referrers":   securityString(tf.Referrers.ValueBool()),
		"jwts":        securityString(tf.JWTs.ValueBool()),
		"ips":         securityString(tf.IPs.ValueBool()),
		"domainMasks": securityString(tf.DomainMasks.ValueBool()),
		"hsts":        securityString(tf.HSTS.ValueBool()),
		"cors":        securityString(tf.CORS.ValueBool()),
	}
}

// Create a new resource.
func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.EndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new order
	endpoint, err := r.client.CreateEndpoint(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint",
			"Could not create endpoint, unexpected error: "+err.Error(),
		)
		return
	}
	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(endpoint.ID)
	plan.Chain = types.StringValue(endpoint.Chain)
	plan.Network = types.StringValue(endpoint.Network)
	plan.HTTPURL = types.StringValue(endpoint.HTTPURL)
	plan.WSSURL = types.StringValue(endpoint.WSSURL)
	plan.Status = types.StringValue(endpoint.Status)
	plan.Multichain = types.BoolValue(endpoint.Multichain)

	// Patch endpoint label if needed
	if plan.Label.ValueString() != "" {
		err = r.client.PatchEndpoint(ctx, plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error patching endpoint label",
				"Could not patch endpoint, unexpected error: "+err.Error(),
			)
			return
		}
		// Keep the label from plan since we just patched it
	} else {
		// Use label from API response if plan didn't have one
		plan.Label = types.StringValue(endpoint.Label)
	}

	// Patch security options
	err = r.client.PatchEndpointSecurity(ctx, plan.ID.ValueString(), buildSecurityOptionsMap(plan.SecurityOptions))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			"Could not patch endpoint security options, unexpected error: "+err.Error(),
		)
		return
	}

	// Read back the full state including security options
	endpoint, err = r.client.GetEndpoint(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	plan.SecurityOptions = mapSecurityOptions(&endpoint.Security.Options)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *endpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.EndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode
	endpoint, err := r.client.GetEndpoint(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite endpoint with refreshed state
	state.ID = types.StringValue(endpoint.ID)
	state.Label = types.StringValue(endpoint.Label)
	state.Chain = types.StringValue(endpoint.Chain)
	state.Network = types.StringValue(endpoint.Network)
	state.HTTPURL = types.StringValue(endpoint.HTTPURL)
	state.WSSURL = types.StringValue(endpoint.WSSURL)
	state.SecurityOptions = mapSecurityOptions(&endpoint.Security.Options)
	state.Status = types.StringValue(endpoint.Status)
	state.Multichain = types.BoolValue(endpoint.Multichain)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.EndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Patch endpoint label
	err := r.client.PatchEndpoint(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Endpoint",
			"Could not update endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Patch security options
	err = r.client.PatchEndpointSecurity(ctx, plan.ID.ValueString(), buildSecurityOptionsMap(plan.SecurityOptions))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			"Could not patch endpoint security options, unexpected error: "+err.Error(),
		)
		return
	}

	// Read back the endpoint to get the full updated state
	endpoint, err := r.client.GetEndpoint(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(endpoint.ID)
	plan.Label = types.StringValue(endpoint.Label)
	plan.Chain = types.StringValue(endpoint.Chain)
	plan.Network = types.StringValue(endpoint.Network)
	plan.HTTPURL = types.StringValue(endpoint.HTTPURL)
	plan.WSSURL = types.StringValue(endpoint.WSSURL)
	plan.SecurityOptions = mapSecurityOptions(&endpoint.Security.Options)
	plan.Status = types.StringValue(endpoint.Status)
	plan.Multichain = types.BoolValue(endpoint.Multichain)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.EndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteEndpoint(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint",
			"Could not delete endpoint, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *endpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *endpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
