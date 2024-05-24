package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page Title</title>
		</head>
		<body>
			<h1>Heading 1</h1>
		</body>
		</html>
	`

	extractor, err := New([]byte(htmlContent))
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}

	title := extractor.Title()
	expectedTitle := "Test Page Title"

	assert.Equal(t, expectedTitle, title)
}
