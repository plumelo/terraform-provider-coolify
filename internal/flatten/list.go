package flatten

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringList(values *[]string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(*values))
	for i, v := range *values {
		elems[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
}
