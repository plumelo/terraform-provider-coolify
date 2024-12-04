package expand

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Time(value types.String) *time.Time {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v, err := time.Parse(time.RFC3339Nano, value.ValueString())
	if err != nil {
		return nil
	}

	return &v
}
