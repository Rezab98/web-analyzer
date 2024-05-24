package pageanalyzer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
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
	htmlVersion := HTMLVersion(pageContent)
	title := htmlExtractor.Title()
	headingTagToTexts := htmlExtractor.HeadingTagToTexts()
	hasLoginForm := htmlExtractor.HasLoginForm()
	allLinks, err := resolveRelativeLinks(htmlExtractor.Links(), pageURL)
	if err != nil {
		logrus.WithError(err).Error("resolveRelativeLinks failed")
	}

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

func HTMLVersion(pageContent []byte) string {
	// Read only the first 1024 bytes to identify the doctype and root elements

	// Helper function to get the minimum of two integers
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	initialContent := bytes.ToLower(pageContent[:min(len(pageContent), 1024)])

	// Convert the relevant portion to a string for easier processing
	initialContentStr := string(initialContent)

	// Check for HTML5 doctype
	if strings.Contains(initialContentStr, "<!doctype html>") {
		return "HTML5"
	}

	// Check for XHTML doctypes
	if strings.Contains(initialContentStr, `<!doctype html public "-//w3c//dtd xhtml 1.0`) {
		return "XHTML 1.0"
	}
	if strings.Contains(initialContentStr, `<!doctype html public "-//w3c//dtd xhtml 1.1`) {
		return "XHTML 1.1"
	}

	// Check for older HTML versions
	if strings.Contains(initialContentStr, "<html>") || strings.Contains(initialContentStr, "<html ") {
		return "HTML4 or Earlier"
	}

	return "Unknown"
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
	base, err := url.Parse(baseURL)
	if err != nil {
		logrus.Errorf("Error parsing baseURL: %v", err)
		return false
	}

	link, err := url.Parse(href)
	if err != nil {
		logrus.Errorf("Error parsing baseURL: %v", err)
		return false
	}

	// Remove www prefix from the href's host for comparison
	linkHost := strings.TrimPrefix(link.Host, "www.")

	// Compare the host of the base URL and the resolved URL
	return base.Host == linkHost
}

// resolveRelativeLinks converts all relative links to absolute links based on the base URL's host.
func resolveRelativeLinks(links []string, baseURL string) ([]string, error) {
	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing baseURL: %v", err)
	}

	var resolvedLinks []string
	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			logrus.Errorf("resolveRelativeLinks: error parsing link: %v", err)
			continue
		}
		resolvedLink := link
		if !parsedLink.IsAbs() {
			resolvedLink = base.ResolveReference(parsedLink).String()
		}
		resolvedLinks = append(resolvedLinks, resolvedLink)

	}

	return resolvedLinks, nil
}
