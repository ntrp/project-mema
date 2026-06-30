import type {
	CustomFormat,
	CustomFormatForm,
	CustomFormatRequest,
	DownloadClient,
	DownloadClientForm,
	DownloadClientRequest,
	Indexer,
	IndexerForm,
	IndexerRequest,
	ManagedUser,
	MediaProfile,
	MediaProfileForm,
	MediaProfileRequest,
	MetadataProvider,
	MetadataProviderForm,
	MetadataProviderRequest,
	LibraryFolderForm,
	LibraryFolderRequest,
	PathMappingForm,
	PathMappingRequest,
	UserCreateRequest,
	UserForm,
	UserUpdateRequest
} from './types';

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

export function emptyMediaProfileForm(): MediaProfileForm {
	return {
		name: '',
		qualityIds: [],
		upgradesAllowed: true,
		upgradeUntilQualityId: undefined,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		targetLanguages: ['english'],
		targetLanguageScores: [{ languageId: 'english', score: 0 }],
		customFormatScores: []
	};
}

export function emptyCustomFormatForm(): CustomFormatForm {
	return {
		name: '',
		includeSpecs: [],
		excludeSpecs: []
	};
}

export function emptyUserForm(): UserForm {
	return {
		username: '',
		password: '',
		role: 'user'
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

export function mediaProfileFormFromProfile(profile: MediaProfile): MediaProfileForm {
	return {
		id: profile.id,
		name: profile.name,
		qualityIds: [...(profile.qualityIds ?? [])],
		upgradesAllowed: profile.upgradesAllowed,
		upgradeUntilQualityId: profile.upgradeUntilQualityId,
		minimumCustomFormatScore: profile.minimumCustomFormatScore,
		upgradeUntilCustomFormatScore: profile.upgradeUntilCustomFormatScore,
		minimumCustomFormatScoreIncrement: profile.minimumCustomFormatScoreIncrement,
		targetLanguages: [...(profile.targetLanguages ?? [])],
		targetLanguageScores: languageScoresFromProfile(profile),
		customFormatScores: (profile.customFormatScores ?? []).map((score) => ({ ...score }))
	};
}

export function customFormatFormFromFormat(format: CustomFormat): CustomFormatForm {
	return {
		id: format.id,
		name: format.name,
		includeSpecs: format.includeSpecs.map((spec) => ({ ...spec })),
		excludeSpecs: format.excludeSpecs.map((spec) => ({ ...spec }))
	};
}

export function userFormFromUser(user: ManagedUser): UserForm {
	return {
		id: user.id,
		username: user.username,
		password: '',
		role: user.role
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

export function normalizeMediaProfileForm(form: MediaProfileForm): MediaProfileRequest {
	const qualityIds = [...new Set(form.qualityIds.map((id) => id.trim()).filter(Boolean))];
	const customFormatScores = form.customFormatScores
		.filter((score) => score.customFormatId)
		.map((score) => ({
			customFormatId: score.customFormatId,
			score: normalizedInteger(score.score)
		}));
	const targetLanguageScores = languageScoresFromForm(form);
	return {
		name: form.name.trim(),
		qualityIds,
		upgradesAllowed: form.upgradesAllowed,
		upgradeUntilQualityId:
			form.upgradeUntilQualityId && qualityIds.includes(form.upgradeUntilQualityId)
				? form.upgradeUntilQualityId
				: undefined,
		minimumCustomFormatScore: normalizedInteger(form.minimumCustomFormatScore),
		upgradeUntilCustomFormatScore: normalizedInteger(form.upgradeUntilCustomFormatScore),
		minimumCustomFormatScoreIncrement: Math.max(
			0,
			normalizedInteger(form.minimumCustomFormatScoreIncrement)
		),
		targetLanguages: targetLanguageScores.map((score) => score.languageId),
		targetLanguageScores,
		customFormatScores
	};
}

function languageScoresFromProfile(profile: MediaProfile) {
	if (profile.targetLanguageScores?.length) {
		return profile.targetLanguageScores.map((score) => ({ ...score }));
	}
	return (profile.targetLanguages ?? []).map((languageId) => ({ languageId, score: 0 }));
}

function languageScoresFromForm(form: MediaProfileForm) {
	const seen = new Set<string>();
	const source = form.targetLanguageScores?.length
		? form.targetLanguageScores
		: form.targetLanguages.map((languageId) => ({ languageId, score: 0 }));
	const scores = [];
	for (const value of source) {
		const languageId = value.languageId.trim();
		if (!languageId || seen.has(languageId)) {
			continue;
		}
		seen.add(languageId);
		scores.push({ languageId, score: normalizedInteger(value.score) });
	}
	return scores;
}

function normalizedInteger(value: number | string | undefined) {
	const parsed = Number(value ?? 0);
	if (!Number.isFinite(parsed)) {
		return 0;
	}
	return Math.trunc(parsed);
}

export function normalizeCustomFormatForm(form: CustomFormatForm): CustomFormatRequest {
	return {
		name: form.name.trim(),
		includeSpecs: normalizeCustomFormatSpecs(form.includeSpecs),
		excludeSpecs: normalizeCustomFormatSpecs(form.excludeSpecs)
	};
}

export function normalizeUserCreateForm(form: UserForm): UserCreateRequest {
	return {
		username: form.username.trim(),
		password: form.password,
		role: form.role
	};
}

export function normalizeUserUpdateForm(form: UserForm): UserUpdateRequest {
	return {
		username: form.username.trim(),
		password: optionalString(form.password),
		role: form.role
	};
}

function optionalString(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}

function normalizeCustomFormatSpecs(specs: CustomFormatRequest['includeSpecs']) {
	return specs
		.map((spec) => ({
			id: spec.id.trim(),
			name: spec.name.trim(),
			type: spec.type,
			value: spec.value.trim(),
			required: spec.required
		}))
		.filter((spec) => spec.id !== '' && spec.name !== '' && spec.value !== '');
}

function parseCategories(value: string) {
	return value
		.split(',')
		.map((item) => Number.parseInt(item.trim(), 10))
		.filter((item) => Number.isInteger(item));
}
