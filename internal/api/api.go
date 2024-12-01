package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const UserAgentPrefix = "terraform-provider-coolify"

var (
	TokenRegex      = regexp.MustCompile(`^\d+\|\w+$`)
	ErrInvalidToken = errors.New("invalid token format")
)

type RetryConfig struct {
	MaxAttempts int64
	MinWait     int64
	MaxWait     int64
}

func NewAPIClient(version, server, apiToken string, retry RetryConfig) (*ClientWithResponses, error) {
	if err := ValidateTokenFormat(apiToken); err != nil {
		return nil, err
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = int(retry.MaxAttempts)
	retryClient.RetryWaitMin = time.Duration(retry.MinWait) * time.Second
	retryClient.RetryWaitMax = time.Duration(retry.MaxWait) * time.Second
	retryClient.Backoff = retryablehttp.DefaultBackoff
	retryClient.Logger = nil

	httpClient := retryClient.StandardClient()

	return NewClientWithResponses(server,
		WithHTTPClient(httpClient),
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
