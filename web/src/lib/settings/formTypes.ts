import type { components } from '$lib/api/generated/schema';

type Schemas = components['schemas'];

export type DownloadClientForm = Schemas['DownloadClientRequest'] & {
	id?: string;
	passwordSet?: boolean;
	apiKeySet?: boolean;
};
export type IndexerForm = Omit<Schemas['IndexerRequest'], 'categories'> & {
	id?: string;
	apiKeySet?: boolean;
	categoriesText: string;
};
export type IndexerProxyForm = Schemas['IndexerProxyRequest'] & { id?: string };
export type MetadataProviderForm = Schemas['MetadataProviderRequest'] & {
	id?: string;
	apiKeySet?: boolean;
	pinSet?: boolean;
	accessTokenSet?: boolean;
};
export type SubtitleProviderForm = Schemas['SubtitleProviderRequest'] & { id?: string };
export type LibraryFolderForm = Schemas['LibraryFolderRequest'];
export type PathMappingForm = Schemas['PathMappingRequest'];
export type MediaProfileForm = Schemas['MediaProfileRequest'] & { id?: string };
export type CustomFormatForm = Schemas['CustomFormatRequest'] & { id?: string };
export type LanguageForm = {
	code: string;
	originalCode?: string;
	displayName: string;
	aliasesText: string;
};
export type UserForm = {
	id?: string;
	username: string;
	password: string;
	role: Schemas['UserRole'];
};

export interface QualityProfileOption {
	id: string;
	name: string;
	isDefault?: boolean;
	audioTargets?: Schemas['MediaProfileAudioTarget'][];
	subtitleTargets?: Schemas['MediaProfileSubtitleTarget'][];
	removeUnwantedAudio?: boolean;
	removeUnwantedSubtitles?: boolean;
}

export type AppView =
	| 'home'
	| 'settings'
	| 'system'
	| 'profile'
	| 'advanced-search'
	| 'metadata-detail'
	| 'media-people'
	| 'person-detail'
	| 'media-collection'
	| 'related-section'
	| 'discover-section'
	| 'discover-movies'
	| 'discover-series';
export type HomeSection =
	| 'discover'
	| 'blacklist'
	| 'requests'
	| 'movies'
	| 'series'
	| 'wanted'
	| 'activity';
export type ActivitySection = 'queue' | 'history' | 'blocklist';
export type SystemSection = 'status' | 'indexing' | 'metadata' | 'jobs' | 'logs' | 'events';
export type SettingsSection =
	| 'general'
	| 'library'
	| 'download-clients'
	| 'indexers'
	| 'quality'
	| 'profiles'
	| 'custom-formats'
	| 'metadata'
	| 'subtitles'
	| 'languages'
	| 'tags'
	| 'users';
export type TagForm = Schemas['TagRequest'] & { id?: string };
export interface SettingsData {
	downloadClients: Schemas['DownloadClient'][];
	indexers: Schemas['Indexer'][];
	indexerSearch: Schemas['IndexerSearchResponse'];
	metadataProviders: Schemas['MetadataProvider'][];
	subtitleProviders: Schemas['SubtitleProvider'][];
	metadataCache: Schemas['MetadataCacheResponse'];
	libraryFolders: Schemas['LibraryFolder'][];
	pathMappings: Schemas['PathMapping'][];
	mediaProfiles: Schemas['MediaProfile'][];
	customFormats: Schemas['CustomFormat'][];
	users: Schemas['ManagedUser'][];
	tags: Schemas['Tag'][];
	languages: Schemas['Language'][];
}

export type IntegrationTestResults = Record<string, Schemas['IntegrationTestResponse'] | undefined>;

export interface ReleaseSearchState {
	loaded: boolean;
	releases: Schemas['ReleaseCandidate'][];
	errors: string[];
}

export type ReleaseSearchResults = Record<string, ReleaseSearchState | undefined>;
