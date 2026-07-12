package providercore

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func ParseHTML(data []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewReader(data))
}

func SelectionText(selection *goquery.Selection) string {
	return selection.First().Text()
}

func SelectionAttr(selection *goquery.Selection, name string) string {
	value, _ := selection.First().Attr(name)
	return value
}
