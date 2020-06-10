package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Make HTTP Get request and parse response as JSON
func GetJSON(ctx context.Context, baseURL string, reply interface{}) error {
	return GetJSONWithHeaders(ctx, baseURL, reply, nil)
}

// Make HTTP Get request and parse response as JSON with custom headers
func GetJSONWithHeaders(ctx context.Context, baseURL string, reply interface{}, headers map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return fmt.Errorf("prepare request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoke request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status %d %s", res.StatusCode, res.Status)
	}
	err = json.NewDecoder(res.Body).Decode(reply)
	if err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// Make HTTP Post request with application/json payload and parse response as JSON
func PostJSON(ctx context.Context, baseURL string, request interface{}, reply interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("prepare request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("invoke request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status %d %s", res.StatusCode, res.Status)
	}
	err = json.NewDecoder(res.Body).Decode(reply)
	if err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
