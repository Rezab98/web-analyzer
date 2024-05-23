package pagedownloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrNotfound is returned when a webpage is not found
var (
	ErrNotfound = errors.New("page not found")
)

type SimpleWebPageDownloader struct {
	client *http.Client
}

func New(client *http.Client) *SimpleWebPageDownloader {
	return &SimpleWebPageDownloader{
		client: client,
	}
}

func (d *SimpleWebPageDownloader) Download(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)

	}
	resp, err := d.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client GET failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			// If the status code is 404 (Not Found), return the specific ErrNotfound error.
			return nil, fmt.Errorf("request GET failed with: %w", ErrNotfound)
		}

		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all data failed: %v", err)
	}

	return body, nil
}
