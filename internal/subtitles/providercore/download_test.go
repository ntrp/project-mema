package providercore

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/ulikunitz/xz"

	"media-manager/internal/subtitles/security"
)

func TestExtractSubtitleSelectsBestArchiveMember(t *testing.T) {
	data := downloadZip(t, map[string]string{
		"sample.forced.srt": "x",
		"sample.en.srt":     "subtitle content",
		"notes.txt":         "ignore",
	})
	member, err := ExtractSubtitle("subs.zip", data, security.ArchiveLimits{MaxMembers: 4, MaxBytes: 1024})
	if err != nil {
		t.Fatalf("ExtractSubtitle failed: %v", err)
	}
	if member.Name != "sample.en.srt" || string(member.Content) != "subtitle content" {
		t.Fatalf("member = %#v", member)
	}
}

func TestExtractSubtitleDecodesGzipAndXZ(t *testing.T) {
	cases := map[string][]byte{
		"subtitle.srt.gz": gzipBytes(t, []byte("gzip subtitle")),
		"subtitle.srt.xz": xzBytes(t, []byte("xz subtitle")),
	}
	for name, data := range cases {
		member, err := ExtractSubtitle(name, data, security.ArchiveLimits{MaxBytes: 1024})
		if err != nil {
			t.Fatalf("ExtractSubtitle(%s) failed: %v", name, err)
		}
		if string(member.Content) == "" || member.Name != "subtitle.srt" {
			t.Fatalf("member = %#v", member)
		}
	}
}

func TestExtractSubtitleRejectsOversizedCompressedPayload(t *testing.T) {
	_, err := ExtractSubtitle("subtitle.srt.gz", gzipBytes(t, []byte("0123456789")), security.ArchiveLimits{MaxBytes: 5})
	if err == nil {
		t.Fatal("expected size error")
	}
}

func TestExtractSubtitleHandlesRawAndInvalidCompressedPayloads(t *testing.T) {
	member, err := ExtractSubtitle("subtitle.srt", []byte("raw"), security.ArchiveLimits{})
	if err != nil {
		t.Fatalf("raw ExtractSubtitle failed: %v", err)
	}
	if member.Name != "subtitle.srt" || string(member.Content) != "raw" {
		t.Fatalf("member = %#v", member)
	}
	if _, err := ExtractSubtitle("subtitle.srt.gz", []byte("not gzip"), security.ArchiveLimits{}); err == nil {
		t.Fatal("expected invalid gzip error")
	}
}

func TestBestSubtitleMemberRejectsArchivesWithoutSubtitles(t *testing.T) {
	_, err := BestSubtitleMember([]security.ArchiveMember{{Name: "readme.txt", Content: []byte("notes")}})
	if err == nil {
		t.Fatal("expected no subtitle error")
	}
	if stripDownloadSuffix("subtitle.srt") != "subtitle.srt" {
		t.Fatalf("unexpected suffix strip")
	}
}

func TestDecompressSingleFallbackHelpers(t *testing.T) {
	member, err := decompressSingle("subtitle.srt.gz", gzipBytes(t, []byte("gzip subtitle")), security.ArchiveLimits{MaxBytes: 1024}, openGzipReader)
	if err != nil {
		t.Fatalf("gzip decompress failed: %v", err)
	}
	if member.Name != "subtitle.srt" || string(member.Content) != "gzip subtitle" {
		t.Fatalf("gzip member = %#v", member)
	}
	member, err = decompressSingle("subtitle.srt.xz", xzBytes(t, []byte("xz subtitle")), security.ArchiveLimits{MaxBytes: 1024}, openXZReader)
	if err != nil {
		t.Fatalf("xz decompress failed: %v", err)
	}
	if member.Name != "subtitle.srt" || string(member.Content) != "xz subtitle" {
		t.Fatalf("xz member = %#v", member)
	}
	if _, err := decompressSingle("subtitle.srt.gz", gzipBytes(t, []byte("0123456789")), security.ArchiveLimits{MaxBytes: 5}, openGzipReader); err == nil {
		t.Fatal("expected fallback size error")
	}
}

func downloadZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip member: %v", err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("write zip member: %v", err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func gzipBytes(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(content); err != nil {
		t.Fatalf("write gzip: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

func xzBytes(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw, err := xz.NewWriter(&buf)
	if err != nil {
		t.Fatalf("create xz: %v", err)
	}
	if _, err := zw.Write(content); err != nil {
		t.Fatalf("write xz: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close xz: %v", err)
	}
	return buf.Bytes()
}
