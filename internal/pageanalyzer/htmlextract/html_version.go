package htmlextract

import (
	"bytes"
	"strings"
)

func (h *HTMLExtractor) HTMLVersion() string {
	// Read only the first 1024 bytes to identify the doctype and root elements
	initialContent := bytes.ToLower(h.pageContent[:min(len(h.pageContent), 1024)])

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

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
