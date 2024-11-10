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
	"terraform-provider-coolify/internal/provider/generated/datasource_team"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &teamDataSource{}
var _ datasource.DataSourceWithConfigure = &teamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

type teamDataSource struct {
	providerData CoolifyProviderData
}

func (d *teamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *teamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_team.TeamDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify team by optional `id`. If no `id` is provided, the team associated with the current API key will be returned."

	// Override the id attribute to make it optional (used to toggle between current team and team by ID)
	if idAttr, ok := resp.Schema.Attributes["id"].(schema.Int64Attribute); ok {
		idAttr.Required = false
		idAttr.Optional = true
		resp.Schema.Attributes["id"] = idAttr
	}

	// Mark sensitive attributes
	sensitiveAttrs := []string{"discord_webhook_url", "smtp_password", "telegram_token", "resend_api_key"}
	for _, attr := range sensitiveAttrs {
		makeDataSourceAttributeSensitive(resp.Schema.Attributes, attr)
	}
}

func (d *teamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_team.TeamModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var team *api.Team
	var teamMembers *[]api.User

	if !plan.Id.IsNull() {
		// Get team by ID
		teamResp, err := d.providerData.client.GetTeamByIdWithResponse(ctx, int(plan.Id.ValueInt64()))
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
		teamResp, err := d.providerData.client.GetCurrentTeamWithResponse(ctx)
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
		teamMembersResponse, err := d.providerData.client.GetMembersByTeamIdWithResponse(ctx, *team.Id)
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
		teamMembers = teamMembersResponse.JSON200
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, team, teamMembers)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *teamDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Team,
	membersResponse *[]api.User,
) datasource_team.TeamModel {
	var elements []attr.Value
	if membersResponse != nil {
		for _, member := range *membersResponse {
			attributes := map[string]attr.Value{
				"created_at":              optionalString(member.CreatedAt),
				"email":                   optionalString(member.Email),
				"email_verified_at":       optionalString(member.EmailVerifiedAt),
				"force_password_reset":    optionalBool(member.ForcePasswordReset),
				"id":                      optionalInt64(member.Id),
				"marketing_emails":        optionalBool(member.MarketingEmails),
				"name":                    optionalString(member.Name),
				"two_factor_confirmed_at": optionalString(member.TwoFactorConfirmedAt),
				"updated_at":              optionalString(member.UpdatedAt),
			}

			data, diag := datasource_team.NewMembersValue(
				datasource_team.MembersValue{}.AttributeTypes(ctx),
				attributes)
			diags.Append(diag...)
			elements = append(elements, data)
		}
	}
	dataList, diag := types.ListValueFrom(ctx, datasource_team.MembersValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return datasource_team.TeamModel{
		CreatedAt:                            optionalString(response.CreatedAt),
		CustomServerLimit:                    optionalString(response.CustomServerLimit),
		Description:                          optionalString(response.Description),
		DiscordEnabled:                       optionalBool(response.DiscordEnabled),
		DiscordNotificationsDatabaseBackups:  optionalBool(response.DiscordNotificationsDatabaseBackups),
		DiscordNotificationsDeployments:      optionalBool(response.DiscordNotificationsDeployments),
		DiscordNotificationsScheduledTasks:   optionalBool(response.DiscordNotificationsScheduledTasks),
		DiscordNotificationsStatusChanges:    optionalBool(response.DiscordNotificationsStatusChanges),
		DiscordNotificationsTest:             optionalBool(response.DiscordNotificationsTest),
		DiscordWebhookUrl:                    optionalString(response.DiscordWebhookUrl),
		Id:                                   optionalInt64(response.Id),
		Members:                              dataList,
		Name:                                 optionalString(response.Name),
		PersonalTeam:                         optionalBool(response.PersonalTeam),
		ResendApiKey:                         optionalString(response.ResendApiKey),
		ResendEnabled:                        optionalBool(response.ResendEnabled),
		ShowBoarding:                         optionalBool(response.ShowBoarding),
		SmtpEnabled:                          optionalBool(response.SmtpEnabled),
		SmtpEncryption:                       optionalString(response.SmtpEncryption),
		SmtpFromAddress:                      optionalString(response.SmtpFromAddress),
		SmtpFromName:                         optionalString(response.SmtpFromName),
		SmtpHost:                             optionalString(response.SmtpHost),
		SmtpNotificationsDatabaseBackups:     optionalBool(response.SmtpNotificationsDatabaseBackups),
		SmtpNotificationsDeployments:         optionalBool(response.SmtpNotificationsDeployments),
		SmtpNotificationsScheduledTasks:      optionalBool(response.SmtpNotificationsScheduledTasks),
		SmtpNotificationsStatusChanges:       optionalBool(response.SmtpNotificationsStatusChanges),
		SmtpNotificationsTest:                optionalBool(response.SmtpNotificationsTest),
		SmtpPassword:                         optionalString(response.SmtpPassword),
		SmtpPort:                             optionalString(response.SmtpPort),
		SmtpRecipients:                       optionalString(response.SmtpRecipients),
		SmtpTimeout:                          optionalString(response.SmtpTimeout),
		SmtpUsername:                         optionalString(response.SmtpUsername),
		TelegramChatId:                       optionalString(response.TelegramChatId),
		TelegramEnabled:                      optionalBool(response.TelegramEnabled),
		TelegramNotificationsDatabaseBackups: optionalBool(response.TelegramNotificationsDatabaseBackups),
		TelegramNotificationsDatabaseBackupsMessageThreadId: optionalString(response.TelegramNotificationsDatabaseBackupsMessageThreadId),
		TelegramNotificationsDeployments:                    optionalBool(response.TelegramNotificationsDeployments),
		TelegramNotificationsDeploymentsMessageThreadId:     optionalString(response.TelegramNotificationsDeploymentsMessageThreadId),
		TelegramNotificationsScheduledTasks:                 optionalBool(response.TelegramNotificationsScheduledTasks),
		TelegramNotificationsScheduledTasksThreadId:         optionalString(response.TelegramNotificationsScheduledTasksThreadId),
		TelegramNotificationsStatusChanges:                  optionalBool(response.TelegramNotificationsStatusChanges),
		TelegramNotificationsStatusChangesMessageThreadId:   optionalString(response.TelegramNotificationsStatusChangesMessageThreadId),
		TelegramNotificationsTest:                           optionalBool(response.TelegramNotificationsTest),
		TelegramNotificationsTestMessageThreadId:            optionalString(response.TelegramNotificationsTestMessageThreadId),
		TelegramToken:                                       optionalString(response.TelegramToken),
		UpdatedAt:                                           optionalString(response.UpdatedAt),
		UseInstanceEmailSettings:                            optionalBool(response.UseInstanceEmailSettings),
	}
}
