package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type CataasClient struct {
	baseURL    string
	httpClient *http.Client
}

type CatImageResponse struct {
	Data        []byte
	ContentType string
	Size        int64
}

func NewCataasClient(baseURL string, timeoutSeconds int) *CataasClient {
	return &CataasClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *CataasClient) GetRandomCat() (*CatImageResponse, error) {
	url := fmt.Sprintf("%s/cat", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cat image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	return &CatImageResponse{
		Data:        data,
		ContentType: contentType,
		Size:        int64(len(data)),
	}, nil
}

func (c *CataasClient) HealthCheck() error {
	url := fmt.Sprintf("%s/cat", c.baseURL)
	
	resp, err := c.httpClient.Head(url)
	if err != nil {
		return fmt.Errorf("api not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	return nil
}