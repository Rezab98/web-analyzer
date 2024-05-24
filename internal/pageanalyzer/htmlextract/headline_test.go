package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadingTagToTexts(t *testing.T) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<h1>Heading 1</h1>
			<h2>Heading 2</h2>
			<h2>Another Heading 2</h2>
			<h3>Heading 3</h3>
		</body>
		</html>
	`

	extractor, err := New([]byte(htmlContent))
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}

	headings := extractor.HeadingTagToTexts()

	expectedHeadings := map[string][]string{
		"h1": {"Heading 1"},
		"h2": {"Heading 2", "Another Heading 2"},
		"h3": {"Heading 3"},
	}

	assert.Equal(t, expectedHeadings, headings)
}
