package provider

import (
	"encoding/base64"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func optionalTime(value *time.Time) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.Format(time.RFC3339Nano))
}

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

// optionalStringListValue converts a list of strings to a ListValue.
func optionalStringListValue(values *[]string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(*values))
	for i, v := range *values {
		elems[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
}

func base64Encode(value *string) *string {
	if value == nil {
		return nil
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(*value))
	return &encoded
}

func base64EncodeAttr(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}
	return base64Encode(value.ValueStringPointer())
}

func base64Decode(value *string) *string {
	if value == nil {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(*value)
	if err != nil {
		return nil
	}
	decodedStr := string(decoded)
	return &decodedStr
}

func base64DecodeAttr(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}
	return base64Decode(value.ValueStringPointer())
}
