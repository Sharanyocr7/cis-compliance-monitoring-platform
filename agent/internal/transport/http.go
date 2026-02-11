package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type HTTPClient struct {
	URL     string
	APIKey  string
	Retries int
	Backoff time.Duration
}

func (c HTTPClient) PostJSON(ctx context.Context, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	for i := 0; i <= c.Retries; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if c.APIKey != "" {
			req.Header.Set("X-API-Key", c.APIKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		if err == nil && resp != nil {
			lastErr = errors.New(resp.Status)
		} else {
			lastErr = err
		}

		// backoff
		if i < c.Retries {
			select {
			case <-time.After(c.Backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return lastErr
}
