package jobs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

func shouldRetryAlternativeRelease(err error, attempt int) bool {
	return errors.Is(err, errRetryAlternativeRelease) && attempt < maxAutomaticGrabAttempts
}

func automaticRetryLimitReached(err error, attempt int) bool {
	return errors.Is(err, errRetryAlternativeRelease) && attempt >= maxAutomaticGrabAttempts
}

func topDecisionRejections(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) []string {
	seen := map[string]struct{}{}
	reasons := []string{}
	for _, release := range releases {
		match := decisions.EvaluateReleaseCandidateInputMatchWithLanguageContext(
			item,
			release,
			profile,
			formats,
			languages,
		)
		if match.Severity != "error" {
			continue
		}
		for _, detail := range match.Details {
			if _, ok := seen[detail]; ok {
				continue
			}
			seen[detail] = struct{}{}
			reasons = append(reasons, detail)
			if len(reasons) == 3 {
				return reasons
			}
		}
	}
	return reasons
}

func unblockedReleaseCandidates(ctx context.Context, settings *storage.SettingsStore, releases []storage.ReleaseCandidateInput) ([]storage.ReleaseCandidateInput, error) {
	filtered := make([]storage.ReleaseCandidateInput, 0, len(releases))
	for _, release := range releases {
		blocked, err := settings.ReleaseCandidateInputBlocked(ctx, release)
		if err != nil {
			return nil, fmt.Errorf("check release blocklist: %w", err)
		}
		if !blocked {
			filtered = append(filtered, release)
		}
	}
	return filtered, nil
}

func downloadClientConfig(client storage.DownloadClient) downloadclients.Config {
	return downloadclients.Config{
		Name:     client.Name,
		Type:     client.Type,
		BaseURL:  client.BaseURL,
		Username: client.Username,
		Password: client.Password,
		APIKey:   client.APIKey,
		Category: client.Category,
	}
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
