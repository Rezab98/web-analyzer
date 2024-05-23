package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Rezab98/web-analyzer/internal/pageanalyzer"
	"github.com/Rezab98/web-analyzer/internal/pagedownloader"
)

type Config struct {
	Port int
	Host string
}

// New creates a new HTTP server and sets up the routes.
func New(cfg *Config, pageDownload *pagedownloader.SimpleWebPageDownloader, pageAnalyzer *pageanalyzer.WebpageAnalyzer) *http.Server {

	router := mux.NewRouter()

	analyzerHandler := NewAnalyzerHandler(pageAnalyzer, pageDownload)

	// Set up the routes
	router.HandleFunc("/", analyzerHandler.showForm).Methods(http.MethodGet)
	router.HandleFunc("/", analyzerHandler.analyzeURL).Methods(http.MethodPost)

	router.Use(LoggingMiddleware)

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: router,
	}
}
