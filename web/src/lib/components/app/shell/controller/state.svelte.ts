import {
	basePrimaryItems,
	settingsPrimaryItem,
	systemPrimaryItem
} from '$lib/components/app/navigation/appNavigation';
import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyLibraryFolderForm,
	emptyMediaProfileForm,
	emptyPathMappingForm,
	emptySubtitleProviderForm,
	emptyUserForm
} from '$lib/settings/forms';
import type {
	ActivitySection,
	AppView,
	CustomFormatForm as CustomFormatFormValue,
	DownloadClientForm as DownloadClientFormValue,
	HomeSection,
	IndexerForm as IndexerFormValue,
	IntegrationTestResults,
	IndexerSearchResponse,
	LanguageForm,
	LibraryFolderForm as LibraryFolderFormValue,
	MediaItem,
	MetadataCacheResponse,
	MediaProfileForm as MediaProfileFormValue,
	MediaSearchResult,
	PathMappingForm,
	SettingsSection,
	SubtitleProviderForm,
	SystemSection,
	TagForm,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';
import type { UserProfile } from '$lib/profile/types';
import { defaultRouteState, type AppRouteState } from './routeState';
import type { PeopleSectionKind, RelatedSectionKind } from './types';
import { initialiseRouteState } from './stateInitializers';

export class AppShellState {
	declare indexerSearch: IndexerSearchResponse;
	declare metadataCache: MetadataCacheResponse;
	declare profile: UserProfile | undefined;
	declare loadingIndexerSearch: boolean;
	declare loadingMetadataCache: boolean;
	declare loadingProfile: boolean;
	authenticated = $state(false);
	loading = $state(true);
	savingDownloadClient = $state(false);
	savingIndexer = $state(false);
	savingMetadataProviderId = $state<string | undefined>();
	savingSubtitleProviderId = $state<string | undefined>();
	savingLibraryFolder = $state(false);
	savingPathMapping = $state(false);
	deletingPathMappingId = $state<string | undefined>();
	savingMediaProfile = $state(false);
	deletingMediaProfileId = $state<string | undefined>();
	savingCustomFormat = $state(false);
	deletingCustomFormatId = $state<string | undefined>();
	savingTag = $state(false);
	deletingTagId = $state<string | undefined>();
	savingLanguage = $state(false);
	deletingLanguageCode = $state<string | undefined>();
	savingUser = $state(false);
	savingProfile = $state(false);
	profileErrorMessage = $state('');
	message = $state('');
	errorMessage = $state('');
	username = $state('admin');
	password = $state('admin');
	currentUser = $state<UserSummary | undefined>();
	downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	libraryFolderForm = $state<LibraryFolderFormValue>(emptyLibraryFolderForm());
	pathMappingForm = $state<PathMappingForm>(emptyPathMappingForm());
	mediaProfileForm = $state<MediaProfileFormValue>(emptyMediaProfileForm());
	customFormatForm = $state<CustomFormatFormValue>(emptyCustomFormatForm());
	tagForm = $state<TagForm>({ name: '' });
	languageForm = $state<LanguageForm>(emptyLanguageForm());
	userForm = $state<UserFormValue>(emptyUserForm());
	subtitleProviderForm = $state<SubtitleProviderForm>(emptySubtitleProviderForm());
	testingIndexerId = $state<string | undefined>();
	testingMetadataProviderId = $state<string | undefined>();
	testingSubtitleProviderId = $state<string | undefined>();
	clearingMetadataCache = $state(false);
	clearingIndexerSearchCache = $state(false);
	savingIndexerSearchSettings = $state(false);
	addingKey = $state<string | undefined>();
	blacklistingKey = $state<string | undefined>();
	removingBlacklistId = $state<string | undefined>();
	savingMediaAction = $state(false);
	activeMediaCandidate = $state<MediaSearchResult | undefined>();
	mediaDeleteCandidate = $state<MediaItem | undefined>();
	approvingRequestId = $state<string | undefined>();
	searchingItemId = $state<string | undefined>();
	scanningMediaItemId = $state<string | undefined>();
	refreshingMetadataItemId = $state<string | undefined>();
	savingMediaItemOptionsId = $state<string | undefined>();
	grabbingKey = $state<string | undefined>();
	deletingMediaItemId = $state<string | undefined>();
	assemblingMediaItemId = $state<string | undefined>();
	reviewingComponentDecisionId = $state<string | undefined>();
	pendingFulfillmentActions = $state<Record<string, number>>({});
	scanningLibraryFolderId = $state<string | undefined>();
	openLibraryFolderId = $state<string | undefined>();
	indexerTests = $state<IntegrationTestResults>({});
	metadataProviderTests = $state<IntegrationTestResults>({});
	subtitleProviderTests = $state<IntegrationTestResults>({});
	activeView = $state<AppView>('home');
	activeHomeSection = $state<HomeSection>('discover');
	activeActivitySection = $state<ActivitySection>('queue');
	activeSettingsSection = $state<SettingsSection>('general');
	activeSystemSection = $state<SystemSection>('status');
	activeDiscoverSectionId = $state<string | undefined>();
	activeDiscoverSubmenuSection = $state<string | undefined>();
	activeRelatedSectionKind = $state<RelatedSectionKind>('recommendations');
	activePeopleSectionKind = $state<PeopleSectionKind>('cast');
	selectedMediaItemId = $state<string | undefined>();
	selectedRequestId = $state<string | undefined>();
	searchQuery = $state('');
	route = $state<AppRouteState>(defaultRouteState());
	isAdmin = $derived(this.currentUser?.role === 'admin');
	activePrimarySection = $derived(
		this.activeView === 'settings'
			? 'settings'
			: this.activeView === 'system'
				? 'system'
				: this.activeView === 'discover-section' ||
					  this.activeView === 'discover-movies' ||
					  this.activeView === 'discover-series'
					? 'discover'
					: this.activeHomeSection === 'movies' ||
						  this.activeHomeSection === 'series' ||
						  this.activeHomeSection === 'wanted'
						? 'library'
						: this.activeHomeSection
	);
	activeSubmenuSection = $derived(
		this.activeView === 'system'
			? this.activeSystemSection
			: this.activeView === 'discover-movies'
				? (this.activeDiscoverSubmenuSection ?? 'movies')
				: this.activeView === 'discover-series'
					? (this.activeDiscoverSubmenuSection ?? 'series')
					: this.activeView === 'discover-section'
						? this.activeDiscoverSectionId
						: this.activePrimarySection === 'library'
							? this.activeHomeSection
							: this.activePrimarySection === 'discover'
								? this.activeHomeSection
								: this.activePrimarySection === 'activity'
									? this.activeActivitySection
									: this.activeSettingsSection
	);
	primaryItems = $derived(
		this.isAdmin
			? [...basePrimaryItems, settingsPrimaryItem, systemPrimaryItem]
			: basePrimaryItems.filter((item) => item.value !== 'blacklist')
	);

	constructor(route: AppRouteState = defaultRouteState()) {
		initialiseRouteState(this, route);
	}
}
