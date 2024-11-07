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
	"terraform-provider-coolify/internal/provider/generated/datasource_servers"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serversDataSource{}
var _ datasource.DataSourceWithConfigure = &serversDataSource{}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{}
}

type serversDataSource struct {
	providerData CoolifyProviderData
}

type serversDataSourceWithFilterModel struct {
	datasource_servers.ServersModel
	Filter []filterBlockModel `tfsdk:"filter"`
}

var serversFilterNames = []string{"id", "uuid", "user", "ip", "name", "description"}

func (d *serversDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

func (d *serversDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_servers.ServersDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify servers."

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": createDatasourceFilter(serversFilterNames),
	}
}
func (d *serversDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *serversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan serversDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.providerData.client.ListServersWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading servers", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading servers",
			fmt.Sprintf("Received %s for servers. Details: %s", listResponse.Status(), listResponse.Body),
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

func (d *serversDataSource) apiToModel(
	ctx context.Context,
	response *[]api.Server,
	filters []filterBlockModel,
) (serversDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var elements []attr.Value

	for _, sv := range *response {
		settings, diag := datasource_servers.NewSettingsValueMust(
			datasource_servers.SettingsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"concurrent_builds":             optionalInt64(sv.Settings.ConcurrentBuilds),
				"created_at":                    optionalString(sv.Settings.CreatedAt),
				"delete_unused_networks":        optionalBool(sv.Settings.DeleteUnusedNetworks),
				"delete_unused_volumes":         optionalBool(sv.Settings.DeleteUnusedVolumes),
				"docker_cleanup_frequency":      optionalString(sv.Settings.DockerCleanupFrequency),
				"docker_cleanup_threshold":      optionalInt64(sv.Settings.DockerCleanupThreshold),
				"dynamic_timeout":               optionalInt64(sv.Settings.DynamicTimeout),
				"force_disabled":                optionalBool(sv.Settings.ForceDisabled),
				"force_server_cleanup":          optionalBool(sv.Settings.ForceServerCleanup),
				"id":                            optionalInt64(sv.Settings.Id),
				"is_build_server":               optionalBool(sv.Settings.IsBuildServer),
				"is_cloudflare_tunnel":          optionalBool(sv.Settings.IsCloudflareTunnel),
				"is_jump_server":                optionalBool(sv.Settings.IsJumpServer),
				"is_logdrain_axiom_enabled":     optionalBool(sv.Settings.IsLogdrainAxiomEnabled),
				"is_logdrain_custom_enabled":    optionalBool(sv.Settings.IsLogdrainCustomEnabled),
				"is_logdrain_highlight_enabled": optionalBool(sv.Settings.IsLogdrainHighlightEnabled),
				"is_logdrain_newrelic_enabled":  optionalBool(sv.Settings.IsLogdrainNewrelicEnabled),
				"is_metrics_enabled":            optionalBool(sv.Settings.IsMetricsEnabled),
				"is_reachable":                  optionalBool(sv.Settings.IsReachable),
				"is_server_api_enabled":         optionalBool(sv.Settings.IsServerApiEnabled),
				"is_swarm_manager":              optionalBool(sv.Settings.IsSwarmManager),
				"is_swarm_worker":               optionalBool(sv.Settings.IsSwarmWorker),
				"is_usable":                     optionalBool(sv.Settings.IsUsable),
				"logdrain_axiom_api_key":        optionalString(sv.Settings.LogdrainAxiomApiKey),
				"logdrain_axiom_dataset_name":   optionalString(sv.Settings.LogdrainAxiomDatasetName),
				"logdrain_custom_config":        optionalString(sv.Settings.LogdrainCustomConfig),
				"logdrain_custom_config_parser": optionalString(sv.Settings.LogdrainCustomConfigParser),
				"logdrain_highlight_project_id": optionalString(sv.Settings.LogdrainHighlightProjectId),
				"logdrain_newrelic_base_uri":    optionalString(sv.Settings.LogdrainNewrelicBaseUri),
				"logdrain_newrelic_license_key": optionalString(sv.Settings.LogdrainNewrelicLicenseKey),
				"metrics_history_days":          optionalInt64(sv.Settings.MetricsHistoryDays),
				"metrics_refresh_rate_seconds":  optionalInt64(sv.Settings.MetricsRefreshRateSeconds),
				"metrics_token":                 optionalString(sv.Settings.MetricsToken),
				"server_id":                     optionalInt64(sv.Settings.ServerId),
				"updated_at":                    optionalString(sv.Settings.UpdatedAt),
				"wildcard_domain":               optionalString(sv.Settings.WildcardDomain),
			},
		).ToObjectValue(ctx)
		diags.Append(diag...)

		attributes := map[string]attr.Value{
			"description":                       optionalString(sv.Description),
			"high_disk_usage_notification_sent": optionalBool(sv.HighDiskUsageNotificationSent),
			"id":                                optionalInt64(sv.Settings.ServerId), // TODO: this should be `id` on root object, upstream spec is wrong
			"ip":                                optionalString(sv.Ip),
			"log_drain_notification_sent":       optionalBool(sv.LogDrainNotificationSent),
			"name":                              optionalString(sv.Name),
			"port":                              optionalString(sv.Port),
			"settings":                          settings,
			"swarm_cluster":                     optionalString(sv.SwarmCluster),
			"unreachable_count":                 optionalInt64(sv.UnreachableCount),
			"unreachable_notification_sent":     optionalBool(sv.UnreachableNotificationSent),
			"user":                              optionalString(sv.User),
			"uuid":                              optionalString(sv.Uuid),
			"validation_logs":                   optionalString(sv.ValidationLogs),
		}

		if !filterOnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_servers.NewServersValue(
			datasource_servers.ServersValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_servers.ServersValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return serversDataSourceWithFilterModel{
		ServersModel: datasource_servers.ServersModel{
			Servers: dataSet,
		},
		Filter: filters,
	}, diags
}
