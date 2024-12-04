package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &teamsDataSource{}
var _ datasource.DataSourceWithConfigure = &teamsDataSource{}

type teamsDataSourceModel struct {
	Teams       []teamDataSourceModel `tfsdk:"teams"`
	WithMembers types.Bool            `tfsdk:"with_members"`
	Filter      []filterBlockModel    `tfsdk:"filter"`
}

func NewTeamsDataSource() datasource.DataSource {
	return &teamsDataSource{}
}

type teamsDataSource struct {
	client *api.ClientWithResponses
}

func (d *teamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *teamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ds := NewTeamDataSource()
	dsResp := datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, &dsResp)

	if attr, ok := dsResp.Schema.Attributes["id"].(schema.Int64Attribute); ok {
		attr.Computed = true
		attr.Optional = false
		attr.Required = false
		dsResp.Schema.Attributes["id"] = attr
	}

	resp.Schema = schema.Schema{
		Description: "Get a list of Coolify teams.",
		Attributes: map[string]schema.Attribute{
			"teams": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: dsResp.Schema.Attributes,
				},
			},
			"with_members": schema.BoolAttribute{
				Description: "Whether to fetch team members. This requires an additional API call per team.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": createDatasourceFilter(teamsFilterNames),
		},
	}
}

func (d *teamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *teamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan teamsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.client.ListTeamsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading teams", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading teams",
			fmt.Sprintf("Received %s for teams. Details: %s", listResponse.Status(), listResponse.Body),
		)
		return
	}

	state, diag := d.apiToModel(ctx, listResponse.JSON200, plan.Filter, plan.WithMembers)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *teamsDataSource) apiToModel(
	ctx context.Context,
	teams *[]api.Team,
	filters []filterBlockModel,
	withMembers types.Bool,
) (teamsDataSourceModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	var teamValues []teamDataSourceModel

	for _, team := range *teams {
		if team.Members == nil && withMembers.ValueBool() {
			// Fetch members separately if requested and not included in team response
			teamMembersResponse, err := d.client.GetMembersByTeamIdWithResponse(ctx, *team.Id)
			if err != nil {
				diags.AddError(
					"Error reading team members",
					err.Error(),
				)
				continue
			}

			if teamMembersResponse.StatusCode() != http.StatusOK {
				diags.AddError(
					"Unexpected HTTP status code reading team members",
					fmt.Sprintf("Received %s for team members. Details: %s", teamMembersResponse.Status(), string(teamMembersResponse.Body)),
				)
				continue
			}

			if teamMembersResponse.JSON200 != nil {
				team.Members = teamMembersResponse.JSON200
			}
		}

		model := teamDataSourceModel{}.FromAPI(&team)

		if !filterOnStruct(ctx, model, filters) {
			continue
		}

		teamValues = append(teamValues, model)
	}

	return teamsDataSourceModel{
		Teams:       teamValues,
		WithMembers: withMembers,
		Filter:      filters,
	}, diags
}
