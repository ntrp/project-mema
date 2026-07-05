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

export function metadataProviderFormFromProvider(provider: MetadataProvider): MetadataProviderForm {
	return {
		id: provider.id,
		name: provider.name,
		type: provider.type,
		baseUrl: provider.baseUrl,
		apiKey: '',
		pin: '',
		accessToken: '',
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
		password: undefined,
		apiKey: undefined,
		enabled: provider.enabled,
		priority: provider.priority
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
