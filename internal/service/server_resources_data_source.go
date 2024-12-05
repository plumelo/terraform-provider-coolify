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
	"terraform-provider-coolify/internal/provider/generated/datasource_server_resources"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serverResourcesDataSource{}
var _ datasource.DataSourceWithConfigure = &serverResourcesDataSource{}

func NewServerResourcesDataSource() datasource.DataSource {
	return &serverResourcesDataSource{}
}

type serverResourcesDataSource struct {
	client *api.ClientWithResponses
}

func (d *serverResourcesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_resources"
}

func (d *serverResourcesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_server_resources.ServerResourcesDataSourceSchema(ctx)
	resp.Schema.Description = "Get resources for a server by `uuid`."
}

func (d *serverResourcesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *serverResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_server_resources.ServerResourcesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetResourcesByServerUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server resources", err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server resources",
			fmt.Sprintf("Received %s for server resources. Details: %s", response.Status(), string(response.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, response.JSON200)
	state.Uuid = plan.Uuid

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serverResourcesDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]struct {
		CreatedAt *string `json:"created_at,omitempty"`
		Id        *int    `json:"id,omitempty"`
		Name      *string `json:"name,omitempty"`
		Status    *string `json:"status,omitempty"`
		Type      *string `json:"type,omitempty"`
		UpdatedAt *string `json:"updated_at,omitempty"`
		Uuid      *string `json:"uuid,omitempty"`
	}, // Yuck. Codegen did not produce a model var for this response.
) datasource_server_resources.ServerResourcesModel {
	var elements []attr.Value

	for _, svRes := range *response {
		attributes := map[string]attr.Value{
			"created_at": flatten.String(svRes.CreatedAt),
			"id":         flatten.Int64(svRes.Id),
			"name":       flatten.String(svRes.Name),
			"status":     flatten.String(svRes.Status),
			"type":       flatten.String(svRes.Type),
			"updated_at": flatten.String(svRes.UpdatedAt),
			"uuid":       flatten.String(svRes.Uuid),
		}

		// todo: add `server_resources` filtering

		data, diag := datasource_server_resources.NewServerResourcesValue(
			datasource_server_resources.ServerResourcesValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_server_resources.ServerResourcesValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return datasource_server_resources.ServerResourcesModel{
		ServerResources: dataSet,
	}
}
