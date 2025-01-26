package service

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/filter"
	"terraform-provider-coolify/internal/flatten"
)

type teamModel struct {
	CreatedAt         types.String      `tfsdk:"created_at"`
	CustomServerLimit types.String      `tfsdk:"custom_server_limit"`
	Description       types.String      `tfsdk:"description"`
	Id                types.Int64       `tfsdk:"id"`
	Members           []teamMemberModel `tfsdk:"members"`
	Name              types.String      `tfsdk:"name"`
	PersonalTeam      types.Bool        `tfsdk:"personal_team"`
	ShowBoarding      types.Bool        `tfsdk:"show_boarding"`
	UpdatedAt         types.String      `tfsdk:"updated_at"`
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

var _ filter.FilterableStructModel = teamModel{}

func (m teamModel) FromAPI(apiModel *api.Team) teamModel {
	var members []teamMemberModel

	if apiModel.Members != nil {
		members = make([]teamMemberModel, len(*apiModel.Members))
		for i, member := range *apiModel.Members {
			members[i] = teamMemberModel{
				CreatedAt:            flatten.String(member.CreatedAt),
				Email:                flatten.String(member.Email),
				EmailVerifiedAt:      flatten.String(member.EmailVerifiedAt),
				ForcePasswordReset:   flatten.Bool(member.ForcePasswordReset),
				Id:                   flatten.Int64(member.Id),
				MarketingEmails:      flatten.Bool(member.MarketingEmails),
				Name:                 flatten.String(member.Name),
				TwoFactorConfirmedAt: flatten.String(member.TwoFactorConfirmedAt),
				UpdatedAt:            flatten.String(member.UpdatedAt),
			}
		}
	}

	return teamModel{
		CreatedAt:         flatten.String(apiModel.CreatedAt),
		CustomServerLimit: flatten.String(apiModel.CustomServerLimit),
		Description:       flatten.String(apiModel.Description),
		Id:                flatten.Int64(apiModel.Id),
		Members:           members,
		Name:              flatten.String(apiModel.Name),
		PersonalTeam:      flatten.Bool(apiModel.PersonalTeam),
		ShowBoarding:      flatten.Bool(apiModel.ShowBoarding),
		UpdatedAt:         flatten.String(apiModel.UpdatedAt),
	}
}

var teamsFilterNames = []string{"name", "description", "id", "personal_team"}

func (m teamModel) FilterAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":   m.Description,
		"id":            m.Id,
		"name":          m.Name,
		"personal_team": m.PersonalTeam,
	}
}
