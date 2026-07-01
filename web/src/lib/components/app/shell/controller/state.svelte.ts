/* global $derived, $state, EventSource */
import {
	basePrimaryItems,
	settingsPrimaryItem,
	systemPrimaryItem
} from '$lib/components/app/navigation/appNavigation';
import { mediaMetadataDetail } from '$lib/components/app/media/mediaDetail';
import { emptyMetadataCache } from '$lib/settings/api';
import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLibraryFolderForm,
	emptyMediaProfileForm,
	emptyPathMappingForm,
	emptyUserForm
} from '$lib/settings/forms';
import type {
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
	IntegrationTestResults,
	LibraryFolder,
	LibraryFolderForm as LibraryFolderFormValue,
	LibraryScan,
	ManagedUser,
	MediaCollection,
	MediaDiscoverSection,
	MediaItem,
	MediaMetadataDetails,
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
	SettingsSection,
	SystemSection,
	Tag,
	TagForm,
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
	savingLibraryFolder = $state(false);
	savingPathMapping = $state(false);
	deletingPathMappingId = $state<string | undefined>();
	savingMediaProfile = $state(false);
	deletingMediaProfileId = $state<string | undefined>();
	savingCustomFormat = $state(false);
	deletingCustomFormatId = $state<string | undefined>();
	savingTag = $state(false);
	deletingTagId = $state<string | undefined>();
	savingUser = $state(false);
	message = $state('');
	errorMessage = $state('');
	username = $state('admin');
	password = $state('admin');
	downloadClients = $state<DownloadClient[]>([]);
	indexers = $state<Indexer[]>([]);
	metadataProviders = $state<MetadataProvider[]>([]);
	metadataCache = $state<MetadataCacheResponse>(emptyMetadataCache());
	libraryFolders = $state<LibraryFolder[]>([]);
	pathMappings = $state<PathMapping[]>([]);
	mediaProfiles = $state<MediaProfile[]>([]);
	customFormats = $state<CustomFormat[]>([]);
	users = $state<ManagedUser[]>([]);
	tags = $state<Tag[]>([]);
	currentUser = $state<UserSummary | undefined>();
	mediaItems = $state<MediaItem[]>([]);
	mediaRequests = $state<MediaRequest[]>([]);
	discoverSections = $state<MediaDiscoverSection[]>([]);
	discoverSection = $state<MediaDiscoverSection | undefined>();
	discoverBlacklist = $state<DiscoverBlacklistItem[]>([]);
	metadataDetail = $state<MediaMetadataDetails | undefined>();
	mediaCollection = $state<MediaCollection | undefined>();
	autocompleteGroups = $state<MediaSearchGroup[]>([]);
	advancedSearchGroups = $state<MediaSearchGroup[]>([]);
	releaseResults = $state<ReleaseSearchResults>({});
	activities = $state<DownloadActivity[]>([]);
	downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	libraryFolderForm = $state<LibraryFolderFormValue>(emptyLibraryFolderForm());
	pathMappingForm = $state<PathMappingForm>(emptyPathMappingForm());
	mediaProfileForm = $state<MediaProfileFormValue>(emptyMediaProfileForm());
	customFormatForm = $state<CustomFormatFormValue>(emptyCustomFormatForm());
	tagForm = $state<TagForm>({ name: '' });
	userForm = $state<UserFormValue>(emptyUserForm());
	testingIndexerId = $state<string | undefined>();
	testingMetadataProviderId = $state<string | undefined>();
	loadingMetadataCache = $state(false);
	clearingMetadataCache = $state(false);
	metadataCachePattern = $state('');
	loadingDiscover = $state(false);
	loadingDiscoverSection = $state(false);
	loadingMoreDiscoverSection = $state(false);
	discoverSectionPage = $state(1);
	discoverSectionHasMore = $state(true);
	loadingBlacklist = $state(false);
	loadingMetadataDetail = $state(false);
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
	grabbingKey = $state<string | undefined>();
	deletingMediaItemId = $state<string | undefined>();
	cancellingActivityId = $state<string | undefined>();
	deletingActivityId = $state<string | undefined>();
	loadingActivity = $state(false);
	scanningLibraryFolderId = $state<string | undefined>();
	libraryScansByFolder = $state<Record<string, LibraryScan>>({});
	openLibraryFolderId = $state<string | undefined>();
	indexerTests = $state<IntegrationTestResults>({});
	metadataProviderTests = $state<IntegrationTestResults>({});
	activeView = $state<AppView>('home');
	activeHomeSection = $state<HomeSection>('discover');
	activeSettingsSection = $state<SettingsSection>('library');
	activeSystemSection = $state<SystemSection>('status');
	activeDiscoverSectionId = $state<string | undefined>();
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
							(item) => item.type === (this.activeHomeSection === 'movies' ? 'movie' : 'series')
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
				: this.activeView === 'discover-section'
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
			: this.activeView === 'discover-section'
				? this.activeDiscoverSectionId
				: this.activePrimarySection === 'library'
					? this.activeHomeSection
					: this.activePrimarySection === 'discover'
						? this.activeHomeSection
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
		this.activeSettingsSection = route.settingsSection;
		this.activeSystemSection = route.systemSection;
		this.activeDiscoverSectionId = route.discoverSectionId;
		this.activeRelatedSectionKind = route.relatedSectionKind;
		this.activePeopleSectionKind = route.peopleSectionKind;
		this.selectedMediaItemId = route.selectedMediaItemId;
		this.selectedRequestId = route.selectedRequestId;
		this.searchQuery = route.advancedQuery;
	}
}
