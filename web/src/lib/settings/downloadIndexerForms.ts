import type {
	DownloadClient,
	DownloadClientForm,
	DownloadClientRequest,
	DownloadClientType,
	Indexer,
	IndexerForm,
	IndexerProtocol,
	IndexerRequest
} from './types';

export function emptyDownloadClientForm(): DownloadClientForm {
	return {
		name: '',
		type: 'transmission',
		protocol: 'torrent',
		baseUrl: '',
		username: '',
		password: '',
		apiKey: '',
		category: '',
		enabled: true,
		priority: 100
	};
}

export function emptyIndexerForm(): IndexerForm {
	return {
		definitionId: 'generic-torznab',
		name: '',
		baseUrl: '',
		apiKey: '',
		categoriesText: '',
		fields: [],
		mediaTypeScopes: undefined,
		tagScopes: [],
		redirect: true,
		appProfileId: 'default',
		preferMagnetUrl: false,
		enabled: true,
		priority: 100
	};
}

export function downloadClientFormFromClient(client: DownloadClient): DownloadClientForm {
	return {
		id: client.id,
		name: client.name,
		type: client.type,
		protocol: client.protocol,
		baseUrl: client.baseUrl,
		username: client.username ?? '',
		password: client.password ?? '',
		apiKey: client.apiKey ?? '',
		passwordSet: client.passwordSet,
		apiKeySet: client.apiKeySet,
		category: client.category ?? '',
		enabled: client.enabled,
		priority: client.priority
	};
}

export function indexerFormFromIndexer(indexer: Indexer): IndexerForm {
	return {
		id: indexer.id,
		definitionId: indexer.definitionId,
		name: indexer.name,
		implementation: indexer.implementation,
		implementationName: indexer.implementationName,
		baseUrl: indexer.baseUrl,
		apiKey: indexer.apiKey ?? '',
		apiKeySet: indexer.apiKeySet,
		categoriesText: (indexer.categories ?? []).join(', '),
		mediaTypeScopes: indexer.mediaTypeScopes ?? [],
		tagScopes: indexer.tagScopes ?? [],
		fields: indexer.fields ?? [],
		redirect: indexer.redirect ?? true,
		appProfileId: indexer.appProfileId ?? 'default',
		minimumSeeders: indexer.minimumSeeders,
		seedRatio: indexer.seedRatio,
		seedTime: indexer.seedTime,
		packSeedTime: indexer.packSeedTime,
		preferMagnetUrl: indexer.preferMagnetUrl ?? false,
		enabled: indexer.enabled,
		priority: indexer.priority
	};
}

export function normalizeDownloadClientForm(form: DownloadClientForm): DownloadClientRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		protocol: downloadClientProtocolForType(form.type),
		baseUrl: form.baseUrl.trim(),
		username: optionalString(form.username),
		password: optionalSavedSecret(form.password, form.passwordSet),
		apiKey: optionalSavedSecret(form.apiKey, form.apiKeySet),
		category: optionalString(form.category),
		enabled: form.enabled,
		priority: form.priority
	};
}

export function downloadClientProtocolForType(type: DownloadClientType): IndexerProtocol {
	return type === 'sabnzbd' ? 'usenet' : 'torrent';
}

export function normalizeIndexerForm(form: IndexerForm): IndexerRequest {
	return {
		definitionId: form.definitionId,
		name: form.name.trim(),
		implementation: form.implementation,
		implementationName: form.implementationName,
		baseUrl: form.baseUrl.trim(),
		apiKey: optionalSavedSecret(form.apiKey, form.apiKeySet),
		categories: parseCategories(form.categoriesText),
		mediaTypeScopes: form.mediaTypeScopes,
		tagScopes: form.tagScopes ?? [],
		fields: form.fields ?? [],
		redirect: form.redirect ?? true,
		appProfileId: form.appProfileId ?? 'default',
		minimumSeeders: form.minimumSeeders,
		seedRatio: form.seedRatio,
		seedTime: form.seedTime,
		packSeedTime: form.packSeedTime,
		preferMagnetUrl: form.preferMagnetUrl ?? false,
		enabled: form.enabled,
		priority: form.priority
	};
}

function optionalString(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}

function optionalSavedSecret(value: string | undefined, saved: boolean | undefined) {
	const trimmed = value?.trim() ?? '';
	if (saved && trimmed === '') {
		return undefined;
	}
	return optionalString(value);
}

function parseCategories(value: string) {
	return value
		.split(',')
		.map((item) => Number.parseInt(item.trim(), 10))
		.filter((item) => Number.isInteger(item));
}
