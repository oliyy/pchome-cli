package prodapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/oliy/pchome-cli/pkg/httpx"
)

const baseURL = "https://ecapi-cdn.pchome.com.tw/ecshop/prodapi/v2/prod"

var idSuffixRe = regexp.MustCompile(`-\d{3}$`)

type Client struct {
	http *httpx.Client
}

func New(http *httpx.Client) *Client {
	return &Client{http: http}
}

func NormalizeID(id string) string {
	return idSuffixRe.ReplaceAllString(id, "")
}

func ProdURL(id string) string {
	return "https://24h.pchome.com.tw/prod/" + NormalizeID(id)
}

type ProductsOptions struct {
	IDs      []string
	Fields   []string
	Callback string
}

// Products fetches prodapi/v2/prod.
//
// Response can be:
//   - map keyed by product id (usually suffixed with -000) when at least one id is valid
//   - empty array [] when no ids are valid
//   - array of product objects on shape drift
//
// This helper always returns a map (possibly empty).
func (c *Client) Products(ctx context.Context, opt ProductsOptions) (map[string]Product, error) {
	ids := normalizeIDs(opt.IDs)
	if len(ids) == 0 {
		return map[string]Product{}, nil
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("id", strings.Join(ids, ","))
	if len(opt.Fields) > 0 {
		q.Set("fields", strings.Join(opt.Fields, ","))
	}
	if opt.Callback != "" {
		q.Set("_callback", opt.Callback)
	}
	u.RawQuery = q.Encode()

	body, err := c.http.GetBytes(ctx, u.String(), "application/json,*/*")
	if err != nil {
		return nil, err
	}

	payload, err := parseJSONOrJSONP(body)
	if err != nil {
		return nil, err
	}

	return decodeProductsPayload(payload)
}

func decodeProductsPayload(payload []byte) (map[string]Product, error) {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(payload, &rawMap); err == nil {
		out := make(map[string]Product, len(rawMap))
		for key, raw := range rawMap {
			var product Product
			if err := json.Unmarshal(raw, &product); err != nil {
				return nil, fmt.Errorf("decode prodapi product %s: %w", NormalizeID(key), err)
			}
			out[key] = product
		}
		return out, nil
	}

	var rawList []json.RawMessage
	if err := json.Unmarshal(payload, &rawList); err == nil {
		out := make(map[string]Product, len(rawList))
		for i, raw := range rawList {
			var product Product
			if err := json.Unmarshal(raw, &product); err != nil {
				return nil, fmt.Errorf("decode prodapi product at index %d: %w", i, err)
			}
			id := NormalizeID(product.Id)
			if id == "" {
				return nil, fmt.Errorf("unexpected prodapi product array item without Id at index %d", i)
			}
			out[id] = product
		}
		return out, nil
	}

	return nil, fmt.Errorf("unexpected prodapi response shape")
}

type ButtonOptions struct {
	IDs      []string
	Fields   []string
	Callback string
}

func (c *Client) Button(ctx context.Context, opt ButtonOptions) ([]ButtonItem, error) {
	ids := normalizeIDs(opt.IDs)
	if len(ids) == 0 {
		return []ButtonItem{}, nil
	}

	// Non-standard URL format (observed): uses `&` rather than `?`.
	// Example: https://.../prod/button&id=<ids>&fields=...&_callback=...
	var b strings.Builder
	b.WriteString(baseURL)
	b.WriteString("/button&id=")
	b.WriteString(url.QueryEscape(strings.Join(ids, ",")))

	if len(opt.Fields) > 0 {
		b.WriteString("&fields=")
		b.WriteString(url.QueryEscape(strings.Join(opt.Fields, ",")))
	}
	if opt.Callback != "" {
		b.WriteString("&_callback=")
		b.WriteString(url.QueryEscape(opt.Callback))
	}

	body, err := c.http.GetBytes(ctx, b.String(), "application/json,*/*")
	if err != nil {
		return nil, err
	}

	payload, err := parseJSONOrJSONP(body)
	if err != nil {
		return nil, err
	}

	var out []ButtonItem
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func normalizeIDs(ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, s := range ids {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func parseJSONOrJSONP(body []byte) ([]byte, error) {
	trim := strings.TrimSpace(string(body))
	if trim == "" {
		return nil, fmt.Errorf("empty response body")
	}

	// Fast path: plain JSON.
	if json.Valid([]byte(trim)) {
		return []byte(trim), nil
	}

	arg, ok := extractFirstJSONPArgument(trim)
	if !ok {
		return nil, fmt.Errorf("unable to parse response as JSON or JSONP")
	}
	if !json.Valid([]byte(arg)) {
		return nil, fmt.Errorf("invalid JSONP payload")
	}
	return []byte(arg), nil
}

func extractFirstJSONPArgument(text string) (string, bool) {
	startParen := strings.IndexByte(text, '(')
	if startParen < 0 {
		return "", false
	}

	i := startParen + 1
	depth := 1
	var quote byte

	for ; i < len(text); i++ {
		ch := text[i]

		if quote != 0 {
			if ch == '\\' {
				i++ // skip escaped char
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}

		if ch == '"' || ch == '\'' {
			quote = ch
			continue
		}

		if ch == '(' {
			depth++
			continue
		}

		if ch == ')' {
			depth--
			if depth == 0 {
				return text[startParen+1 : i], true
			}
		}
	}

	return "", false
}
