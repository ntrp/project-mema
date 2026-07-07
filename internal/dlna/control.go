package dlna

import (
	"context"

	"media-manager/internal/dlna/soap"
	"media-manager/internal/dlna/ssdp"
)

func SOAPDispatcher() *soap.Dispatcher {
	dispatcher := soap.NewDispatcher()
	for _, prefix := range []string{"", "/dlna"} {
		dispatcher.Register(prefix+"/control/content-directory", ssdp.ContentDir, contentDirectoryActions())
		dispatcher.Register(prefix+"/control/connection-manager", ssdp.Connection, connectionManagerActions())
		dispatcher.Register(prefix+"/control/media-receiver-registrar", "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1", registrarActions())
	}
	return dispatcher
}

func contentDirectoryActions() map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetSearchCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SearchCaps": "dc:title,upnp:class,upnp:genre,dc:creator"}, nil
		},
		"GetSortCapabilities": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"SortCaps": "dc:title,dc:date"}, nil
		},
		"GetSystemUpdateID": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Id": "0"}, nil
		},
		"Browse": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return nil, soap.Error{Code: 501, Description: "Action Failed"}
		},
		"Search": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return nil, soap.Error{Code: 501, Description: "Action Failed"}
		},
	}
}

func connectionManagerActions() map[string]soap.HandlerFunc {
	return map[string]soap.HandlerFunc{
		"GetProtocolInfo": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Source": "", "Sink": ""}, nil
		},
		"GetCurrentConnectionIDs": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"ConnectionIDs": ""}, nil
		},
		"GetCurrentConnectionInfo": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			if _, err := soap.RequiredArg(args, "ConnectionID"); err != nil {
				return nil, err
			}
			return nil, soap.Error{Code: 706, Description: "No Such Connection"}
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
