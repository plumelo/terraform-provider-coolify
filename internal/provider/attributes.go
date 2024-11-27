package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

type modelWithAttributes = interface {
	// Attributes is required for filtering
	Attributes() map[string]attr.Value
	// AttributeTypes is required for List/Set type parsing
	AttributeTypes() map[string]attr.Type
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

func setResourceDefaultValue(attributes map[string]resource_schema.Attribute, attrName string, defaultValue interface{}) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case resource_schema.StringAttribute:
		typedAttr.Computed = true
		typedAttr.Optional = true
		if defaultValue == nil {
			typedAttr.Computed = false
		} else if strVal, ok := defaultValue.(string); ok {
			typedAttr.Default = stringdefault.StaticString(strVal)
		}
		attributes[attrName] = typedAttr
	case resource_schema.BoolAttribute:
		typedAttr.Computed = true
		typedAttr.Optional = true
		if defaultValue == nil {
			typedAttr.Computed = false
		} else if boolVal, ok := defaultValue.(bool); ok {
			typedAttr.Default = booldefault.StaticBool(boolVal)
		}
		attributes[attrName] = typedAttr
	case resource_schema.Int64Attribute:
		typedAttr.Computed = true
		typedAttr.Optional = true
		if defaultValue == nil {
			typedAttr.Computed = false
		} else if intVal, ok := defaultValue.(int64); ok {
			typedAttr.Default = int64default.StaticInt64(intVal)
		}
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}
