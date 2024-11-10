package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_private_keys"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &privateKeysDataSource{}
var _ datasource.DataSourceWithConfigure = &privateKeysDataSource{}

func NewPrivateKeysDataSource() datasource.DataSource {
	return &privateKeysDataSource{}
}

type privateKeysDataSource struct {
	providerData CoolifyProviderData
}

type privateKeysDataSourceWithFilterModel struct {
	datasource_private_keys.PrivateKeysModel
	Filter []filterBlockModel `tfsdk:"filter"`
}

var privateKeysFilterNames = []string{"name", "description", "team_id", "is_git_related"}

func (d *privateKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_keys"
}

func (d *privateKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_private_keys.PrivateKeysDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify private keys."

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": createDatasourceFilter(privateKeysFilterNames),
	}

	// Mark sensitive attributes
	if privateKeysSet, ok := resp.Schema.Attributes["private_keys"].(schema.SetNestedAttribute); ok {
		sensitiveAttrs := []string{"private_key"}
		for _, attr := range sensitiveAttrs {
			makeDataSourceAttributeSensitive(privateKeysSet.NestedObject.Attributes, attr)
		}
	}
}

func (d *privateKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *privateKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan privateKeysDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.providerData.client.ListPrivateKeysWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private keys", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private keys. Details: %s", listResponse.Status(), listResponse.Body),
		)
		return
	}

	state, diag := d.apiToModel(ctx, listResponse.JSON200, plan.Filter)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *privateKeysDataSource) apiToModel(
	ctx context.Context,
	privateKeys *[]api.PrivateKey,
	filters []filterBlockModel,
) (privateKeysDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var elements []attr.Value

	for _, pk := range *privateKeys {
		attributes := map[string]attr.Value{
			"created_at":     types.StringValue(*pk.CreatedAt),
			"description":    types.StringValue(*pk.Description),
			"id":             types.Int64Value(int64(*pk.Id)),
			"is_git_related": types.BoolValue(*pk.IsGitRelated),
			"name":           types.StringValue(*pk.Name),
			"private_key":    types.StringValue(*pk.PrivateKey),
			"team_id":        types.Int64Value(int64(*pk.TeamId)),
			"updated_at":     types.StringValue(*pk.UpdatedAt),
			"uuid":           types.StringValue(*pk.Uuid),
		}
		if !filterOnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_private_keys.NewPrivateKeysValue(
			datasource_private_keys.PrivateKeysValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_private_keys.PrivateKeysValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return privateKeysDataSourceWithFilterModel{
		PrivateKeysModel: datasource_private_keys.PrivateKeysModel{
			PrivateKeys: dataSet,
		},
		Filter: filters,
	}, diags
}
