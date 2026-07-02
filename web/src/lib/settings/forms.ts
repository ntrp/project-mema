import type {
	DownloadClient,
	DownloadClientForm,
	DownloadClientRequest,
	Indexer,
	IndexerForm,
	IndexerRequest,
	MetadataProvider,
	MetadataProviderForm,
	MetadataProviderRequest,
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
		name: '',
		type: 'torznab',
		baseUrl: '',
		apiKey: '',
		categoriesText: '',
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
		name: indexer.name,
		type: indexer.type,
		baseUrl: indexer.baseUrl,
		apiKey: indexer.apiKey ?? '',
		categoriesText: (indexer.categories ?? []).join(', '),
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

export function normalizeDownloadClientForm(form: DownloadClientForm): DownloadClientRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		username: optionalString(form.username),
		password: optionalString(form.password),
		apiKey: optionalString(form.apiKey),
		category: optionalString(form.category),
		enabled: form.enabled,
		priority: form.priority
	};
}

export function normalizeIndexerForm(form: IndexerForm): IndexerRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		apiKey: optionalString(form.apiKey),
		categories: parseCategories(form.categoriesText),
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
