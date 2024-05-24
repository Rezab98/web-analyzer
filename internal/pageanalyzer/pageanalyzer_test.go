package pageanalyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLVersion(t *testing.T) {
	tests := []struct {
		name        string
		htmlContent []byte
		expected    string
	}{
		{
			name: "HTML5",
			htmlContent: []byte(`
				<!DOCTYPE html>
				<html>
				<head>
					<title>Test</title>
				</head>
				<body></body>
				</html>
			`),
			expected: "HTML5",
		},
		{
			name: "XHTML 1.0",
			htmlContent: []byte(`
				<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
				<html xmlns="http://www.w3.org/1999/xhtml">
				<head>
					<title>Test</title>
				</head>
				<body></body>
				</html>
			`),
			expected: "XHTML 1.0",
		},
		{
			name: "XHTML 1.1",
			htmlContent: []byte(`
				<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
				<html xmlns="http://www.w3.org/1999/xhtml">
				<head>
					<title>Test</title>
				</head>
				<body></body>
				</html>
			`),
			expected: "XHTML 1.1",
		},
		{
			name: "HTML4",
			htmlContent: []byte(`
				<html>
				<head>
					<title>Test</title>
				</head>
				<body></body>
				</html>
			`),
			expected: "HTML4 or Earlier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTMLVersion(tt.htmlContent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		baseURL  string
		expected bool
	}{
		{
			name:     "Internal Link",
			href:     "https://example.com/path/to/resource",
			baseURL:  "https://example.com",
			expected: true,
		},
		{
			name:     "External Link",
			href:     "https://other.com/path/to/resource",
			baseURL:  "https://example.com",
			expected: false,
		},
		{
			name:     "Internal Link with WWW",
			href:     "https://www.example.com/path/to/resource",
			baseURL:  "https://example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInternalLink(tt.href, tt.baseURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolveRelativeLinks(t *testing.T) {
	baseURL := "https://example.com"
	links := []string{
		"/path/to/resource",
		"https://other.com/path/to/resource",
		"relative/path",
	}

	expected := []string{
		"https://example.com/path/to/resource",
		"https://other.com/path/to/resource",
		"https://example.com/relative/path",
	}

	resolvedLinks, err := resolveRelativeLinks(links, baseURL)
	if err != nil {
		t.Fatalf("resolveRelativeLinks failed: %v", err)
	}

	assert.Equal(t, expected, resolvedLinks)
}
