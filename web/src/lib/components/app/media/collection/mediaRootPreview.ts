import type { FileNamingSettings, LibraryFolder, MediaItem } from '$lib/settings/types';
import { defaultFileNamingTemplates } from '$lib/settings/fileNamingTemplates';

export function mediaRootPreview(
	item: MediaItem,
	folder: LibraryFolder | undefined,
	settings: FileNamingSettings | undefined
) {
	if (!folder) return '-';
	const template = folderTemplate(item, settings);
	const segment = sanitizePathSegment(renderMediaTemplate(template, item));
	return `${folder.path}/${segment}`;
}

export function mediaRootWarning(
	item: MediaItem,
	folder: LibraryFolder | undefined,
	settings: FileNamingSettings | undefined
) {
	const current = item.mediaFolderPath?.trim();
	if (!current || !folder) return undefined;
	const expected = mediaRootPreview(item, folder, settings);
	if (expected === '-' || normalizedRoot(current) === normalizedRoot(expected)) return undefined;
	return { expected };
}

function folderTemplate(item: MediaItem, settings: FileNamingSettings | undefined) {
	if (item.type === 'movie') {
		return settings?.movieFolderFormat ?? defaultFileNamingTemplates.movieFolderFormat;
	}
	return settings?.seriesFolderFormat ?? defaultFileNamingTemplates.seriesFolderFormat;
}

function renderMediaTemplate(template: string, item: MediaItem) {
	const title = item.title;
	const year = item.year ? `${item.year}` : '';
	const values: Record<string, string> = {
		movie_title: title,
		release_year: year,
		series_title: title,
		year
	};
	return template
		.replace(/\{([^{}]+)\}/g, (token, rawKey: string) => values[tokenName(rawKey)] ?? token)
		.replace(/\s+/g, ' ')
		.trim();
}

function tokenName(key: string) {
	const legacy: Record<string, string> = {
		'Movie Title': 'movie_title',
		'Release Year': 'release_year',
		'Series Title': 'series_title',
		Year: 'year'
	};
	return legacy[key] ?? key.toLowerCase().replace(/[ -]/g, '_');
}

function sanitizePathSegment(value: string) {
	const sanitized = value
		.replace(/[\\/]/g, ' ')
		.replace(/:/g, ' -')
		.replace(/[*?"<>|]/g, '')
		.trim()
		.replace(/^\.+|\.+$/g, '');
	return sanitized || 'Untitled';
}

function normalizedRoot(path: string) {
	const normalized = path.trim().replace(/[/\\]+$/g, '');
	return normalized || path.trim();
}
