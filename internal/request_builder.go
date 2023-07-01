package internal

import (
	"context"
	"net/http"
)

type RequestBuilder interface {
	Build(
		ctx context.Context,
		method string,
		url string,
	) (*http.Request, error)
}

type HTTPRequestBuilder struct{}

func NewRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{}
}

func (b *HTTPRequestBuilder) Build(
	ctx context.Context,
	method string,
	url string,
) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
