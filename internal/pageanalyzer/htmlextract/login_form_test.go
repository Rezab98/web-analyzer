package htmlextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasLoginForm(t *testing.T) {
	htmlContentWithLogin := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<form>
				<input type="text" name="username">
				<input type="password" name="password">
				<input type="submit" value="Login">
			</form>
		</body>
		</html>
	`

	htmlContentWithoutLogin := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<form>
				<input type="text" name="username">
				<input type="submit" value="Submit">
			</form>
		</body>
		</html>
	`

	extractorWithLogin, err := New([]byte(htmlContentWithLogin))
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}

	extractorWithoutLogin, err := New([]byte(htmlContentWithoutLogin))
	if err != nil {
		t.Fatalf("failed to create extractor: %v", err)
	}

	assert.True(t, extractorWithLogin.HasLoginForm())
	assert.False(t, extractorWithoutLogin.HasLoginForm())
}
