package htmlextract

import (
	"github.com/PuerkitoBio/goquery"
)

func (h *HTMLExtractor) HasLoginForm() bool {
	var hasLoginForm bool

	h.goQueryDoc.Find("input[type='password']").Each(func(index int, item *goquery.Selection) {
		hasLoginForm = true
	})

	return hasLoginForm
}
