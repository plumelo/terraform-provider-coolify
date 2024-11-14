package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

const UserAgentPrefix = "terraform-provider-coolify"

var (
	TokenRegex      = regexp.MustCompile(`^\d+\|\w+$`)
	ErrInvalidToken = errors.New("invalid token format")
)

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
	if !TokenRegex.MatchString(token) {
		return ErrInvalidToken
	}
	return nil
}
