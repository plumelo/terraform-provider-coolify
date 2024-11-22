package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func TestOptionalStringListValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]string
		expected types.List
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: types.ListNull(types.StringType),
		},
		{
			name:     "empty array input",
			input:    &[]string{},
			expected: types.ListValueMust(types.StringType, []attr.Value{}),
		},
		{
			name:  "populated array input",
			input: &[]string{"one", "two", "three"},
			expected: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("one"),
				types.StringValue("two"),
				types.StringValue("three"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, optionalStringListValue(tt.input))
		})
	}
}

func TestOptionalTime(t *testing.T) {
	now := time.Now()
	formattedTime := now.Format(time.RFC3339Nano)

	tests := []struct {
		name     string
		input    *time.Time
		expected types.String
	}{
		{"nil input", nil, types.StringNull()},
		{"valid time input", &now, types.StringValue(formattedTime)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, optionalTime(tt.input))
		})
	}
}
func TestBase64EncodeAttr(t *testing.T) {
	tests := []struct {
		name     string
		input    types.String
		expected *string
	}{
		{"null value", types.StringNull(), nil},
		{"unknown value", types.StringUnknown(), nil},
		{"empty string", types.StringValue(""), &[]string{""}[0]},
		{"simple string", types.StringValue("hello"), &[]string{"aGVsbG8="}[0]},
		{"string with special characters", types.StringValue("hello@world!123"), &[]string{"aGVsbG9Ad29ybGQhMTIz"}[0]},
		{"unicode string", types.StringValue("こんにちは"), &[]string{"44GT44KT44Gr44Gh44Gv"}[0]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base64EncodeAttr(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestBase64DecodeAttr(t *testing.T) {
	tests := []struct {
		name     string
		input    types.String
		expected *string
	}{
		{"null value", types.StringNull(), nil},
		{"unknown value", types.StringUnknown(), nil},
		{"invalid base64", types.StringValue("!@#$"), nil},
		{"simple string", types.StringValue("YWJj"), &[]string{"abc"}[0]},
		{"empty string", types.StringValue(""), &[]string{""}[0]},
		{"string with special characters", types.StringValue("IUAjJCVeJiooKV8r"), &[]string{"!@#$%^&*()_+"}[0]},
		{"unicode string", types.StringValue("5pel5pys6Kqe"), &[]string{"日本語"}[0]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base64DecodeAttr(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}
