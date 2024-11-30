package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
)

type commonDatabaseModel struct {
	Description             types.String `tfsdk:"description"`
	DestinationUuid         types.String `tfsdk:"destination_uuid"`
	EnvironmentName         types.String `tfsdk:"environment_name"`
	Image                   types.String `tfsdk:"image"`
	InstantDeploy           types.Bool   `tfsdk:"instant_deploy"`
	IsPublic                types.Bool   `tfsdk:"is_public"`
	LimitsCpuShares         types.Int64  `tfsdk:"limits_cpu_shares"`
	LimitsCpus              types.String `tfsdk:"limits_cpus"`
	LimitsCpuset            types.String `tfsdk:"limits_cpuset"`
	LimitsMemory            types.String `tfsdk:"limits_memory"`
	LimitsMemoryReservation types.String `tfsdk:"limits_memory_reservation"`
	LimitsMemorySwap        types.String `tfsdk:"limits_memory_swap"`
	LimitsMemorySwappiness  types.Int64  `tfsdk:"limits_memory_swappiness"`
	Name                    types.String `tfsdk:"name"`
	ProjectUuid             types.String `tfsdk:"project_uuid"`
	PublicPort              types.Int64  `tfsdk:"public_port"`
	ServerUuid              types.String `tfsdk:"server_uuid"`
	Uuid                    types.String `tfsdk:"uuid"`
	InternalDbUrl           types.String `tfsdk:"internal_db_url"`
}

func (m commonDatabaseModel) CommonSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the database",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "UUID of the destination if the server has multiple destinations",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Default:       stringdefault.StaticString(""),
			},
			"environment_name": schema.StringAttribute{
				Required:      true,
				Description:   "Name of the environment",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"image": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Docker Image of the database",
				Default:     stringdefault.StaticString("postgres:16-alpine"),
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Instant deploy the database",
				Default:     booldefault.StaticBool(false),
			},
			"is_public": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Is the database public?",
				Default:     booldefault.StaticBool(false),
			},
			"limits_cpu_shares": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU shares of the database",
				Default:     int64default.StaticInt64(1024),
			},
			"limits_cpus": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU limit of the database",
				Default:     stringdefault.StaticString("0"),
			},
			"limits_cpuset": schema.StringAttribute{
				Optional:    true,
				Description: "CPU set of the database",
			},
			"limits_memory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory limit of the database",
				Default:     stringdefault.StaticString("0"),
			},
			"limits_memory_reservation": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory reservation of the database",
				Default:     stringdefault.StaticString("0"),
			},
			"limits_memory_swap": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory swap limit of the database",
				Default:     stringdefault.StaticString("0"),
			},
			"limits_memory_swappiness": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory swappiness of the database",
				Default:     int64default.StaticInt64(60),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the database",
			},
			"project_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the project",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"public_port": schema.Int64Attribute{
				Optional:    true,
				Description: "Public port of the database",
			},
			"server_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the server",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"uuid": schema.StringAttribute{
				Computed:      true,
				Description:   "UUID of the database.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"internal_db_url": schema.StringAttribute{
				Computed:      true,
				Sensitive:     true,
				Description:   "Internal URL of the database.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (m commonDatabaseModel) FromAPI(apiModel *api.Database, state commonDatabaseModel) commonDatabaseModel {
	db, _ := apiModel.AsDatabaseCommon()

	return commonDatabaseModel{
		Uuid:                    types.StringValue(db.Uuid),
		Name:                    optionalString(db.Name),
		Description:             optionalString(db.Description),
		ServerUuid:              state.ServerUuid, // Values not returned by API, so use the plan value
		ProjectUuid:             state.ProjectUuid,
		EnvironmentName:         state.EnvironmentName,
		DestinationUuid:         state.DestinationUuid,
		InstantDeploy:           state.InstantDeploy,
		InternalDbUrl:           optionalString(db.InternalDbUrl),
		Image:                   optionalString(db.Image),
		IsPublic:                optionalBool(db.IsPublic),
		PublicPort:              optionalInt64(db.PublicPort),
		LimitsCpuShares:         optionalInt64(db.LimitsCpuShares),
		LimitsCpus:              optionalString(db.LimitsCpus),
		LimitsCpuset:            optionalString(db.LimitsCpuset),
		LimitsMemory:            optionalString(db.LimitsMemory),
		LimitsMemoryReservation: optionalString(db.LimitsMemoryReservation),
		LimitsMemorySwap:        optionalString(db.LimitsMemorySwap),
		LimitsMemorySwappiness:  optionalInt64(db.LimitsMemorySwappiness),
	}
}
