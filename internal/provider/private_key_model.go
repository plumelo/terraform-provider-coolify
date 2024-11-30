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

var _ filterableStructModel = privateKeyModel{}

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

var privateKeysFilterNames = []string{"name", "description", "team_id", "is_git_related"}

func (m privateKeyModel) FilterAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":    m.Description,
		"is_git_related": m.IsGitRelated,
		"name":           m.Name,
		"team_id":        m.TeamId,
	}
}
