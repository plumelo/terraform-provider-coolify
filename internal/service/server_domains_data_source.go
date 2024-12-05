package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/datasource_server_domains"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serverDomainsDataSource{}
var _ datasource.DataSourceWithConfigure = &serverDomainsDataSource{}

func NewServerDomainsDataSource() datasource.DataSource {
	return &serverDomainsDataSource{}
}

type serverDomainsDataSource struct {
	client *api.ClientWithResponses
}

func (d *serverDomainsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_domains"
}

func (d *serverDomainsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_server_domains.ServerDomainsDataSourceSchema(ctx)
	resp.Schema.Description = "Get domains for a server by `uuid`."
}

func (d *serverDomainsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *serverDomainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_server_domains.ServerDomainsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetDomainsByServerUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server domains", err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server domains",
			fmt.Sprintf("Received %s for server domains. Details: %s", response.Status(), string(response.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, response.JSON200)
	state.Uuid = plan.Uuid

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serverDomainsDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]struct {
		Domains *[]string `json:"domains,omitempty"`
		Ip      *string   `json:"ip,omitempty"`
	}, // Yuck. Codegen did not produce a model var for this response.
) datasource_server_domains.ServerDomainsModel {
	var elements []attr.Value

	for _, svDomain := range *response {
		attributes := map[string]attr.Value{
			"domains": flatten.StringList(svDomain.Domains),
			"ip":      flatten.String(svDomain.Ip),
		}

		// todo: add `server_domains` filtering

		data, diag := datasource_server_domains.NewServerDomainsValue(
			datasource_server_domains.ServerDomainsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_server_domains.ServerDomainsValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return datasource_server_domains.ServerDomainsModel{
		ServerDomains: dataSet,
	}
}
