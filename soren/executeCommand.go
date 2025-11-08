package soren

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
)

// APIResponse is the common wrapper for all API calls.
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func (c *Client) executeCommand(ctx context.Context, apiPath string, payload any) (*APIResponse, error) {
	fullURL := c.baseURL + apiPath

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	if err := writer.WriteField("data", string(jsonBytes)); err != nil {
		return nil, fmt.Errorf("failed to write 'data' field: %w", err)
	}

	if err := writer.WriteField("key", c.apiKey); err != nil {
		return nil, fmt.Errorf("failed to write 'key' field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to form request for %s: %w", fullURL, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to complete request to %s: %w", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request to %s failed with status: %s", fullURL, resp.Status)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response wrapper for %s: %w", fullURL, err)
	}

	if !apiResponse.Success {
		return nil, errors.New("API command failed: " + apiResponse.Error)
	}

	return &apiResponse, nil
}
