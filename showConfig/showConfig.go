package soren

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"mime/multipart"
	"bytes"

	"github.com/aevrex/soren/client" 
)

type ShowConfigPayload struct {
	Op   string   `json:"op"`
	Path []string `json:"path"`
}

type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func (c *Client) ShowConfig (ctx context.Context, path []string) (*APIResponse, error) {
	configURL := c.baseURL + "/retrieve"

	var requestBody bytes.Buffer 
	writer := multipart.NewWriter(&requestBody)

	payload := ShowConfigPayload{
		Op: "showConfig",
		Path: path,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal JSON payload: %w", err)
	}

    dataField, err := writer.CreateFormField("data")
    if err != nil {
        return nil, fmt.Errorf("Failed to create data field: %s", err)
    }

    if _, err := dataField.Write(jsonBytes); err != nil {
        return nil, fmt.Errorf("Failed to write 'data' field: %w", err)
    }

    keyField, err := writer.CreateFormField("key")
    if err != nil {
        return nil, fmt.Errorf("Failed to create 'key' form field: %w", err)
    }

    if _, err := keyField.Write([]byte(c.apiKey)); err != nil {
        return nil, fmt.Errorf("Failed to write 'key' field: %w", err)
    }

    if err := writer.Close(); err != nil {
        return nil, fmt.Errorf("Failed to close writer: %s", err)
    }

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, configURL, &requestBody) 
	if err != nil {
		return nil, fmt.Errorf("Failed to form request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to complete request: %w", err)
	}
	defer resp.Body.Close()


	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("Failed to decode API response wrapper: %w", err)
	}

	if !apiResponse.Success {
		return nil, errors.New("API command failed: " + apiResponse.Error)
	}

	return &APIResponse, nil
}