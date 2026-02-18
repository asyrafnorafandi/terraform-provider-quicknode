package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// JSON Models.
type EndpointSecurityOptionsModel struct {
	Tokens      bool `json:"tokens"`
	Referrers   bool `json:"referrers"`
	JWTs        bool `json:"jwts"`
	IPs         bool `json:"ips"`
	DomainMasks bool `json:"domainMasks"`
	HSTS        bool `json:"hsts"`
	CORS        bool `json:"cors"`
}

type EndpointSecurityTokensModel struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type EndpointSecurityModel struct {
	Options EndpointSecurityOptionsModel  `json:"options"`
	Tokens  []EndpointSecurityTokensModel `json:"tokens,omitempty"`
}

type EndpointRateLimitsModel struct {
	RateLimitByIP bool  `json:"rate_limit_by_ip"`
	Account       int64 `json:"account"`
	RPS           int64 `json:"rps"`
	RPD           int64 `json:"rpd"`
	RPM           int64 `json:"rpm"`
}

type EndpointModel struct {
	ID         string                  `json:"id"`
	Label      string                  `json:"label"`
	Chain      string                  `json:"chain"`
	Network    string                  `json:"network"`
	HTTPURL    string                  `json:"http_url"`
	WSSURL     string                  `json:"wss_url"`
	Security   EndpointSecurityModel   `json:"security"`
	Status     string                  `json:"status"`
	RateLimits EndpointRateLimitsModel `json:"rate_limits"`
	Tags       []string                `json:"tags"`
	Multichain bool                    `json:"multichain"`
}

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
	// Tags            types.List                    `tfsdk:"tags"`
	Multichain types.Bool `tfsdk:"multichain"`
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
