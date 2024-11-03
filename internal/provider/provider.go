package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &CoolifyProvider{}
var _ provider.ProviderWithFunctions = &CoolifyProvider{}
var _ provider.ProviderWithEphemeralResources = &CoolifyProvider{}

// CoolifyProvider defines the provider implementation.
type CoolifyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CoolifyProviderData struct {
	endpoint string
	client   *api.ClientWithResponses
}

type CoolifyProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CoolifyProvider{
			version: version,
		}
	}
}

func (p *CoolifyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "coolify"
	resp.Version = p.version
}

func (p *CoolifyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	hasEnvToken := os.Getenv("COOLIFY_TOKEN") != ""
	resp.Schema = schema.Schema{
		MarkdownDescription: "" +
			"The \"coolify\" provider facilitates interaction with resources supported by [Coolify](https://coolify.io/). " +
			"Before using this provider, you must configure it with your credentials, typically by setting the environment variable `COOLIFY_TOKEN`. " +
			"For instructions on obtaining an API token, refer to Coolify's [API documentation](https://coolify.io/docs/api-reference/authorization#generate).",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Coolify endpoint. If not set, checks env for `COOLIFY_ENDPOINT`. Default: `https://app.coolify.io/api/v1`",
			},
			"token": schema.StringAttribute{
				Required:            !hasEnvToken,
				Optional:            hasEnvToken,
				Sensitive:           true,
				MarkdownDescription: "Coolify token. If not set, checks env for `COOLIFY_TOKEN`.",
			},
		},
	}
}

func (p *CoolifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CoolifyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiEndpoint string
	if !data.Endpoint.IsNull() {
		apiEndpoint = data.Endpoint.ValueString()
	} else if apiEndpointFromEnv, found := os.LookupEnv("COOLIFY_ENDPOINT"); found {
		apiEndpoint = apiEndpointFromEnv
	} else {
		apiEndpoint = "https://app.coolify.io/api/v1"
	}

	if apiEndpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "Failed to configure client", "No API Endpoint provided")
	}

	var apiToken string
	if !data.Token.IsNull() {
		apiToken = data.Token.ValueString()
	} else {
		if apiTokenFromEnv, found := os.LookupEnv("COOLIFY_TOKEN"); found {
			apiToken = apiTokenFromEnv
		}
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(path.Root("token"), "Failed to configure client", "No token provided")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := api.NewAPIClient(p.version, apiEndpoint, apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create API client",
			err.Error(),
		)
		return
	}

	// GET /version
	versionResp, err := client.VersionWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to connect to Coolify API",
			err.Error(),
		)
		return
	}

	if versionResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code API client",
			fmt.Sprintf("Received %s creating API client. Details: %s", versionResp.Status(), versionResp.Body),
		)
		return
	}

	currentVersion := string(versionResp.Body)
	minVersion := "4.0.0-beta.360"

	if !isVersionCompatible(currentVersion, minVersion) {
		resp.Diagnostics.AddError(
			"Unsupported API version",
			fmt.Sprintf("The Coolify API version %s is not supported. The minimum supported version is %s", currentVersion, minVersion),
		)
		return
	}

	tflog.Info(ctx, "Successfully connected to Coolify API", map[string]interface{}{"version": currentVersion})

	providerData := &CoolifyProviderData{
		endpoint: apiEndpoint,
		client:   client,
	}

	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

func (p *CoolifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPrivateKeyResource,
	}
}

func (p *CoolifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *CoolifyProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *CoolifyProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}
