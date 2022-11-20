package http

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Client represents default http client that can be used to send requests to third party services
type Client struct {
	base http.Client
}

// Post send posts request to url with additional headers
func (c *Client) Post(ctx context.Context, url string, req []byte) ([]byte, error) {
	reqBody := bytes.NewBuffer(req)

	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	return executeRequest(c, request)
}

// Get send request to url with requestID headers
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url,
		http.NoBody)
	if err != nil {
		return nil, err
	}

	return executeRequest(c, req)
}

// executeRequest contains utils logic of request execution
func executeRequest(c *Client, r *http.Request) ([]byte, error) {
	resp, err := c.base.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http request failed with status %v, error: %v", resp.StatusCode, string(body))
	}

	return body, nil
}
