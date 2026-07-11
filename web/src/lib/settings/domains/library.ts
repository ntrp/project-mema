import { client } from '$lib/api/client';
import { normalizeLibraryFolderForm, normalizePathMappingForm } from '../forms';
import type {
	LibraryFolderForm,
	LibraryFolderOption,
	LibraryFolderOptionListResponse,
	LibraryMediaKind,
	LibraryScanImportRequest,
	LibraryScanItemMatchRequest,
	PathMappingForm
} from '../types';

export async function saveLibraryFolder(form: LibraryFolderForm) {
	const { data, error } = await client.POST('/settings/library/folders', {
		body: normalizeLibraryFolderForm(form)
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library scan was not returned');
	}
	return data;
}

export async function listLibraryFolderOptions(
	path?: string
): Promise<LibraryFolderOptionListResponse> {
	const { data, error } = await client.GET('/settings/library/folder-options', {
		params: { query: { path } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Folder options were not returned');
	}
	return data;
}

export async function createLibraryFolderOption(
	parentPath: string,
	name: string
): Promise<LibraryFolderOption> {
	const { data, error } = await client.POST('/settings/library/folder-options', {
		body: { parentPath, name }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Folder was not returned');
	}
	return data;
}

export async function deleteLibraryFolder(id: string) {
	const { error } = await client.DELETE('/settings/library/folders/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function scanLibraryFolder(id: string) {
	const { data, error } = await client.POST('/settings/library/folders/{id}/scan', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library scan was not returned');
	}
	return data;
}

export async function savePathMapping(form: PathMappingForm) {
	const { data, error } = await client.POST('/settings/library/path-mappings', {
		body: normalizePathMappingForm(form)
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Path mapping was not returned');
	}
	return data;
}

export async function deletePathMapping(id: string) {
	const { error } = await client.DELETE('/settings/library/path-mappings/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function getLibraryScan(id: string) {
	const { data, error } = await client.GET('/settings/library/scans/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library scan was not returned');
	}
	return data;
}

export async function matchLibraryScanItem(
	scanId: string,
	itemId: string,
	request: LibraryScanItemMatchRequest
) {
	const { data, error } = await client.POST('/settings/library/scans/{id}/items/{itemId}/match', {
		params: { path: { id: scanId, itemId } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library match was not returned');
	}
	return data;
}

export async function importLibraryScanItems(scanId: string, request: LibraryScanImportRequest) {
	const { data, error } = await client.POST('/settings/library/scans/{id}/import', {
		params: { path: { id: scanId } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library import result was not returned');
	}
	return data;
}

export async function resetLibraryScanItemImport(scanId: string, itemId: string) {
	const { data, error } = await client.POST('/settings/library/scans/{id}/items/{itemId}/reset', {
		params: { path: { id: scanId, itemId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library import reset result was not returned');
	}
	return data;
}

export function mediaTypeForLibraryKind(kind: LibraryMediaKind) {
	return kind === 'series' || kind === 'anime_series' ? 'serie' : 'movie';
}
