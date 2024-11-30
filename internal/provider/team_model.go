package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
)

type teamModel struct {
	CreatedAt                                           types.String      `tfsdk:"created_at"`
	CustomServerLimit                                   types.String      `tfsdk:"custom_server_limit"`
	Description                                         types.String      `tfsdk:"description"`
	DiscordEnabled                                      types.Bool        `tfsdk:"discord_enabled"`
	DiscordNotificationsDatabaseBackups                 types.Bool        `tfsdk:"discord_notifications_database_backups"`
	DiscordNotificationsDeployments                     types.Bool        `tfsdk:"discord_notifications_deployments"`
	DiscordNotificationsScheduledTasks                  types.Bool        `tfsdk:"discord_notifications_scheduled_tasks"`
	DiscordNotificationsServerDiskUsage                 types.Bool        `tfsdk:"discord_notifications_server_disk_usage"`
	DiscordNotificationsStatusChanges                   types.Bool        `tfsdk:"discord_notifications_status_changes"`
	DiscordNotificationsTest                            types.Bool        `tfsdk:"discord_notifications_test"`
	DiscordWebhookUrl                                   types.String      `tfsdk:"discord_webhook_url"`
	Id                                                  types.Int64       `tfsdk:"id"`
	Members                                             []teamMemberModel `tfsdk:"members"`
	Name                                                types.String      `tfsdk:"name"`
	PersonalTeam                                        types.Bool        `tfsdk:"personal_team"`
	ResendApiKey                                        types.String      `tfsdk:"resend_api_key"`
	ResendEnabled                                       types.Bool        `tfsdk:"resend_enabled"`
	ShowBoarding                                        types.Bool        `tfsdk:"show_boarding"`
	SmtpEnabled                                         types.Bool        `tfsdk:"smtp_enabled"`
	SmtpEncryption                                      types.String      `tfsdk:"smtp_encryption"`
	SmtpFromAddress                                     types.String      `tfsdk:"smtp_from_address"`
	SmtpFromName                                        types.String      `tfsdk:"smtp_from_name"`
	SmtpHost                                            types.String      `tfsdk:"smtp_host"`
	SmtpNotificationsDatabaseBackups                    types.Bool        `tfsdk:"smtp_notifications_database_backups"`
	SmtpNotificationsDeployments                        types.Bool        `tfsdk:"smtp_notifications_deployments"`
	SmtpNotificationsScheduledTasks                     types.Bool        `tfsdk:"smtp_notifications_scheduled_tasks"`
	SmtpNotificationsServerDiskUsage                    types.Bool        `tfsdk:"smtp_notifications_server_disk_usage"`
	SmtpNotificationsStatusChanges                      types.Bool        `tfsdk:"smtp_notifications_status_changes"`
	SmtpNotificationsTest                               types.Bool        `tfsdk:"smtp_notifications_test"`
	SmtpPassword                                        types.String      `tfsdk:"smtp_password"`
	SmtpPort                                            types.String      `tfsdk:"smtp_port"`
	SmtpRecipients                                      types.String      `tfsdk:"smtp_recipients"`
	SmtpTimeout                                         types.String      `tfsdk:"smtp_timeout"`
	SmtpUsername                                        types.String      `tfsdk:"smtp_username"`
	TelegramChatId                                      types.String      `tfsdk:"telegram_chat_id"`
	TelegramEnabled                                     types.Bool        `tfsdk:"telegram_enabled"`
	TelegramNotificationsDatabaseBackups                types.Bool        `tfsdk:"telegram_notifications_database_backups"`
	TelegramNotificationsDatabaseBackupsMessageThreadId types.String      `tfsdk:"telegram_notifications_database_backups_message_thread_id"`
	TelegramNotificationsDeployments                    types.Bool        `tfsdk:"telegram_notifications_deployments"`
	TelegramNotificationsDeploymentsMessageThreadId     types.String      `tfsdk:"telegram_notifications_deployments_message_thread_id"`
	TelegramNotificationsScheduledTasks                 types.Bool        `tfsdk:"telegram_notifications_scheduled_tasks"`
	TelegramNotificationsScheduledTasksThreadId         types.String      `tfsdk:"telegram_notifications_scheduled_tasks_thread_id"`
	TelegramNotificationsStatusChanges                  types.Bool        `tfsdk:"telegram_notifications_status_changes"`
	TelegramNotificationsStatusChangesMessageThreadId   types.String      `tfsdk:"telegram_notifications_status_changes_message_thread_id"`
	TelegramNotificationsTest                           types.Bool        `tfsdk:"telegram_notifications_test"`
	TelegramNotificationsTestMessageThreadId            types.String      `tfsdk:"telegram_notifications_test_message_thread_id"`
	TelegramToken                                       types.String      `tfsdk:"telegram_token"`
	UpdatedAt                                           types.String      `tfsdk:"updated_at"`
	UseInstanceEmailSettings                            types.Bool        `tfsdk:"use_instance_email_settings"`
}

type teamMemberModel struct {
	CreatedAt            types.String `tfsdk:"created_at"`
	Email                types.String `tfsdk:"email"`
	EmailVerifiedAt      types.String `tfsdk:"email_verified_at"`
	ForcePasswordReset   types.Bool   `tfsdk:"force_password_reset"`
	Id                   types.Int64  `tfsdk:"id"`
	MarketingEmails      types.Bool   `tfsdk:"marketing_emails"`
	Name                 types.String `tfsdk:"name"`
	TwoFactorConfirmedAt types.String `tfsdk:"two_factor_confirmed_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

var _ filterableStructModel = teamModel{}

func (m teamModel) FromAPI(apiModel *api.Team) teamModel {
	var members []teamMemberModel

	if apiModel.Members != nil {
		members = make([]teamMemberModel, len(*apiModel.Members))
		for i, member := range *apiModel.Members {
			members[i] = teamMemberModel{
				CreatedAt:            optionalString(member.CreatedAt),
				Email:                optionalString(member.Email),
				EmailVerifiedAt:      optionalString(member.EmailVerifiedAt),
				ForcePasswordReset:   optionalBool(member.ForcePasswordReset),
				Id:                   optionalInt64(member.Id),
				MarketingEmails:      optionalBool(member.MarketingEmails),
				Name:                 optionalString(member.Name),
				TwoFactorConfirmedAt: optionalString(member.TwoFactorConfirmedAt),
				UpdatedAt:            optionalString(member.UpdatedAt),
			}
		}
	}

	return teamModel{
		CreatedAt:                            optionalString(apiModel.CreatedAt),
		CustomServerLimit:                    optionalString(apiModel.CustomServerLimit),
		Description:                          optionalString(apiModel.Description),
		DiscordEnabled:                       optionalBool(apiModel.DiscordEnabled),
		DiscordNotificationsDatabaseBackups:  optionalBool(apiModel.DiscordNotificationsDatabaseBackups),
		DiscordNotificationsDeployments:      optionalBool(apiModel.DiscordNotificationsDeployments),
		DiscordNotificationsScheduledTasks:   optionalBool(apiModel.DiscordNotificationsScheduledTasks),
		DiscordNotificationsServerDiskUsage:  optionalBool(apiModel.DiscordNotificationsServerDiskUsage),
		DiscordNotificationsStatusChanges:    optionalBool(apiModel.DiscordNotificationsStatusChanges),
		DiscordNotificationsTest:             optionalBool(apiModel.DiscordNotificationsTest),
		DiscordWebhookUrl:                    optionalString(apiModel.DiscordWebhookUrl),
		Id:                                   optionalInt64(apiModel.Id),
		Members:                              members,
		Name:                                 optionalString(apiModel.Name),
		PersonalTeam:                         optionalBool(apiModel.PersonalTeam),
		ResendApiKey:                         optionalString(apiModel.ResendApiKey),
		ResendEnabled:                        optionalBool(apiModel.ResendEnabled),
		ShowBoarding:                         optionalBool(apiModel.ShowBoarding),
		SmtpEnabled:                          optionalBool(apiModel.SmtpEnabled),
		SmtpEncryption:                       optionalString(apiModel.SmtpEncryption),
		SmtpFromAddress:                      optionalString(apiModel.SmtpFromAddress),
		SmtpFromName:                         optionalString(apiModel.SmtpFromName),
		SmtpHost:                             optionalString(apiModel.SmtpHost),
		SmtpNotificationsDatabaseBackups:     optionalBool(apiModel.SmtpNotificationsDatabaseBackups),
		SmtpNotificationsDeployments:         optionalBool(apiModel.SmtpNotificationsDeployments),
		SmtpNotificationsScheduledTasks:      optionalBool(apiModel.SmtpNotificationsScheduledTasks),
		SmtpNotificationsServerDiskUsage:     optionalBool(apiModel.SmtpNotificationsServerDiskUsage),
		SmtpNotificationsStatusChanges:       optionalBool(apiModel.SmtpNotificationsStatusChanges),
		SmtpNotificationsTest:                optionalBool(apiModel.SmtpNotificationsTest),
		SmtpPassword:                         optionalString(apiModel.SmtpPassword),
		SmtpPort:                             optionalString(apiModel.SmtpPort),
		SmtpRecipients:                       optionalString(apiModel.SmtpRecipients),
		SmtpTimeout:                          optionalString(apiModel.SmtpTimeout),
		SmtpUsername:                         optionalString(apiModel.SmtpUsername),
		TelegramChatId:                       optionalString(apiModel.TelegramChatId),
		TelegramEnabled:                      optionalBool(apiModel.TelegramEnabled),
		TelegramNotificationsDatabaseBackups: optionalBool(apiModel.TelegramNotificationsDatabaseBackups),
		TelegramNotificationsDatabaseBackupsMessageThreadId: optionalString(apiModel.TelegramNotificationsDatabaseBackupsMessageThreadId),
		TelegramNotificationsDeployments:                    optionalBool(apiModel.TelegramNotificationsDeployments),
		TelegramNotificationsDeploymentsMessageThreadId:     optionalString(apiModel.TelegramNotificationsDeploymentsMessageThreadId),
		TelegramNotificationsScheduledTasks:                 optionalBool(apiModel.TelegramNotificationsScheduledTasks),
		TelegramNotificationsScheduledTasksThreadId:         optionalString(apiModel.TelegramNotificationsScheduledTasksThreadId),
		TelegramNotificationsStatusChanges:                  optionalBool(apiModel.TelegramNotificationsStatusChanges),
		TelegramNotificationsStatusChangesMessageThreadId:   optionalString(apiModel.TelegramNotificationsStatusChangesMessageThreadId),
		TelegramNotificationsTest:                           optionalBool(apiModel.TelegramNotificationsTest),
		TelegramNotificationsTestMessageThreadId:            optionalString(apiModel.TelegramNotificationsTestMessageThreadId),
		TelegramToken:                                       optionalString(apiModel.TelegramToken),
		UpdatedAt:                                           optionalString(apiModel.UpdatedAt),
		UseInstanceEmailSettings:                            optionalBool(apiModel.UseInstanceEmailSettings),
	}
}

var teamsFilterNames = []string{"name", "description", "id", "discord_enabled", "resend_enabled", "smtp_enabled", "telegram_enabled"}

func (m teamModel) FilterAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":      m.Description,
		"discord_enabled":  m.DiscordEnabled,
		"id":               m.Id,
		"name":             m.Name,
		"personal_team":    m.PersonalTeam,
		"resend_enabled":   m.ResendEnabled,
		"smtp_enabled":     m.SmtpEnabled,
		"telegram_enabled": m.TelegramEnabled,
	}
}
