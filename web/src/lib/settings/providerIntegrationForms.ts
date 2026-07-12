import type {
	MetadataProvider,
	MetadataProviderForm,
	MetadataProviderRequest,
	SubtitleProvider,
	SubtitleProviderCatalogEntry,
	SubtitleProviderForm,
	SubtitleProviderRequest,
	SubtitleProviderSettingValue
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
	type: SubtitleProviderRequest['type'] = 'opensubtitlescom',
	entry?: SubtitleProviderCatalogEntry
): SubtitleProviderForm {
	const mock = type === 'mock';
	const settings = defaultSubtitleSettings(entry);
	const baseUrl = stringSetting(settings.baseUrl) ?? (mock ? 'mock://subtitles' : '');
	return {
		name: entry?.displayName ?? (mock ? 'Mock Subtitles' : 'OpenSubtitles'),
		type,
		catalogKey: entry?.key ?? type,
		baseUrl,
		username: stringSetting(settings.username) ?? '',
		password: '',
		apiKey: '',
		settings,
		secretSettings: {},
		clearSecretFields: [],
		enabled: entry?.runtimeStatus === 'supported',
		priority: mock ? 900 : 100,
		mockSubtitles: [],
		secretFieldsSet: [],
		runtimeStatus: entry?.runtimeStatus,
		runtimeMessage: entry?.runtimeMessage
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
		catalogKey: provider.catalogKey,
		baseUrl: provider.baseUrl,
		username: provider.username ?? '',
		password: '',
		apiKey: '',
		settings: { ...provider.settings },
		secretSettings: {},
		clearSecretFields: [],
		enabled: provider.enabled,
		priority: provider.priority,
		apiKeySet: provider.apiKeySet,
		passwordSet: provider.passwordSet,
		secretFieldsSet: provider.secretFieldsSet,
		runtimeStatus: provider.runtimeStatus,
		runtimeMessage: provider.runtimeMessage,
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
	const settings = normalizeSettings(form.settings ?? {});
	const secretSettings = normalizeSecretSettings(form.secretSettings ?? {});
	const clearSecretFields = (form.clearSecretFields ?? []).filter((field) => field.trim() !== '');
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: optionalString(form.baseUrl),
		username: optionalString(form.username),
		password: optionalString(form.password),
		apiKey: optionalString(form.apiKey),
		settings,
		secretSettings: Object.keys(secretSettings).length ? secretSettings : undefined,
		clearSecretFields: clearSecretFields.length ? clearSecretFields : undefined,
		enabled: form.enabled,
		priority: form.priority,
		mockSubtitles: form.type === 'mock' ? normalizedMockSubtitles(form) : undefined
	};
}

function defaultSubtitleSettings(
	entry: SubtitleProviderCatalogEntry | undefined
): Record<string, SubtitleProviderSettingValue> {
	if (!entry) return {};
	return Object.fromEntries(
		entry.fields
			.filter((field) => field.persisted && field.options?.[0] && field.type !== 'password')
			.map((field) => [field.key, { stringValue: field.options?.[0] }])
	);
}

function normalizeSettings(settings: Record<string, SubtitleProviderSettingValue>) {
	return Object.fromEntries(
		Object.entries(settings).filter(([, value]) =>
			Boolean(
				value.stringValue?.trim() ||
					value.numberValue !== undefined ||
					value.booleanValue !== undefined ||
					value.stringValues?.length
			)
		)
	);
}

function normalizeSecretSettings(settings: Record<string, string>) {
	return Object.fromEntries(
		Object.entries(settings)
			.map(([key, value]) => [key, value.trim()] as const)
			.filter(([, value]) => value !== '')
	);
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

function stringSetting(value: SubtitleProviderSettingValue | undefined) {
	return value?.stringValue;
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
