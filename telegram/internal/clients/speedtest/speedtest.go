package speedtest

import (
	"encoding/json"
	"fmt"
	"net/http"
)


const(
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

func (c *Client) SpeedTest(chatID int64) (*SpeedTestResult, error) {
	url := fmt.Sprintf("%s%s?chat_id=%d", c.addr, SpeedTestMethod, chatID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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