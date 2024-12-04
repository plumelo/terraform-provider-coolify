package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type mockProviderData struct {
	Value string
}

func TestProviderDataFromDataSourceConfigureRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		providerData  any
		expected      bool
		expectError   bool
		expectedValue string
	}{
		{"NilProviderData", nil, false, false, ""},
		{"ValidProviderData", mockProviderData{Value: "test"}, true, false, "test"},
		{"InvalidProviderData", "invalid", false, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := datasource.ConfigureRequest{ProviderData: tt.providerData}
			resp := &datasource.ConfigureResponse{Diagnostics: diag.Diagnostics{}}
			var out mockProviderData

			got := ProviderDataFromDataSourceConfigureRequest(req, &out, resp)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}

			if tt.expectError && len(resp.Diagnostics) == 0 {
				t.Error("expected error diagnostics, got none")
			}

			if !tt.expectError && len(resp.Diagnostics) > 0 {
				t.Error("expected no error diagnostics, got some")
			}

			if tt.expected && out.Value != tt.expectedValue {
				t.Errorf("expected value %s, got %s", tt.expectedValue, out.Value)
			}
		})
	}
}

func TestProviderDataFromResourceConfigureRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		providerData  any
		expected      bool
		expectError   bool
		expectedValue string
	}{
		{"NilProviderData", nil, false, false, ""},
		{"ValidProviderData", mockProviderData{Value: "test"}, true, false, "test"},
		{"InvalidProviderData", "invalid", false, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := resource.ConfigureRequest{ProviderData: tt.providerData}
			resp := &resource.ConfigureResponse{Diagnostics: diag.Diagnostics{}}
			var out mockProviderData

			got := ProviderDataFromResourceConfigureRequest(req, &out, resp)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}

			if tt.expectError && len(resp.Diagnostics) == 0 {
				t.Error("expected error diagnostics, got none")
			}

			if !tt.expectError && len(resp.Diagnostics) > 0 {
				t.Error("expected no error diagnostics, got some")
			}

			if tt.expected && out.Value != tt.expectedValue {
				t.Errorf("expected value %s, got %s", tt.expectedValue, out.Value)
			}
		})
	}
}
