package htmlextract

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func (h *HTMLExtractor) HeadingTagToTexts() map[string][]string {
	headingTagToTexts := make(map[string][]string)

	for i := 1; i <= 6; i++ {
		headingTag := fmt.Sprintf("h%d", i)
		h.goQueryDoc.Find(headingTag).Each(func(index int, item *goquery.Selection) {
			headingTagToTexts[headingTag] = append(headingTagToTexts[headingTag], item.Text())
		})
	}

	return headingTagToTexts
}
