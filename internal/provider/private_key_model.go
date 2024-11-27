package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
)

type privateKeyModel struct {
	Description  types.String `tfsdk:"description"`
	Fingerprint  types.String `tfsdk:"fingerprint"`
	Id           types.Int64  `tfsdk:"id"`
	IsGitRelated types.Bool   `tfsdk:"is_git_related"`
	Name         types.String `tfsdk:"name"`
	PrivateKey   types.String `tfsdk:"private_key"`
	TeamId       types.Int64  `tfsdk:"team_id"`
	Uuid         types.String `tfsdk:"uuid"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

var _ modelWithAttributes = privateKeyModel{}

type privateKeyResourceModel = privateKeyModel
type privateKeyDataSourceModel = privateKeyModel
type privateKeysDataSourceModel struct {
	PrivateKeys types.Set          `tfsdk:"private_keys"`
	Filter      []filterBlockModel `tfsdk:"filter"`
}

func (m privateKeyModel) FromAPI(apiModel *api.PrivateKey) privateKeyModel {
	return privateKeyModel{
		Description:  optionalString(apiModel.Description),
		Fingerprint:  optionalString(apiModel.Fingerprint),
		Id:           optionalInt64(apiModel.Id),
		IsGitRelated: optionalBool(apiModel.IsGitRelated),
		Name:         optionalString(apiModel.Name),
		PrivateKey:   optionalString(apiModel.PrivateKey),
		TeamId:       optionalInt64(apiModel.TeamId),
		Uuid:         optionalString(apiModel.Uuid),
		CreatedAt:    optionalString(apiModel.CreatedAt),
		UpdatedAt:    optionalString(apiModel.UpdatedAt),
	}
}

// Helpers required for types.Set

func (m privateKeyModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description":    types.StringType,
		"fingerprint":    types.StringType,
		"id":             types.Int64Type,
		"is_git_related": types.BoolType,
		"name":           types.StringType,
		"private_key":    types.StringType,
		"team_id":        types.Int64Type,
		"uuid":           types.StringType,
		"created_at":     types.StringType,
		"updated_at":     types.StringType,
	}
}

func (m privateKeyModel) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":    m.Description,
		"fingerprint":    m.Fingerprint,
		"id":             m.Id,
		"is_git_related": m.IsGitRelated,
		"name":           m.Name,
		"private_key":    m.PrivateKey,
		"team_id":        m.TeamId,
		"uuid":           m.Uuid,
		"created_at":     m.CreatedAt,
		"updated_at":     m.UpdatedAt,
	}
}
