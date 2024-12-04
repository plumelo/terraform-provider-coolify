package flatten

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	input := true
	assert.Equal(t, types.BoolNull(), Bool(nil))
	assert.Equal(t, types.BoolValue(input), Bool(&input))
}
