package resource_server

import "github.com/hashicorp/terraform-plugin-framework/attr"

func (v *SettingsValue) SetKnown() {
	v.state = attr.ValueStateKnown
}
