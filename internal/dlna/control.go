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
		dispatcher.Register(prefix+"/control/content-directory", ssdp.ContentDir, contentDirectoryActions(tree))
		dispatcher.Register(prefix+"/control/connection-manager", ssdp.Connection, connectionManagerActions())
		dispatcher.Register(prefix+"/control/media-receiver-registrar", "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1", registrarActions())
	}
	return dispatcher
}

func contentDirectoryActions(tree *content.Tree) map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetSearchCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SearchCaps": "dc:title,upnp:class,upnp:genre,dc:creator,dc:date"}, nil
		},
		"GetSortCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SortCaps": "dc:title,dc:date"}, nil
		},
		"GetSystemUpdateID": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Id": "0"}, nil
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
			payload, err := content.RenderDIDL(response.Objects, nil)
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
			payload, err := content.RenderDIDL(response.Objects, nil)
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

func connectionManagerActions() map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetProtocolInfo": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Source": SourceProtocolInfo(), "Sink": ""}, nil
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
