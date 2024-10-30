package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestOptionalString(t *testing.T) {
	input := "test"
	assert.Equal(t, types.StringNull(), optionalString(nil))
	assert.Equal(t, types.StringValue(input), optionalString(&input))
}

func TestOptionalInt64(t *testing.T) {
	input := 42
	assert.Equal(t, types.Int64Null(), optionalInt64(nil))
	assert.Equal(t, types.Int64Value(int64(input)), optionalInt64(&input))
}

func TestOptionalBool(t *testing.T) {
	input := true
	assert.Equal(t, types.BoolNull(), optionalBool(nil))
	assert.Equal(t, types.BoolValue(input), optionalBool(&input))
}
