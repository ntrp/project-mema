package providercore

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

const openSubtitlesChunkSize int64 = 64 * 1024

func BuildSearchRequest(item storage.MediaItem, languageID string, filePath string) subtitles.SearchRequest {
	request := subtitles.SearchRequest{
		MediaType:    item.Type,
		Title:        item.Title,
		LanguageID:   strings.TrimSpace(languageID),
		Year:         item.Year,
		FilePath:     strings.TrimSpace(filePath),
		MediaContext: BuildMediaContext(item, filePath),
	}
	return request
}

func BuildMediaContext(item storage.MediaItem, filePath string) subtitles.MediaContext {
	return subtitles.MediaContext{
		ExternalIDs:        externalIDs(item, ""),
		SeasonExternalIDs:  externalIDs(item, "season"),
		EpisodeExternalIDs: externalIDs(item, "episode"),
		Aliases:            aliases(item),
		EpisodeNumbering:   episodeNumbering(item),
		File:               fileContext(item, filePath),
		Provenance:         provenance(item),
	}
}

func externalIDs(item storage.MediaItem, entityType string) map[string]string {
	ids := map[string]string{}
	if entityType == "" && item.ExternalProvider != nil && item.ExternalID != nil {
		addID(ids, *item.ExternalProvider, *item.ExternalID)
	}
	for _, mapping := range item.ProviderMappings {
		if entityType != "" && mapping.EntityType != entityType {
			continue
		}
		if entityType == "" && mapping.EntityType != "" && mapping.EntityType != "media" {
			continue
		}
		addID(ids, mapping.ProviderName, mapping.ExternalID)
	}
	return ids
}

func addID(ids map[string]string, provider string, id string) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	id = strings.TrimSpace(id)
	if provider != "" && id != "" {
		ids[provider] = id
	}
}

func aliases(item storage.MediaItem) []subtitles.MediaAlias {
	values := make([]subtitles.MediaAlias, 0, len(item.Aliases))
	for _, alias := range item.Aliases {
		value := strings.TrimSpace(alias.Alias)
		if value == "" {
			continue
		}
		values = append(values, subtitles.MediaAlias{
			Value:        value,
			Language:     stringValue(alias.Language),
			Kind:         alias.Kind,
			ProviderName: stringValue(alias.ProviderName),
		})
	}
	return values
}

func episodeNumbering(item storage.MediaItem) []subtitles.EpisodeNumbering {
	values := make([]subtitles.EpisodeNumbering, 0, len(item.EpisodeNumbering))
	for _, numbering := range item.EpisodeNumbering {
		values = append(values, subtitles.EpisodeNumbering{
			ProviderName:    numbering.ProviderName,
			NumberingScheme: numbering.NumberingScheme,
			SeasonNumber:    numbering.SeasonNumber,
			EpisodeNumber:   numbering.EpisodeNumber,
			AbsoluteNumber:  numbering.AbsoluteNumber,
		})
	}
	return values
}

func fileContext(item storage.MediaItem, filePath string) subtitles.FileContext {
	filePath = strings.TrimSpace(filePath)
	name := filepath.Base(filePath)
	ctx := subtitles.FileContext{
		Path:      filePath,
		Name:      name,
		BaseName:  strings.TrimSuffix(name, filepath.Ext(name)),
		Extension: strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), "."),
		Hashes:    map[string]string{},
	}
	if size, ok := fileSize(item, filePath); ok {
		ctx.SizeBytes = size
	}
	if stat, err := os.Stat(filePath); err == nil && !stat.IsDir() {
		ctx.SizeBytes = stat.Size()
	}
	return ctx
}

func fileSize(item storage.MediaItem, filePath string) (int64, bool) {
	for _, fact := range item.FileFacts {
		if fact.FilePath == filePath && fact.SizeBytes != nil {
			return *fact.SizeBytes, true
		}
	}
	return 0, false
}

func ComputeFileHashes(filePath string, names ...string) (map[string]string, error) {
	if len(names) == 0 {
		names = []string{"sha256", "opensubtitles"}
	}
	hashes := map[string]string{}
	for _, name := range names {
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "sha256":
			hash, err := fileSHA256(filePath)
			if err != nil {
				return nil, err
			}
			hashes["sha256"] = hash
		case "opensubtitles":
			stat, err := os.Stat(filePath)
			if err != nil {
				return nil, err
			}
			hash, err := openSubtitlesHash(filePath, stat.Size())
			if err != nil {
				return nil, err
			}
			hashes["opensubtitles"] = hash
		}
	}
	return hashes, nil
}

func fileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func openSubtitlesHash(filePath string, size int64) (string, error) {
	if size < openSubtitlesChunkSize*2 {
		return "", fmt.Errorf("file too small for opensubtitles hash")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	sum := uint64(size)
	if err := hashChunk(file, &sum); err != nil {
		return "", err
	}
	if _, err := file.Seek(-openSubtitlesChunkSize, io.SeekEnd); err != nil {
		return "", err
	}
	if err := hashChunk(file, &sum); err != nil {
		return "", err
	}
	return fmt.Sprintf("%016x", sum), nil
}

func hashChunk(reader io.Reader, sum *uint64) error {
	buffer := make([]byte, 8)
	for i := int64(0); i < openSubtitlesChunkSize/8; i++ {
		if _, err := io.ReadFull(reader, buffer); err != nil {
			return err
		}
		*sum += binary.LittleEndian.Uint64(buffer)
	}
	return nil
}

func provenance(item storage.MediaItem) []subtitles.ReleaseProvenance {
	values := []subtitles.ReleaseProvenance{}
	for _, source := range item.ComponentSources {
		for _, candidate := range []string{source.SourceFilePath, stringValue(source.ReleaseID), stringValue(source.SourceMetadata)} {
			if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") {
				values = append(values, subtitles.ReleaseProvenance{Source: source.SourceRole, InfoURL: candidate})
			}
		}
	}
	return values
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
