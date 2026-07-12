package providercore

import "testing"

func TestHTMLHelpersParseTextAndAttributes(t *testing.T) {
	doc, err := ParseHTML([]byte(`<html><body><a href="/subtitle"> Subtitle </a></body></html>`))
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}
	link := doc.Find("a")
	if SelectionText(link) != " Subtitle " {
		t.Fatalf("text = %q", SelectionText(link))
	}
	if SelectionAttr(link, "href") != "/subtitle" {
		t.Fatalf("href = %q", SelectionAttr(link, "href"))
	}
	if SelectionAttr(link, "missing") != "" {
		t.Fatalf("missing attr = %q", SelectionAttr(link, "missing"))
	}
}
