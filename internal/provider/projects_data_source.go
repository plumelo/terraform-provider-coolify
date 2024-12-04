package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_projects"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &projectsDataSource{}
var _ datasource.DataSourceWithConfigure = &projectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type projectsDataSource struct {
	providerData CoolifyProviderData
}

type projectsDataSourceWithFilterModel struct {
	datasource_projects.ProjectsModel
	Filter []filterBlockModel `tfsdk:"filter"`
}

var projectsFilterNames = []string{"id", "uuid", "name", "description"}

func (d *projectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_projects.ProjectsDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify projects." +
		"\nNOTE: Environments are not returned due to an API bug. Combine with `coolify_project` as a temporary workaround."

	// todo: Coolify API bug, environments are not returned
	if projectsSet, ok := resp.Schema.Attributes["projects"].(schema.SetNestedAttribute); ok {
		if envAttr, ok := projectsSet.NestedObject.Attributes["environments"].(schema.ListNestedAttribute); ok {
			envAttr.DeprecationMessage = "The environments field is currently not functional due to an API bug. Use coolify_project data source instead."
			envAttr.Description = "This field is currently not populated due to a Coolify API bug."
			envAttr.MarkdownDescription = "*" + envAttr.Description + "*"
			envAttr.NestedObject.Attributes = map[string]schema.Attribute{}

			projectsSet.NestedObject.Attributes["environments"] = envAttr
			resp.Schema.Attributes["projects"] = projectsSet
		}
	}

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": createDatasourceFilter(projectsFilterNames),
	}
}

func (d *projectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan projectsDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.providerData.Client.ListProjectsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading projects", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading projects",
			fmt.Sprintf("Received %s for projects. Details: %s", listResponse.Status(), listResponse.Body),
		)
		return
	}

	state, diag := d.apiToModel(ctx, listResponse.JSON200, plan.Filter)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *projectsDataSource) apiToModel(
	ctx context.Context,
	response *[]api.Project,
	filters []filterBlockModel,
) (projectsDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var projects []attr.Value

	for _, project := range *response {
		var envs []attr.Value

		// todo: Coolify API bug, environments are not returned
		if project.Environments != nil {
			for _, env := range *project.Environments {
				attributes := map[string]attr.Value{
					"created_at":  optionalString(env.CreatedAt),
					"description": optionalString(env.Description),
					"id":          optionalInt64(env.Id),
					"name":        optionalString(env.Name),
					"project_id":  optionalInt64(env.ProjectId),
					"updated_at":  optionalString(env.UpdatedAt),
				}

				data, diag := datasource_projects.NewEnvironmentsValue(
					datasource_projects.EnvironmentsValue{}.AttributeTypes(ctx),
					attributes)
				diags.Append(diag...)
				envs = append(envs, data)
			}
		}

		envsList, diag := types.ListValueFrom(ctx, datasource_projects.EnvironmentsValue{}.Type(ctx), envs)
		diags.Append(diag...)

		attributes := map[string]attr.Value{
			"description":  optionalString(project.Description),
			"environments": envsList,
			"id":           optionalInt64(project.Id),
			"name":         optionalString(project.Name),
			"uuid":         optionalString(project.Uuid),
		}

		if !filterOnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_projects.NewProjectsValue(
			datasource_projects.ProjectsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		projects = append(projects, data)
	}

	dataSet, diag := types.SetValue(datasource_projects.ProjectsValue{}.Type(ctx), projects)
	diags.Append(diag...)

	return projectsDataSourceWithFilterModel{
		ProjectsModel: datasource_projects.ProjectsModel{
			Projects: dataSet,
		},
		Filter: filters,
	}, diags
}
