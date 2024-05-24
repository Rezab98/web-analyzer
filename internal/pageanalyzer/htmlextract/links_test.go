package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinks(t *testing.T) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<a href="http://example.com">Example</a>
			<a href="/relative">Relative</a>
			<a href="http://another.com">Another Example</a>
		</body>
		</html>
	`

	extractor, err := New([]byte(htmlContent))
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}

	links := extractor.Links()

	expectedLinks := []string{
		"http://example.com",
		"/relative",
		"http://another.com",
	}

	assert.ElementsMatch(t, expectedLinks, links)
}
