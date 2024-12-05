package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	res_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

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
			result := mergeResourceSchemas(tt.schemas...)
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
			result := mergeDataSourceSchemas(tt.schemas...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// generateAttrTypesFromStruct is a helper function for doing attribute comparisons during testing
func generateAttrTypesFromStruct(t *testing.T, structType any) map[string]attr.Type {
	t.Helper()

	attrTypes := make(map[string]attr.Type)
	refT := reflect.TypeOf(structType)

	for i := 0; i < refT.NumField(); i++ {
		field := refT.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" {
			continue
		}

		if field.Type.Kind() == reflect.Slice {
			attrTypes[tag] = types.ListType{ElemType: types.ObjectType{}}
			continue
		}

		switch field.Type {
		case reflect.TypeOf(types.Dynamic{}):
			attrTypes[tag] = types.DynamicType
		case reflect.TypeOf(types.String{}):
			attrTypes[tag] = types.StringType
		case reflect.TypeOf(types.Int32{}):
			attrTypes[tag] = types.Int32Type
		case reflect.TypeOf(types.Int64{}):
			attrTypes[tag] = types.Int64Type
		case reflect.TypeOf(types.Float32{}):
			attrTypes[tag] = types.Float32Type
		case reflect.TypeOf(types.Float64{}):
			attrTypes[tag] = types.Float64Type
		case reflect.TypeOf(types.Number{}):
			attrTypes[tag] = types.NumberType
		case reflect.TypeOf(types.Bool{}):
			attrTypes[tag] = types.BoolType
		case reflect.TypeOf(types.Set{}):
			attrTypes[tag] = types.SetType{}
		case reflect.TypeOf(types.List{}):
			attrTypes[tag] = types.ListType{}
		case reflect.TypeOf(types.Map{}):
			attrTypes[tag] = types.MapType{}
		case reflect.TypeOf(types.Object{}):
			attrTypes[tag] = types.ObjectType{}
		default:
			panic(fmt.Sprintf("unsupported type %s", field.Type))
		}
	}

	return attrTypes
}
