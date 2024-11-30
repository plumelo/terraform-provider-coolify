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

func TestCombineResourceSchemas(t *testing.T) {
	tests := []struct {
		name     string
		schemas  []res_schema.Schema
		expected res_schema.Schema
	}{
		{
			name:    "empty schemas",
			schemas: []res_schema.Schema{},
			expected: res_schema.Schema{
				Attributes: map[string]res_schema.Attribute{},
				Blocks:     map[string]res_schema.Block{},
			},
		},
		{
			name: "single schema",
			schemas: []res_schema.Schema{
				{
					Description: "Test schema",
					Attributes: map[string]res_schema.Attribute{
						"attr1": res_schema.StringAttribute{Description: "Attribute 1"},
					},
					Blocks: map[string]res_schema.Block{
						"block1": res_schema.SingleNestedBlock{Description: "Block 1"},
					},
				},
			},
			expected: res_schema.Schema{
				Description: "Test schema",
				Attributes: map[string]res_schema.Attribute{
					"attr1": res_schema.StringAttribute{Description: "Attribute 1"},
				},
				Blocks: map[string]res_schema.Block{
					"block1": res_schema.SingleNestedBlock{Description: "Block 1"},
				},
			},
		},
		{
			name: "multiple schemas with overlapping fields",
			schemas: []res_schema.Schema{
				{
					Description: "Schema 1",
					Attributes: map[string]res_schema.Attribute{
						"attr1": res_schema.StringAttribute{Description: "Attribute 1"},
					},
				},
				{
					Description: "Schema 2",
					Attributes: map[string]res_schema.Attribute{
						"attr1": res_schema.StringAttribute{Description: "Attribute 1 Override"},
						"attr2": res_schema.StringAttribute{Description: "Attribute 2"},
					},
				},
			},
			expected: res_schema.Schema{
				Description: "Schema 2",
				Attributes: map[string]res_schema.Attribute{
					"attr1": res_schema.StringAttribute{Description: "Attribute 1 Override"},
					"attr2": res_schema.StringAttribute{Description: "Attribute 2"},
				},
				Blocks: map[string]res_schema.Block{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineResourceSchemas(tt.schemas...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCombineDataSourceSchemas(t *testing.T) {
	tests := []struct {
		name     string
		schemas  []ds_schema.Schema
		expected ds_schema.Schema
	}{
		{
			name:    "empty schemas",
			schemas: []ds_schema.Schema{},
			expected: ds_schema.Schema{
				Attributes: map[string]ds_schema.Attribute{},
				Blocks:     map[string]ds_schema.Block{},
			},
		},
		{
			name: "single schema",
			schemas: []ds_schema.Schema{
				{
					Description: "Test schema",
					Attributes: map[string]ds_schema.Attribute{
						"attr1": ds_schema.StringAttribute{Description: "Attribute 1"},
					},
					Blocks: map[string]ds_schema.Block{
						"block1": ds_schema.SingleNestedBlock{Description: "Block 1"},
					},
				},
			},
			expected: ds_schema.Schema{
				Description: "Test schema",
				Attributes: map[string]ds_schema.Attribute{
					"attr1": ds_schema.StringAttribute{Description: "Attribute 1"},
				},
				Blocks: map[string]ds_schema.Block{
					"block1": ds_schema.SingleNestedBlock{Description: "Block 1"},
				},
			},
		},
		{
			name: "multiple schemas with overlapping fields",
			schemas: []ds_schema.Schema{
				{
					Description: "Schema 1",
					Attributes: map[string]ds_schema.Attribute{
						"attr1": ds_schema.StringAttribute{Description: "Attribute 1"},
					},
				},
				{
					Description: "Schema 2",
					Attributes: map[string]ds_schema.Attribute{
						"attr1": ds_schema.StringAttribute{Description: "Attribute 1 Override"},
						"attr2": ds_schema.StringAttribute{Description: "Attribute 2"},
					},
				},
			},
			expected: ds_schema.Schema{
				Description: "Schema 2",
				Attributes: map[string]ds_schema.Attribute{
					"attr1": ds_schema.StringAttribute{Description: "Attribute 1 Override"},
					"attr2": ds_schema.StringAttribute{Description: "Attribute 2"},
				},
				Blocks: map[string]ds_schema.Block{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineDataSourceSchemas(tt.schemas...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
