package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func TestMakeResourceAttributeRequired(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]resource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
		{
			name: "bool attribute",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:    "bool_attr",
			expectedErr: "",
		},
		{
			name: "int64 attribute",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:    "int64_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeResourceAttributeRequired(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case resource_schema.StringAttribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				case resource_schema.BoolAttribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				case resource_schema.Int64Attribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				}
			}
		})
	}
}

func TestMakeResourceAttributeSensitive(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]resource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeResourceAttributeSensitive(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case datasource_schema.StringAttribute:
					assert.True(t, typedAttr.Sensitive)
				}
			}
		})
	}
}

func TestMakeDataSourceAttributeSensitive(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]datasource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]datasource_schema.Attribute{
				"existing_attr": datasource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]datasource_schema.Attribute{
				"unsupported_attr": datasource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]datasource_schema.Attribute{
				"string_attr": datasource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeDataSourceAttributeSensitive(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case datasource_schema.StringAttribute:
					assert.True(t, typedAttr.Sensitive)
				}
			}
		})
	}
}

