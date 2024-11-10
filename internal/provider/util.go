package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// optionalStringListValue converts a list of strings to a ListValue.
func optionalStringListValue(values *[]string) basetypes.ListValue {
	if values == nil {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(*values))
	for i, v := range *values {
		elems[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
}

func makeResourceAttributeRequired(
	attributes map[string]resource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case resource_schema.StringAttribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	case resource_schema.BoolAttribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	case resource_schema.Int64Attribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}

func makeResourceAttributeSensitive(
	attributes map[string]resource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case resource_schema.StringAttribute:
		typedAttr.Sensitive = true
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}

func makeDataSourceAttributeSensitive(
	attributes map[string]datasource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case datasource_schema.StringAttribute:
		typedAttr.Sensitive = true
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}
