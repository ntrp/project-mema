import { client } from '$lib/api/client';
import type {
	MediaFileHistoryResponse,
	MediaFilePreviewInfo,
	MediaFileTrackDeleteRequest,
	MediaItemSubtitleListResponse,
	MediaItemSubtitleSelectionRequest,
	MediaRenameApplyResponse,
	MediaRenamePreviewResponse
} from '$lib/settings/types';

export async function deleteMediaItemFile(id: string, path: string) {
	const { data, error } = await client.POST('/media/items/{id}/files/delete', {
		params: { path: { id } },
		body: { path }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function deleteMediaItemFileTrack(id: string, request: MediaFileTrackDeleteRequest) {
	const { data, error } = await client.POST('/media/items/{id}/files/tracks/delete', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function getMediaItemFilePreviewInfo(
	id: string,
	path: string
): Promise<MediaFilePreviewInfo> {
	const { data, error } = await client.GET('/media/items/{id}/files/preview-info', {
		params: { path: { id }, query: { path } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media file preview info was not returned');
	}
	return data;
}

export async function listMediaItemSubtitles(id: string): Promise<MediaItemSubtitleListResponse> {
	const { data, error } = await client.GET('/media/items/{id}/subtitles', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { subtitles: [] };
}

export async function deleteMediaItemSubtitle(id: string, subtitleId: string) {
	const { data, error } = await client.DELETE('/media/items/{id}/subtitles/{subtitleId}', {
		params: { path: { id, subtitleId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function updateMediaItemSubtitle(
	id: string,
	subtitleId: string,
	request: MediaItemSubtitleSelectionRequest
) {
	const { data, error } = await client.PUT('/media/items/{id}/subtitles/{subtitleId}', {
		params: { path: { id, subtitleId } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function listMediaFileHistory(id: string): Promise<MediaFileHistoryResponse> {
	const { data, error } = await client.GET('/media/items/{id}/file-history', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { entries: [] };
}

export async function previewMediaRename(id: string): Promise<MediaRenamePreviewResponse> {
	const { data, error } = await client.GET('/media/items/{id}/rename-preview', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { rows: [] };
}

export async function applyMediaRename(
	id: string,
	currentPaths?: string[]
): Promise<MediaRenameApplyResponse> {
	const { data, error } = await client.POST('/media/items/{id}/rename-apply', {
		params: { path: { id } },
		body: currentPaths ? { currentPaths } : undefined
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { rows: [], appliedCount: 0, skippedCount: 0, failedCount: 0 };
}
