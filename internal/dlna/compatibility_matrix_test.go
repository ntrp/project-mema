package dlna

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna/content"
)

type compatibilityFixture struct {
	Name           string `json:"name"`
	Delivery       string `json:"delivery"`
	Container      string `json:"container"`
	VideoCodec     string `json:"videoCodec"`
	AudioCodec     string `json:"audioCodec"`
	SubtitleFormat string `json:"subtitleFormat"`
	ImageFormat    string `json:"imageFormat"`
	Language       string `json:"language"`
	Height         int32  `json:"height"`
	HDR            bool   `json:"hdr"`
}

func TestCompatibilityMatrixProfiles(t *testing.T) {
	tests := []struct {
		client     string
		userAgent  string
		profileID  string
		wantFirst  string
		noEventing bool
	}{
		{client: "VLC", userAgent: "VLC/3.0.20 LibVLC", profileID: "vlc", wantFirst: "video/mp4"},
		{client: "Kodi", userAgent: "Kodi/21.0", profileID: "kodi", wantFirst: "video/mp4"},
		{client: "Samsung", userAgent: "Samsung Tizen DLNADOC/1.50", profileID: "samsung", wantFirst: "video/mp4"},
		{client: "LG", userAgent: "LG webOS TV", profileID: "lg", wantFirst: "video/mp4"},
		{client: "Sony", userAgent: "Sony BRAVIA", profileID: "sony", wantFirst: "video/mp4"},
		{client: "BubbleUPnP", userAgent: "BubbleUPnP/4.3", profileID: "bubbleupnp", wantFirst: "video/mp4"},
		{
			client: "iOS/tvOS", userAgent: "Chromecast", profileID: "chromecast",
			wantFirst: "application/vnd.apple.mpegurl", noEventing: true,
		},
	}

	for _, test := range tests {
		t.Run(test.client, func(t *testing.T) {
			profile := MatchRendererProfile(RendererRequest{UserAgent: test.userAgent}, nil)
			requireProfileBasics(t, profile, test.profileID, test.wantFirst, test.noEventing)
		})
	}
}

func TestCompatibilitySeededCatalogRepresentativeFamilies(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	fixtures := compatibilityFixtures(t)
	tests := []struct {
		name       string
		request    RendererRequest
		profileID  string
		fixture    string
		wantMode   delivery.Mode
		wantProto  delivery.Protocol
		wantFirst  string
		noEventing bool
	}{
		{name: "LG", request: RendererRequest{UserAgent: "LG webOS TV"}, profileID: "lg", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Samsung", request: RendererRequest{UserAgent: "Samsung Tizen DLNADOC"}, profileID: "samsung", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Sony", request: RendererRequest{UserAgent: "Sony BRAVIA"}, profileID: "sony", fixture: "mpegts-h264-eac3-tv", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Panasonic", request: RendererRequest{UserAgent: "Panasonic VIERA"}, profileID: "panasonic-tv", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Philips", request: RendererRequest{UserAgent: "Philips PUS 6500"}, profileID: "philips-tv", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Roku", request: RendererRequest{UserAgent: "Roku TV"}, profileID: "roku-tv", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Chromecast", request: RendererRequest{Headers: http.Header{"X-Cast": []string{"Google Cast"}}}, profileID: "chromecast", fixture: "hdr-hevc-dts-remux", wantMode: delivery.ModeTranscode, wantProto: delivery.ProtocolHLS, wantFirst: "application/vnd.apple.mpegurl", noEventing: true},
		{name: "VLC", request: RendererRequest{UserAgent: "VLC/3.0.20 LibVLC"}, profileID: "vlc", fixture: "avi-mpeg2-ac3-legacy", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "BubbleUPnP", request: RendererRequest{FriendlyName: "BubbleUPnP Renderer"}, profileID: "bubbleupnp", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Kodi", request: RendererRequest{FriendlyName: "Kodi"}, profileID: "kodi", fixture: "hdr-hevc-dts-remux", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Windows Media Player", request: RendererRequest{UserAgent: "Windows Media Player"}, profileID: "windows-media-player", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Xbox", request: RendererRequest{UserAgent: "Xbox One"}, profileID: "xbox-one", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
		{name: "Android", request: RendererRequest{UserAgent: "Android AOSP DLNA"}, profileID: "android-generic", fixture: "baseline-mp4-h264-aac", wantMode: delivery.ModeDirect, wantProto: delivery.ProtocolFile, wantFirst: "video/mp4"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := manager.ExplainRendererProfile(ctx, test.request)
			requireProfileBasics(t, match.Profile, test.profileID, test.wantFirst, test.noEventing)
			decision := EvaluateRendererCapability(match.Profile, probeFromFixture(fixtures[test.fixture]))
			if decision.Decision.Mode != test.wantMode || decision.Decision.DeliveryProtocol != test.wantProto {
				t.Fatalf("decision = %#v, trace=%#v", decision.Decision, decision.Trace)
			}
		})
	}
}

func TestCompatibilityMediaFixturesCoverCodecMatrix(t *testing.T) {
	fixtures := compatibilityFixtures(t)
	seen := map[string]bool{}
	for _, item := range fixtures {
		markFixtureCoverage(seen, item)
	}
	for _, want := range []string{
		"delivery:direct", "delivery:remux", "delivery:transcode", "delivery:subtitle", "delivery:image",
		"container:mp4", "container:mkv", "container:avi", "container:mpegts",
		"subtitle:srt", "subtitle:vtt", "subtitle:ass", "subtitle:ssa",
		"image:jpeg", "image:png", "codec:aac", "codec:ac3", "codec:eac3", "codec:dts",
		"codec:h264", "codec:hevc", "codec:av1", "hdr",
	} {
		if !seen[want] {
			t.Fatalf("fixture missing %s: %#v", want, fixtures)
		}
	}
}

func TestCompatibilityDIDLGoldenFamilies(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	tests := []struct {
		name      string
		request   RendererRequest
		contains  []string
		forbidden []string
	}{
		{
			name: "vlc", request: RendererRequest{UserAgent: "VLC/3.0.20"},
			contains: []string{"application/x-subrip", "text/vtt", "albumArtURI", "upnp:genre"},
		},
		{
			name: "chromecast", request: RendererRequest{UserAgent: "Chromecast"},
			contains:  []string{"text/vtt"},
			forbidden: []string{"application/x-subrip", "albumArtURI", "dc:date", "upnp:genre"},
		},
		{
			name: "samsung", request: RendererRequest{UserAgent: "Samsung Tizen DLNADOC"},
			contains:  []string{"application/x-subrip", "albumArtURI"},
			forbidden: []string{"text/vtt", "dc:date", "upnp:genre"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := manager.ExplainRendererProfile(ctx, test.request).Profile
			payload := renderCompatibilityDIDL(t, profile)
			requireContainsAll(t, payload, test.contains)
			requireContainsNone(t, payload, test.forbidden)
		})
	}
}

func compatibilityFixtures(t *testing.T) map[string]compatibilityFixture {
	t.Helper()
	payload, err := os.ReadFile("testdata/compatibility_media.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Cases []compatibilityFixture `json:"cases"`
	}
	if err := json.Unmarshal(payload, &fixture); err != nil {
		t.Fatal(err)
	}
	fixtures := map[string]compatibilityFixture{}
	for _, item := range fixture.Cases {
		if item.Name == "" {
			t.Fatalf("fixture missing name: %#v", item)
		}
		if _, exists := fixtures[item.Name]; exists {
			t.Fatalf("duplicate fixture %q", item.Name)
		}
		fixtures[item.Name] = item
	}
	return fixtures
}

func requireProfileBasics(t *testing.T, profile RendererProfile, wantID string, wantFirst string, noEventing bool) {
	t.Helper()
	if profile.ID != wantID {
		t.Fatalf("profile = %q, want %q", profile.ID, wantID)
	}
	protocols := SourceProtocolInfosForProfile(profile)
	if len(protocols) == 0 || !strings.Contains(protocols[0], wantFirst) {
		t.Fatalf("protocols[0] = %q, want %q first", protocols[0], wantFirst)
	}
	if profile.DisableEventing != noEventing {
		t.Fatalf("DisableEventing = %v, want %v", profile.DisableEventing, noEventing)
	}
}

func probeFromFixture(fixture compatibilityFixture) delivery.ProbeResult {
	format := containerFormatName(fixture.Container)
	tracks := []delivery.Track{}
	if fixture.VideoCodec != "" {
		video := delivery.Track{Type: delivery.TrackVideo, Codec: stringPtr(fixture.VideoCodec)}
		if fixture.Height > 0 {
			video.Height = int32Ptr(fixture.Height)
		}
		if fixture.HDR {
			video.Profile = stringPtr("main10")
			video.PixelFormat = stringPtr("yuv420p10le")
		}
		tracks = append(tracks, video)
	}
	if fixture.AudioCodec != "" {
		tracks = append(tracks, delivery.Track{Type: delivery.TrackAudio, Codec: stringPtr(fixture.AudioCodec)})
	}
	return delivery.ProbeResult{Container: delivery.Container{FormatName: &format}, Tracks: tracks}
}

func containerFormatName(container string) string {
	switch container {
	case "mkv", "webm":
		return "matroska,webm"
	case "jpg":
		return "jpeg"
	default:
		return container
	}
}

func markFixtureCoverage(seen map[string]bool, item compatibilityFixture) {
	seen["delivery:"+item.Delivery] = true
	seen["container:"+item.Container] = true
	if item.SubtitleFormat != "" {
		seen["subtitle:"+item.SubtitleFormat] = true
	}
	if item.ImageFormat != "" {
		seen["image:"+item.ImageFormat] = true
	}
	if item.VideoCodec != "" {
		seen["codec:"+item.VideoCodec] = true
	}
	if item.AudioCodec != "" {
		seen["codec:"+item.AudioCodec] = true
	}
	if item.HDR {
		seen["hdr"] = true
	}
}

func renderCompatibilityDIDL(t *testing.T, profile RendererProfile) string {
	t.Helper()
	date := "2026-07-08"
	artwork := "http://127.0.0.1:18080/dlna/artwork/movie"
	object := content.Object{
		ID: "movie:1", ParentID: "movies", Title: "Scenario Movie", Kind: content.ObjectItem,
		Class: "object.item.videoItem.movie", Date: &date, Genres: []string{"Drama"}, Artwork: &artwork,
		Subtitles: []content.Subtitle{
			{URL: "http://127.0.0.1:18080/dlna/subtitle/movie.srt", Format: "srt"},
			{URL: "http://127.0.0.1:18080/dlna/subtitle/movie.vtt", Format: "vtt"},
			{URL: "http://127.0.0.1:18080/dlna/subtitle/movie.ass", Format: "ass"},
		},
	}
	payload, err := content.RenderDIDLWithOptions([]content.Object{object}, nil, DIDLOptionsForProfile(profile))
	if err != nil {
		t.Fatal(err)
	}
	return string(payload)
}

func requireContainsAll(t *testing.T, value string, wants []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(value, want) {
			t.Fatalf("missing %q:\n%s", want, value)
		}
	}
}

func requireContainsNone(t *testing.T, value string, forbidden []string) {
	t.Helper()
	for _, item := range forbidden {
		if strings.Contains(value, item) {
			t.Fatalf("contains %q:\n%s", item, value)
		}
	}
}

func int32Ptr(value int32) *int32 {
	return &value
}
