package metadata

import "context"

type Provider interface {
	Search(ctx context.Context, request SearchRequest) ([]SearchResult, error)
	Details(ctx context.Context, request DetailsRequest) (Details, error)
}

type DiscoverProvider interface {
	Discover(ctx context.Context, request DiscoverRequest) ([]SearchResult, error)
}

type tmdbProvider struct {
	service *Service
	config  Config
}

func (p tmdbProvider) Search(ctx context.Context, request SearchRequest) ([]SearchResult, error) {
	if request.MediaType != "movie" && request.MediaType != "serie" {
		return nil, nil
	}
	return p.service.searchTMDB(ctx, p.config, request)
}

func (p tmdbProvider) Discover(ctx context.Context, request DiscoverRequest) ([]SearchResult, error) {
	return p.service.discoverTMDB(ctx, p.config, request)
}

func (p tmdbProvider) Details(ctx context.Context, request DetailsRequest) (Details, error) {
	return p.service.detailsTMDB(ctx, p.config, request)
}

type tvdbProvider struct {
	service *Service
	config  Config
}

func (p tvdbProvider) Search(ctx context.Context, request SearchRequest) ([]SearchResult, error) {
	return p.service.searchTVDB(ctx, p.config, request)
}

func (p tvdbProvider) Details(ctx context.Context, request DetailsRequest) (Details, error) {
	return p.service.detailsTVDB(ctx, p.config, request)
}
