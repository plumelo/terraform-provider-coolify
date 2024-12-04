package flatten

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestStringList(t *testing.T) {
	values := []string{"a", "b", "c"}
	assert.Equal(t, types.ListNull(types.StringType), StringList(nil))
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	}), StringList(&values))
}
