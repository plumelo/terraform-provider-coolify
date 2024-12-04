package expand

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTfInt64(t *testing.T) {
	input := 42
	assert.Nil(t, Int64(types.Int64Null()))
	assert.Nil(t, Int64(types.Int64Unknown()))
	assert.Equal(t, &input, Int64(types.Int64Value(int64(input))))
}
