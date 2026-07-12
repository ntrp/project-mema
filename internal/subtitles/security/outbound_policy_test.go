package security

import (
	"errors"
	"testing"
)

func TestValidateProviderURLAllowsOpenSubtitlesHosts(t *testing.T) {
	if err := ValidateProviderURL("opensubtitlescom", "https://api.opensubtitles.com/api/v1/subtitles", false); err != nil {
		t.Fatalf("expected API URL to be allowed: %v", err)
	}
	if err := ValidateProviderURL("opensubtitlescom", "https://dl.opensubtitles.com/en/download/file", true); err != nil {
		t.Fatalf("expected download URL to be allowed: %v", err)
	}
}

func TestValidateProviderURLRejectsArbitraryDownloadHost(t *testing.T) {
	err := ValidateProviderURL("opensubtitlescom", "https://evil.example/download.srt", true)
	if !errors.Is(err, ErrOutboundURLBlocked) {
		t.Fatalf("expected blocked outbound URL, got %v", err)
	}
}

func TestValidateProviderURLRejectsPrivateHosts(t *testing.T) {
	for _, rawURL := range []string{
		"http://127.0.0.1/subtitle.srt",
		"http://10.1.2.3/subtitle.srt",
		"http://169.254.1.1/subtitle.srt",
		"http://localhost/subtitle.srt",
	} {
		err := ValidateProviderURL("opensubtitlescom", rawURL, false)
		if !errors.Is(err, ErrOutboundURLBlocked) {
			t.Fatalf("expected %s to be blocked, got %v", rawURL, err)
		}
	}
}

func TestValidateProviderURLAllowsLocalOnlyWhenCatalogAllowsIt(t *testing.T) {
	if err := ValidateProviderURL("whisperai", "http://127.0.0.1:9000/transcribe", false); err != nil {
		t.Fatalf("expected whisper local endpoint to be allowed: %v", err)
	}
}

func TestValidateProviderURLFailsClosedForUnknownDownloadHosts(t *testing.T) {
	err := ValidateProviderURL("addic7ed", "https://addic7ed.example/archive.zip", true)
	if !errors.Is(err, ErrDownloadHostClosed) {
		t.Fatalf("expected fail-closed download host error, got %v", err)
	}
}

func TestValidateRedirectRejectsHostChange(t *testing.T) {
	err := ValidateRedirect("opensubtitlescom", "https://dl.opensubtitles.com/file.srt", "https://attacker.example/file.srt", true)
	if !errors.Is(err, ErrOutboundURLBlocked) {
		t.Fatalf("expected redirect host change to be blocked, got %v", err)
	}
}
