package util_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-coolify/internal/provider/util"
)

type TestProviderData struct {
	Value string
}

func TestProviderDataFromDataSourceConfigureRequest(t *testing.T) {
	tests := []struct {
		name     string
		req      datasource.ConfigureRequest
		wantData TestProviderData
		wantBool bool
		wantDiag bool
	}{
		{
			name:     "Valid provider data",
			req:      datasource.ConfigureRequest{ProviderData: &TestProviderData{Value: "test"}},
			wantData: TestProviderData{Value: "test"},
			wantBool: true,
		},
		{
			name:     "Nil provider data",
			req:      datasource.ConfigureRequest{ProviderData: nil},
			wantBool: false,
		},
		{
			name:     "Invalid provider data type",
			req:      datasource.ConfigureRequest{ProviderData: "invalid"},
			wantDiag: true,
		},
		{
			name:     "Empty provider data",
			req:      datasource.ConfigureRequest{ProviderData: &TestProviderData{}},
			wantBool: true,
		},
		{
			name:     "Different provider data type",
			req:      datasource.ConfigureRequest{ProviderData: &struct{ OtherValue int }{OtherValue: 42}},
			wantDiag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotData TestProviderData
			resp := &datasource.ConfigureResponse{Diagnostics: diag.Diagnostics{}}

			gotBool := util.ProviderDataFromDataSourceConfigureRequest(tt.req, &gotData, resp)

			if gotBool != tt.wantBool {
				t.Errorf("ProviderDataFromDataSourceConfigureRequest() returned %v, want %v", gotBool, tt.wantBool)
			}

			if gotData != tt.wantData {
				t.Errorf("ProviderDataFromDataSourceConfigureRequest() set data to %v, want %v", gotData, tt.wantData)
			}

			if (len(resp.Diagnostics) > 0) != tt.wantDiag {
				t.Errorf("ProviderDataFromDataSourceConfigureRequest() diagnostic presence: got %v, want %v", len(resp.Diagnostics) > 0, tt.wantDiag)
			}
		})
	}
}

func TestProviderDataFromResourceConfigureRequest(t *testing.T) {
	tests := []struct {
		name     string
		req      resource.ConfigureRequest
		wantData TestProviderData
		wantBool bool
		wantDiag bool
	}{
		{
			name:     "Valid resource provider data",
			req:      resource.ConfigureRequest{ProviderData: &TestProviderData{Value: "resource_test"}},
			wantData: TestProviderData{Value: "resource_test"},
			wantBool: true,
		},
		{
			name:     "Nil resource provider data",
			req:      resource.ConfigureRequest{ProviderData: nil},
			wantBool: false,
		},
		{
			name:     "Invalid resource provider data type",
			req:      resource.ConfigureRequest{ProviderData: "invalid_resource"},
			wantDiag: true,
		},
		{
			name:     "Empty resource provider data",
			req:      resource.ConfigureRequest{ProviderData: &TestProviderData{}},
			wantBool: true,
		},
		{
			name:     "Different resource provider data type",
			req:      resource.ConfigureRequest{ProviderData: &struct{ ResourceValue int }{ResourceValue: 100}},
			wantDiag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotData TestProviderData
			resp := &resource.ConfigureResponse{Diagnostics: diag.Diagnostics{}}

			gotBool := util.ProviderDataFromResourceConfigureRequest(tt.req, &gotData, resp)

			if gotBool != tt.wantBool {
				t.Errorf("ProviderDataFromResourceConfigureRequest() returned %v, want %v", gotBool, tt.wantBool)
			}

			if gotData != tt.wantData {
				t.Errorf("ProviderDataFromResourceConfigureRequest() set data to %v, want %v", gotData, tt.wantData)
			}

			if (len(resp.Diagnostics) > 0) != tt.wantDiag {
				t.Errorf("ProviderDataFromResourceConfigureRequest() diagnostic presence: got %v, want %v", len(resp.Diagnostics) > 0, tt.wantDiag)
			}
		})
	}
}
