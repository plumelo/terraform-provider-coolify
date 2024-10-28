package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/resource_private_key"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &privateKeyResource{}
	_ resource.ResourceWithConfigure   = &privateKeyResource{}
	_ resource.ResourceWithImportState = &privateKeyResource{}
)

func NewPrivateKeyResource() resource.Resource {
	return &privateKeyResource{}
}

type privateKeyResource struct {
	providerData CoolifyProviderData
}

func (r *privateKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (r *privateKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_private_key.PrivateKeyResourceSchema(ctx)
}

func (r *privateKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *privateKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating private key", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})
	createResp, err := r.providerData.client.CreatePrivateKeyWithResponse(ctx, api.CreatePrivateKeyJSONRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating private key",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating private key",
			fmt.Sprintf("Received %s creating private key. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading private key", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_private_key.PrivateKeyModel
	var state resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	if uuid == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// Update API call logic
	tflog.Debug(ctx, "Updating private key", map[string]interface{}{
		"uuid": uuid,
	})
	updateResp, err := r.providerData.client.UpdatePrivateKeyWithResponse(ctx, uuid, api.UpdatePrivateKeyJSONRequestBody{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating private key: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating private key",
			fmt.Sprintf("Received %s updating private key: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	// TODO: BUG: All computed fields are being recalculated, even if they are not updated. May have to do with this.
	data := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_private_key.PrivateKeyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting private key", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.providerData.client.DeletePrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete private key, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting private key",
			fmt.Sprintf("Received %s deleting private key: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *privateKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *privateKeyResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) resource_private_key.PrivateKeyModel {
	readResp, err := r.providerData.client.GetPrivateKeyByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading private key: uuid=%s", uuid),
			err.Error(),
		)
		return resource_private_key.PrivateKeyModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return resource_private_key.PrivateKeyModel{}
	}

	return r.ApiToModel(ctx, diags, readResp.JSON200)
}

func (r *privateKeyResource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.PrivateKey,
) resource_private_key.PrivateKeyModel {
	return resource_private_key.PrivateKeyModel{
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
