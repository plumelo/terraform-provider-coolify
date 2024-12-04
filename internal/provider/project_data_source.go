package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_project"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &projectDataSource{}
var _ datasource.DataSourceWithConfigure = &projectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	providerData CoolifyProviderData
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_project.ProjectDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify project by `uuid`."
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_project.ProjectModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.providerData.Client.GetProjectByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project", err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading project",
			fmt.Sprintf("Received %s for project. Details: %s", response.Status(), string(response.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, response.JSON200)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *projectDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Project,
) datasource_project.ProjectModel {
	var elements []attr.Value
	for _, env := range *response.Environments {
		attributes := map[string]attr.Value{
			"created_at":  optionalString(env.CreatedAt),
			"description": optionalString(env.Description),
			"id":          optionalInt64(env.Id),
			"name":        optionalString(env.Name),
			"project_id":  optionalInt64(env.ProjectId),
			"updated_at":  optionalString(env.UpdatedAt),
		}

		data, diag := datasource_project.NewEnvironmentsValue(
			datasource_project.EnvironmentsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}
	dataList, diag := types.ListValueFrom(ctx, datasource_project.EnvironmentsValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return datasource_project.ProjectModel{
		Description:  optionalString(response.Description),
		Environments: dataList,
		Id:           optionalInt64(response.Id),
		Name:         optionalString(response.Name),
		Uuid:         optionalString(response.Uuid),
	}
}
