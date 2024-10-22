package provider_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"terraform-provider-coolify/internal/provider"
)

const (
	accEndpoint = "http://localhost:8000/api/v1"
	accToken    = "1|Ey3eX9TNOeUv7W1E5XX6Uf4OJxgq9TPcIFHf7yDbe09e497d"

	providerConfig = `
	provider "coolify" {
		endpoint = "` + accEndpoint + `"
		token 	= "` + accToken + `"
	}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"coolify": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func providerConfigDynamicValue(config map[string]interface{}) (tfprotov6.DynamicValue, error) {
	providerConfigTypes := map[string]tftypes.Type{
		"endpoint": tftypes.String,
		"token":    tftypes.String,
	}
	providerConfigObjectType := tftypes.Object{AttributeTypes: providerConfigTypes}

	providerConfigObjectValue := tftypes.NewValue(providerConfigObjectType, map[string]tftypes.Value{
		"endpoint": tftypes.NewValue(tftypes.String, config["endpoint"]),
		"token":    tftypes.NewValue(tftypes.String, config["token"]),
	})

	value, err := tfprotov6.NewDynamicValue(providerConfigObjectType, providerConfigObjectValue)
	if err != nil {
		err = fmt.Errorf("failed to create dynamic value: %w", err)
	}

	return value, err
}

func TestProtocol6ProviderServerSchemaVersion(t *testing.T) {
	t.Parallel()

	providerServer, err := testAccProtoV6ProviderFactories["coolify"]()
	require.NotNil(t, providerServer)
	require.NoError(t, err)

	resp, err := providerServer.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	require.NotNil(t, resp.Provider)
	require.NoError(t, err)
	assert.Empty(t, resp.Diagnostics)

	assert.EqualValues(t, 0, resp.Provider.Version)
}

func TestProtocol6ProviderServerConfigure(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.SkipNow() // Skip if not running acceptance tests
	}

	tests := map[string]struct {
		config          map[string]interface{}
		env             map[string]string
		expectedSuccess bool
	}{
		"config: endpoint": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
			},
			expectedSuccess: false,
		},
		"config: token": {
			config: map[string]interface{}{
				"token": accToken,
			},
			expectedSuccess: false,
		},
		"config: endpoint,token": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
				"token":    accToken,
			},
			expectedSuccess: true,
		},
		"config: endpoint(invalid),token": {
			config: map[string]interface{}{
				"endpoint": "url://an invalid url %/",
				"token":    accToken,
			},
			expectedSuccess: false,
		},
		"env: endpoint": {
			env: map[string]string{
				"COOLIFY_ENDPOINT": accEndpoint,
			},
			expectedSuccess: false,
		},
		"env: endpoint,token": {
			env: map[string]string{
				"COOLIFY_ENDPOINT": accEndpoint,
				"COOLIFY_TOKEN":    accToken,
			},
			expectedSuccess: true,
		},
		"config: endpoint env: token": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
			},
			env: map[string]string{
				"COOLIFY_TOKEN": accToken,
			},
			expectedSuccess: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for key, value := range test.env {
				t.Setenv(key, value)
			}

			providerServer, err := testAccProtoV6ProviderFactories["coolify"]()
			require.NotNil(t, providerServer)
			require.NoError(t, err)

			providerConfigValue, err := providerConfigDynamicValue(test.config)
			require.NotNil(t, providerConfigValue)
			require.NoError(t, err)

			resp, err := providerServer.ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{
				Config: &providerConfigValue,
			})
			require.NotNil(t, resp)
			require.NoError(t, err)

			if test.expectedSuccess {
				assert.Empty(t, resp.Diagnostics)
			} else {
				assert.NotEmpty(t, resp.Diagnostics)
			}
		})
	}
}

// ---------------------

// func TestAccPreCheck(t *testing.T) {

// 	if v := os.Getenv("COOLIFY_ENDPOINT"); v == "" {
// 		t.Fatal("COOLIFY_ENDPOINT must be set for acceptance tests")
// 	}

// 	if v := os.Getenv("COOLIFY_TOKEN"); v == "" {
// 		t.Fatal("COOLIFY_TOKEN must be set for acceptance tests")
// 	}
// }

// func GetTestAccEndpoint() string {
// 	return os.Getenv("COOLIFY_ENDPOINT")
// }

// func GetTestAccToken() string {
// 	return os.Getenv("COOLIFY_TOKEN")
// }

const TestAccNamePrefix = "tf-acc"

func GetRandomResourceName(resType string) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s-%s-%s", TestAccNamePrefix, resType, string(b))
}
