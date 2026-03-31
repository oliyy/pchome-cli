package hermes

import (
	"context"
	"encoding/json"

	"github.com/oliy/pchome-cli/pkg/httpx"
)

const (
	GoodsRankURL   = "https://apih.pcloud.tw/hermes/api/goods/rank"
	DefaultToken   = "JeFbL9ypvv"
	DefaultRecPos  = "bsim"
	DefaultRecType = "ClickStream"
)

type Client struct {
	http *httpx.Client
}

func New(http *httpx.Client) *Client {
	return &Client{http: http}
}

func (c *Client) GoodsRank(ctx context.Context, req GoodsRankRequest) (*GoodsRankResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	body, err := c.http.PostJSONBytes(ctx, GoodsRankURL, b, "application/json")
	if err != nil {
		return nil, err
	}

	var out GoodsRankResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
