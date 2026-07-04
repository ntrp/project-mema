package cardigann

import (
	"context"
)

func (s *Engine) releaseFromHTMLRow(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	baseCtx cardigannContext,
	searchURL string,
	row *htmlRow,
) (Release, bool, error) {
	rowCtx := baseCtx
	rowCtx.Result = map[string]any{}
	release := baseCardigannRelease(config)
	for _, field := range def.Search.Fields {
		value, found, err := cardigannHTMLValue(row, field.Selector, rowCtx, !isOptionalCardigannField(field.Name, field.Selector))
		if err != nil {
			return Release{}, false, cardigannFieldError(field.Name, err)
		}
		if !found {
			continue
		}
		rowCtx.Result[field.Name] = value
		applyCardigannField(&release, field.Name, value, searchURL)
	}
	if def.Download != nil && release.DownloadURL != "" {
		resolved, err := s.resolveCardigannDownload(ctx, def, config, rowCtx, release.DownloadURL)
		if err == nil && resolved != "" {
			release.DownloadURL = resolved
		}
	}
	finalized, ok := finalizeCardigannRelease(release)
	return finalized, ok, nil
}

func (s *Engine) releaseFromJSONRow(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	baseCtx cardigannContext,
	searchURL string,
	row jsonRow,
) (Release, bool, error) {
	rowCtx := baseCtx
	rowCtx.Result = map[string]any{}
	release := baseCardigannRelease(config)
	for _, field := range def.Search.Fields {
		value, found, err := cardigannJSONValue(row, field.Selector, rowCtx, !isOptionalCardigannField(field.Name, field.Selector))
		if err != nil {
			return Release{}, false, cardigannFieldError(field.Name, err)
		}
		if !found {
			continue
		}
		rowCtx.Result[field.Name] = value
		applyCardigannField(&release, field.Name, value, searchURL)
	}
	if def.Download != nil && release.DownloadURL != "" {
		resolved, err := s.resolveCardigannDownload(ctx, def, config, rowCtx, release.DownloadURL)
		if err == nil && resolved != "" {
			release.DownloadURL = resolved
		}
	}
	finalized, ok := finalizeCardigannRelease(release)
	return finalized, ok, nil
}
