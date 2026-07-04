package httpapi

import (
	"context"

	"media-manager/internal/downloadrouting"
)

func (s *Server) enabledDownloadClientProtocols(ctx context.Context) (map[string]struct{}, error) {
	clients, err := s.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		return nil, err
	}
	return downloadrouting.Protocols(clients), nil
}

func releaseGrabDisabledReason(protocol string, protocols map[string]struct{}) *string {
	if downloadrouting.HasProtocol(protocols, protocol) {
		return nil
	}
	message := downloadrouting.MissingClientMessage(protocol)
	return &message
}
