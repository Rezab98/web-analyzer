package analyzer

import (
	"context"
)

type Downloader interface {
	Download(ctx context.Context, url string) ([]byte, error)
}

type WebpageAnalyzer struct {
	downloader Downloader
}

func New(d Downloader) *WebpageAnalyzer {
	return &WebpageAnalyzer{downloader: d}
}

type AnalysisResult struct {
	URL               string
	Error             string
	HTMLVersion       string
	Title             string
	HeadingTexts      map[string][]string
	HeadingCounts     map[string]int
	InternalLinks     int
	ExternalLinks     int
	InternalLinksList []string
	ExternalLinksList []string
	InaccessibleLinks int
	HasLoginForm      bool
}

func (a *WebpageAnalyzer) Analyze(ctx context.Context, url string) (*AnalysisResult, error) {
	// return a temp AnalysisResult with nil values

	AnalysisResult := &AnalysisResult{
		URL:               url,
		Error:             "",
		HTMLVersion:       "",
		Title:             "",
		HeadingTexts:      map[string][]string{},
		HeadingCounts:     map[string]int{},
		InternalLinks:     0,
		ExternalLinks:     0,
		InternalLinksList: []string{},
		ExternalLinksList: []string{},
		InaccessibleLinks: 0,
		HasLoginForm:      false,
	}

	return AnalysisResult, nil

}
