package soren

import (
	"context"
)

type ShowConfigPayload struct {
	Op   string   `json:"op"`
	Path []string `json:"path"`
}

func (c *Client) ShowConfig(ctx context.Context, path []string) (map[string]any, error) {
	payload := ShowConfigPayload{
		Op:   "getInfo",
		Path: path,
	}

	var result map[string]any

	err := c.Call(ctx, "/retrieve", payload, &result)
	if err != nil {
		return nil, err
	}

	// 3. Return the fully populated dynamic map
	return result, nil
}
