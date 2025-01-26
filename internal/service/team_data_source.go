package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &teamDataSource{}
var _ datasource.DataSourceWithConfigure = &teamDataSource{}

type teamDataSourceModel = teamModel

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

type teamDataSource struct {
	client *api.ClientWithResponses
}

func (d *teamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *teamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a Coolify team by optional `id`. If no `id` is provided, the team associated with the current API key will be returned.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the team was created.",
			},
			"custom_server_limit": schema.StringAttribute{
				Computed:    true,
				Description: "The custom server limit.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the team.",
			},
			"id": schema.Int64Attribute{
				Required:    false,
				Optional:    true,
				Description: "The unique identifier of the team.",
			},
			"members": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user was created.",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Sensitive:   false, // todo: mark sensitive?
							Description: "The user email.",
						},
						"email_verified_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user email was verified.",
						},
						"force_password_reset": schema.BoolAttribute{
							Computed:    true,
							Description: "The flag to force the user to reset the password.",
						},
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The user identifier in the database.",
						},
						"marketing_emails": schema.BoolAttribute{
							Computed:    true,
							Description: "The flag to receive marketing emails.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The user name.",
						},
						"two_factor_confirmed_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user two factor was confirmed.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user was updated.",
						},
					},
				},
				Computed:    true,
				Description: "The members of the team.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the team.",
			},
			"personal_team": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the team is personal or not.",
			},

			"show_boarding": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to show the boarding screen or not.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the team was last updated.",
			},
		},
	}
}

func (d *teamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan teamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var team *api.Team
	var teamMembers *[]api.User

	if !plan.Id.IsNull() {
		// Get team by ID
		teamResp, err := d.client.GetTeamByIdWithResponse(ctx, int(plan.Id.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team", err.Error(),
			)
			return
		}

		if teamResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team",
				fmt.Sprintf("Received %s for team. Details: %s", teamResp.Status(), string(teamResp.Body)),
			)
			return
		}

		team = teamResp.JSON200
		teamMembers = team.Members
	} else {
		// Get current team
		teamResp, err := d.client.GetCurrentTeamWithResponse(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team", err.Error(),
			)
			return
		}

		if teamResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team",
				fmt.Sprintf("Received %s for team. Details: %s", teamResp.Status(), string(teamResp.Body)),
			)
			return
		}

		team = teamResp.JSON200
		teamMembers = team.Members
	}

	// If the API did not provide members, we need to fetch them separately
	// TODO: Coolify API inconsistency: Spec says members should be returned, but it is not.
	if teamMembers == nil {
		teamMembersResponse, err := d.client.GetMembersByTeamIdWithResponse(ctx, *team.Id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team members", err.Error(),
			)
			return
		}

		if teamMembersResponse.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team members",
				fmt.Sprintf("Received %s for team members. Details: %s", teamMembersResponse.Status(), string(teamMembersResponse.Body)),
			)
			return
		}
		team.Members = teamMembersResponse.JSON200
	}

	state := teamDataSourceModel{}.FromAPI(team)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
