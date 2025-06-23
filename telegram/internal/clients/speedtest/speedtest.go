package speedtest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SpeedTestMethod = "/speedtest"
)

type Client struct {
	client http.Client
	addr   string
}

func New(addr string) *Client {
	return &Client{
		client: http.Client{},
		addr:   addr,
	}
}

func (c *Client) SpeedTest(chatID int64, appID int, appHash, proxy string) (*SpeedTestResult, error) {
	baseURL := fmt.Sprintf("%s%s", c.addr, SpeedTestMethod)
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("chat_id", fmt.Sprintf("%d", chatID))
	if appID == 0 {
		q.Add("app_id", "")
	} else {
		q.Add("app_id", fmt.Sprintf("%d", appID))
	}
	q.Add("app_hash", appHash)
	if proxy != "" {
		q.Add("proxy", proxy)
	}
	q.Add("mb", "10")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var result SpeedTestResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
