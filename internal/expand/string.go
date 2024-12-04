package expand

import "github.com/hashicorp/terraform-plugin-framework/types"

func String(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	return value.ValueStringPointer()
}

func RequiredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}

	return value.ValueString()
}
