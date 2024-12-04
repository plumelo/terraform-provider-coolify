package flatten

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Time(value *time.Time) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(value.Format(time.RFC3339Nano))
}
