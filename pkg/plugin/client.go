package plugin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client

	addr string
}

func (c *Client) Command(cmd string) (string, error) {
	resp, err := c.client.Post(c.addr, "application/text", bytes.NewBufferString(cmd))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("plugin returned status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func NewClient(addr string, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		addr:   addr,
		client: client,
	}
}
