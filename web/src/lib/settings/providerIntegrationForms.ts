import type {
	MetadataProvider,
	MetadataProviderForm,
	MetadataProviderRequest,
	SubtitleProvider,
	SubtitleProviderForm,
	SubtitleProviderRequest
} from './types';

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

export function emptySubtitleProviderForm(
	type: SubtitleProviderRequest['type'] = 'opensubtitles'
): SubtitleProviderForm {
	const mock = type === 'mock';
	return {
		name: mock ? 'Mock Subtitles' : 'OpenSubtitles',
		type,
		baseUrl: mock ? 'mock://subtitles' : 'https://api.opensubtitles.com',
		username: '',
		password: '',
		apiKey: '',
		enabled: true,
		priority: mock ? 900 : 100,
		mockSubtitles: []
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
		apiKeySet: provider.apiKeySet,
		pinSet: provider.pinSet,
		accessTokenSet: provider.accessTokenSet,
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
		password: provider.password ?? '',
		apiKey: provider.apiKey ?? '',
		enabled: provider.enabled,
		priority: provider.priority,
		mockSubtitles: provider.mockSubtitles.map(({ title, languageId, format }) => ({
			title,
			languageId,
			format
		}))
	};
}

export function normalizeMetadataProviderForm(form: MetadataProviderForm): MetadataProviderRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		apiKey: optionalSavedSecret(form.apiKey, form.apiKeySet),
		pin: optionalSavedSecret(form.pin, form.pinSet),
		accessToken: optionalSavedSecret(form.accessToken, form.accessTokenSet),
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
		priority: form.priority,
		mockSubtitles: form.type === 'mock' ? normalizedMockSubtitles(form) : undefined
	};
}

function normalizedMockSubtitles(
	form: SubtitleProviderForm
): SubtitleProviderRequest['mockSubtitles'] {
	return (form.mockSubtitles ?? [])
		.map((row) => ({
			title: row.title.trim(),
			languageId: row.languageId.trim(),
			format: row.format.trim()
		}))
		.filter((row) => row.title !== '' && row.languageId !== '' && row.format !== '');
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
