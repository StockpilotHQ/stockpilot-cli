package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURL = "https://api.stockpilot.dev"

type Client struct {
	ClientID     string
	ClientSecret string
	HTTP         *http.Client
}

func New(clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTP:         &http.Client{},
	}
}

func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	endpoint := BaseURL + path
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	return c.do(req)
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.request("POST", path, body)
}

func (c *Client) Patch(path string, body any) ([]byte, error) {
	return c.request("PATCH", path, body)
}

func (c *Client) Put(path string, body any) ([]byte, error) {
	return c.request("PUT", path, body)
}

func (c *Client) Delete(path string, body any) ([]byte, error) {
	return c.request("DELETE", path, body)
}

func (c *Client) request(method, path string, body any) ([]byte, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, BaseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req)
	return c.do(req)
}

func (c *Client) setAuth(req *http.Request) {
	req.Header.Set("X-CLIENT-ID", c.ClientID)
	req.Header.Set("X-CLIENT-SECRET", c.ClientSecret)
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var apiErr struct {
			Detail string `json:"detail"`
			Error  string `json:"error"`
		}
		if jsonErr := json.Unmarshal(data, &apiErr); jsonErr == nil {
			msg := apiErr.Detail
			if msg == "" {
				msg = apiErr.Error
			}
			if msg != "" {
				return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, msg)
			}
		}
		return nil, fmt.Errorf("API error %d", resp.StatusCode)
	}
	return data, nil
}
