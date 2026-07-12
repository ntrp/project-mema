package titulky

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestCandidateRequiresTitulkyDownloadShape(t *testing.T) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<table><tr data-language="cs"><td class="title">Film 2024</td><td><a href="/download.php?id=12">Stahnout</a></td></tr></table>`))
	cand, ok := candidate(doc.Find("tr"), "https://www.titulky.com/", "")
	if !ok || cand.LanguageID != "cs" || cand.ReleaseName != "Film 2024" || cand.SourceURL != "https://www.titulky.com/download.php?id=12" {
		t.Fatalf("candidate=%#v ok=%v", cand, ok)
	}
	bad, _ := goquery.NewDocumentFromReader(strings.NewReader(`<a href="#">noop</a>`))
	if _, ok := candidate(bad.Find("a"), "https://www.titulky.com/", ""); ok {
		t.Fatal("expected non-download link to be rejected")
	}
}
