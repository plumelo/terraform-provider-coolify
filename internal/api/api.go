package api

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

const UserAgentPrefix = "terraform-provider-coolify"

func NewAPIClient(version, server, apiToken string) (*ClientWithResponses, error) {
	if err := ValidateTokenFormat(apiToken); err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	return NewClientWithResponses(server,
		WithHTTPClient(&httpClient),
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apiToken)
			req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", UserAgentPrefix, version))
			req.Header.Set("Accept", "application/json")
			return nil
		}),
	)
}

func ValidateTokenFormat(token string) error {
	matched, _ := regexp.MatchString(`^\d+\|\w+$`, token)
	if !matched {
		return fmt.Errorf("invalid token format")
	}
	return nil
}
