package pageanalyzer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/Rezab98/web-analyzer/internal/pageanalyzer/htmlextract"
	"github.com/Rezab98/web-analyzer/pkg/slicetools"
)

type WebpageAnalyzer struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *WebpageAnalyzer {
	return &WebpageAnalyzer{httpClient: httpClient}
}

type Result struct {
	HTMLVersion          string
	Title                string
	HeadingTagToTexts    map[string][]string
	HasLoginForm         bool
	InternalLinks        []string
	ExternalLinks        []string
	InaccessibleLinksNum int
}

func (w *WebpageAnalyzer) Analyze(ctx context.Context, pageURL string, pageContent []byte) (*Result, error) {
	// Initialize html extractor
	htmlExtractor, err := htmlextract.New(pageContent)
	if err != nil {
		return nil, fmt.Errorf("initialize html extractor failed: %v", err)
	}

	// Extracts information from the HTML document
	htmlVersion := htmlExtractor.HTMLVersion()
	title := htmlExtractor.Title()
	headingTagToTexts := htmlExtractor.HeadingTagToTexts()
	hasLoginForm := htmlExtractor.HasLoginForm()
	allLinks := htmlExtractor.Links()

	// Get inaccessible links number
	inaccessibleLinksNum := w.CountInaccessibleLinks(ctx, allLinks)

	// Extract internal links
	internalLinks := slicetools.Filter(
		allLinks, func(link string) bool {
			return isInternalLink(link, pageURL)
		},
	)

	// Extract external links
	externalLinks := slicetools.Filter(
		allLinks, func(link string) bool {
			return !isInternalLink(link, pageURL)
		},
	)

	return &Result{
		HTMLVersion:          htmlVersion,
		Title:                title,
		HeadingTagToTexts:    headingTagToTexts,
		InternalLinks:        internalLinks,
		ExternalLinks:        externalLinks,
		InaccessibleLinksNum: inaccessibleLinksNum,
		HasLoginForm:         hasLoginForm,
	}, nil
}

// CountInaccessibleLinks sends concurrent HEAD requests to each link and returns the count of inaccessible links.
func (w *WebpageAnalyzer) CountInaccessibleLinks(ctx context.Context, links []string) int {
	var (
		inaccessibleLinkNum int
		lock                sync.Mutex
	)

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)

		go func(link string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if !w.isLinkAccessible(ctx, link) {
				lock.Lock()
				inaccessibleLinkNum++
				lock.Unlock()
			}
		}(link)
	}

	wg.Wait()

	return inaccessibleLinkNum
}

func (w *WebpageAnalyzer) isLinkAccessible(ctx context.Context, link string) bool {

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, link, nil /* body */)
	if err != nil {
		logrus.WithError(err).Error("new request with context failed")

		return false
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Send request failed")

		return false
	}
	defer resp.Body.Close()

	inaccessibleStatusCodes := []int{
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}

	return !slices.Contains(inaccessibleStatusCodes, resp.StatusCode)
}

func isInternalLink(href, baseURL string) bool {
	return strings.HasPrefix(href, baseURL) || !strings.HasPrefix(href, "http")
}
