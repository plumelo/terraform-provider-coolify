package provider

import (
	"encoding/base64"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	res_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func tfStringToOptionalString(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	return value.ValueStringPointer()
}

func tfStringToRequiredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

func tfBoolToOptionalBool(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	return value.ValueBoolPointer()
}

func tfBoolToRequiredBool(value types.Bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return false
	}
	return value.ValueBool()
}

func tfInt64ToOptionalInt(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := int(value.ValueInt64())
	return &v
}

func tfInt64ToRequiredInt(value types.Int64) int {
	if value.IsNull() || value.IsUnknown() {
		return 0
	}
	return int(value.ValueInt64())
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

// mergeResourceSchemas combines multiple resource schemas by merging their attributes and blocks.
// If an attribute or block exists in multiple schemas, the last one takes precedence.
func mergeResourceSchemas(schemas ...res_schema.Schema) res_schema.Schema {
	result := res_schema.Schema{
		Attributes: make(map[string]res_schema.Attribute),
		Blocks:     make(map[string]res_schema.Block),
	}

	// Merge/overwrite attributes and blocks from all schemas
	for _, s := range schemas {
		result.Description = s.Description
		result.MarkdownDescription = s.MarkdownDescription
		result.DeprecationMessage = s.DeprecationMessage
		result.Version = s.Version
		for name, attr := range s.Attributes {
			result.Attributes[name] = attr
		}
		for name, block := range s.Blocks {
			result.Blocks[name] = block
		}
	}

	return result
}

// mergeDataSourceSchemas combines multiple datasource schemas by merging their attributes and blocks.
// If an attribute or block exists in multiple schemas, the last one takes precedence.
func mergeDataSourceSchemas(schemas ...ds_schema.Schema) ds_schema.Schema {
	result := ds_schema.Schema{
		Attributes: make(map[string]ds_schema.Attribute),
		Blocks:     make(map[string]ds_schema.Block),
	}

	// Merge/overwrite attributes and blocks from all schemas
	for _, s := range schemas {
		result.Description = s.Description
		result.MarkdownDescription = s.MarkdownDescription
		result.DeprecationMessage = s.DeprecationMessage
		for name, attr := range s.Attributes {
			result.Attributes[name] = attr
		}
		for name, block := range s.Blocks {
			result.Blocks[name] = block
		}
	}

	return result
}
