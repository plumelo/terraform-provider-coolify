package flatten

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringSet(values *[]string) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, len(*values))
	for i, v := range *values {
		elems[i] = types.StringValue(v)
	}

	return types.SetValueMust(types.StringType, elems)
}
