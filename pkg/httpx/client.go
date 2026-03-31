package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Options struct {
	Timeout   time.Duration
	UserAgent string
	RPS       float64
	Burst     int
}

type Client struct {
	http      *http.Client
	limiter   *rate.Limiter
	userAgent string
}

func New(opts Options) *Client {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	rps := opts.RPS
	if rps <= 0 {
		rps = 2
	}
	burst := opts.Burst
	if burst <= 0 {
		burst = 2
	}

	ua := opts.UserAgent
	if ua == "" {
		ua = "pchome/0 (https://github.com/oliy/pchome-cli)"
	}

	return &Client{
		http: &http.Client{
			Timeout: timeout,
		},
		limiter:   rate.NewLimiter(rate.Limit(rps), burst),
		userAgent: ua,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.limiter != nil {
		if err := c.limiter.Wait(req.Context()); err != nil {
			return nil, err
		}
	}
	return c.http.Do(req)
}

type HTTPError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("HTTP %d for %s", e.StatusCode, e.URL)
	}
	return fmt.Sprintf("HTTP %d for %s: %s", e.StatusCode, e.URL, e.Body)
}

func (c *Client) GetBytes(ctx context.Context, url string, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: res.StatusCode,
			URL:        url,
			Body:       truncateBytes(body, 4000),
		}
	}
	return body, nil
}

func (c *Client) PostJSONBytes(ctx context.Context, url string, jsonBody []byte, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: res.StatusCode,
			URL:        url,
			Body:       truncateBytes(body, 4000),
		}
	}
	return body, nil
}

func truncateBytes(b []byte, max int) string {
	if len(b) <= max {
		return string(bytes.TrimSpace(b))
	}
	return string(bytes.TrimSpace(b[:max])) + "..."
}
