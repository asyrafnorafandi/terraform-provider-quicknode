package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type EndpointModel struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Chain   string `json:"chain"`
	Network string `json:"network"`
	HTTPURL string `json:"http_url"`
	WSSURL  string `json:"wss_url"`
}

type EndpointResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Label   types.String `tfsdk:"label"`
	Chain   types.String `tfsdk:"chain"`
	Network types.String `tfsdk:"network"`
	HTTPURL types.String `tfsdk:"http_url"`
	WSSURL  types.String `tfsdk:"wss_url"`
}
