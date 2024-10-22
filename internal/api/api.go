package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const DefaultServerURL = "https://app.coolify.io/api/v1"
const UserAgentPrefix = "terraform-provider-coolify"

type APIClient struct {
	*ClientWithResponses
	httpClient *http.Client
}

func NewAPIClient(version, server, apiToken string) *APIClient {
	httpClient := http.Client{}
	apiClient, err := NewClientWithResponses(server,
		WithHTTPClient(&httpClient),
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apiToken)
			req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", UserAgentPrefix, version))
			req.Header.Set("Accept", "application/json")
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	return &APIClient{
		httpClient:          &httpClient,
		ClientWithResponses: apiClient,
	}
}
