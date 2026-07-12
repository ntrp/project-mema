package security

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"testing"

	"github.com/ulikunitz/xz"
)

func TestReadArchiveReadsZipInMemory(t *testing.T) {
	data := zipBytes(t, map[string]string{"movie.en.srt": "subtitle"})
	members, err := ReadArchive("subs.zip", data, ArchiveLimits{MaxMembers: 4, MaxBytes: 1024})
	if err != nil {
		t.Fatalf("ReadArchive failed: %v", err)
	}
	if len(members) != 1 || members[0].Name != "movie.en.srt" || string(members[0].Content) != "subtitle" {
		t.Fatalf("unexpected members: %#v", members)
	}
}

func TestReadArchiveRejectsTraversal(t *testing.T) {
	data := zipBytes(t, map[string]string{"../escape.srt": "bad"})
	_, err := ReadArchive("subs.zip", data, ArchiveLimits{MaxMembers: 4, MaxBytes: 1024})
	if !errors.Is(err, ErrUnsafeArchive) {
		t.Fatalf("expected unsafe archive, got %v", err)
	}
}

func TestReadArchiveRejectsNestedArchive(t *testing.T) {
	data := zipBytes(t, map[string]string{"nested.zip": "PK"})
	_, err := ReadArchive("subs.zip", data, ArchiveLimits{MaxMembers: 4, MaxBytes: 1024})
	if !errors.Is(err, ErrUnsafeArchive) {
		t.Fatalf("expected unsafe nested archive, got %v", err)
	}
}

func TestReadArchiveRejectsOversizedZip(t *testing.T) {
	data := zipBytes(t, map[string]string{"large.srt": "0123456789"})
	_, err := ReadArchive("subs.zip", data, ArchiveLimits{MaxMembers: 4, MaxBytes: 5})
	if !errors.Is(err, ErrUnsafeArchive) {
		t.Fatalf("expected unsafe oversized archive, got %v", err)
	}
}

func TestReadArchiveRejectsSymlink(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	header := &zip.FileHeader{Name: "link.srt", Method: zip.Deflate}
	header.SetMode(os.ModeSymlink | 0o777)
	w, err := zw.CreateHeader(header)
	if err != nil {
		t.Fatalf("create symlink header: %v", err)
	}
	if _, err := w.Write([]byte("target")); err != nil {
		t.Fatalf("write symlink header: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	_, err = ReadArchive("subs.zip", buf.Bytes(), ArchiveLimits{MaxMembers: 4, MaxBytes: 1024})
	if !errors.Is(err, ErrUnsafeArchive) {
		t.Fatalf("expected unsafe symlink archive, got %v", err)
	}
}

func TestReadArchiveRejectsTooManyMembers(t *testing.T) {
	data := zipBytes(t, map[string]string{"a.srt": "a", "b.srt": "b"})
	_, err := ReadArchive("subs.zip", data, ArchiveLimits{MaxMembers: 1, MaxBytes: 1024})
	if !errors.Is(err, ErrUnsafeArchive) {
		t.Fatalf("expected unsafe member count, got %v", err)
	}
}

func TestReadArchiveReadsGzipAndXZ(t *testing.T) {
	cases := map[string][]byte{
		"movie.en.srt.gz": gzipArchiveBytes(t, []byte("gzip subtitle")),
		"movie.en.srt.xz": xzArchiveBytes(t, []byte("xz subtitle")),
	}
	for name, data := range cases {
		members, err := ReadArchive(name, data, ArchiveLimits{MaxMembers: 1, MaxBytes: 1024})
		if err != nil {
			t.Fatalf("ReadArchive(%s) failed: %v", name, err)
		}
		if len(members) != 1 || members[0].Name != "movie.en.srt" || len(members[0].Content) == 0 {
			t.Fatalf("members = %#v", members)
		}
	}
}

func TestReadArchiveRecognizesMalformedRARAsRAR(t *testing.T) {
	_, err := ReadArchive("subs.rar", []byte("Rar!\x1a\x07\x00bad"), ArchiveLimits{})
	if err == nil || errors.Is(err, ErrUnsupportedArchive) {
		t.Fatalf("expected rar parse error, got %v", err)
	}
}

func TestReadArchiveRejectsUnsupportedFormats(t *testing.T) {
	_, err := ReadArchive("subs.7z", []byte("not supported"), ArchiveLimits{})
	if !errors.Is(err, ErrUnsupportedArchive) {
		t.Fatalf("expected unsupported format, got %v", err)
	}
}

func gzipArchiveBytes(t *testing.T, content []byte) []byte {
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

func xzArchiveBytes(t *testing.T, content []byte) []byte {
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

func zipBytes(t *testing.T, files map[string]string) []byte {
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
