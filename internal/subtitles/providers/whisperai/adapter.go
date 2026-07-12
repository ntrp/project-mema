package whisperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
)

const defaultBaseURL = "http://localhost:9000"

type adapter struct{}
type ref struct{ Task, AudioLanguage, InputLanguage, ForceAudioStream, FilePath string }
type detectResponse struct {
	LanguageCode     string `json:"language_code"`
	DetectedLanguage string `json:"detected_language"`
}

type probeData struct {
	Streams []struct {
		Index int               `json:"index"`
		Tags  map[string]string `json:"tags"`
	} `json:"streams"`
	Packets []struct {
		StreamIndex int    `json:"stream_index"`
		PTSTime     string `json:"pts_time"`
	} `json:"packets"`
}

func init() { providers.Register("whisperai", adapter{}) }

func (adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if providercore.NewConfig(config).BaseURL("") == "" {
		return fmt.Errorf("%w: baseUrl is required", providercore.ErrProviderPrerequisiteMissing)
	}
	if err := runCheck(ctx, config, "ffprobe", "-version"); err != nil {
		return err
	}
	if err := runCheck(ctx, config, "ffmpeg", "-version"); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint(config, "/"), nil)
	if err != nil {
		return err
	}
	resp, err := service.DoProviderRequest(req, "whisperai", false)
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	if err != nil {
		return nil
	}
	if resp.StatusCode >= 500 {
		return fmt.Errorf("%w: Whisper endpoint status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return nil
}

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	file := firstPath(request.FilePath, request.MediaContext.File.Path)
	if file == "" {
		return nil, fmt.Errorf("%w: media file path is required", providercore.ErrProviderPrerequisiteMissing)
	}
	want := normalizeLang(request.LanguageID)
	if want == "" {
		want = "eng"
	}
	audio := firstAudioLanguage(config)
	force := ""
	if audio == "" || ambiguous[audio] {
		detected, err := detectLanguage(ctx, service, config, file)
		if err != nil || detected == "" {
			return nil, nil
		}
		audio = detected
	}
	task := "transcribe"
	if audio != want {
		task = "translate"
	}
	if task == "translate" && want != "eng" {
		return nil, nil
	}
	input := alpha3ToAlpha2[audio]
	if input == "" && want == "eng" {
		input, task = "en", "transcribe"
	}
	if input == "" {
		return nil, nil
	}
	r := ref{Task: task, AudioLanguage: audio, InputLanguage: input, ForceAudioStream: force, FilePath: file}
	b, _ := json.Marshal(r)
	return []providercore.Candidate{{ProviderName: "whisperai", LanguageID: want, Format: "srt", ReleaseName: fmt.Sprintf("%s %s audio -> %s SRT", task, audio, want), SourceURL: endpoint(config, "/asr"), SourceRef: string(b)}}, nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	ctx, cancel := timeoutContext(ctx, config, "transcriptionTimeoutSeconds", 3600)
	defer cancel()
	var r ref
	if err := json.Unmarshal([]byte(candidate.SourceRef), &r); err != nil {
		return providercore.Download{}, err
	}
	audio, err := encodeAudio(ctx, config, r.FilePath, r.ForceAudioStream)
	if err != nil {
		return providercore.Download{}, err
	}
	body, ctype, err := multipartBody("audio_file", filepath.Base(r.FilePath)+".raw", audio)
	if err != nil {
		return providercore.Download{}, err
	}
	u, _ := url.Parse(endpoint(config, "/asr"))
	q := u.Query()
	q.Set("task", r.Task)
	q.Set("language", r.InputLanguage)
	q.Set("output", "srt")
	q.Set("encode", "false")
	if providercore.NewConfig(config).BoolSetting("passVideoName") {
		q.Set("video_file", r.FilePath)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return providercore.Download{}, err
	}
	req.Header.Set("Content-Type", ctype)
	resp, err := service.DoProviderRequest(req, "whisperai", true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20))
	if err != nil {
		return providercore.Download{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("Whisper ASR status %d", resp.StatusCode)
	}
	return providercore.Download{Content: data, URL: u.String()}, nil
}

func detectLanguage(ctx context.Context, service providercore.Service, config providercore.Config, file string) (string, error) {
	ctx, cancel := timeoutContext(ctx, config, "responseTimeoutSeconds", 30)
	defer cancel()
	audio, err := encodeAudio(ctx, config, file, "")
	if err != nil {
		return "", err
	}
	body, ctype, err := multipartBody("audio_file", filepath.Base(file)+".raw", audio)
	if err != nil {
		return "", err
	}
	u, _ := url.Parse(endpoint(config, "/detect-language"))
	q := u.Query()
	q.Set("encode", "false")
	if providercore.NewConfig(config).BoolSetting("passVideoName") {
		q.Set("video_file", file)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", ctype)
	resp, err := service.DoProviderRequest(req, "whisperai", false)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var parsed detectResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&parsed); err != nil {
		return "", nil
	}
	code := strings.ToLower(strings.TrimSpace(parsed.LanguageCode))
	if code == "" || code == "und" {
		return "", nil
	}
	if len(code) == 2 {
		code = alpha2ToAlpha3[code]
	}
	if mapped := languageMapping[code]; mapped != "" {
		code = mapped
	}
	return code, nil
}

func encodeAudio(ctx context.Context, config providercore.Config, file, lang string) ([]byte, error) {
	filter := "aresample=async=1"
	if delay := audioDelay(ctx, config, file, lang); delay > 20 {
		filter = fmt.Sprintf("adelay=%d|%d,%s", delay, delay, filter)
	} else if delay < -20 {
		filter = fmt.Sprintf("atrim=start=%.3f,asetpts=PTS-STARTPTS,%s", float64(-delay)/1000, filter)
	}
	logLevel := strings.TrimSpace(providercore.NewConfig(config).StringSetting("logLevel"))
	if logLevel == "" {
		logLevel = "error"
	}
	args := []string{"-nostdin", "-loglevel", logLevel, "-i", file, "-map", "0:a:0", "-f", "s16le", "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", "-af", filter, "-"}
	if lang != "" {
		args[6] = "0:a:m:language:" + iso6392(lang)
	}
	out, err := command(ctx, config, "ffmpeg", args...)
	if err != nil {
		return nil, fmt.Errorf("%w: ffmpeg failed: %v", providercore.ErrProviderPrerequisiteMissing, err)
	}
	return out, nil
}

func audioDelay(ctx context.Context, config providercore.Config, file, lang string) int {
	out, err := command(ctx, config, "ffprobe", "-v", "error", "-select_streams", "a", "-read_intervals", "%+30", "-show_entries", "stream=index:stream_tags=language:packet=stream_index,pts_time", "-of", "json", file)
	if err != nil {
		return 0
	}
	var p probeData
	if json.Unmarshal(out, &p) != nil || len(p.Streams) == 0 {
		return 0
	}
	target := p.Streams[0].Index
	if lang != "" {
		short := lang[:2]
		for _, s := range p.Streams {
			if strings.Contains(strings.ToLower(s.Tags["language"]), short) {
				target = s.Index
				break
			}
		}
	}
	for _, pkt := range p.Packets {
		if pkt.StreamIndex == target {
			f, _ := strconv.ParseFloat(pkt.PTSTime, 64)
			return int(f * 1000)
		}
	}
	return 0
}

func runCheck(ctx context.Context, config providercore.Config, name string, args ...string) error {
	if _, err := command(ctx, config, name, args...); err != nil {
		return fmt.Errorf("%w: %s is required", providercore.ErrProviderPrerequisiteMissing, name)
	}
	return nil
}

func command(ctx context.Context, config providercore.Config, name string, args ...string) ([]byte, error) {
	if config.CommandRunner == nil {
		return nil, fmt.Errorf("no command runner")
	}
	return config.CommandRunner(ctx, name, args...)
}
func multipartBody(field, filename string, data []byte) (io.Reader, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		return nil, "", err
	}
	if _, err = part.Write(data); err != nil {
		return nil, "", err
	}
	err = w.Close()
	return &b, w.FormDataContentType(), err
}
func endpoint(config providercore.Config, p string) string {
	return strings.TrimRight(providercore.NewConfig(config).BaseURL(defaultBaseURL), "/") + p
}
