package provider

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

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

func TestSetResourceDefaultValue(t *testing.T) {
	tests := []struct {
		name         string
		attributes   map[string]resource_schema.Attribute
		attrName     string
		defaultValue interface{}
		expectedErr  string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:     "missing_attr",
			defaultValue: "default",
			expectedErr:  "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:     "unsupported_attr",
			defaultValue: nil,
			expectedErr:  "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: "default",
			expectedErr:  "",
		},
		{
			name: "bool attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:     "bool_attr",
			defaultValue: true,
			expectedErr:  "",
		},
		{
			name: "int64 attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:     "int64_attr",
			defaultValue: int64(42),
			expectedErr:  "",
		},
		{
			name: "string attribute with wrong default type",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: 123, // wrong type
			expectedErr:  "",
		},
		{
			name: "string attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
		{
			name: "bool attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:     "bool_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
		{
			name: "int64 attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:     "int64_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setResourceDefaultValue(tt.attributes, tt.attrName, tt.defaultValue)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case resource_schema.StringAttribute:
					if strVal, ok := tt.defaultValue.(string); ok {
						assert.Equal(t, stringdefault.StaticString(strVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				case resource_schema.BoolAttribute:
					if boolVal, ok := tt.defaultValue.(bool); ok {
						assert.Equal(t, booldefault.StaticBool(boolVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				case resource_schema.Int64Attribute:
					if intVal, ok := tt.defaultValue.(int64); ok {
						assert.Equal(t, int64default.StaticInt64(intVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				default:
					t.Errorf("Unexpected attribute type")
				}
			}
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
