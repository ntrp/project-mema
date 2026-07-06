import {
	basePrimaryItems,
	settingsPrimaryItem,
	systemPrimaryItem
} from '$lib/components/app/navigation/appNavigation';
import { mediaMetadataDetail } from '$lib/components/app/media/detail/mediaDetail';
import { emptyIndexerSearch, emptyMetadataCache } from '$lib/settings/api';
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
	CustomFormat,
	CustomFormatForm as CustomFormatFormValue,
	DiscoverBlacklistItem,
	DownloadActivity,
	DownloadClient,
	DownloadClientForm as DownloadClientFormValue,
	HomeSection,
	Indexer,
	IndexerForm as IndexerFormValue,
	IndexerSearchResponse,
	IntegrationTestResults,
	Language,
	LanguageForm,
	LibraryFolder,
	LibraryFolderForm as LibraryFolderFormValue,
	LibraryScan,
	ManagedUser,
	MediaCollection,
	MediaDiscoverSection,
	MediaItem,
	MediaMetadataDetails,
	PersonDetails,
	MediaProfile,
	MediaProfileForm as MediaProfileFormValue,
	MediaRequest,
	MediaSearchGroup,
	MediaSearchResult,
	MetadataCacheResponse,
	MetadataProvider,
	PathMapping,
	PathMappingForm,
	ReleaseSearchResults,
	ReleaseBlocklistItem,
	SettingsSection,
	SubtitleProvider,
	SubtitleProviderForm,
	SystemSection,
	Tag,
	TagForm,
	UserProfile,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';
import { relatedSectionFromDetail } from './discoverFilters';
import { defaultRouteState, type AppRouteState } from './routeState';
import type { PeopleSectionKind, RelatedSectionKind } from './types';

export class AppShellState {
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
	loadingProfile = $state(false);
	savingProfile = $state(false);
	profileErrorMessage = $state('');
	message = $state('');
	errorMessage = $state('');
	username = $state('admin');
	password = $state('admin');
	downloadClients = $state<DownloadClient[]>([]);
	indexers = $state<Indexer[]>([]);
	indexerSearch = $state<IndexerSearchResponse>(emptyIndexerSearch());
	metadataProviders = $state<MetadataProvider[]>([]);
	subtitleProviders = $state<SubtitleProvider[]>([]);
	metadataCache = $state<MetadataCacheResponse>(emptyMetadataCache());
	libraryFolders = $state<LibraryFolder[]>([]);
	pathMappings = $state<PathMapping[]>([]);
	mediaProfiles = $state<MediaProfile[]>([]);
	customFormats = $state<CustomFormat[]>([]);
	users = $state<ManagedUser[]>([]);
	profile = $state<UserProfile | undefined>();
	tags = $state<Tag[]>([]);
	languages = $state<Language[]>([]);
	currentUser = $state<UserSummary | undefined>();
	mediaItems = $state<MediaItem[]>([]);
	mediaRequests = $state<MediaRequest[]>([]);
	discoverSections = $state<MediaDiscoverSection[]>([]);
	discoverSection = $state<MediaDiscoverSection | undefined>();
	discoverBlacklist = $state<DiscoverBlacklistItem[]>([]);
	metadataDetail = $state<MediaMetadataDetails | undefined>();
	personDetail = $state<PersonDetails | undefined>();
	mediaCollection = $state<MediaCollection | undefined>();
	autocompleteGroups = $state<MediaSearchGroup[]>([]);
	advancedSearchGroups = $state<MediaSearchGroup[]>([]);
	releaseResults = $state<ReleaseSearchResults>({});
	activities = $state<DownloadActivity[]>([]);
	releaseBlocklist = $state<ReleaseBlocklistItem[]>([]);
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
	loadingMetadataCache = $state(false);
	clearingMetadataCache = $state(false);
	loadingIndexerSearch = $state(false);
	clearingIndexerSearchCache = $state(false);
	savingIndexerSearchSettings = $state(false);
	loadingDiscover = $state(false);
	loadingDiscoverSection = $state(false);
	loadingMoreDiscoverSection = $state(false);
	discoverSectionPage = $state(1);
	discoverSectionHasMore = $state(true);
	loadingBlacklist = $state(false);
	loadingMediaItems = $state(false);
	loadingMetadataDetail = $state(false);
	loadingPersonDetail = $state(false);
	loadingMediaCollection = $state(false);
	loadingAutocomplete = $state(false);
	searchingAdvanced = $state(false);
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
	cancellingActivityId = $state<string | undefined>();
	deletingActivityId = $state<string | undefined>();
	deletingReleaseBlocklistId = $state<string | undefined>();
	clearingReleaseBlocklist = $state(false);
	loadingActivity = $state(false);
	scanningLibraryFolderId = $state<string | undefined>();
	libraryScansByFolder = $state<Record<string, LibraryScan>>({});
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
	eventSource: EventSource | undefined;

	mediaPeopleDetail = $derived(
		this.metadataDetail ??
			(this.selectedMediaItemId
				? this.mediaItems
						.filter(
							(item) => item.type === (this.activeHomeSection === 'movies' ? 'movie' : 'serie')
						)
						.find((item) => item.id === this.selectedMediaItemId)
				: undefined)
	);
	mediaPeopleMetadataDetail = $derived(
		this.mediaPeopleDetail && 'id' in this.mediaPeopleDetail
			? mediaMetadataDetail(this.mediaPeopleDetail)
			: this.mediaPeopleDetail
	);
	relatedMediaSection = $derived(
		relatedSectionFromDetail(
			this.metadataDetail,
			this.activeRelatedSectionKind,
			this.discoverBlacklist
		)
	);
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
		this.route = route;
		this.activeView = route.view;
		this.activeHomeSection = route.homeSection;
		this.activeActivitySection = route.activitySection;
		this.activeSettingsSection = route.settingsSection;
		this.activeSystemSection = route.systemSection;
		this.activeDiscoverSectionId = route.discoverSectionId;
		this.activeDiscoverSubmenuSection = route.discoverSubmenuSection;
		this.activeRelatedSectionKind = route.relatedSectionKind;
		this.activePeopleSectionKind = route.peopleSectionKind;
		this.selectedMediaItemId = route.selectedMediaItemId;
		this.selectedRequestId = route.selectedRequestId;
		this.searchQuery = route.advancedQuery;
	}
}
