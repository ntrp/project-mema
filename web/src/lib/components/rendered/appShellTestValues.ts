import type {
	IndexerSearchResponse,
	ManagedUser,
	MediaProfile,
	MetadataProvider,
	Tag,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';

export function userSummary(): UserSummary {
	return {
		id: 'user-1',
		username: 'scenario-admin',
		role: 'admin'
	};
}

export function managedUser(overrides: Partial<ManagedUser> = {}): ManagedUser {
	return {
		...userSummary(),
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function userForm(overrides: Partial<UserFormValue> = {}): UserFormValue {
	return {
		username: 'scenario-user',
		password: 'long-password',
		role: 'user',
		...overrides
	};
}

export function tag(overrides: Partial<Tag> = {}): Tag {
	return {
		id: 'tag-1',
		name: 'scenario-tag',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function mediaProfile(overrides: Partial<MediaProfile> = {}): MediaProfile {
	return {
		id: 'profile-1',
		name: 'Scenario Profile',
		qualityIds: ['webdl-1080p'],
		upgradesAllowed: true,
		upgradeUntilQualityId: 'webdl-1080p',
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 50,
		minimumCustomFormatScoreIncrement: 1,
		removeNonEnabledLanguages: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		targetLanguages: ['english'],
		targetLanguageScores: [{ languageId: 'english', score: 0, required: false }],
		customFormatScores: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function metadataProvider(overrides: Partial<MetadataProvider> = {}): MetadataProvider {
	return {
		id: 'metadata-1',
		name: 'Scenario Metadata',
		type: 'tmdb',
		enabled: true,
		baseUrl: 'https://metadata.example.test',
		apiKey: 'scenario-key',
		priority: 1,
		settings: {},
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	} as MetadataProvider;
}

export function emptyIndexerSearch(): IndexerSearchResponse {
	return {
		settings: { cacheDurationMinutes: 60, historyRetentionDays: 14 },
		stats: { totalEntries: 0, activeEntries: 0, expiredEntries: 0, indexerCount: 0 },
		cacheEntries: [],
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 }
	};
}
