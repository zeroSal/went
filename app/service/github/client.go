package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"went/app/model/github"
	"went/app/service/http"
)

const baseURL = "https://api.github.com"

var _ ClientInterface = (*Client)(nil)

type Client struct {
	httpClient http.ClientInterface
	baseURL  string
}

func NewClient(httpClient http.ClientInterface) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:   baseURL,
	}
}

func (c *Client) ListReleases(
	ctx context.Context,
	owner,
	repo string,
) ([]github.Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", c.baseURL, owner, repo)

	c.httpClient.AddHeader("Accept", "application/vnd.github+json")
	c.httpClient.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch releases: status %d", resp.StatusCode)
	}

	var releases []github.Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}

	return releases, nil
}

func (c *Client) DownloadAsset(
	ctx context.Context,
	downloadURL string,
	writer io.Writer,
) error {
	c.httpClient.AddHeader("Accept", "application/octet-stream")
	c.httpClient.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	err := c.httpClient.Download(downloadURL, writer)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	return nil
}

func (c *Client) DownloadAssetBytes(
	ctx context.Context,
	downloadURL string,
) ([]byte, error) {
	c.httpClient.AddHeader("Accept", "application/octet-stream")
	c.httpClient.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	var buf bytes.Buffer
	err := c.httpClient.Download(downloadURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}

	return buf.Bytes(), nil
}