package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Rezab98/web-analyzer/internal/pageanalyzer"
	"github.com/Rezab98/web-analyzer/internal/pagedownloader"
)

type AnalyzerHandler struct {
	pageAnalyzer   *pageanalyzer.WebpageAnalyzer
	pageDownloader *pagedownloader.SimpleWebPageDownloader

	template *template.Template
}

func NewAnalyzerHandler(pageAnalyzer *pageanalyzer.WebpageAnalyzer, pageDownloader *pagedownloader.SimpleWebPageDownloader) *AnalyzerHandler {
	template := template.Must(template.ParseFiles("templates/form.html", "templates/results.html"))

	return &AnalyzerHandler{
		pageAnalyzer:   pageAnalyzer,
		pageDownloader: pageDownloader,

		template: template,
	}
}

func (h *AnalyzerHandler) showForm(w http.ResponseWriter, r *http.Request) {
	if err := h.template.ExecuteTemplate(w, "form.html", nil); err != nil {
		handleHTTPError(w, r,
			"An error occurred rendering template. Please try again later.",
			http.StatusInternalServerError,
			fmt.Errorf("handler.analyzeURL: error rendering template: %v", err),
		)
		return
	}
}

type TemplateData struct {
	URL                  string
	Error                string
	HTMLVersion          string
	Title                string
	HeadingTagToTexts    map[string][]string
	HeadingTagToTextsNum map[string]int
	InternalLinksNum     int
	ExternalLinksNum     int
	InternalLinks        []string
	ExternalLinks        []string
	InaccessibleLinksNum int
	HasLoginForm         bool
}

func (h *AnalyzerHandler) analyzeURL(w http.ResponseWriter, r *http.Request) {
	urlStr := r.FormValue("url")

	// EnsureURLIsValid checks and corrects the URL format, adding "https://" if missing.
	url, err := validateURL(urlStr)
	if err != nil {
		handleHTTPError(w, r,
			fmt.Sprintf("Invalid URL: %v", err),
			http.StatusBadRequest,
			nil,
		)
		return
	}

	downloadCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pageContent, err := h.pageDownloader.Download(downloadCtx, urlStr)
	if err != nil {
		if errors.Is(err, pagedownloader.ErrNotfound) {
			handleHTTPError(w, r,
				"page doesn't exist",
				http.StatusNotFound,
				nil,
			)
			return
		}

		handleHTTPError(w, r,
			"An error occurred while downloading the Page. Please try again later.",
			http.StatusInternalServerError,
			fmt.Errorf("download page failed: %v", err),
		)
		return
	}

	analyzerCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pageAnalyzedResult, err := h.pageAnalyzer.Analyze(analyzerCtx, url, pageContent)
	if err != nil {
		handleHTTPError(w, r,
			"An error occurred while analyzing the page. Please try again later.",
			http.StatusInternalServerError,
			fmt.Errorf("analyze page failed: %v", err),
		)
		return
	}

	// Prepare template data
	headingTagToTextNum := make(map[string]int)
	for tag, texts := range pageAnalyzedResult.HeadingTagToTexts {
		headingTagToTextNum[tag] = len(texts)
	}

	templateData := TemplateData{
		URL:          url,
		HTMLVersion:  pageAnalyzedResult.HTMLVersion,
		Title:        pageAnalyzedResult.Title,
		HasLoginForm: pageAnalyzedResult.HasLoginForm,

		HeadingTagToTexts:    pageAnalyzedResult.HeadingTagToTexts,
		HeadingTagToTextsNum: headingTagToTextNum,
		InternalLinks:        pageAnalyzedResult.InternalLinks,
		InternalLinksNum:     len(pageAnalyzedResult.InternalLinks),

		ExternalLinks:    pageAnalyzedResult.ExternalLinks,
		ExternalLinksNum: len(pageAnalyzedResult.ExternalLinks),

		InaccessibleLinksNum: pageAnalyzedResult.InaccessibleLinksNum,
	}

	// Execute template
	if err := h.template.ExecuteTemplate(w, "results.html", templateData); err != nil {
		handleHTTPError(w, r,
			"An error occurred rendering template. Please try again later.",
			http.StatusInternalServerError,
			fmt.Errorf("execute template failed: %v", err),
		)
		return
	}
}

// validateURL checks if the given string is a valid URL using regex.
func validateURL(urlStr string) (string, error) {
	// Define the URL pattern.
	const urlPattern = `^(http|https)://[a-zA-Z0-9]+([\-\.]{1}[a-zA-Z0-9]+)*\.[a-zA-Z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

	urlStr = strings.TrimSpace(urlStr)

	// Compile the regex.
	regex := regexp.MustCompile(urlPattern)

	// Check if the URL matches the regex pattern.
	if !regex.MatchString(urlStr) {
		return "", fmt.Errorf("invalid URL format")
	}

	return urlStr, nil
}

// handleHTTPError wil return proper http error on the responseWriter and log error with request scope information
func handleHTTPError(w http.ResponseWriter, r *http.Request, msg string, statusCode int, causeErr error) {
	logrus.WithError(causeErr).WithFields(logrus.Fields{
		"path":         r.URL.Path,
		"method":       r.Method,
		"responseCode": statusCode,
		"responseMsg":  msg,
	}).Error("request failed")

	http.Error(w, msg, statusCode)
}
