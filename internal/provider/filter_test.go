package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAttributeValueToString(t *testing.T) {
	tests := []struct {
		name     string
		input    attr.Value
		expected string
	}{
		{"StringValue", types.StringValue("test"), "test"},
		{"BoolValue", types.BoolValue(true), "true"},
		{"Int64Value", types.Int64Value(42), "42"},
		{"Int32Value", types.Int32Value(32), "32"},
		{"Float64Value", types.Float64Value(3.14), "3.140000"},
		{"Float32Value", types.Float32Value(1.23), "1.230000"},
		{"NumberValue", types.NumberValue(big.NewFloat(1.23)), "1.230000"},
		{"DynamicValue-String", types.DynamicValue(types.StringValue("test")), "test"},
		{"DynamicValue-Int64", types.DynamicValue(types.Int64Value(42)), "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attributeValueToString(tt.input)
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
				"field1": types.StringValue("value1"),
			},
			filters:  []filterBlockModel{},
			expected: true,
		},
		{
			name: "MatchingFilter",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
			},
			expected: true,
		},
		{
			name: "NonMatchingFilter",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value2")}),
				},
			},
			expected: false,
		},
		{
			name: "MissingAttribute",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field2"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
			},
			expected: false,
		},
		{
			name: "MultipleFilters",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
				"field2": types.StringValue("value2"),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
				{
					Name:   types.StringValue("field2"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value2")}),
				},
			},
			expected: true,
		},
		{
			name: "MultipleFiltersNonMatching",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
				"field2": types.StringValue("value2"),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
				{
					Name:   types.StringValue("field2"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value3")}),
				},
			},
			expected: false,
		},
		{
			name: "UnsupportedAttributeType",
			attributes: map[string]attr.Value{
				"field1": types.ListValueMust(types.StringType, nil),
			},
			filters: []filterBlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
			},
			expected: false,
		},
		{
			name: "MultipleValuesInFilter_OR_Logic",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name: types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("value2"),
						types.StringValue("value1"),
						types.StringValue("value3"),
					}),
				},
			},
			expected: true,
		},
		{
			name: "MultipleValuesInFilter_NoMatch",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []filterBlockModel{
				{
					Name: types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("value2"),
						types.StringValue("value3"),
					}),
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
