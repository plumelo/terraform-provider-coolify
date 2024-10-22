package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-coolify/internal/api"
)

const MOCK_TOKEN = "valid-token"

// MockHandler simulates API responses based on the request.
func MockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+MOCK_TOKEN {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if r.Header.Get("User-Agent") != api.UserAgentPrefix+"/test" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/version":
		// Simulate a successful response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`9.9.9-beta.999`))
	default:
		// Simulate a 404 not found
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestAPIClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockHandler))
	defer mockServer.Close()

	client := api.NewAPIClient("test", mockServer.URL, MOCK_TOKEN)

	resp, err := client.VersionWithResponse(context.Background())
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.HTTPResponse.StatusCode)
	}

}
