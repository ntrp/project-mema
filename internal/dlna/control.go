package dlna

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/delivery"
	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/soap"
	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"
)

func (m *Manager) SOAPDispatcher() *soap.Dispatcher {
	dispatcher := soap.NewDispatcher()
	tree := m.contentTree()
	for _, prefix := range []string{"", "/dlna"} {
		dispatcher.Register(prefix+"/control/content-directory", ssdp.ContentDir, m.diagnosticActions(contentDirectoryActions(tree, m.baseURL, m.events.UpdateID, m.rendererProfileFromContext)))
		dispatcher.Register(prefix+"/control/connection-manager", ssdp.Connection, m.diagnosticActions(connectionManagerActions(m.rendererProfileFromContext)))
		dispatcher.Register(prefix+"/control/media-receiver-registrar", "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1", m.diagnosticActions(registrarActions()))
	}
	return dispatcher
}

func (m *Manager) diagnosticActions(actions map[string]soap.HandlerFunc) map[string]soap.HandlerFunc {
	wrapped := make(map[string]soap.HandlerFunc, len(actions))
	for name, handler := range actions {
		actionName := name
		next := handler
		wrapped[actionName] = func(ctx context.Context, args map[string]string) (map[string]string, error) {
			values, err := next(ctx, args)
			m.recordSOAPAction(ctx, actionName, args, err)
			return values, err
		}
	}
	return wrapped
}

func contentDirectoryActions(tree *content.Tree, baseURL string, updateID func() int, profileForContext func(context.Context) RendererProfile) map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetSearchCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SearchCaps": "dc:title,upnp:class,upnp:genre,dc:creator,dc:date"}, nil
		},
		"GetSortCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SortCaps": "dc:title,dc:date"}, nil
		},
		"GetSystemUpdateID": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Id": strconv.Itoa(updateID())}, nil
		},
		"Browse": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			request, err := content.ParseBrowseRequest(args)
			if err != nil {
				return nil, soap.InvalidArgs(err.Error())
			}
			response, err := tree.Browse(ctx, request)
			if errors.Is(err, content.ErrObjectNotFound) {
				return nil, soap.Error{Code: 701, Description: "No Such Object"}
			}
			if err != nil {
				return nil, soap.InvalidArgs(err.Error())
			}
			responseBaseURL := contentActionBaseURL(ctx, baseURL)
			resources := contentResources(response.Objects, responseBaseURL, profileForContext(ctx))
			objects := content.ApplySubtitleURLs(responseBaseURL, content.ApplyArtworkURLs(responseBaseURL, response.Objects))
			response.UpdateID = updateID()
			payload, err := content.RenderDIDL(objects, resources)
			if err != nil {
				return nil, err
			}
			return map[string]string{
				"Result":         string(payload),
				"NumberReturned": strconv.Itoa(response.NumberReturned),
				"TotalMatches":   strconv.Itoa(response.TotalMatches),
				"UpdateID":       strconv.Itoa(response.UpdateID),
			}, nil
		},
		"Search": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			request, err := content.ParseSearchRequest(args)
			if err != nil {
				return nil, soap.InvalidArgs(err.Error())
			}
			response, err := tree.Search(ctx, request)
			if errors.Is(err, content.ErrObjectNotFound) {
				return nil, soap.Error{Code: 701, Description: "No Such Object"}
			}
			if err != nil {
				return nil, soap.InvalidArgs(err.Error())
			}
			responseBaseURL := contentActionBaseURL(ctx, baseURL)
			resources := contentResources(response.Objects, responseBaseURL, profileForContext(ctx))
			objects := content.ApplySubtitleURLs(responseBaseURL, content.ApplyArtworkURLs(responseBaseURL, response.Objects))
			response.UpdateID = updateID()
			payload, err := content.RenderDIDL(objects, resources)
			if err != nil {
				return nil, err
			}
			return map[string]string{
				"Result":         string(payload),
				"NumberReturned": strconv.Itoa(response.NumberReturned),
				"TotalMatches":   strconv.Itoa(response.TotalMatches),
				"UpdateID":       strconv.Itoa(response.UpdateID),
			}, nil
		},
	}
}

func contentResources(objects []content.Object, baseURL string, profile RendererProfile) map[string][]content.Resource {
	resources := map[string][]content.Resource{}
	for _, object := range objects {
		if object.FilePath == "" {
			continue
		}
		probe := delivery.Probe(object.FilePath)
		probe = probeWithPathContainer(probe, object.FilePath)
		directResourceURL := content.ResourceURL(baseURL, object)
		var size *int64
		if info, err := os.Stat(object.FilePath); err == nil && !info.IsDir() {
			value := info.Size()
			size = &value
		}
		itemResources := []content.Resource{}
		if profile.ID == "lg" && audioNeedsTranscode(probe) {
			itemResources = append(itemResources, content.ResourceFromDelivery(content.ResourceInput{
				URL:       directResourceURL,
				SizeBytes: size,
				Probe:     probe,
				Decision:  directDecision(),
			}))
			resources[object.ID] = itemResources
			continue
		}
		itemResources = append(itemResources, content.ResourceFromDelivery(content.ResourceInput{
			URL:       directResourceURL,
			SizeBytes: size,
			Probe:     probe,
			Decision:  directDecision(),
		}))
		decision := delivery.DecisionFromTracks(object.FilePath, probe.Tracks, nil, DeliveryClientProfile(profile))
		resourceURL := resourceURLForDecision(directResourceURL, decision)
		if decision.DeliveryProtocol == delivery.ProtocolHLS && !profile.AvoidHLS {
			itemResources = append(itemResources, content.ResourceFromDelivery(content.ResourceInput{
				URL:       resourceURL,
				SizeBytes: size,
				Probe:     probe,
				Decision:  decision,
			}))
		}
		resources[object.ID] = itemResources
	}
	return resources
}

func directDecision() delivery.Decision {
	return delivery.Decision{DeliveryProtocol: delivery.ProtocolFile, Mode: delivery.ModeDirect}
}

func matroskaAudioTranscodeDecision(probe delivery.ProbeResult) delivery.Decision {
	decision := delivery.DecisionFromTracks("", probe.Tracks, nil, delivery.ClientBrowser)
	if decision.Plan.VideoCodec == "" {
		decision.Plan.VideoCodec = "copy"
	}
	return delivery.Decision{
		DeliveryProtocol: delivery.ProtocolFile,
		Mode:             delivery.ModeTranscode,
		Plan: delivery.TranscodePlan{
			VideoCodec: decision.Plan.VideoCodec,
			AudioCodec: "aac",
		},
		Reasons: append(decision.Reasons, "lg_audio_codec_not_supported"),
	}
}

func matroskaTranscodeProbe(probe delivery.ProbeResult) delivery.ProbeResult {
	probe.Container.FormatName = stringPtr("matroska,webm")
	return probe
}

func audioNeedsTranscode(probe delivery.ProbeResult) bool {
	audio := delivery.FirstTrackByType(probe.Tracks, delivery.TrackAudio, nil)
	if audio == nil || audio.Codec == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*audio.Codec)) {
	case "", "aac", "ac3", "eac3", "mp3", "mp2":
		return false
	default:
		return true
	}
}

func stringPtr(value string) *string {
	return &value
}

func probeWithPathContainer(probe delivery.ProbeResult, filePath string) delivery.ProbeResult {
	if probe.Container.FormatName != nil && strings.TrimSpace(*probe.Container.FormatName) != "" {
		return probe
	}
	format := formatNameFromPath(filePath)
	if format == "" {
		return probe
	}
	probe.Container.FormatName = &format
	return probe
}

func formatNameFromPath(filePath string) string {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".mkv", ".webm":
		return "matroska,webm"
	case ".mp4", ".m4v", ".mov":
		return "mov,mp4,m4a,3gp,3g2,mj2"
	case ".ts", ".m2ts":
		return "mpegts"
	default:
		return ""
	}
}

func resourceURLForDecision(resourceURL string, decision delivery.Decision) string {
	if decision.DeliveryProtocol != delivery.ProtocolHLS {
		return resourceURL
	}
	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return resourceURL
	}
	values := parsed.Query()
	values.Set("mode", "hls")
	parsed.RawQuery = values.Encode()
	return parsed.String()
}

func contentActionBaseURL(ctx context.Context, fallback string) string {
	if request, ok := soap.RequestFromContext(ctx); ok {
		return requestBaseURL(request)
	}
	return fallback
}

func (m *Manager) contentTree() *content.Tree {
	if m.source != nil {
		return content.NewTree(m.source)
	}
	if m.store == nil {
		return content.NewTree(emptyLibrarySource{})
	}
	return content.NewTree(m.store)
}

type emptyLibrarySource struct{}

func (emptyLibrarySource) ListMediaItems(context.Context) ([]storage.MediaItem, error) {
	return []storage.MediaItem{}, nil
}

func connectionManagerActions(profileForContext func(context.Context) RendererProfile) map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetProtocolInfo": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Source": SourceProtocolInfoForProfile(profileForContext(ctx)), "Sink": ""}, nil
		},
		"GetCurrentConnectionIDs": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"ConnectionIDs": "0"}, nil
		},
		"GetCurrentConnectionInfo": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			connectionID, err := soap.RequiredArg(args, "ConnectionID")
			if err != nil {
				return nil, err
			}
			response, err := CurrentConnectionInfo(connectionID)
			if connectionErr, ok := err.(connectionError); ok {
				return nil, soap.Error{Code: connectionErr.code, Description: connectionErr.description}
			}
			return response, err
		},
	}
}

func registrarActions() map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"IsAuthorized": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Result": "1"}, nil
		},
		"IsValidated": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Result": "1"}, nil
		},
	}
}
