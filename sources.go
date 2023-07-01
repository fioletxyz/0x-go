package zerox

import (
	"context"
	"net/http"
)

type (
	SourcesResponse map[string][]string
)

func (c *Client) GetSources(
	ctx context.Context,
) (SourcesResponse, error) {
	var response SourcesResponse

	urlSuffix := "swap/v1/sources"

	req, err := c.newRequest(ctx, http.MethodGet, c.config.URL+urlSuffix)
	if err != nil {
		return SourcesResponse{}, err
	}

	err = c.sendRequest(req, &response)
	if err != nil {
		return SourcesResponse{}, err
	}
	return response, nil
}
