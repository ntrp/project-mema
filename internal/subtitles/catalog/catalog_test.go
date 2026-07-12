package catalog

import "testing"

const provenanceCommit = "e54edd769b7062280118a14aa0fef3808829714d"

func TestCatalogContainsExactPickerKeys(t *testing.T) {
	entries := MustAll()
	if len(entries) != 59 {
		t.Fatalf("expected 59 picker entries, got %d", len(entries))
	}
	expected := map[string]bool{}
	for _, key := range []string{"addic7ed", "animekalesi", "animetosho", "animesubinfo", "avistaz", "bayflix", "assrt", "betaseries", "bsplayer", "cinemaz", "gestdown", "greeksubs", "greeksubtitles", "hdbits", "jimaku", "hosszupuska", "karagarga", "ktuvit", "pipocas", "legendasdivx", "legendasnet", "napiprojekt", "napisy24", "nekur", "opensubtitlescom", "prijevodionline", "regielive", "soustitreseu", "subclub", "subdl", "subf2m", "subsource", "subsarr", "subssabbz", "subs4free", "subs4series", "subscenter", "subsro", "subsunacs", "subsynchro", "subtis", "subtitrarinoi", "subtitriid", "subtitulamostv", "subx", "supersubtitles", "titlovi", "titrari", "titulky", "turkcealtyaziorg", "tvsubtitles", "whisperai", "wizdom", "xsubs", "yavkanet", "yifysubtitles", "vladoonmooo", "zimuku", "mock"} {
		expected[key] = true
	}
	for _, entry := range entries {
		if !expected[entry.Key] {
			t.Fatalf("unexpected key %q", entry.Key)
		}
		delete(expected, entry.Key)
	}
	if len(expected) > 0 {
		t.Fatalf("missing keys: %#v", expected)
	}
	if _, ok := Lookup("opensubtitles"); ok {
		t.Fatalf("legacy opensubtitles alias must not be a picker entry")
	}
}

func TestCatalogProvenanceAndRuntimeStatus(t *testing.T) {
	for _, entry := range MustAll() {
		if entry.Key != "mock" && entry.ProvenanceCommit != provenanceCommit {
			t.Fatalf("%s provenance = %q", entry.Key, entry.ProvenanceCommit)
		}
		if entry.RuntimeStatus == "" || entry.RuntimeMessage == "" {
			t.Fatalf("%s runtime status/message must be honest", entry.Key)
		}
	}
	for _, key := range []string{"opensubtitlescom", "mock"} {
		entry, ok := Lookup(key)
		if !ok || entry.RuntimeStatus != RuntimeSupported {
			t.Fatalf("%s should be supported, got %#v", key, entry.RuntimeStatus)
		}
	}
	entry, ok := Lookup("bsplayer")
	if !ok || entry.RuntimeStatus != RuntimeCatalogOnly {
		t.Fatalf("bsplayer should be catalog-only, got %#v", entry.RuntimeStatus)
	}
}

func TestOpenSubtitlesFieldNormalization(t *testing.T) {
	entry, ok := Lookup("opensubtitlescom")
	if !ok {
		t.Fatal("opensubtitlescom entry missing")
	}
	fields := map[string]Field{}
	for _, field := range entry.Fields {
		fields[field.Key] = field
		if field.Type == FieldAction && field.Persisted {
			t.Fatalf("action field %q must not be persisted", field.Key)
		}
	}
	assertField(t, fields["baseUrl"], FieldText, false, "base_url")
	assertField(t, fields["username"], FieldText, false, "username")
	assertField(t, fields["password"], FieldPassword, true, "password")
	assertField(t, fields["apiKey"], FieldPassword, true, "api_key")
}

func assertField(t *testing.T, field Field, fieldType FieldType, secret bool, semantic string) {
	t.Helper()
	if field.Type != fieldType || field.Secret != secret || !field.Persisted || field.SemanticKey != semantic {
		t.Fatalf("unexpected field normalization: %#v", field)
	}
}
