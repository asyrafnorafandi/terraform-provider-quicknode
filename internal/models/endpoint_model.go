package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// Terraform Models.
type SecurityOptionsResourceModel struct {
	Tokens      types.Bool `tfsdk:"tokens"`
	Referrers   types.Bool `tfsdk:"referrers"`
	JWTs        types.Bool `tfsdk:"jwts"`
	IPs         types.Bool `tfsdk:"ips"`
	DomainMasks types.Bool `tfsdk:"domain_masks"`
	HSTS        types.Bool `tfsdk:"hsts"`
	CORS        types.Bool `tfsdk:"cors"`
}

type EndpointResourceModel struct {
	ID              types.String                  `tfsdk:"id"`
	Label           types.String                  `tfsdk:"label"`
	Chain           types.String                  `tfsdk:"chain"`
	Network         types.String                  `tfsdk:"network"`
	HTTPURL         types.String                  `tfsdk:"http_url"`
	WSSURL          types.String                  `tfsdk:"wss_url"`
	SecurityOptions *SecurityOptionsResourceModel `tfsdk:"security_options"`
	Status          types.String                  `tfsdk:"status"`
	Tags            types.List                    `tfsdk:"tags"` // element type: types.StringType
	Multichain      types.Bool                    `tfsdk:"multichain"`
}

type EndpointsDataSourceModel struct {
	Limit     types.Int64 `tfsdk:"limit"`
	Offset    types.Int64 `tfsdk:"offset"`
	Endpoints []struct {
		ID      types.String `tfsdk:"id"`
		Label   types.String `tfsdk:"label"`
		Chain   types.String `tfsdk:"chain"`
		Network types.String `tfsdk:"network"`
		HTTPURL types.String `tfsdk:"http_url"`
		WSSURL  types.String `tfsdk:"wss_url"`
	} `tfsdk:"endpoints"`
}

type EndpointWhitelistIPResourceModel struct {
	ID         types.String `tfsdk:"id"`
	IP         types.String `tfsdk:"ip"`
	EndpointID types.String `tfsdk:"endpoint_id"`
}
