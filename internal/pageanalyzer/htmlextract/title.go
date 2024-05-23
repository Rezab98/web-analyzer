package htmlextract

func (h *HTMLExtractor) Title() string {
	return h.goQueryDoc.Find("title").Text()
}
