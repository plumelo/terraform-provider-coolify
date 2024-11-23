package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-coolify/internal/api"
)

const MOCK_TOKEN = "1|validToken"

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

	retryConfig := api.RetryConfig{
		MaxAttempts: 1,
		MinWait:     1,
		MaxWait:     1,
	}

	// Test with valid token
	client, err := api.NewAPIClient("test", mockServer.URL, MOCK_TOKEN, retryConfig)
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	resp, err := client.VersionWithResponse(context.Background())
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.HTTPResponse.StatusCode)
	}

	// Test with invalid token
	invalidToken := "invalid_token"
	_, err = api.NewAPIClient("test", mockServer.URL, invalidToken, retryConfig)
	if err == nil {
		t.Fatalf("Expected error when creating API client with invalid token, got none")
	}
}

func TestValidateTokenFormat(t *testing.T) {
	tests := []struct {
		token   string
		wantErr bool
	}{
		{"12345|validToken", false},
		{"7|anotherValidToken", false},
		{"invalid_token", true},
		{"12345|", true},
		{"|token", true},
		{"12345|token with spaces", true},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			err := api.ValidateTokenFormat(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTokenFormat(%q) error = %v, wantErr %v", tt.token, err, tt.wantErr)
			}
		})
	}
}
