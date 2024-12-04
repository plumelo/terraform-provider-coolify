package flatten

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	now := time.Now()
	formattedTime := now.Format(time.RFC3339Nano)
	assert.Equal(t, types.StringNull(), Time(nil))
	assert.Equal(t, types.StringValue(formattedTime), Time(&now))
}
