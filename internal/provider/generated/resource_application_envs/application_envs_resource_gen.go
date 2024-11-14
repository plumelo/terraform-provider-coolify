// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_application_envs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ApplicationEnvsResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"is_build_time": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the environment variable is used in build time.",
				MarkdownDescription: "The flag to indicate if the environment variable is used in build time.",
			},
			"is_literal": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the environment variable is a literal, nothing espaced.",
				MarkdownDescription: "The flag to indicate if the environment variable is a literal, nothing espaced.",
			},
			"is_multiline": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the environment variable is multiline.",
				MarkdownDescription: "The flag to indicate if the environment variable is multiline.",
			},
			"is_preview": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the environment variable is used in preview deployments.",
				MarkdownDescription: "The flag to indicate if the environment variable is used in preview deployments.",
			},
			"is_shown_once": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The flag to indicate if the environment variable's value is shown on the UI.",
				MarkdownDescription: "The flag to indicate if the environment variable's value is shown on the UI.",
			},
			"key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The key of the environment variable.",
				MarkdownDescription: "The key of the environment variable.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				Description:         "UUID of the application.",
				MarkdownDescription: "UUID of the application.",
			},
			"value": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The value of the environment variable.",
				MarkdownDescription: "The value of the environment variable.",
			},
		},
	}
}

type ApplicationEnvsModel struct {
	IsBuildTime types.Bool   `tfsdk:"is_build_time"`
	IsLiteral   types.Bool   `tfsdk:"is_literal"`
	IsMultiline types.Bool   `tfsdk:"is_multiline"`
	IsPreview   types.Bool   `tfsdk:"is_preview"`
	IsShownOnce types.Bool   `tfsdk:"is_shown_once"`
	Key         types.String `tfsdk:"key"`
	Uuid        types.String `tfsdk:"uuid"`
	Value       types.String `tfsdk:"value"`
}
