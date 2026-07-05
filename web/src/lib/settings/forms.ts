import type {
	DownloadClient,
	DownloadClientForm,
	DownloadClientRequest,
	DownloadClientType,
	Indexer,
	IndexerForm,
	IndexerProtocol,
	IndexerRequest,
	MetadataProvider,
	MetadataProviderForm,
	MetadataProviderRequest,
	SubtitleProvider,
	SubtitleProviderForm,
	SubtitleProviderRequest,
	LibraryFolderForm,
	LibraryFolderRequest,
	PathMappingForm,
	PathMappingRequest
} from './types';

export {
	customFormatFormFromFormat,
	emptyCustomFormatForm,
	normalizeCustomFormatForm
} from './customFormatForms';
export {
	emptyLanguageForm,
	languageFormFromLanguage,
	normalizeLanguageForm,
	normalizeLanguageUpdateForm
} from './languageForms';
export {
	emptyMediaProfileForm,
	mediaProfileFormFromProfile,
	normalizeMediaProfileForm
} from './mediaProfileForms';
export {
	emptyUserForm,
	normalizeUserCreateForm,
	normalizeUserUpdateForm,
	userFormFromUser
} from './userForms';

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

export function emptyMetadataProviderForm(): MetadataProviderForm {
	return {
		name: '',
		type: 'tmdb',
		baseUrl: 'https://api.themoviedb.org/3',
		apiKey: '',
		pin: '',
		accessToken: '',
		enabled: true,
		priority: 100
	};
}

export function emptySubtitleProviderForm(): SubtitleProviderForm {
	return {
		name: 'OpenSubtitles',
		type: 'opensubtitles',
		baseUrl: 'https://api.opensubtitles.com',
		username: '',
		password: '',
		apiKey: '',
		enabled: true,
		priority: 100
	};
}

export function emptyLibraryFolderForm(): LibraryFolderForm {
	return {
		path: ''
	};
}

export function emptyPathMappingForm(): PathMappingForm {
	return {
		clientPath: '',
		appPath: ''
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

export function metadataProviderFormFromProvider(provider: MetadataProvider): MetadataProviderForm {
	return {
		id: provider.id,
		name: provider.name,
		type: provider.type,
		baseUrl: provider.baseUrl,
		apiKey: provider.apiKey ?? '',
		pin: provider.pin ?? '',
		accessToken: provider.accessToken ?? '',
		enabled: provider.enabled,
		priority: provider.priority
	};
}

export function subtitleProviderFormFromProvider(provider: SubtitleProvider): SubtitleProviderForm {
	return {
		id: provider.id,
		name: provider.name,
		type: provider.type,
		baseUrl: provider.baseUrl,
		username: provider.username ?? '',
		password: undefined,
		apiKey: undefined,
		enabled: provider.enabled,
		priority: provider.priority
	};
}

export function normalizeDownloadClientForm(form: DownloadClientForm): DownloadClientRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		protocol: downloadClientProtocolForType(form.type),
		baseUrl: form.baseUrl.trim(),
		username: optionalString(form.username),
		password: optionalString(form.password),
		apiKey: optionalString(form.apiKey),
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
		apiKey: optionalString(form.apiKey),
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

export function normalizeMetadataProviderForm(form: MetadataProviderForm): MetadataProviderRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		apiKey: optionalString(form.apiKey),
		pin: optionalString(form.pin),
		accessToken: optionalString(form.accessToken),
		enabled: form.enabled,
		priority: form.priority
	};
}

export function normalizeSubtitleProviderForm(form: SubtitleProviderForm): SubtitleProviderRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		username: optionalString(form.username),
		password: optionalString(form.password),
		apiKey: optionalString(form.apiKey),
		enabled: form.enabled,
		priority: form.priority
	};
}

export function normalizeLibraryFolderForm(form: LibraryFolderForm): LibraryFolderRequest {
	return {
		path: form.path.trim()
	};
}

export function normalizePathMappingForm(form: PathMappingForm): PathMappingRequest {
	return {
		clientPath: form.clientPath.trim(),
		appPath: form.appPath.trim()
	};
}

function optionalString(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}

function parseCategories(value: string) {
	return value
		.split(',')
		.map((item) => Number.parseInt(item.trim(), 10))
		.filter((item) => Number.isInteger(item));
}
