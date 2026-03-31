package onsale

import (
	"context"
	"encoding/json"

	"github.com/oliy/pchome-cli/pkg/httpx"
)

const URL = "https://ecapi-cdn.pchome.com.tw/fsapi/cms/onsale"

type Client struct {
	http *httpx.Client
}

func New(http *httpx.Client) *Client {
	return &Client{http: http}
}

func (c *Client) Fetch(ctx context.Context) (*Response, error) {
	body, err := c.http.GetBytes(ctx, URL, "application/json")
	if err != nil {
		return nil, err
	}
	var out Response
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
