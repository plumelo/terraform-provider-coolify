package testutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generateAttrTypesFromStruct is a helper function for doing attribute comparisons during testing
func GenerateAttrTypesFromStruct(t *testing.T, structType any) map[string]attr.Type {
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
