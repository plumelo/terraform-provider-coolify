package flatten

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	input := "test"
	assert.Equal(t, types.StringNull(), String(nil))
	assert.Equal(t, types.StringValue(input), String(&input))
}

func TestRequiredString(t *testing.T) {
	input := "test"
	assert.Equal(t, types.StringNull(), RequiredString(""))
	assert.Equal(t, types.StringValue(input), RequiredString(input))
}
