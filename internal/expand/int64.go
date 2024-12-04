package expand

import "github.com/hashicorp/terraform-plugin-framework/types"

func Int64(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := int(value.ValueInt64())
	return &v
}
