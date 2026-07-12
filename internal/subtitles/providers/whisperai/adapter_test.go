package whisperai

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type httpService struct{ client *http.Client }

func (s httpService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) {
	return s.client.Do(req)
}

func TestSearchDetectsLanguageAndBuildsTranscribeCandidate(t *testing.T) {
	var detectPath, contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		detectPath = r.URL.Path
		contentType = r.Header.Get("Content-Type")
		_, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"language_code":"es","detected_language":"Spanish"}`))
	}))
	defer server.Close()
	config := providercore.Config{BaseURL: server.URL, CommandRunner: whisperRunner(t)}
	got, err := adapter{}.Search(context.Background(), httpService{server.Client()}, config, providercore.SearchRequest{MediaType: "movie", FilePath: "/video.mkv", LanguageID: "es"})
	if err != nil {
		t.Fatal(err)
	}
	if detectPath != "/detect-language" || !strings.HasPrefix(contentType, "multipart/form-data") {
		t.Fatalf("detect request = %s %s", detectPath, contentType)
	}
	if len(got) != 1 || got[0].LanguageID != "spa" || !strings.Contains(got[0].ReleaseName, "transcribe spa") {
		t.Fatalf("candidates = %#v", got)
	}
}

func TestSearchOnlyTranslatesToEnglish(t *testing.T) {
	config := providercore.Config{BaseURL: "http://localhost", Settings: map[string]providercore.SettingValue{"audioLanguage": strSetting("spa")}, CommandRunner: whisperRunner(t)}
	got, err := adapter{}.Search(context.Background(), httpService{http.DefaultClient}, config, providercore.SearchRequest{FilePath: "/video.mkv", LanguageID: "fr"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("non-English translation candidates = %#v", got)
	}
	got, err = adapter{}.Search(context.Background(), httpService{http.DefaultClient}, config, providercore.SearchRequest{FilePath: "/video.mkv", LanguageID: "en"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !strings.Contains(got[0].ReleaseName, "translate spa audio -> eng") {
		t.Fatalf("English candidate = %#v", got)
	}
}

func TestDownloadPostsMultipartASRWithSettings(t *testing.T) {
	var path, query, contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.RawQuery
		contentType = r.Header.Get("Content-Type")
		_, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nHola"))
	}))
	defer server.Close()
	config := providercore.Config{BaseURL: server.URL, Settings: map[string]providercore.SettingValue{"passVideoName": boolSetting(true), "responseTimeoutSeconds": numberSetting(5), "transcriptionTimeoutSeconds": numberSetting(10), "logLevel": strSetting("warning")}, CommandRunner: whisperRunner(t)}
	candidate := providercore.Candidate{SourceRef: `{"Task":"transcribe","AudioLanguage":"spa","InputLanguage":"es","FilePath":"/video.mkv"}`}
	got, err := adapter{}.Download(context.Background(), httpService{server.Client()}, config, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/asr" || !strings.Contains(query, "task=transcribe") || !strings.Contains(query, "language=es") || !strings.Contains(query, "video_file=%2Fvideo.mkv") {
		t.Fatalf("asr request = %s?%s", path, query)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data") || !strings.Contains(string(got.Content), "Hola") {
		t.Fatalf("content type/content = %q %q", contentType, got.Content)
	}
}

func whisperRunner(t *testing.T) providercore.CommandRunner {
	t.Helper()
	return func(_ context.Context, name string, args ...string) ([]byte, error) {
		joined := strings.Join(args, " ")
		switch name {
		case "ffprobe":
			return []byte(`{"streams":[{"index":1,"tags":{"language":"spa"}}],"packets":[{"stream_index":1,"pts_time":"0.000"}]}`), nil
		case "ffmpeg":
			if !strings.Contains(joined, "-f s16le") || !strings.Contains(joined, "-loglevel") {
				t.Fatalf("ffmpeg args = %v", args)
			}
			return []byte("audio-bytes"), nil
		default:
			t.Fatalf("unexpected command %s", name)
		}
		return nil, nil
	}
}

func strSetting(s string) providercore.SettingValue {
	return providercore.SettingValue{StringValue: &s}
}
func boolSetting(b bool) providercore.SettingValue {
	return providercore.SettingValue{BooleanValue: &b}
}
func numberSetting(n float64) providercore.SettingValue {
	return providercore.SettingValue{NumberValue: &n}
}
