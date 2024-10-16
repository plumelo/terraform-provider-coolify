package resource_private_key

// contentfulManagement "github.com/cysp/terraform-provider-contentful/internal/contentful-management-go"
// "github.com/cysp/terraform-provider-contentful/internal/provider/util"
// "github.com/go-faster/jx"
// "github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
// "github.com/hashicorp/terraform-plugin-framework/diag"
// "github.com/hashicorp/terraform-plugin-framework/path"

// type CreateRequest = api.Eb4780acaa990c594cdbe8ffa80b4fb0JSONRequestBody

// func (model *PrivateKeyModel) ToCreateReq() (CreateRequest, diag.Diagnostics) {
// 	diags := diag.Diagnostics{}

// 	req := CreateRequest{}

// 	switch {
// 	case model.Uuid.IsUnknown():
// 		diags.AddAttributeWarning(path.Root("uuid"), "Failed to update app installation parameters", "Parameters are unknown")
// 	// case model.Uuid.IsNull():
// 	default:
// 		// appInstallationParametersValue := contentfulManagement.PutAppInstallationReqParameters{}
// 		// diags.Append(model.Parameters.Unmarshal(&appInstallationParametersValue)...)
// 		// req.Parameters.SetTo(appInstallationParametersValue)
// 	}

// 	return req, diags
// }

// func (model *PrivateKeyModel) ReadFromResponse(appInstallation *contentfulManagement.AppInstallation) {
// 	// // SpaceId, EnvironmentId and AppDefinitionId are all already known
// 	// if parameters, ok := appInstallation.Parameters.Get(); ok {
// 	// 	encoder := jx.Encoder{}
// 	// 	util.EncodeJxRawMapOrdered(&encoder, parameters)
// 	// 	model.Parameters = jsontypes.NewNormalizedValue(encoder.String())
// 	// } else {
// 	// 	model.Parameters = jsontypes.NewNormalizedNull()
// 	// }
// }
