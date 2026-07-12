package httpapi

import "media-manager/internal/subtitles/catalog"

func subtitleProviderCatalogEntry(entry catalog.Entry) SubtitleProviderCatalogEntry {
	return SubtitleProviderCatalogEntry{
		Key:              entry.Key,
		DisplayName:      entry.DisplayName,
		ProvenanceCommit: subtitleCatalogStringPtr(entry.ProvenanceCommit),
		RuntimeStatus:    SubtitleProviderRuntimeStatus(entry.RuntimeStatus),
		RuntimeMessage:   entry.RuntimeMessage,
		MediaTypes:       entry.MediaTypes,
		Dependencies:     subtitleProviderDependencies(entry.Dependencies),
		Warning:          subtitleCatalogStringPtr(entry.Warning),
		OutboundPolicy:   subtitleProviderOutboundPolicy(entry.OutboundPolicy),
		Fields:           subtitleProviderFields(entry.Fields),
	}
}

func subtitleProviderDependencies(deps catalog.Dependencies) SubtitleProviderDependencies {
	return SubtitleProviderDependencies{
		Captcha:           optionalTrueBool(deps.Captcha),
		AntiCaptcha:       optionalTrueBool(deps.AntiCaptcha),
		Archive:           optionalTrueBool(deps.Archive),
		Ffmpeg:            optionalTrueBool(deps.FFmpeg),
		Ffprobe:           optionalTrueBool(deps.FFprobe),
		Anidb:             optionalTrueBool(deps.AniDB),
		ArrHistory:        optionalTrueBool(deps.ArrHistory),
		LocalHttpEndpoint: optionalTrueBool(deps.LocalHTTPEndpoint),
	}
}

func subtitleProviderOutboundPolicy(policy catalog.OutboundPolicy) SubtitleProviderOutboundPolicy {
	return SubtitleProviderOutboundPolicy{
		AllowedBaseHosts:     optionalStrings(policy.AllowedBaseHosts),
		AllowedDownloadHosts: optionalStrings(policy.AllowedDownloadHosts),
		AllowLocalHosts:      optionalTrueBool(policy.AllowLocalHosts),
	}
}

func subtitleProviderFields(fields []catalog.Field) []SubtitleProviderField {
	result := make([]SubtitleProviderField, 0, len(fields))
	for _, field := range fields {
		result = append(result, SubtitleProviderField{
			Key:         field.Key,
			Label:       field.Label,
			Type:        SubtitleProviderFieldType(field.Type),
			Secret:      optionalTrueBool(field.Secret),
			Required:    optionalTrueBool(field.Required),
			Persisted:   field.Persisted,
			SemanticKey: subtitleCatalogStringPtr(field.SemanticKey),
			Options:     optionalStrings(field.Options),
		})
	}
	return result
}

func optionalTrueBool(value bool) *bool {
	if !value {
		return nil
	}
	return &value
}

func subtitleCatalogStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func optionalStrings(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	return &values
}
