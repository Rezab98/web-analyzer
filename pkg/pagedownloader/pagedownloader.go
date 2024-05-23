package pagedownloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrNotfound is returned when a webpage is not found
var ErrNotfound = errors.New("page not found")

type SimpleWebPageDownloader struct {
	client *http.Client
}

func New() *SimpleWebPageDownloader {
	return &SimpleWebPageDownloader{
		client: http.DefaultClient,
	}
}

// Download downloads the webpage content from the specified URL.
func (d *SimpleWebPageDownloader) Download(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := d.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If the status code is 404 (Not Found), return the specific ErrNotfound error.
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("request GET failed with: %w", ErrNotfound)
		}
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all data failed: %w", err)
	}

	return body, nil
}
