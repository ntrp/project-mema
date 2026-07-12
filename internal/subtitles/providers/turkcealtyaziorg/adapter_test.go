package turkcealtyaziorg

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestCandidateUsesTurkishDownloadRoute(t *testing.T) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<div class="item"><span class="title">Dizi S01E02</span><a href="/subtitle/123/indirmek.html">indir</a></div>`))
	cand, ok := candidate(doc.Find("div.item"), "https://turkcealtyazi.org/find.php", "en")
	if !ok || cand.LanguageID != "tr" || !strings.Contains(cand.SourceURL, "/indirmek.html") {
		t.Fatalf("candidate=%#v ok=%v", cand, ok)
	}
	bad, _ := goquery.NewDocumentFromReader(strings.NewReader(`<a href="/movie/123">details</a>`))
	if _, ok := candidate(bad.Find("a"), "https://turkcealtyazi.org/find.php", "en"); ok {
		t.Fatal("expected non-subtitle route to be rejected")
	}
}
