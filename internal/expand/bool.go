package expand

import "github.com/hashicorp/terraform-plugin-framework/types"

func Bool(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueBool()
	return &v
}
