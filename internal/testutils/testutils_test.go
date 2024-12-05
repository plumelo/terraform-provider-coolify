package testutils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAttrTypesFromStruct(t *testing.T) {
	type testStruct struct {
		String    types.String  `tfsdk:"string_field"`
		Int64     types.Int64   `tfsdk:"int64_field"`
		Int32     types.Int32   `tfsdk:"int32_field"`
		Float64   types.Float64 `tfsdk:"float64_field"`
		Float32   types.Float32 `tfsdk:"float32_field"`
		Number    types.Number  `tfsdk:"number_field"`
		Bool      types.Bool    `tfsdk:"bool_field"`
		List      types.List    `tfsdk:"list_field"`
		Map       types.Map     `tfsdk:"map_field"`
		Set       types.Set     `tfsdk:"set_field"`
		Object    types.Object  `tfsdk:"object_field"`
		Dynamic   types.Dynamic `tfsdk:"dynamic_field"`
		NoTag     types.String
		SliceType []string `tfsdk:"slice_field"`
	}

	type unsupportedStruct struct {
		Custom chan int `tfsdk:"custom_field"`
	}

	tests := []struct {
		name          string
		input         any
		expected      map[string]attr.Type
		expectedPanic bool
	}{
		{
			name:  "all_supported_types",
			input: testStruct{},
			expected: map[string]attr.Type{
				"string_field":  types.StringType,
				"int64_field":   types.Int64Type,
				"int32_field":   types.Int32Type,
				"float64_field": types.Float64Type,
				"float32_field": types.Float32Type,
				"number_field":  types.NumberType,
				"bool_field":    types.BoolType,
				"list_field":    types.ListType{},
				"map_field":     types.MapType{},
				"set_field":     types.SetType{},
				"object_field":  types.ObjectType{},
				"dynamic_field": types.DynamicType,
				"slice_field":   types.ListType{ElemType: types.ObjectType{}},
			},
		},
		{
			name:     "empty_struct",
			input:    struct{}{},
			expected: map[string]attr.Type{},
		},
		{
			name: "struct_with_single_field",
			input: struct {
				Name types.String `tfsdk:"name"`
			}{},
			expected: map[string]attr.Type{
				"name": types.StringType,
			},
		},
		{
			name:          "unsupported_type",
			input:         unsupportedStruct{},
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedPanic {
				assert.Panics(t, func() {
					GenerateAttrTypesFromStruct(t, tt.input)
				}, "function should panic with unsupported type")
				return
			}

			got := GenerateAttrTypesFromStruct(t, tt.input)
			assert.Equal(t, tt.expected, got, "attribute types should match")
		})
	}
}
