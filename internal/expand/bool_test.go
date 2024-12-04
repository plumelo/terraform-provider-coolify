package expand

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTfBool(t *testing.T) {
	input := types.BoolValue(true)
	assert.Nil(t, Bool(types.BoolNull()))
	assert.Nil(t, Bool(types.BoolUnknown()))
	assert.Equal(t, input.ValueBool(), *Bool(input))
}
