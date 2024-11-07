package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_server"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serverDataSource{}
var _ datasource.DataSourceWithConfigure = &serverDataSource{}

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{}
}

type serverDataSource struct {
	providerData CoolifyProviderData
}

func (d *serverDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *serverDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_server.ServerDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify server by `uuid`."
}

func (d *serverDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_server.ServerModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.providerData.client.GetServerByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server", err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server. Details: %s", response.Status(), string(response.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serverDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Server,
) datasource_server.ServerModel {

	settings := datasource_server.NewSettingsValueMust(
		datasource_server.SettingsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"concurrent_builds":             optionalInt64(response.Settings.ConcurrentBuilds),
			"created_at":                    optionalString(response.Settings.CreatedAt),
			"delete_unused_networks":        optionalBool(response.Settings.DeleteUnusedNetworks),
			"delete_unused_volumes":         optionalBool(response.Settings.DeleteUnusedVolumes),
			"docker_cleanup_frequency":      optionalString(response.Settings.DockerCleanupFrequency),
			"docker_cleanup_threshold":      optionalInt64(response.Settings.DockerCleanupThreshold),
			"dynamic_timeout":               optionalInt64(response.Settings.DynamicTimeout),
			"force_disabled":                optionalBool(response.Settings.ForceDisabled),
			"force_server_cleanup":          optionalBool(response.Settings.ForceServerCleanup),
			"id":                            optionalInt64(response.Settings.Id),
			"is_build_server":               optionalBool(response.Settings.IsBuildServer),
			"is_cloudflare_tunnel":          optionalBool(response.Settings.IsCloudflareTunnel),
			"is_jump_server":                optionalBool(response.Settings.IsJumpServer),
			"is_logdrain_axiom_enabled":     optionalBool(response.Settings.IsLogdrainAxiomEnabled),
			"is_logdrain_custom_enabled":    optionalBool(response.Settings.IsLogdrainCustomEnabled),
			"is_logdrain_highlight_enabled": optionalBool(response.Settings.IsLogdrainHighlightEnabled),
			"is_logdrain_newrelic_enabled":  optionalBool(response.Settings.IsLogdrainNewrelicEnabled),
			"is_metrics_enabled":            optionalBool(response.Settings.IsMetricsEnabled),
			"is_reachable":                  optionalBool(response.Settings.IsReachable),
			"is_server_api_enabled":         optionalBool(response.Settings.IsServerApiEnabled),
			"is_swarm_manager":              optionalBool(response.Settings.IsSwarmManager),
			"is_swarm_worker":               optionalBool(response.Settings.IsSwarmWorker),
			"is_usable":                     optionalBool(response.Settings.IsUsable),
			"logdrain_axiom_api_key":        optionalString(response.Settings.LogdrainAxiomApiKey),
			"logdrain_axiom_dataset_name":   optionalString(response.Settings.LogdrainAxiomDatasetName),
			"logdrain_custom_config":        optionalString(response.Settings.LogdrainCustomConfig),
			"logdrain_custom_config_parser": optionalString(response.Settings.LogdrainCustomConfigParser),
			"logdrain_highlight_project_id": optionalString(response.Settings.LogdrainHighlightProjectId),
			"logdrain_newrelic_base_uri":    optionalString(response.Settings.LogdrainNewrelicBaseUri),
			"logdrain_newrelic_license_key": optionalString(response.Settings.LogdrainNewrelicLicenseKey),
			"metrics_history_days":          optionalInt64(response.Settings.MetricsHistoryDays),
			"metrics_refresh_rate_seconds":  optionalInt64(response.Settings.MetricsRefreshRateSeconds),
			"metrics_token":                 optionalString(response.Settings.MetricsToken),
			"server_id":                     optionalInt64(response.Settings.ServerId),
			"updated_at":                    optionalString(response.Settings.UpdatedAt),
			"wildcard_domain":               optionalString(response.Settings.WildcardDomain),
		},
	)

	return datasource_server.ServerModel{
		Description:                   optionalString(response.Description),
		HighDiskUsageNotificationSent: optionalBool(response.HighDiskUsageNotificationSent),
		Id:                            optionalInt64(response.Settings.ServerId), // TODO: this should be `id` on root object, upstream spec is wrong
		Ip:                            optionalString(response.Ip),
		LogDrainNotificationSent:      optionalBool(response.LogDrainNotificationSent),
		Name:                          optionalString(response.Name),
		Port:                          optionalString(response.Port),
		Settings:                      settings,
		SwarmCluster:                  optionalString(response.SwarmCluster),
		UnreachableCount:              optionalInt64(response.UnreachableCount),
		UnreachableNotificationSent:   optionalBool(response.UnreachableNotificationSent),
		User:                          optionalString(response.User),
		Uuid:                          optionalString(response.Uuid),
		ValidationLogs:                optionalString(response.ValidationLogs),
	}
}
