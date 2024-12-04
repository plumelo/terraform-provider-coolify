package flatten

import "github.com/hashicorp/terraform-plugin-framework/types"

func Int64(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*value))
}
