package flatten

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestStringSet(t *testing.T) {
	values := []string{"a", "b", "c"}
	assert.Equal(t, types.SetNull(types.StringType), StringSet(nil))
	assert.Equal(t, types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	}), StringSet(&values))
}
