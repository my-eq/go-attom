package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// NewRequest constructs an HTTP request relative to the client's base URL.
//
// The endpoint must be a relative path without leading scheme. Query parameters
// are optional and will be URL-encoded. The Accept header defaults to
// application/json when not already provided.
func (c *Client) NewRequest(ctx context.Context, method, endpoint string, query url.Values, body io.Reader) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", c.baseURL, err)
	}

	trimmed := strings.TrimLeft(strings.TrimSpace(endpoint), "/")
	rel := &url.URL{Path: trimmed}
	if query != nil {
		rel.RawQuery = query.Encode()
	}

	finalURL := base.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, method, finalURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
