package api

import (
	"context"
	"log"
	"net/http"
)

const DefaultServerURL = "https://app.coolify.io/api/v1"

type APIClient struct {
	*ClientWithResponses
	httpClient *http.Client
}

func NewAPIClient(server string, apiToken string) *APIClient {
	httpClient := http.Client{}
	apiClient, err := NewClientWithResponses(server,
		WithHTTPClient(&httpClient),
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apiToken)
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

// func withAuthentication(apiKey string) RequestEditorFn {
// 	return func(ctx context.Context, req *http.Request) error {
// 		req.Header.Set("Authorization", "Bearer "+apiKey)
// 		return nil
// 	}
// }
// func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
// 	// req.Header.Set("Authorization", "Bearer " + c.apiKey)
// 	// return c.Client.Do(req)
// }

// func withRetries(maxRetries int) RequestEditorFn {
// 	return func(ctx context.Context, req *http.Request) error {
// 		var resp *http.Response
// 		var err error
// 		for i := 0; i < maxRetries; i++ {
// 			resp, err = http.DefaultClient.Do(req)
// 			if err == nil && resp.StatusCode < 500 {
// 				return nil
// 			}
// 			time.Sleep(time.Second * time.Duration(i+1))
// 		}
// 		return err
// 	}
// }
