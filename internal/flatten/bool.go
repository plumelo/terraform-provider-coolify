package flatten

import "github.com/hashicorp/terraform-plugin-framework/types"

func Bool(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*value)
}
