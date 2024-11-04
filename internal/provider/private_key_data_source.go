package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_private_key"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &privateKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &privateKeyDataSource{}

func NewPrivateKeyDataSource() datasource.DataSource {
	return &privateKeyDataSource{}
}

type privateKeyDataSource struct {
	providerData CoolifyProviderData
}

func (d *privateKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (d *privateKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_private_key.PrivateKeyDataSourceSchema(ctx)
	resp.Schema.Description = "Get a single Coolify private key by UUID."

	// Override the private_key.private_key attribute to make it sensitive
	if privateKeyAttr, ok := resp.Schema.Attributes["private_key"].(schema.StringAttribute); ok {
		privateKeyAttr.Sensitive = true
		resp.Schema.Attributes["private_key"] = privateKeyAttr
	}
}

func (d *privateKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *privateKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_private_key.PrivateKeyModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	privateKey, err := d.providerData.client.GetPrivateKeyByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private key", err.Error(),
		)
		return
	}

	if privateKey.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key. Details: %s", privateKey.Status(), privateKey.Body),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, privateKey.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *privateKeyDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.PrivateKey,
) datasource_private_key.PrivateKeyModel {
	return datasource_private_key.PrivateKeyModel{
		Id:           optionalInt64(response.Id),
		Uuid:         optionalString(response.Uuid),
		Name:         optionalString(response.Name),
		Description:  optionalString(response.Description),
		PrivateKey:   optionalString(response.PrivateKey),
		IsGitRelated: optionalBool(response.IsGitRelated),
		TeamId:       optionalInt64(response.TeamId),
		CreatedAt:    optionalString(response.CreatedAt),
		UpdatedAt:    optionalString(response.UpdatedAt),
	}
}
