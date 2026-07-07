package dlna

import (
	"context"
	"errors"
	"strconv"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/soap"
	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"
)

func (m *Manager) SOAPDispatcher() *soap.Dispatcher {
	dispatcher := soap.NewDispatcher()
	tree := m.contentTree()
	for _, prefix := range []string{"", "/dlna"} {
		dispatcher.Register(prefix+"/control/content-directory", ssdp.ContentDir, m.diagnosticActions(contentDirectoryActions(tree, m.baseURL, m.events.UpdateID)))
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
			m.recordSOAPAction(ctx, actionName, err)
			return values, err
		}
	}
	return wrapped
}

func contentDirectoryActions(tree *content.Tree, baseURL string, updateID func() int) map[string]soap.HandlerFunc {
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
			objects := content.ApplySubtitleURLs(baseURL, content.ApplyArtworkURLs(baseURL, response.Objects))
			response.UpdateID = updateID()
			payload, err := content.RenderDIDL(objects, nil)
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
			objects := content.ApplySubtitleURLs(baseURL, content.ApplyArtworkURLs(baseURL, response.Objects))
			response.UpdateID = updateID()
			payload, err := content.RenderDIDL(objects, nil)
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
