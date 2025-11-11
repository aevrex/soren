package soren

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// Call executes an API command and unmarshals the response data into the result variable.
func Call[T any](c *Client, ctx context.Context, apiPath string, payload any, result *T) error {
	apiResp, err := c.executeCommand(ctx, apiPath, payload)
	if err != nil {
		return err
	}

	prettyJSON, indentErr := json.MarshalIndent(apiResp.Data, "", "  ")
	if indentErr != nil {
		log.Printf("Warning: Failed to pretty-print API response data: %v", indentErr)
	} else {
		log.Printf("API Response Data (Path: %s):\n%s", apiPath, string(prettyJSON))
	}

	if err := json.Unmarshal(apiResp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal API data into target type: %w", err)
	}

	return nil
}
