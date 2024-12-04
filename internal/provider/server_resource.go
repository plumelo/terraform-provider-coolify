package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/resource_server"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

func NewServerResource() resource.Resource {
	return &serverResource{}
}

type serverResource struct {
	client *api.ClientWithResponses
}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_server.ServerResourceSchema(ctx)
	resp.Schema.Description = "Create, read, update, and delete a Coolify server resource." +
		"\n**NOTE:** This resource is not fully implemented and may not work as expected because the Coolify API is incomplete."

	requiredAttrs := []string{"name", "private_key_uuid", "ip", "instant_validate"}
	for _, attr := range requiredAttrs {
		makeResourceAttributeRequired(resp.Schema.Attributes, attr)
	}
}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_server.ServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating server", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})
	createResp, err := r.client.CreateServerWithResponse(ctx, api.CreateServerJSONRequestBody{
		Description:     plan.Description.ValueStringPointer(),
		Name:            plan.Name.ValueStringPointer(),
		InstantValidate: plan.InstantValidate.ValueBoolPointer(),
		Ip:              plan.Ip.ValueStringPointer(),
		IsBuildServer:   plan.IsBuildServer.ValueBoolPointer(),
		Port: func() *int {
			if plan.Port.IsUnknown() || plan.Port.IsNull() {
				return nil
			}
			value := int(*plan.Port.ValueInt64Pointer())
			return &value
		}(),
		PrivateKeyUuid: plan.PrivateKeyUuid.ValueStringPointer(),
		User:           plan.User.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating server",
			fmt.Sprintf("Received %s creating server. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid)
	r.copyMissingAttributes(&plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_server.ServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading server", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	r.copyMissingAttributes(&state, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_server.ServerModel
	var state resource_server.ServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	if uuid == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// Update API call logic
	tflog.Debug(ctx, "Updating server", map[string]interface{}{
		"uuid": uuid,
	})
	updateResp, err := r.client.UpdateServerByUuidWithResponse(ctx, uuid, api.UpdateServerByUuidJSONRequestBody{
		Description:     plan.Description.ValueStringPointer(),
		Name:            plan.Name.ValueStringPointer(),
		InstantValidate: plan.InstantValidate.ValueBoolPointer(),
		Ip:              plan.Ip.ValueStringPointer(),
		IsBuildServer:   plan.IsBuildServer.ValueBoolPointer(),
		Port: func() *int { // todo: make a reusable fn for these inline conversions
			if plan.Port.IsUnknown() || plan.Port.IsNull() {
				return nil
			}
			value := int(*plan.Port.ValueInt64Pointer())
			return &value
		}(),
		PrivateKeyUuid: plan.PrivateKeyUuid.ValueStringPointer(),
		User: func() *string {
			if plan.User.IsUnknown() {
				return nil
			}
			return plan.User.ValueStringPointer()
		}(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating server: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating server",
			fmt.Sprintf("Received %s updating server: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid)
	r.copyMissingAttributes(&plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_server.ServerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting server", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting server",
			fmt.Sprintf("Received %s deleting server: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *serverResource) copyMissingAttributes(
	plan *resource_server.ServerModel,
	data *resource_server.ServerModel,
) {
	// Values that are not returned in API response
	data.InstantValidate = plan.InstantValidate
	data.PrivateKeyUuid = plan.PrivateKeyUuid

	if plan.PrivateKeyUuid.IsNull() {
		data.PrivateKeyUuid = types.StringValue("")
	}

	// Values that are incorrectly mapped in API
	data.Id = data.Settings.ServerId
}

func (r *serverResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) resource_server.ServerModel {
	readResp, err := r.client.GetServerByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading server: uuid=%s", uuid),
			err.Error(),
		)
		return resource_server.ServerModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return resource_server.ServerModel{}
	}

	return r.ApiToModel(ctx, diags, readResp.JSON200)
}

func (r *serverResource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Server,
) resource_server.ServerModel {
	settings := resource_server.NewSettingsValueMust(
		resource_server.SettingsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"concurrent_builds":                     optionalInt64(response.Settings.ConcurrentBuilds),
			"created_at":                            optionalString(response.Settings.CreatedAt),
			"delete_unused_networks":                optionalBool(response.Settings.DeleteUnusedNetworks),
			"delete_unused_volumes":                 optionalBool(response.Settings.DeleteUnusedVolumes),
			"docker_cleanup_frequency":              optionalString(response.Settings.DockerCleanupFrequency),
			"docker_cleanup_threshold":              optionalInt64(response.Settings.DockerCleanupThreshold),
			"dynamic_timeout":                       optionalInt64(response.Settings.DynamicTimeout),
			"force_disabled":                        optionalBool(response.Settings.ForceDisabled),
			"force_server_cleanup":                  optionalBool(response.Settings.ForceServerCleanup),
			"id":                                    optionalInt64(response.Settings.Id),
			"is_build_server":                       optionalBool(response.Settings.IsBuildServer),
			"is_cloudflare_tunnel":                  optionalBool(response.Settings.IsCloudflareTunnel),
			"is_jump_server":                        optionalBool(response.Settings.IsJumpServer),
			"is_logdrain_axiom_enabled":             optionalBool(response.Settings.IsLogdrainAxiomEnabled),
			"is_logdrain_custom_enabled":            optionalBool(response.Settings.IsLogdrainCustomEnabled),
			"is_logdrain_highlight_enabled":         optionalBool(response.Settings.IsLogdrainHighlightEnabled),
			"is_logdrain_newrelic_enabled":          optionalBool(response.Settings.IsLogdrainNewrelicEnabled),
			"is_metrics_enabled":                    optionalBool(response.Settings.IsMetricsEnabled),
			"is_reachable":                          optionalBool(response.Settings.IsReachable),
			"is_sentinel_enabled":                   optionalBool(response.Settings.IsSentinelEnabled),
			"is_swarm_manager":                      optionalBool(response.Settings.IsSwarmManager),
			"is_swarm_worker":                       optionalBool(response.Settings.IsSwarmWorker),
			"is_usable":                             optionalBool(response.Settings.IsUsable),
			"logdrain_axiom_api_key":                optionalString(response.Settings.LogdrainAxiomApiKey),
			"logdrain_axiom_dataset_name":           optionalString(response.Settings.LogdrainAxiomDatasetName),
			"logdrain_custom_config":                optionalString(response.Settings.LogdrainCustomConfig),
			"logdrain_custom_config_parser":         optionalString(response.Settings.LogdrainCustomConfigParser),
			"logdrain_highlight_project_id":         optionalString(response.Settings.LogdrainHighlightProjectId),
			"logdrain_newrelic_base_uri":            optionalString(response.Settings.LogdrainNewrelicBaseUri),
			"logdrain_newrelic_license_key":         optionalString(response.Settings.LogdrainNewrelicLicenseKey),
			"sentinel_metrics_history_days":         optionalInt64(response.Settings.SentinelMetricsHistoryDays),
			"sentinel_metrics_refresh_rate_seconds": optionalInt64(response.Settings.SentinelMetricsRefreshRateSeconds),
			"sentinel_token":                        optionalString(response.Settings.SentinelToken),
			"server_id":                             optionalInt64(response.Settings.ServerId),
			"updated_at":                            optionalString(response.Settings.UpdatedAt),
			"wildcard_domain":                       optionalString(response.Settings.WildcardDomain),
		},
	)

	return resource_server.ServerModel{
		Description:                   optionalString(response.Description),
		HighDiskUsageNotificationSent: optionalBool(response.HighDiskUsageNotificationSent), // missing
		Id:                            optionalInt64(response.Id),
		Ip:                            optionalString(response.Ip),
		IsBuildServer:                 optionalBool(response.Settings.IsBuildServer),
		LogDrainNotificationSent:      optionalBool(response.LogDrainNotificationSent),
		Name:                          optionalString(response.Name),
		Port:                          optionalInt64(response.Port),
		SwarmCluster:                  optionalString(response.SwarmCluster),
		UnreachableCount:              optionalInt64(response.UnreachableCount),
		UnreachableNotificationSent:   optionalBool(response.UnreachableNotificationSent),
		User:                          optionalString(response.User),
		Uuid:                          optionalString(response.Uuid),
		ValidationLogs:                optionalString(response.ValidationLogs),

		// Proxy:                         resource_server.NewProxyValueUnknown(),
		ProxyType:       optionalString((*string)(response.ProxyType)), // enum value
		PrivateKeyUuid:  types.StringUnknown(),
		InstantValidate: types.BoolUnknown(),
		Settings:        settings,
	}
}
