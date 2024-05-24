package htmlextract

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type HTMLExtractor struct {
	goQueryDoc *goquery.Document
}

func New(htmlContent []byte) (*HTMLExtractor, error) {
	goQueryDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("goquery new document from reader failed: %v", err)
	}

	return &HTMLExtractor{goQueryDoc: goQueryDoc}, nil
}
