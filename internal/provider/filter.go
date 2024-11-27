package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type filterBlockModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// createDatasourceFilter creates a filter block for a datasource schema.
func createDatasourceFilter(allowedFields []string) schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "Filter results by values",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: fmt.Sprintf("Name of the field to filter on. Valid names are `%s`", strings.Join(allowedFields, "`, `")),
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(allowedFields...),
					},
				},
				"values": schema.ListAttribute{
					Required:            true,
					MarkdownDescription: "List of values to match against - if any value matches, the filter is satisfied (**OR** operation). Non-string values will be converted to strings if possible, ie `true` -> `\"true\"`",
					ElementType:         types.StringType,
				},
			},
		},
	}
}

// filterOnAttributes filters a map of attributes based on a list of filter blocks.
func filterOnAttributes(attributes map[string]attr.Value, filters []filterBlockModel) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		if attr, ok := attributes[filter.Name.ValueString()]; ok {
			attrValueString := attributeValueToString(attr)

			filterValues := []string{}
			filter.Values.ElementsAs(context.Background(), &filterValues, false)
			if !slices.Contains(filterValues, attrValueString) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// attributeValueToString converts any supported attribute value to its string representation.
func attributeValueToString(value attr.Value) string {
	switch v := value.(type) {
	case types.String:
		return v.ValueString()
	case types.Bool:
		return fmt.Sprintf("%t", v.ValueBool())
	case types.Int64:
		return fmt.Sprintf("%d", v.ValueInt64())
	case types.Int32:
		return fmt.Sprintf("%d", v.ValueInt32())
	case types.Float64:
		return fmt.Sprintf("%f", v.ValueFloat64())
	case types.Float32:
		return fmt.Sprintf("%f", v.ValueFloat32())
	case types.Number:
		return fmt.Sprintf("%f", v.ValueBigFloat())
	case types.Dynamic:
		if underlyingValue := v.UnderlyingValue(); underlyingValue != nil {
			return attributeValueToString(underlyingValue)
		}
	}

	// Fall back to Terraform's string representation
	return value.String()
}

type structWithAttributeMap interface {
	Attributes() map[string]attr.Value
}

func filterOnStruct(
	ctx context.Context,
	item structWithAttributeMap,
	filters []filterBlockModel,
) bool {
	if len(filters) == 0 {
		return true
	}

	attributes := item.Attributes()

	for _, filter := range filters {
		filterName := filter.Name.ValueString()
		filterValues := []string{}
		filter.Values.ElementsAs(ctx, &filterValues, false)

		attrValue, ok := attributes[filterName]
		if !ok {
			return false
		}

		attrValueString := attributeValueToString(attrValue)
		if !slices.Contains(filterValues, attrValueString) {
			return false
		}
	}
	return true
}
