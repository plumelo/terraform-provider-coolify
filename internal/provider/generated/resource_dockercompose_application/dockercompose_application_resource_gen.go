// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_dockercompose_application

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func DockercomposeApplicationResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_directory": schema.StringAttribute{
				Computed:            true,
				Description:         "Base directory for all commands.",
				MarkdownDescription: "Base directory for all commands.",
			},
			"build_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Build command.",
				MarkdownDescription: "Build command.",
			},
			"build_pack": schema.StringAttribute{
				Computed:            true,
				Description:         "Build pack.",
				MarkdownDescription: "Build pack.",
			},
			"compose_parsing_version": schema.StringAttribute{
				Computed:            true,
				Description:         "How Coolify parse the compose file.",
				MarkdownDescription: "How Coolify parse the compose file.",
			},
			"config_hash": schema.StringAttribute{
				Computed:            true,
				Description:         "Configuration hash.",
				MarkdownDescription: "Configuration hash.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The date and time when the application was created.",
				MarkdownDescription: "The date and time when the application was created.",
			},
			"custom_docker_run_options": schema.StringAttribute{
				Computed:            true,
				Description:         "Custom docker run options.",
				MarkdownDescription: "Custom docker run options.",
			},
			"custom_healthcheck_found": schema.BoolAttribute{
				Computed:            true,
				Description:         "Custom healthcheck found.",
				MarkdownDescription: "Custom healthcheck found.",
			},
			"custom_labels": schema.StringAttribute{
				Computed:            true,
				Description:         "Custom labels.",
				MarkdownDescription: "Custom labels.",
			},
			"custom_nginx_configuration": schema.StringAttribute{
				Computed:            true,
				Description:         "Custom Nginx configuration base64 encoded.",
				MarkdownDescription: "Custom Nginx configuration base64 encoded.",
			},
			"deleted_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The date and time when the application was deleted.",
				MarkdownDescription: "The date and time when the application was deleted.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The application description.",
				MarkdownDescription: "The application description.",
			},
			"destination_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Destination identifier.",
				MarkdownDescription: "Destination identifier.",
			},
			"destination_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Destination type.",
				MarkdownDescription: "Destination type.",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The destination UUID if the server has more than one destinations.",
				MarkdownDescription: "The destination UUID if the server has more than one destinations.",
			},
			"docker_compose": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker compose content. Used for docker compose build pack.",
				MarkdownDescription: "Docker compose content. Used for docker compose build pack.",
			},
			"docker_compose_custom_build_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker compose custom build command.",
				MarkdownDescription: "Docker compose custom build command.",
			},
			"docker_compose_custom_start_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker compose custom start command.",
				MarkdownDescription: "Docker compose custom start command.",
			},
			"docker_compose_domains": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker compose domains.",
				MarkdownDescription: "Docker compose domains.",
			},
			"docker_compose_location": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker compose location.",
				MarkdownDescription: "Docker compose location.",
			},
			"docker_compose_raw": schema.StringAttribute{
				Required:            true,
				Description:         "The Docker Compose raw content.",
				MarkdownDescription: "The Docker Compose raw content.",
			},
			"docker_registry_image_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker registry image name.",
				MarkdownDescription: "Docker registry image name.",
			},
			"docker_registry_image_tag": schema.StringAttribute{
				Computed:            true,
				Description:         "Docker registry image tag.",
				MarkdownDescription: "Docker registry image tag.",
			},
			"dockerfile": schema.StringAttribute{
				Computed:            true,
				Description:         "Dockerfile content. Used for dockerfile build pack.",
				MarkdownDescription: "Dockerfile content. Used for dockerfile build pack.",
			},
			"dockerfile_location": schema.StringAttribute{
				Computed:            true,
				Description:         "Dockerfile location.",
				MarkdownDescription: "Dockerfile location.",
			},
			"dockerfile_target_build": schema.StringAttribute{
				Computed:            true,
				Description:         "Dockerfile target build.",
				MarkdownDescription: "Dockerfile target build.",
			},
			"environment_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Environment identifier.",
				MarkdownDescription: "Environment identifier.",
			},
			"environment_name": schema.StringAttribute{
				Required:            true,
				Description:         "The environment name. You need to provide at least one of environment_name or environment_uuid.",
				MarkdownDescription: "The environment name. You need to provide at least one of environment_name or environment_uuid.",
			},
			"environment_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "The environment UUID. You need to provide at least one of environment_name or environment_uuid.",
				MarkdownDescription: "The environment UUID. You need to provide at least one of environment_name or environment_uuid.",
			},
			"fqdn": schema.StringAttribute{
				Computed:            true,
				Description:         "The application domains.",
				MarkdownDescription: "The application domains.",
			},
			"git_branch": schema.StringAttribute{
				Computed:            true,
				Description:         "Git branch.",
				MarkdownDescription: "Git branch.",
			},
			"git_commit_sha": schema.StringAttribute{
				Computed:            true,
				Description:         "Git commit SHA.",
				MarkdownDescription: "Git commit SHA.",
			},
			"git_full_url": schema.StringAttribute{
				Computed:            true,
				Description:         "Git full URL.",
				MarkdownDescription: "Git full URL.",
			},
			"git_repository": schema.StringAttribute{
				Computed:            true,
				Description:         "Git repository URL.",
				MarkdownDescription: "Git repository URL.",
			},
			"health_check_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Health check enabled.",
				MarkdownDescription: "Health check enabled.",
			},
			"health_check_host": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check host.",
				MarkdownDescription: "Health check host.",
			},
			"health_check_interval": schema.Int64Attribute{
				Computed:            true,
				Description:         "Health check interval in seconds.",
				MarkdownDescription: "Health check interval in seconds.",
			},
			"health_check_method": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check method.",
				MarkdownDescription: "Health check method.",
			},
			"health_check_path": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check path.",
				MarkdownDescription: "Health check path.",
			},
			"health_check_port": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check port.",
				MarkdownDescription: "Health check port.",
			},
			"health_check_response_text": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check response text.",
				MarkdownDescription: "Health check response text.",
			},
			"health_check_retries": schema.Int64Attribute{
				Computed:            true,
				Description:         "Health check retries count.",
				MarkdownDescription: "Health check retries count.",
			},
			"health_check_return_code": schema.Int64Attribute{
				Computed:            true,
				Description:         "Health check return code.",
				MarkdownDescription: "Health check return code.",
			},
			"health_check_scheme": schema.StringAttribute{
				Computed:            true,
				Description:         "Health check scheme.",
				MarkdownDescription: "Health check scheme.",
			},
			"health_check_start_period": schema.Int64Attribute{
				Computed:            true,
				Description:         "Health check start period in seconds.",
				MarkdownDescription: "Health check start period in seconds.",
			},
			"health_check_timeout": schema.Int64Attribute{
				Computed:            true,
				Description:         "Health check timeout in seconds.",
				MarkdownDescription: "Health check timeout in seconds.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The application identifier in the database.",
				MarkdownDescription: "The application identifier in the database.",
			},
			"install_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Install command.",
				MarkdownDescription: "Install command.",
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the application should be deployed instantly.",
				MarkdownDescription: "The flag to indicate if the application should be deployed instantly.",
			},
			"limits_cpu_shares": schema.Int64Attribute{
				Computed:            true,
				Description:         "CPU shares.",
				MarkdownDescription: "CPU shares.",
			},
			"limits_cpus": schema.StringAttribute{
				Computed:            true,
				Description:         "CPU limit.",
				MarkdownDescription: "CPU limit.",
			},
			"limits_cpuset": schema.StringAttribute{
				Computed:            true,
				Description:         "CPU set.",
				MarkdownDescription: "CPU set.",
			},
			"limits_memory": schema.StringAttribute{
				Computed:            true,
				Description:         "Memory limit.",
				MarkdownDescription: "Memory limit.",
			},
			"limits_memory_reservation": schema.StringAttribute{
				Computed:            true,
				Description:         "Memory reservation.",
				MarkdownDescription: "Memory reservation.",
			},
			"limits_memory_swap": schema.StringAttribute{
				Computed:            true,
				Description:         "Memory swap limit.",
				MarkdownDescription: "Memory swap limit.",
			},
			"limits_memory_swappiness": schema.Int64Attribute{
				Computed:            true,
				Description:         "Memory swappiness.",
				MarkdownDescription: "Memory swappiness.",
			},
			"manual_webhook_secret_bitbucket": schema.StringAttribute{
				Computed:            true,
				Description:         "Manual webhook secret for Bitbucket.",
				MarkdownDescription: "Manual webhook secret for Bitbucket.",
			},
			"manual_webhook_secret_gitea": schema.StringAttribute{
				Computed:            true,
				Description:         "Manual webhook secret for Gitea.",
				MarkdownDescription: "Manual webhook secret for Gitea.",
			},
			"manual_webhook_secret_github": schema.StringAttribute{
				Computed:            true,
				Description:         "Manual webhook secret for GitHub.",
				MarkdownDescription: "Manual webhook secret for GitHub.",
			},
			"manual_webhook_secret_gitlab": schema.StringAttribute{
				Computed:            true,
				Description:         "Manual webhook secret for GitLab.",
				MarkdownDescription: "Manual webhook secret for GitLab.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The application name.",
				MarkdownDescription: "The application name.",
			},
			"ports_exposes": schema.StringAttribute{
				Computed:            true,
				Description:         "Ports exposes.",
				MarkdownDescription: "Ports exposes.",
			},
			"ports_mappings": schema.StringAttribute{
				Computed:            true,
				Description:         "Ports mappings.",
				MarkdownDescription: "Ports mappings.",
			},
			"post_deployment_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Post deployment command.",
				MarkdownDescription: "Post deployment command.",
			},
			"post_deployment_command_container": schema.StringAttribute{
				Computed:            true,
				Description:         "Post deployment command container.",
				MarkdownDescription: "Post deployment command container.",
			},
			"pre_deployment_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Pre deployment command.",
				MarkdownDescription: "Pre deployment command.",
			},
			"pre_deployment_command_container": schema.StringAttribute{
				Computed:            true,
				Description:         "Pre deployment command container.",
				MarkdownDescription: "Pre deployment command container.",
			},
			"preview_url_template": schema.StringAttribute{
				Computed:            true,
				Description:         "Preview URL template.",
				MarkdownDescription: "Preview URL template.",
			},
			"private_key_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Private key identifier.",
				MarkdownDescription: "Private key identifier.",
			},
			"project_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "The project UUID.",
				MarkdownDescription: "The project UUID.",
			},
			"publish_directory": schema.StringAttribute{
				Computed:            true,
				Description:         "Publish directory.",
				MarkdownDescription: "Publish directory.",
			},
			"redirect": schema.StringAttribute{
				Computed:            true,
				Description:         "How to set redirect with Traefik / Caddy. www<->non-www.",
				MarkdownDescription: "How to set redirect with Traefik / Caddy. www<->non-www.",
			},
			"repository_project_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The repository project identifier.",
				MarkdownDescription: "The repository project identifier.",
			},
			"server_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "The server UUID.",
				MarkdownDescription: "The server UUID.",
			},
			"source_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Source identifier.",
				MarkdownDescription: "Source identifier.",
			},
			"start_command": schema.StringAttribute{
				Computed:            true,
				Description:         "Start command.",
				MarkdownDescription: "Start command.",
			},
			"static_image": schema.StringAttribute{
				Computed:            true,
				Description:         "Static image used when static site is deployed.",
				MarkdownDescription: "Static image used when static site is deployed.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "Application status.",
				MarkdownDescription: "Application status.",
			},
			"swarm_placement_constraints": schema.StringAttribute{
				Computed:            true,
				Description:         "Swarm placement constraints. Only used for swarm deployments.",
				MarkdownDescription: "Swarm placement constraints. Only used for swarm deployments.",
			},
			"swarm_replicas": schema.Int64Attribute{
				Computed:            true,
				Description:         "Swarm replicas. Only used for swarm deployments.",
				MarkdownDescription: "Swarm replicas. Only used for swarm deployments.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The date and time when the application was last updated.",
				MarkdownDescription: "The date and time when the application was last updated.",
			},
			"use_build_server": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Use build server.",
				MarkdownDescription: "Use build server.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				Description:         "The application UUID.",
				MarkdownDescription: "The application UUID.",
			},
			"watch_paths": schema.StringAttribute{
				Computed:            true,
				Description:         "Watch paths.",
				MarkdownDescription: "Watch paths.",
			},
		},
	}
}

type DockercomposeApplicationModel struct {
	BaseDirectory                   types.String `tfsdk:"base_directory"`
	BuildCommand                    types.String `tfsdk:"build_command"`
	BuildPack                       types.String `tfsdk:"build_pack"`
	ComposeParsingVersion           types.String `tfsdk:"compose_parsing_version"`
	ConfigHash                      types.String `tfsdk:"config_hash"`
	CreatedAt                       types.String `tfsdk:"created_at"`
	CustomDockerRunOptions          types.String `tfsdk:"custom_docker_run_options"`
	CustomHealthcheckFound          types.Bool   `tfsdk:"custom_healthcheck_found"`
	CustomLabels                    types.String `tfsdk:"custom_labels"`
	CustomNginxConfiguration        types.String `tfsdk:"custom_nginx_configuration"`
	DeletedAt                       types.String `tfsdk:"deleted_at"`
	Description                     types.String `tfsdk:"description"`
	DestinationId                   types.Int64  `tfsdk:"destination_id"`
	DestinationType                 types.String `tfsdk:"destination_type"`
	DestinationUuid                 types.String `tfsdk:"destination_uuid"`
	DockerCompose                   types.String `tfsdk:"docker_compose"`
	DockerComposeCustomBuildCommand types.String `tfsdk:"docker_compose_custom_build_command"`
	DockerComposeCustomStartCommand types.String `tfsdk:"docker_compose_custom_start_command"`
	DockerComposeDomains            types.String `tfsdk:"docker_compose_domains"`
	DockerComposeLocation           types.String `tfsdk:"docker_compose_location"`
	DockerComposeRaw                types.String `tfsdk:"docker_compose_raw"`
	DockerRegistryImageName         types.String `tfsdk:"docker_registry_image_name"`
	DockerRegistryImageTag          types.String `tfsdk:"docker_registry_image_tag"`
	Dockerfile                      types.String `tfsdk:"dockerfile"`
	DockerfileLocation              types.String `tfsdk:"dockerfile_location"`
	DockerfileTargetBuild           types.String `tfsdk:"dockerfile_target_build"`
	EnvironmentId                   types.Int64  `tfsdk:"environment_id"`
	EnvironmentName                 types.String `tfsdk:"environment_name"`
	EnvironmentUuid                 types.String `tfsdk:"environment_uuid"`
	Fqdn                            types.String `tfsdk:"fqdn"`
	GitBranch                       types.String `tfsdk:"git_branch"`
	GitCommitSha                    types.String `tfsdk:"git_commit_sha"`
	GitFullUrl                      types.String `tfsdk:"git_full_url"`
	GitRepository                   types.String `tfsdk:"git_repository"`
	HealthCheckEnabled              types.Bool   `tfsdk:"health_check_enabled"`
	HealthCheckHost                 types.String `tfsdk:"health_check_host"`
	HealthCheckInterval             types.Int64  `tfsdk:"health_check_interval"`
	HealthCheckMethod               types.String `tfsdk:"health_check_method"`
	HealthCheckPath                 types.String `tfsdk:"health_check_path"`
	HealthCheckPort                 types.String `tfsdk:"health_check_port"`
	HealthCheckResponseText         types.String `tfsdk:"health_check_response_text"`
	HealthCheckRetries              types.Int64  `tfsdk:"health_check_retries"`
	HealthCheckReturnCode           types.Int64  `tfsdk:"health_check_return_code"`
	HealthCheckScheme               types.String `tfsdk:"health_check_scheme"`
	HealthCheckStartPeriod          types.Int64  `tfsdk:"health_check_start_period"`
	HealthCheckTimeout              types.Int64  `tfsdk:"health_check_timeout"`
	Id                              types.Int64  `tfsdk:"id"`
	InstallCommand                  types.String `tfsdk:"install_command"`
	InstantDeploy                   types.Bool   `tfsdk:"instant_deploy"`
	LimitsCpuShares                 types.Int64  `tfsdk:"limits_cpu_shares"`
	LimitsCpus                      types.String `tfsdk:"limits_cpus"`
	LimitsCpuset                    types.String `tfsdk:"limits_cpuset"`
	LimitsMemory                    types.String `tfsdk:"limits_memory"`
	LimitsMemoryReservation         types.String `tfsdk:"limits_memory_reservation"`
	LimitsMemorySwap                types.String `tfsdk:"limits_memory_swap"`
	LimitsMemorySwappiness          types.Int64  `tfsdk:"limits_memory_swappiness"`
	ManualWebhookSecretBitbucket    types.String `tfsdk:"manual_webhook_secret_bitbucket"`
	ManualWebhookSecretGitea        types.String `tfsdk:"manual_webhook_secret_gitea"`
	ManualWebhookSecretGithub       types.String `tfsdk:"manual_webhook_secret_github"`
	ManualWebhookSecretGitlab       types.String `tfsdk:"manual_webhook_secret_gitlab"`
	Name                            types.String `tfsdk:"name"`
	PortsExposes                    types.String `tfsdk:"ports_exposes"`
	PortsMappings                   types.String `tfsdk:"ports_mappings"`
	PostDeploymentCommand           types.String `tfsdk:"post_deployment_command"`
	PostDeploymentCommandContainer  types.String `tfsdk:"post_deployment_command_container"`
	PreDeploymentCommand            types.String `tfsdk:"pre_deployment_command"`
	PreDeploymentCommandContainer   types.String `tfsdk:"pre_deployment_command_container"`
	PreviewUrlTemplate              types.String `tfsdk:"preview_url_template"`
	PrivateKeyId                    types.Int64  `tfsdk:"private_key_id"`
	ProjectUuid                     types.String `tfsdk:"project_uuid"`
	PublishDirectory                types.String `tfsdk:"publish_directory"`
	Redirect                        types.String `tfsdk:"redirect"`
	RepositoryProjectId             types.Int64  `tfsdk:"repository_project_id"`
	ServerUuid                      types.String `tfsdk:"server_uuid"`
	SourceId                        types.Int64  `tfsdk:"source_id"`
	StartCommand                    types.String `tfsdk:"start_command"`
	StaticImage                     types.String `tfsdk:"static_image"`
	Status                          types.String `tfsdk:"status"`
	SwarmPlacementConstraints       types.String `tfsdk:"swarm_placement_constraints"`
	SwarmReplicas                   types.Int64  `tfsdk:"swarm_replicas"`
	UpdatedAt                       types.String `tfsdk:"updated_at"`
	UseBuildServer                  types.Bool   `tfsdk:"use_build_server"`
	Uuid                            types.String `tfsdk:"uuid"`
	WatchPaths                      types.String `tfsdk:"watch_paths"`
}
