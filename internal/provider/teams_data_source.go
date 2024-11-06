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
	"terraform-provider-coolify/internal/provider/generated/datasource_teams"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &teamsDataSource{}
var _ datasource.DataSourceWithConfigure = &teamsDataSource{}

func NewTeamsDataSource() datasource.DataSource {
	return &teamsDataSource{}
}

type teamsDataSource struct {
	providerData CoolifyProviderData
}

type teamsDataSourceWithFilterModel struct {
	datasource_teams.TeamsModel
	Filter      []filterBlockModel `tfsdk:"filter"`
	WithMembers types.Bool         `tfsdk:"with_members"`
}

var teamsFilterNames = []string{"name", "description", "id", "discord_enabled", "resend_enabled", "smtp_enabled", "telegram_enabled"}

func (d *teamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *teamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_teams.TeamsDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify teams."

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": createDatasourceFilter(teamsFilterNames),
	}

	// Add with_members attribute
	resp.Schema.Attributes["with_members"] = schema.BoolAttribute{
		Optional:    true,
		Description: "Whether to fetch team members. This requires an additional API call per team.",
	}

	// Make sensitive fields sensitive
	if teamsSet, ok := resp.Schema.Attributes["teams"].(schema.SetNestedAttribute); ok {
		sensitiveFields := []string{"discord_webhook_url", "smtp_password", "telegram_token", "resend_api_key"}
		for _, field := range sensitiveFields {
			if attr, ok := teamsSet.NestedObject.Attributes[field].(schema.StringAttribute); ok {
				attr.Sensitive = true
				teamsSet.NestedObject.Attributes[field] = attr
			}
		}
		resp.Schema.Attributes["teams"] = teamsSet
	}
}

func (d *teamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *teamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan teamsDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.providerData.client.ListTeamsWithResponse(ctx)
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
) (teamsDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var elements []attr.Value

	for _, team := range *teams {
		var members []attr.Value
		if team.Members != nil {
			members = convertMembersToAttrValues(ctx, *team.Members)
		} else if withMembers.ValueBool() {
			// Fetch members separately if requested and not included in team response
			teamMembersResponse, err := d.providerData.client.GetMembersByTeamIdWithResponse(ctx, *team.Id)
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
				members = convertMembersToAttrValues(ctx, *teamMembersResponse.JSON200)
			}
		}

		attributes := map[string]attr.Value{
			"created_at":                             optionalString(team.CreatedAt),
			"custom_server_limit":                    optionalString(team.CustomServerLimit),
			"description":                            optionalString(team.Description),
			"discord_enabled":                        optionalBool(team.DiscordEnabled),
			"discord_notifications_database_backups": optionalBool(team.DiscordNotificationsDatabaseBackups),
			"discord_notifications_deployments":      optionalBool(team.DiscordNotificationsDeployments),
			"discord_notifications_scheduled_tasks":  optionalBool(team.DiscordNotificationsScheduledTasks),
			"discord_notifications_status_changes":   optionalBool(team.DiscordNotificationsStatusChanges),
			"discord_notifications_test":             optionalBool(team.DiscordNotificationsTest),
			"discord_webhook_url":                    optionalString(team.DiscordWebhookUrl),
			"id":                                     optionalInt64(team.Id),
			"members": types.ListValueMust(
				datasource_teams.MembersValue{}.Type(ctx),
				members,
			),
			"name":                                    optionalString(team.Name),
			"personal_team":                           optionalBool(team.PersonalTeam),
			"resend_api_key":                          optionalString(team.ResendApiKey),
			"resend_enabled":                          optionalBool(team.ResendEnabled),
			"show_boarding":                           optionalBool(team.ShowBoarding),
			"smtp_enabled":                            optionalBool(team.SmtpEnabled),
			"smtp_encryption":                         optionalString(team.SmtpEncryption),
			"smtp_from_address":                       optionalString(team.SmtpFromAddress),
			"smtp_from_name":                          optionalString(team.SmtpFromName),
			"smtp_host":                               optionalString(team.SmtpHost),
			"smtp_notifications_database_backups":     optionalBool(team.SmtpNotificationsDatabaseBackups),
			"smtp_notifications_deployments":          optionalBool(team.SmtpNotificationsDeployments),
			"smtp_notifications_scheduled_tasks":      optionalBool(team.SmtpNotificationsScheduledTasks),
			"smtp_notifications_status_changes":       optionalBool(team.SmtpNotificationsStatusChanges),
			"smtp_notifications_test":                 optionalBool(team.SmtpNotificationsTest),
			"smtp_password":                           optionalString(team.SmtpPassword),
			"smtp_port":                               optionalString(team.SmtpPort),
			"smtp_recipients":                         optionalString(team.SmtpRecipients),
			"smtp_timeout":                            optionalString(team.SmtpTimeout),
			"smtp_username":                           optionalString(team.SmtpUsername),
			"telegram_chat_id":                        optionalString(team.TelegramChatId),
			"telegram_enabled":                        optionalBool(team.TelegramEnabled),
			"telegram_notifications_database_backups": optionalBool(team.TelegramNotificationsDatabaseBackups),
			"telegram_notifications_database_backups_message_thread_id": optionalString(team.TelegramNotificationsDatabaseBackupsMessageThreadId),
			"telegram_notifications_deployments":                        optionalBool(team.TelegramNotificationsDeployments),
			"telegram_notifications_deployments_message_thread_id":      optionalString(team.TelegramNotificationsDeploymentsMessageThreadId),
			"telegram_notifications_scheduled_tasks":                    optionalBool(team.TelegramNotificationsScheduledTasks),
			"telegram_notifications_scheduled_tasks_thread_id":          optionalString(team.TelegramNotificationsScheduledTasksThreadId),
			"telegram_notifications_status_changes":                     optionalBool(team.TelegramNotificationsStatusChanges),
			"telegram_notifications_status_changes_message_thread_id":   optionalString(team.TelegramNotificationsStatusChangesMessageThreadId),
			"telegram_notifications_test":                               optionalBool(team.TelegramNotificationsTest),
			"telegram_notifications_test_message_thread_id":             optionalString(team.TelegramNotificationsTestMessageThreadId),
			"telegram_token":              optionalString(team.TelegramToken),
			"updated_at":                  optionalString(team.UpdatedAt),
			"use_instance_email_settings": optionalBool(team.UseInstanceEmailSettings),
		}

		if !filterOnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_teams.NewTeamsValue(
			datasource_teams.TeamsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_teams.TeamsValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return teamsDataSourceWithFilterModel{
		TeamsModel: datasource_teams.TeamsModel{
			Teams: dataSet,
		},
		Filter:      filters,
		WithMembers: withMembers,
	}, diags
}

func convertMembersToAttrValues(ctx context.Context, members []api.User) []attr.Value {
	var values []attr.Value
	for _, member := range members {
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
		data, _ := datasource_teams.NewMembersValue(
			datasource_teams.MembersValue{}.AttributeTypes(ctx),
			attributes)
		values = append(values, data)
	}
	return values
}
