package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (updates []Update, err error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can`t doRequest: %w", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can`t Unmarshall: %w", err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	data, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	var response struct {
		Ok    bool   `json:"ok"`
		Error string `json:"description"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("can't parse Telegram response: %w", err)
	}

	if !response.Ok {
		return fmt.Errorf("Telegram API error: %s", response.Error)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() {
		if err != nil {
			log.Printf("Telegram request error: %v", err)
		}
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	log.Printf("Making request to URL: %s", u.String())
	log.Printf("Query parameters: %v", query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	log.Printf("Response status code: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(body))

	return body, nil
}
