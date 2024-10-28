package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func optionalString(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}

func optionalInt64(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

func optionalBool(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

// MARK: Old helpers - to be removed once deemed unnecessary

// type ValueType interface {
// 	IsUnknown() bool
// 	IsNull() bool
// }

// func assignGeneric[S string | int | int64 | bool, T ValueType](source *S, target T) {
// 	// if target.IsUnknown() {
// 	switch v := any(target).(type) {
// 	case *basetypes.StringValue:
// 		if strSource, ok := any(source).(*string); ok {
// 			if strSource != nil {
// 				*v = types.StringValue(*strSource)
// 			} else {
// 				*v = types.StringNull()
// 			}
// 		}
// 	case *basetypes.Int64Value:
// 		if intSource, ok := any(source).(*int); ok {
// 			if intSource != nil {
// 				*v = types.Int64Value(int64(*intSource))
// 			} else {
// 				*v = types.Int64Null()
// 			}
// 		}
// 	case *basetypes.BoolValue:
// 		if boolSource, ok := any(source).(*bool); ok {
// 			if boolSource != nil {
// 				*v = types.BoolValue(*boolSource)
// 			} else {
// 				*v = types.BoolNull()
// 			}
// 		}
// 		// case *basetypes.Float64Value:
// 		// 	if floatSource, ok := any(source).(*float32); ok {
// 		// 		if floatSource != nil {
// 		// 			*v = types.Float64Value(float64(*floatSource))
// 		// 		} else {
// 		// 			*v = types.Float64Null()
// 		// 		}
// 		// 	}
// 		// }
// 	}
// 	// }
// }

func assignStr(source *string, target *basetypes.StringValue) {
	if source != nil {
		*target = types.StringValue(*source)
	} else {
		*target = types.StringNull()
	}
}

func attrStr(source *string) attr.Value {
	if source != nil {
		return types.StringValue(*source)
	}

	return types.StringNull()
}

func assignInt(source *int, target *basetypes.Int64Value) {
	if source != nil {
		*target = types.Int64Value(int64(*source))
	} else {
		*target = types.Int64Null()
	}
}

func attrInt(source *int) attr.Value {
	if source != nil {
		return types.Int64Value(int64(*source))
	}

	return types.Int64Null()
}

// func assignFloat(source *float32, target *basetypes.Float64Value) {
// 	if source != nil {
// 		*target = types.Float64Value(float64(*source))
// 	} else {
// 		*target = types.Float64Null()
// 	}
// }

// func attrFloat(source *float32) attr.Value {
// 	if source != nil {
// 		return types.Float64Value(float64(*source))
// 	}

// 	return types.Float64Null()
// }

func assignBool(source *bool, target *basetypes.BoolValue) {
	if source != nil {
		*target = types.BoolValue(bool(*source))
	} else {
		*target = types.BoolNull()
	}
}

// func attrBool(source *bool) attr.Value {
// 	if source != nil {
// 		return types.BoolValue(bool(*source))
// 	}

// 	return types.BoolNull()
// }

func intPointerToInt64Pointer(val *int) *int64 {
	if val == nil {
		return nil
	}

	v := int64(*val)
	return &v
}

func int64ToIntPointer(val basetypes.Int64Value) *int {
	if val.IsNull() {
		return nil
	}

	v := int(val.ValueInt64())
	return &v
}
