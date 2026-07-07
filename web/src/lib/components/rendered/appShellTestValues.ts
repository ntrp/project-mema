import type {
	IndexerSearchResponse,
	ManagedUser,
	MediaProfile,
	MetadataProvider,
	SubtitleProvider,
	Tag,
	UserProfile,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';

export function userSummary(): UserSummary {
	return {
		id: 'user-1',
		username: 'scenario-admin',
		displayName: 'Scenario Admin',
		pictureUrl: '',
		role: 'admin'
	};
}

export function userProfile(overrides: Partial<UserProfile> = {}): UserProfile {
	return {
		id: 'user-1',
		username: 'scenario-admin',
		displayName: 'Scenario Admin',
		pictureUrl: 'https://example.test/profile.png',
		role: 'admin',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
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
		isDefault: false,
		qualityIds: ['webdl-1080p'],
		upgradesAllowed: true,
		upgradeUntilQualityId: 'webdl-1080p',
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 50,
		minimumCustomFormatScoreIncrement: 1,
		finalContainer: 'mkv',
		removeUnwantedAudio: false,
		audioLossyTranscodePolicy: 'disabled',
		removeUnwantedSubtitles: false,
		subtitlePreferredMode: 'mixed',
		allowSubtitleReleaseFallback: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		videoTarget: {},
		audioTargets: [
			{
				languageId: 'english',
				score: 0
			}
		],
		subtitleTargets: [{ languageId: 'english', score: 0 }],
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

export function subtitleProvider(overrides: Partial<SubtitleProvider> = {}): SubtitleProvider {
	return {
		id: 'subtitle-1',
		name: 'OpenSubtitles',
		type: 'opensubtitles',
		enabled: true,
		baseUrl: 'https://api.opensubtitles.com',
		username: 'scenario-user',
		apiKey: 'scenario-key',
		password: 'scenario-password',
		apiKeySet: true,
		passwordSet: true,
		priority: 100,
		mockSubtitles: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function emptyIndexerSearch(): IndexerSearchResponse {
	return {
		settings: {
			cacheDurationMinutes: 60,
			historyRetentionDays: 14,
			automaticBlocklistExpiryDays: 7
		},
		stats: { totalEntries: 0, activeEntries: 0, expiredEntries: 0, indexerCount: 0 },
		cacheEntries: [],
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 }
	};
}
