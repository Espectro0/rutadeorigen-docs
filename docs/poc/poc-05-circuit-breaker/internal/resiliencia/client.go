package resiliencia

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type Client struct {
	breaker *gobreaker.CircuitBreaker
	http    *http.Client
	url     string
}

func NewClient(url string, breaker *gobreaker.CircuitBreaker) *Client {
	return &Client{
		breaker: breaker,
		http:    &http.Client{Timeout: 2 * time.Second},
		url:     url,
	}
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("respuesta inesperada: %d", resp.StatusCode)
		}

		return "pong", nil
	})
	if err != nil {
		return "", err
	}

	return result.(string), nil
}
