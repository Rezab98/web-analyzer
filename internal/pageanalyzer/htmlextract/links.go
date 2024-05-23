package htmlextract

import (
	"github.com/PuerkitoBio/goquery"
)

func (h *HTMLExtractor) Links() []string {
	var links []string

	h.goQueryDoc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	return links
}
