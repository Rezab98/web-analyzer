package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/Rezab98/web-analyzer/internal/analyzer"
)

type AnalyzerHandler struct {
	pageAnalyzer *analyzer.WebpageAnalyzer
	template     *template.Template
}

func NewAnalyzerHandler(pageAnalyzer *analyzer.WebpageAnalyzer) *AnalyzerHandler {
	template := template.Must(template.ParseFiles("templates/form.html", "templates/results.html"))

	return &AnalyzerHandler{
		pageAnalyzer: pageAnalyzer,
		template:     template,
	}
}

// showForm displays the webpage analysis form.
func (h *AnalyzerHandler) showForm(w http.ResponseWriter, req *http.Request) {
	if err := h.template.ExecuteTemplate(w, "form.html", nil); err != nil {
		logrus.Errorf("handler.analyzeURL: error rendering template: %v", err)
		// Return an internal server error to the client if an error occurs during template rendering
		http.Error(w, "An error occurred when rendering template", http.StatusInternalServerError)
	}
}

// analyzeURL handles webpage analysis requests.
func (h *AnalyzerHandler) analyzeURL(w http.ResponseWriter, req *http.Request) {
	urlStr := req.FormValue("url")

	// Validate and potentially correct the URL format
	url, err := validateURL(urlStr)
	if err != nil {
		logrus.Errorf("handler.analyzeURL: invalid URL: %v", err)
		// Return a bad request error to the client if the URL is invalid
		http.Error(w, fmt.Sprintf("Invalid URL: %v", err), http.StatusBadRequest)
		return
	}

	result, err := h.pageAnalyzer.Analyze(req.Context(), url)
	if err != nil {
		logrus.Errorf("handler.analyzeURL: error analyzing URL: %v", err)
		// Return an internal server error to the client if an error occurs during analysis
		http.Error(w, "An error occurred while analyzing the URL. Please try again later.", http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "results.html", result); err != nil {
		logrus.Errorf("handler.analyzeURL: error rendering template: %v", err)
		// Return an internal server error to the client if an error occurs during template rendering
		http.Error(w, "An error occurred when rendering template", http.StatusInternalServerError)
		return
	}
}

// validateURL checks if the given string is a valid URL using a regular expression.
func validateURL(urlStr string) (string, error) {
	const urlPattern = `^(http|https)://[a-zA-Z0-9]+([\-\.]{1}[a-zA-Z0-9]+)*\.[a-zA-Z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

	urlStr = strings.TrimSpace(urlStr)

	regex := regexp.MustCompile(urlPattern)

	if !regex.MatchString(urlStr) {
		// Wrap the error in a fmt.Errorf to provide a more detailed error message
		return "", fmt.Errorf("invalid URL format")
	}

	// Return the URL string if it is valid
	return urlStr, nil
}
