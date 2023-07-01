package zerox

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	utils "github.com/Daniil675/fiolet-playground-backend/internal"
)

type Client struct {
	config         Config
	requestBuilder utils.RequestBuilder
	Err            ApiErrorResponse
}

func NewClient(apiKey, nodeURL string) *Client {
	config := GetConfig(apiKey, nodeURL)
	return NewClientWithConfig(config)
}

func NewClientWithConfig(c Config) *Client {
	return &Client{
		config:         c,
		requestBuilder: utils.NewRequestBuilder(),
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := c.requestBuilder.Build(ctx, method, url)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if isFailureStatusCode(res) {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) setHeaders(req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	apiKey := req.Header.Get("0x-api-key")
	if apiKey == "" {
		req.Header.Set("0x-api-key", c.config.AuthToken)
	}
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	err := json.NewDecoder(resp.Body).Decode(&c.Err.Message)
	if err != nil {
		return err
	}

	c.Err.HTTPStatusCode = resp.StatusCode
	return &c.Err
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	if result, ok := v.(*string); ok {
		return decodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
