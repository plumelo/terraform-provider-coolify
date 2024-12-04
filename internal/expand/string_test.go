package expand

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTfString(t *testing.T) {
	input := types.StringValue("test")
	assert.Nil(t, String(types.StringNull()))
	assert.Nil(t, String(types.StringUnknown()))
	assert.Equal(t, input.ValueString(), *String(input))
}

func TestTfRequiredString(t *testing.T) {
	input := types.StringValue("test")
	assert.Equal(t, "", RequiredString(types.StringNull()))
	assert.Equal(t, "", RequiredString(types.StringUnknown()))
	assert.Equal(t, input.ValueString(), RequiredString(input))
}
