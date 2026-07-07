package storage

import (
	"context"
	"net"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"
)

const (
	DefaultDLNAFriendlyName            = "Mema"
	DefaultDLNAAnnounceIntervalSeconds = int32(1800)
	DefaultDLNARendererProfile         = "generic"
)

var DefaultDLNAAllowedCIDRs = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

type DLNASettings struct {
	Enabled                 bool
	FriendlyName            string
	Interfaces              []string
	AllowedCIDRs            []string
	AnnounceIntervalSeconds int32
	TranscodeEnabled        bool
	ThumbnailsEnabled       bool
	SubtitlesEnabled        bool
	DefaultRendererProfile  string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type DLNASettingsInput struct {
	Enabled                 bool
	FriendlyName            string
	Interfaces              []string
	AllowedCIDRs            []string
	AnnounceIntervalSeconds int32
	TranscodeEnabled        bool
	ThumbnailsEnabled       bool
	SubtitlesEnabled        bool
	DefaultRendererProfile  string
}

func (s *SettingsStore) GetDLNASettings(ctx context.Context) (DLNASettings, error) {
	row, err := storagegen.New(s.pool).GetDLNASettings(ctx, dlnaDefaultParams())
	if err != nil {
		return DLNASettings{}, err
	}
	return dlnaSettingsFromGetRow(row), nil
}

func (s *SettingsStore) UpdateDLNASettings(ctx context.Context, input DLNASettingsInput) (DLNASettings, error) {
	input = normalizeDLNASettings(input)
	if err := validateDLNASettings(input); err != nil {
		return DLNASettings{}, err
	}
	row, err := storagegen.New(s.pool).UpdateDLNASettings(ctx, storagegen.UpdateDLNASettingsParams{
		Enabled:                 input.Enabled,
		FriendlyName:            input.FriendlyName,
		Interfaces:              input.Interfaces,
		AllowedCidrs:            input.AllowedCIDRs,
		AnnounceIntervalSeconds: input.AnnounceIntervalSeconds,
		TranscodeEnabled:        input.TranscodeEnabled,
		ThumbnailsEnabled:       input.ThumbnailsEnabled,
		SubtitlesEnabled:        input.SubtitlesEnabled,
		DefaultRendererProfile:  input.DefaultRendererProfile,
	})
	if err != nil {
		return DLNASettings{}, err
	}
	return dlnaSettingsFromUpdateRow(row), nil
}

func dlnaDefaultParams() storagegen.GetDLNASettingsParams {
	return storagegen.GetDLNASettingsParams{
		Enabled:                 false,
		FriendlyName:            DefaultDLNAFriendlyName,
		Interfaces:              []string{},
		AllowedCidrs:            append([]string{}, DefaultDLNAAllowedCIDRs...),
		AnnounceIntervalSeconds: DefaultDLNAAnnounceIntervalSeconds,
		TranscodeEnabled:        true,
		ThumbnailsEnabled:       true,
		SubtitlesEnabled:        true,
		DefaultRendererProfile:  DefaultDLNARendererProfile,
	}
}

func normalizeDLNASettings(input DLNASettingsInput) DLNASettingsInput {
	input.FriendlyName = strings.TrimSpace(input.FriendlyName)
	if input.FriendlyName == "" {
		input.FriendlyName = DefaultDLNAFriendlyName
	}
	input.Interfaces = normalizedStringList(input.Interfaces)
	input.AllowedCIDRs = normalizedStringList(input.AllowedCIDRs)
	if len(input.AllowedCIDRs) == 0 {
		input.AllowedCIDRs = append([]string{}, DefaultDLNAAllowedCIDRs...)
	}
	if input.AnnounceIntervalSeconds == 0 {
		input.AnnounceIntervalSeconds = DefaultDLNAAnnounceIntervalSeconds
	}
	input.DefaultRendererProfile = strings.TrimSpace(input.DefaultRendererProfile)
	if input.DefaultRendererProfile == "" {
		input.DefaultRendererProfile = DefaultDLNARendererProfile
	}
	return input
}

func validateDLNASettings(input DLNASettingsInput) error {
	if len(input.FriendlyName) > 120 || len(input.DefaultRendererProfile) > 80 {
		return ErrInvalidInput
	}
	if input.AnnounceIntervalSeconds < 60 || input.AnnounceIntervalSeconds > 86400 {
		return ErrInvalidInput
	}
	for _, name := range input.Interfaces {
		if _, err := net.InterfaceByName(name); err != nil {
			return ErrInvalidInput
		}
	}
	for _, cidr := range input.AllowedCIDRs {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return ErrInvalidInput
		}
	}
	return nil
}

func normalizedStringList(values []string) []string {
	results := []string{}
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		results = append(results, value)
	}
	return results
}

func dlnaSettingsFromGetRow(row storagegen.GetDLNASettingsRow) DLNASettings {
	return DLNASettings{
		Enabled:                 row.Enabled,
		FriendlyName:            row.FriendlyName,
		Interfaces:              append([]string{}, row.Interfaces...),
		AllowedCIDRs:            append([]string{}, row.AllowedCidrs...),
		AnnounceIntervalSeconds: row.AnnounceIntervalSeconds,
		TranscodeEnabled:        row.TranscodeEnabled,
		ThumbnailsEnabled:       row.ThumbnailsEnabled,
		SubtitlesEnabled:        row.SubtitlesEnabled,
		DefaultRendererProfile:  row.DefaultRendererProfile,
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               row.UpdatedAt,
	}
}

func dlnaSettingsFromUpdateRow(row storagegen.UpdateDLNASettingsRow) DLNASettings {
	return DLNASettings{
		Enabled:                 row.Enabled,
		FriendlyName:            row.FriendlyName,
		Interfaces:              append([]string{}, row.Interfaces...),
		AllowedCIDRs:            append([]string{}, row.AllowedCidrs...),
		AnnounceIntervalSeconds: row.AnnounceIntervalSeconds,
		TranscodeEnabled:        row.TranscodeEnabled,
		ThumbnailsEnabled:       row.ThumbnailsEnabled,
		SubtitlesEnabled:        row.SubtitlesEnabled,
		DefaultRendererProfile:  row.DefaultRendererProfile,
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               row.UpdatedAt,
	}
}
