package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/oliy/pchome-cli/pkg/httpx"
)

const baseURL = "https://ecshweb.pchome.com.tw/search/v4.3/all"

type Client struct {
	http *httpx.Client
}

func New(http *httpx.Client) *Client {
	return &Client{http: http}
}

type ResultsParams struct {
	Q         string
	CateID    string
	Page      int
	PageCount int
	Sort      string
	Attr      string
	Price     string
}

func (c *Client) Results(ctx context.Context, p ResultsParams) (*ResultsResponse, error) {
	if p.Q == "" {
		return nil, fmt.Errorf("q is required")
	}

	u, err := url.Parse(baseURL + "/results")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", p.Q)
	if p.CateID != "" {
		q.Set("cateid", p.CateID)
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.PageCount > 0 {
		q.Set("pageCount", strconv.Itoa(p.PageCount))
	}
	if p.Sort != "" {
		q.Set("sort", p.Sort)
	}
	if p.Attr != "" {
		q.Set("attr", p.Attr)
	}
	if p.Price != "" {
		q.Set("price", p.Price)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}

	var out ResultsResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	// Some invalid param combos return "{}" with HTTP 200.
	if out.TotalRows == 0 && out.TotalPage == 0 && len(out.Prods) == 0 && out.Q == "" {
		return nil, fmt.Errorf("search API returned empty object (likely invalid params)")
	}
	return &out, nil
}

func (c *Client) Categories(ctx context.Context, qText, attr, price string) ([]CategoryNode, error) {
	u, err := url.Parse(baseURL + "/categories")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", qText)
	if attr != "" {
		q.Set("attr", attr)
	}
	if price != "" {
		q.Set("price", price)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}
	var out []CategoryNode
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) PCategories(ctx context.Context, qText, cateID, attr, price string) ([]PCategoryNode, error) {
	u, err := url.Parse(baseURL + "/pcategories")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", qText)
	if cateID != "" {
		q.Set("cateid", cateID)
	}
	if attr != "" {
		q.Set("attr", attr)
	}
	if price != "" {
		q.Set("price", price)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}
	var out []PCategoryNode
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Brands(ctx context.Context, qText, cateID, attr, price string) ([]BrandFacet, error) {
	if qText == "" {
		return nil, fmt.Errorf("q is required")
	}
	if cateID == "" {
		return nil, fmt.Errorf("cateid is required")
	}
	u, err := url.Parse(baseURL + "/brands")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", qText)
	q.Set("cateid", cateID)
	if attr != "" {
		q.Set("attr", attr)
	}
	if price != "" {
		q.Set("price", price)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}
	var out []BrandFacet
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GroupAttributes(ctx context.Context, qText, cateID, attr, price string) ([]GroupAttributeGroup, error) {
	if qText == "" {
		return nil, fmt.Errorf("q is required")
	}
	u, err := url.Parse(baseURL + "/groupAttributes")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", qText)
	if cateID != "" {
		q.Set("cateid", cateID)
	}
	if attr != "" {
		q.Set("attr", attr)
	}
	if price != "" {
		q.Set("price", price)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}
	var out []GroupAttributeGroup
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) SuggestWords(ctx context.Context, qText string) ([]SuggestWord, error) {
	if qText == "" {
		return nil, fmt.Errorf("q is required")
	}
	u, err := url.Parse(baseURL + "/suggestwords")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", qText)
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json")
	if err != nil {
		return nil, err
	}
	var out []SuggestWord
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
