import { relativePath } from '$lib/components/app/media/files/mediaFilePath';
import { displayLanguage } from '$lib/settings/languageDisplay';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

export type MediaFileOtherFile = MediaFileRow['otherFiles'][number];

export function otherFileTypeLabel(type: MediaFileOtherFile['type']) {
	switch (type) {
		case 'subtitle':
			return 'Subtitle';
		case 'metadata':
			return 'Metadata';
		default:
			return 'Unknown';
	}
}

export function otherFileStatusLabel(status: MediaFileOtherFile['status']) {
	return status === 'missing' ? 'Missing' : 'Available';
}

export function otherFileLanguageLabel(file: MediaFileOtherFile) {
	return file.type === 'subtitle' ? displayLanguage(file.language) : '-';
}

export function otherFileDisplayPath(row: MediaFileRow, file: MediaFileOtherFile) {
	return relativePath(row.path ? row.path.replace(/[^/]+$/, '') : undefined, file.path);
}
