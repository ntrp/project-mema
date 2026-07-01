import type { LibraryScanImportRow } from '$lib/components/settings/libraryScanImport';
import type {
	CustomFormat,
	CustomFormatForm as CustomFormatFormValue,
	DownloadClient,
	DownloadClientForm as DownloadClientFormValue,
	Indexer,
	IndexerForm as IndexerFormValue,
	IntegrationTestResponse,
	IntegrationTestResults,
	LibraryFolder,
	LibraryFolderForm as LibraryFolderFormValue,
	LibraryMediaKind,
	LibraryScan,
	ManagedUser,
	MediaProfile,
	MediaProfileForm as MediaProfileFormValue,
	MediaSearchResult,
	MetadataCacheResponse,
	MetadataProvider,
	MetadataProviderForm as MetadataProviderFormValue,
	PathMapping,
	PathMappingForm,
	SettingsSection,
	Tag,
	TagForm,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';

export const staticSettingsSections = [
	'metadata',
	'quality',
	'profiles',
	'custom-formats',
	'file-naming',
	'tags'
] as const;

export function isStaticSettingsSection(section: SettingsSection) {
	return staticSettingsSections.includes(section as (typeof staticSettingsSections)[number]);
}

export interface SettingsAreaProps {
	activeSection: SettingsSection;
	downloadClients: DownloadClient[];
	indexers: Indexer[];
	metadataProviders: MetadataProvider[];
	metadataCache: MetadataCacheResponse;
	libraryFolders: LibraryFolder[];
	pathMappings: PathMapping[];
	mediaProfiles: MediaProfile[];
	customFormats: CustomFormat[];
	users: ManagedUser[];
	tags: Tag[];
	currentUser?: UserSummary;
	libraryScansByFolder: Record<string, LibraryScan>;
	openLibraryFolderId?: string;
	downloadForm: DownloadClientFormValue;
	indexerForm: IndexerFormValue;
	libraryFolderForm: LibraryFolderFormValue;
	pathMappingForm: PathMappingForm;
	mediaProfileForm: MediaProfileFormValue;
	customFormatForm: CustomFormatFormValue;
	tagForm: TagForm;
	userForm: UserFormValue;
	savingDownloadClient: boolean;
	savingIndexer: boolean;
	savingMetadataProviderId?: string;
	loadingMetadataCache: boolean;
	clearingMetadataCache: boolean;
	metadataCachePattern: string;
	savingLibraryFolder: boolean;
	savingPathMapping: boolean;
	deletingPathMappingId?: string;
	savingMediaProfile: boolean;
	deletingMediaProfileId?: string;
	savingCustomFormat: boolean;
	deletingCustomFormatId?: string;
	savingTag: boolean;
	deletingTagId?: string;
	savingUser: boolean;
	scanningLibraryFolderId?: string;
	testingIndexerId?: string;
	testingMetadataProviderId?: string;
	indexerTests: IntegrationTestResults;
	metadataProviderTests: IntegrationTestResults;
	onSaveDownloadClient: (_event: SubmitEvent) => void | Promise<void>;
	onTestDownloadClientConfig: (_form: DownloadClientFormValue) => Promise<IntegrationTestResponse>;
	onSaveIndexer: (_event: SubmitEvent) => void | Promise<void>;
	onSaveMetadataProvider: (_form: MetadataProviderFormValue) => void | Promise<void>;
	onRefreshMetadataCache: () => void | Promise<void>;
	onClearMetadataCache: () => void | Promise<void>;
	onClearMetadataCachePattern: (_event: SubmitEvent) => void | Promise<void>;
	onSaveLibraryFolder: (_event: SubmitEvent) => void | Promise<void>;
	onScanLibraryFolder: (_id: string) => void | Promise<void>;
	onSavePathMapping: (_event: SubmitEvent) => void | Promise<void>;
	onSaveMediaProfile: (_event: SubmitEvent) => void | Promise<void>;
	onSaveCustomFormat: (_event: SubmitEvent) => void | Promise<void>;
	onImportCustomFormat: (_format: CustomFormatFormValue) => void | Promise<void>;
	onSaveTag: (_event: SubmitEvent) => void | Promise<void>;
	onSaveUser: (_event: SubmitEvent) => void | Promise<void>;
	onCancelDownloadClient: () => void;
	onCancelIndexer: () => void;
	onCancelMediaProfile: () => void;
	onCancelCustomFormat: () => void;
	onCancelTag: () => void;
	onCancelUser: () => void;
	onEditDownloadClient: (_client: DownloadClient) => void;
	onEditIndexer: (_indexer: Indexer) => void;
	onEditMediaProfile: (_profile: MediaProfile) => void;
	onEditCustomFormat: (_format: CustomFormat) => void;
	onEditTag: (_tag: Tag) => void;
	onEditUser: (_user: ManagedUser) => void;
	onDeleteDownloadClient: (_id: string) => void | Promise<void>;
	onDeleteIndexer: (_id: string) => void | Promise<void>;
	onDeleteLibraryFolder: (_id: string) => void | Promise<void>;
	onDeletePathMapping: (_id: string) => void | Promise<void>;
	onDeleteMediaProfile: (_id: string) => void | Promise<void>;
	onDeleteCustomFormat: (_id: string) => void | Promise<void>;
	onDeleteTag: (_id: string) => void | Promise<void>;
	onDeleteUser: (_id: string) => void | Promise<void>;
	onTestIndexer: (_id: string) => void | Promise<void>;
	onTestMetadataProvider: (_id: string) => void | Promise<void>;
	onSearchLibraryMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
	onImportLibraryScanRows: (_scan: LibraryScan, _rows: LibraryScanImportRow[]) => Promise<void>;
}
