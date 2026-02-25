// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asyrafnorafandi/terraform-provider-quicknode/internal/api"
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
func securityString(val bool) api.UpdateSecurityOptionsJSONBodyOptionsTokens {
	if val {
		return "enabled"
	}
	return "disabled"
}

// securityOptionsResponse is used to parse security options from the raw response body,
// because the OpenAPI spec doesn't include hsts/cors in the response schema.
type securityOptionsResponse struct {
	Options struct {
		Tokens      *bool `json:"tokens"`
		Referrers   *bool `json:"referrers"`
		JWTs        *bool `json:"jwts"`
		IPs         *bool `json:"ips"`
		DomainMasks *bool `json:"domainMasks"`
		HSTS        *bool `json:"hsts"`
		CORS        *bool `json:"cors"`
	} `json:"options"`
}

// parseSecurityOptions extracts security options from the raw JSON body, including
// hsts and cors which are not defined in the OpenAPI spec response schema.
func parseSecurityOptions(body []byte) *models.SecurityOptionsResourceModel {
	var envelope struct {
		Data struct {
			Security securityOptionsResponse `json:"security"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return defaultSecurityOptions()
	}
	opts := envelope.Data.Security.Options
	return &models.SecurityOptionsResourceModel{
		Tokens:      types.BoolPointerValue(opts.Tokens),
		Referrers:   types.BoolPointerValue(opts.Referrers),
		JWTs:        types.BoolPointerValue(opts.JWTs),
		IPs:         types.BoolPointerValue(opts.IPs),
		DomainMasks: types.BoolPointerValue(opts.DomainMasks),
		HSTS:        types.BoolPointerValue(opts.HSTS),
		CORS:        types.BoolPointerValue(opts.CORS),
	}
}

func defaultSecurityOptions() *models.SecurityOptionsResourceModel {
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

// buildSecurityOptionsBody builds the request body for updating security options.
func buildSecurityOptionsBody(tf *models.SecurityOptionsResourceModel) api.UpdateSecurityOptionsJSONRequestBody {
	if tf == nil {
		tokens := api.UpdateSecurityOptionsJSONBodyOptionsTokens("enabled")
		referrers := api.UpdateSecurityOptionsJSONBodyOptionsReferrers("disabled")
		jwts := api.UpdateSecurityOptionsJSONBodyOptionsJwts("disabled")
		ips := api.UpdateSecurityOptionsJSONBodyOptionsIps("disabled")
		domainMasks := api.UpdateSecurityOptionsJSONBodyOptionsDomainMasks("disabled")
		hsts := api.UpdateSecurityOptionsJSONBodyOptionsHsts("disabled")
		cors := api.UpdateSecurityOptionsJSONBodyOptionsCors("enabled")
		return api.UpdateSecurityOptionsJSONRequestBody{
			Options: struct {
				Cors           *api.UpdateSecurityOptionsJSONBodyOptionsCors           `json:"cors,omitempty"`
				DomainMasks    *api.UpdateSecurityOptionsJSONBodyOptionsDomainMasks    `json:"domainMasks,omitempty"`
				Hsts           *api.UpdateSecurityOptionsJSONBodyOptionsHsts           `json:"hsts,omitempty"`
				IpCustomHeader *api.UpdateSecurityOptionsJSONBodyOptionsIpCustomHeader `json:"ipCustomHeader,omitempty"`
				Ips            *api.UpdateSecurityOptionsJSONBodyOptionsIps            `json:"ips,omitempty"`
				Jwts           *api.UpdateSecurityOptionsJSONBodyOptionsJwts           `json:"jwts,omitempty"`
				Referrers      *api.UpdateSecurityOptionsJSONBodyOptionsReferrers      `json:"referrers,omitempty"`
				RequestFilters *api.UpdateSecurityOptionsJSONBodyOptionsRequestFilters `json:"requestFilters,omitempty"`
				Tokens         *api.UpdateSecurityOptionsJSONBodyOptionsTokens         `json:"tokens,omitempty"`
			}{
				Tokens:      &tokens,
				Referrers:   &referrers,
				Jwts:        &jwts,
				Ips:         &ips,
				DomainMasks: &domainMasks,
				Hsts:        &hsts,
				Cors:        &cors,
			},
		}
	}

	tokens := securityString(tf.Tokens.ValueBool())
	referrers := api.UpdateSecurityOptionsJSONBodyOptionsReferrers(securityString(tf.Referrers.ValueBool()))
	jwts := api.UpdateSecurityOptionsJSONBodyOptionsJwts(securityString(tf.JWTs.ValueBool()))
	ips := api.UpdateSecurityOptionsJSONBodyOptionsIps(securityString(tf.IPs.ValueBool()))
	domainMasks := api.UpdateSecurityOptionsJSONBodyOptionsDomainMasks(securityString(tf.DomainMasks.ValueBool()))
	hsts := api.UpdateSecurityOptionsJSONBodyOptionsHsts(securityString(tf.HSTS.ValueBool()))
	cors := api.UpdateSecurityOptionsJSONBodyOptionsCors(securityString(tf.CORS.ValueBool()))
	return api.UpdateSecurityOptionsJSONRequestBody{
		Options: struct {
			Cors           *api.UpdateSecurityOptionsJSONBodyOptionsCors           `json:"cors,omitempty"`
			DomainMasks    *api.UpdateSecurityOptionsJSONBodyOptionsDomainMasks    `json:"domainMasks,omitempty"`
			Hsts           *api.UpdateSecurityOptionsJSONBodyOptionsHsts           `json:"hsts,omitempty"`
			IpCustomHeader *api.UpdateSecurityOptionsJSONBodyOptionsIpCustomHeader `json:"ipCustomHeader,omitempty"`
			Ips            *api.UpdateSecurityOptionsJSONBodyOptionsIps            `json:"ips,omitempty"`
			Jwts           *api.UpdateSecurityOptionsJSONBodyOptionsJwts           `json:"jwts,omitempty"`
			Referrers      *api.UpdateSecurityOptionsJSONBodyOptionsReferrers      `json:"referrers,omitempty"`
			RequestFilters *api.UpdateSecurityOptionsJSONBodyOptionsRequestFilters `json:"requestFilters,omitempty"`
			Tokens         *api.UpdateSecurityOptionsJSONBodyOptionsTokens         `json:"tokens,omitempty"`
		}{
			Tokens:      &tokens,
			Referrers:   &referrers,
			Jwts:        &jwts,
			Ips:         &ips,
			DomainMasks: &domainMasks,
			Hsts:        &hsts,
			Cors:        &cors,
		},
	}
}

// mapSingleEndpointToState maps a SingleEndpoint from the API to the Terraform resource model.
func mapSingleEndpointToState(endpoint *api.SingleEndpoint, body []byte) models.EndpointResourceModel {
	label := ""
	if endpoint.Label != nil {
		label = *endpoint.Label
	}
	wssURL := ""
	if endpoint.WssUrl != nil {
		wssURL = *endpoint.WssUrl
	}
	status := ""
	if endpoint.Status != nil {
		status = *endpoint.Status
	}
	multichain := false
	if endpoint.Multichain != nil {
		multichain = *endpoint.Multichain
	}

	state := models.EndpointResourceModel{
		ID:              types.StringValue(endpoint.Id),
		Label:           types.StringValue(label),
		Chain:           types.StringValue(endpoint.Chain),
		Network:         types.StringValue(endpoint.Network),
		HTTPURL:         types.StringValue(endpoint.HttpUrl),
		WSSURL:          types.StringValue(wssURL),
		SecurityOptions: parseSecurityOptions(body),
		Status:          types.StringValue(status),
		Multichain:      types.BoolValue(multichain),
	}
	return state
}

// Create a new resource.
func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan models.EndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chain := plan.Chain.ValueString()
	network := plan.Network.ValueString()

	// Create new endpoint.
	createResp, err := r.client.API.CreateEndpointWithResponse(ctx, api.CreateEndpointJSONRequestBody{
		Chain:   &chain,
		Network: &network,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint",
			"Could not create endpoint, unexpected error: "+err.Error(),
		)
		return
	}
	if createResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating endpoint",
			fmt.Sprintf("API returned status %d: %s", createResp.StatusCode(), string(createResp.Body)),
		)
		return
	}

	endpoint := createResp.JSON200.Data
	plan.ID = types.StringValue(endpoint.Id)
	plan.Chain = types.StringValue(endpoint.Chain)
	plan.Network = types.StringValue(endpoint.Network)
	plan.HTTPURL = types.StringValue(endpoint.HttpUrl)
	if endpoint.WssUrl != nil {
		plan.WSSURL = types.StringValue(*endpoint.WssUrl)
	}
	if endpoint.Status != nil {
		plan.Status = types.StringValue(*endpoint.Status)
	}
	plan.Multichain = types.BoolValue(*endpoint.Multichain)

	// Patch endpoint label if needed.
	if plan.Label.ValueString() != "" {
		label := plan.Label.ValueString()
		updateResp, err := r.client.API.UpdateEndpointWithResponse(ctx, endpoint.Id, api.UpdateEndpointJSONRequestBody{
			Label: &label,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error patching endpoint label",
				"Could not patch endpoint, unexpected error: "+err.Error(),
			)
			return
		}
		if updateResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Error patching endpoint label",
				fmt.Sprintf("API returned status %d: %s", updateResp.StatusCode(), string(updateResp.Body)),
			)
			return
		}
	} else if endpoint.Label != nil {
		plan.Label = types.StringValue(*endpoint.Label)
	}

	// Patch security options.
	secBody := buildSecurityOptionsBody(plan.SecurityOptions)
	secResp, err := r.client.API.UpdateSecurityOptionsWithResponse(ctx, plan.ID.ValueString(), secBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			"Could not patch endpoint security options, unexpected error: "+err.Error(),
		)
		return
	}
	if secResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			fmt.Sprintf("API returned status %d: %s", secResp.StatusCode(), string(secResp.Body)),
		)
		return
	}

	// Read back the full state including security options.
	showResp, err := r.client.API.ShowEndpointWithResponse(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
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
	plan.SecurityOptions = parseSecurityOptions(showResp.Body)

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *endpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state models.EndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from QuickNode.
	showResp, err := r.client.API.ShowEndpointWithResponse(ctx, state.ID.ValueString())
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
	state = mapSingleEndpointToState(endpoint, showResp.Body)

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan models.EndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Patch endpoint label.
	label := plan.Label.ValueString()
	updateResp, err := r.client.API.UpdateEndpointWithResponse(ctx, plan.ID.ValueString(), api.UpdateEndpointJSONRequestBody{
		Label: &label,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Endpoint",
			"Could not update endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Endpoint",
			fmt.Sprintf("API returned status %d: %s", updateResp.StatusCode(), string(updateResp.Body)),
		)
		return
	}

	// Patch security options.
	secBody := buildSecurityOptionsBody(plan.SecurityOptions)
	secResp, err := r.client.API.UpdateSecurityOptionsWithResponse(ctx, plan.ID.ValueString(), secBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			"Could not patch endpoint security options, unexpected error: "+err.Error(),
		)
		return
	}
	if secResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error patching endpoint security options",
			fmt.Sprintf("API returned status %d: %s", secResp.StatusCode(), string(secResp.Body)),
		)
		return
	}

	// Read back the endpoint to get the full updated state.
	showResp, err := r.client.API.ShowEndpointWithResponse(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Endpoint",
			"Could not read QuickNode endpoint ID "+plan.ID.ValueString()+": "+err.Error(),
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
	plan = mapSingleEndpointToState(endpoint, showResp.Body)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state models.EndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing endpoint.
	deleteResp, err := r.client.API.ArchiveEndpointWithResponse(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint",
			"Could not delete endpoint, unexpected error: "+err.Error(),
		)
		return
	}
	if deleteResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Endpoint",
			fmt.Sprintf("API returned status %d: %s", deleteResp.StatusCode(), string(deleteResp.Body)),
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
	// Retrieve import ID and save to id attribute.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
