package filter

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func TestOnAttributes(t *testing.T) {
	tests := []struct {
		name       string
		attributes map[string]attr.Value
		filters    []BlockModel
		expected   bool
	}{
		{
			name: "NoFilters",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters:  []BlockModel{},
			expected: true,
		},
		{
			name: "MatchingFilter",
			attributes: map[string]attr.Value{
				"field1": types.StringValue("value1"),
			},
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			filters: []BlockModel{
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
			result := OnAttributes(tt.attributes, tt.filters)
			assert.Equal(t, tt.expected, result)
		})
	}
}

type mockStruct struct {
	attributes map[string]attr.Value
}

func (m mockStruct) FilterAttributes() map[string]attr.Value {
	return m.attributes
}

func TestFilterOnStruct(t *testing.T) {
	tests := []struct {
		name     string
		item     FilterableStructModel
		filters  []BlockModel
		expected bool
	}{
		{
			name: "NoFilters",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters:  []BlockModel{},
			expected: true,
		},
		{
			name: "MatchingFilter",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters: []BlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
			},
			expected: true,
		},
		{
			name: "NonMatchingFilter",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters: []BlockModel{
				{
					Name:   types.StringValue("field1"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value2")}),
				},
			},
			expected: false,
		},
		{
			name: "MissingAttribute",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters: []BlockModel{
				{
					Name:   types.StringValue("field2"),
					Values: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("value1")}),
				},
			},
			expected: false,
		},
		{
			name: "MultipleFilters",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
					"field2": types.StringValue("value2"),
				},
			},
			filters: []BlockModel{
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
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
					"field2": types.StringValue("value2"),
				},
			},
			filters: []BlockModel{
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
			name: "MultipleValuesInFilter_OR_Logic",
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters: []BlockModel{
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
			item: mockStruct{
				attributes: map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			},
			filters: []BlockModel{
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
			result := OnStruct(context.Background(), tt.item, tt.filters)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateDatasourceFilter(t *testing.T) {
	allowedFields := []string{"field1", "field2", "field3"}
	block := CreateDatasourceFilter(allowedFields)

	listBlock, ok := block.(schema.ListNestedBlock)
	assert.True(t, ok, "Expected block to be a ListNestedBlock")

	attributes := listBlock.NestedObject.Attributes
	assert.Contains(t, attributes, "name", "Block should contain 'name' attribute")
	assert.Contains(t, attributes, "values", "Block should contain 'values' attribute")
}
