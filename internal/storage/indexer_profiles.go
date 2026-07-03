package storage

type IndexerAppProfile struct {
	ID                      string
	Name                    string
	EnableRSS               bool
	EnableAutomaticSearch   bool
	EnableInteractiveSearch bool
}

func DefaultIndexerAppProfiles() []IndexerAppProfile {
	return []IndexerAppProfile{
		{
			ID:                      "default",
			Name:                    "Default",
			EnableRSS:               true,
			EnableAutomaticSearch:   true,
			EnableInteractiveSearch: true,
		},
	}
}
