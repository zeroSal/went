package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

var _ ClientInterface = (*Client)(nil)

type Client struct {
	client  *http.Client
	headers Header
	timeout time.Duration
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		client:  httpClient,
		headers: make(Header),
		timeout: defaultTimeout,
	}
}

func (c *Client) SetHeaders(headers Header) {
	c.headers = make(Header)
	for k, v := range headers {
		c.headers[k] = v
	}
}

func (c *Client) RemoveHeader(key string) {
	delete(c.headers, key)
}

func (c *Client) AddHeader(key, value string) {
	c.headers[key] = value
}

func (c *Client) SetTimeout(d time.Duration) {
	c.timeout = d
}

func (c *Client) Get(url string) (*Response, error) {
	return c.doRequest(http.MethodGet, url, nil)
}

func (c *Client) Post(url string, body io.Reader) (*Response, error) {
	return c.doRequest(http.MethodPost, url, body)
}

func (c *Client) Put(url string, body io.Reader) (*Response, error) {
	return c.doRequest(http.MethodPut, url, body)
}

func (c *Client) Patch(url string, body io.Reader) (*Response, error) {
	return c.doRequest(http.MethodPatch, url, body)
}

func (c *Client) Delete(url string) (*Response, error) {
	return c.doRequest(http.MethodDelete, url, nil)
}

func (c *Client) Download(url string, writer io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}

func (c *Client) doRequest(method, url string, body io.Reader) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{
		Timeout: c.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	responseHeader := make(Header)
	for k, v := range resp.Header {
		if len(v) > 0 {
			responseHeader[k] = v[0]
		}
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       bytes.NewReader(bodyBytes),
		Header:     responseHeader,
	}, nil
}