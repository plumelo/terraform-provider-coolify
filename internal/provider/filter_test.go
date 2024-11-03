package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestAttributeValueToString(t *testing.T) {
	tests := []struct {
		name     string
		input    attr.Value
		expected string
		err      error
	}{
		{"StringValue", basetypes.NewStringValue("test"), "test", nil},
		{"BoolValue", basetypes.NewBoolValue(true), "true", nil},
		{"Int64Value", basetypes.NewInt64Value(42), "42", nil},
		{"Int32Value", basetypes.NewInt32Value(32), "32", nil},
		{"Float64Value", basetypes.NewFloat64Value(3.14), "3.140000", nil},
		{"Float32Value", basetypes.NewFloat32Value(1.23), "1.230000", nil},
		{"UnsupportedType", basetypes.ListValue{}, "", fmt.Errorf("unsupported attribute type: %T", basetypes.ListValue{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := attributeValueToString(tt.input)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterOnAttributes(t *testing.T) {
	tests := []struct {
		name       string
		attributes map[string]attr.Value
		filters    []filterBlockModel
		expected   bool
	}{
		{
			name: "NoFilters",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
			},
			filters:  []filterBlockModel{},
			expected: true,
		},
		{
			name: "MatchingFilter",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field1"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value1")}),
				},
			},
			expected: true,
		},
		{
			name: "NonMatchingFilter",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field1"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value2")}),
				},
			},
			expected: false,
		},
		{
			name: "MissingAttribute",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field2"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value1")}),
				},
			},
			expected: false,
		},
		{
			name: "MultipleFilters",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
				"field2": basetypes.NewStringValue("value2"),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field1"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value1")}),
				},
				{
					Name:   basetypes.NewStringValue("field2"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value2")}),
				},
			},
			expected: true,
		},
		{
			name: "MultipleFiltersNonMatching",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewStringValue("value1"),
				"field2": basetypes.NewStringValue("value2"),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field1"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value1")}),
				},
				{
					Name:   basetypes.NewStringValue("field2"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value3")}),
				},
			},
			expected: false,
		},
		{
			name: "UnsupportedAttributeType",
			attributes: map[string]attr.Value{
				"field1": basetypes.NewListValueMust(types.StringType, nil),
			},
			filters: []filterBlockModel{
				{
					Name:   basetypes.NewStringValue("field1"),
					Values: basetypes.NewListValueMust(types.StringType, []attr.Value{basetypes.NewStringValue("value1")}),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOnAttributes(tt.attributes, tt.filters)
			assert.Equal(t, tt.expected, result)
		})
	}
}
