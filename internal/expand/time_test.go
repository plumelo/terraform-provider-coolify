package expand

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTfTime(t *testing.T) {
	now := time.Now()
	formattedTime := now.Format(time.RFC3339Nano)
	assert.Nil(t, Time(types.StringNull()))
	assert.Nil(t, Time(types.StringUnknown()))
	assert.Nil(t, Time(types.StringValue("invalid-time-format")))
	assert.Equal(t, formattedTime, Time(types.StringValue(formattedTime)).Format(time.RFC3339Nano))
}
