package flatten

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestInt64(t *testing.T) {
	input := 42
	assert.Equal(t, types.Int64Null(), Int64(nil))
	assert.Equal(t, types.Int64Value(int64(input)), Int64(&input))
}
