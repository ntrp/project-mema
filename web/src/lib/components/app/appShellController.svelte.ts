/* global $derived, $state, EventSource, MessageEvent */
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';

import {
	basePrimaryItems,
	settingsPrimaryItem,
	settingsSectionHref,
	systemPrimaryItem,
	systemSectionHref
} from '$lib/components/app/appNavigation';
import {
	advancedSearchMedia as advancedSearchMediaRequest,
	addDiscoverBlacklistItem as addDiscoverBlacklistItemRequest,
	approveMediaRequest as approveMediaRequestRequest,
	autocompleteMedia as autocompleteMediaRequest,
	cancelDownloadActivity as cancelDownloadActivityRequest,
	clearMetadataCache as clearMetadataCacheRequest,
	clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
	createMediaItem as createMediaItemRequest,
	createMediaRequest as createMediaRequestRequest,
	currentSession as currentSessionRequest,
	deleteCustomFormat as deleteCustomFormatRequest,
	deleteDiscoverBlacklistItem as deleteDiscoverBlacklistItemRequest,
	deleteDownloadActivity as deleteDownloadActivityRequest,
	deleteDownloadClient as deleteDownloadClientRequest,
	deleteIndexer as deleteIndexerRequest,
	deleteMediaItemFile as deleteMediaItemFileRequest,
	deleteLibraryFolder as deleteLibraryFolderRequest,
	deleteMediaItem as deleteMediaItemRequest,
	deleteMediaProfile as deleteMediaProfileRequest,
	deletePathMapping as deletePathMappingRequest,
	deleteTag as deleteTagRequest,
	deleteUser as deleteUserRequest,
	emptyMetadataCache,
	enqueueMediaAutomaticSearch as enqueueMediaAutomaticSearchRequest,
	enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
	getMediaCollection as getMediaCollectionRequest,
	getMediaMetadataDetails as getMediaMetadataDetailsRequest,
	getMetadataCache as getMetadataCacheRequest,
	grabMediaRelease as grabMediaReleaseRequest,
	listDownloadActivity as listDownloadActivityRequest,
	listDiscoverBlacklist as listDiscoverBlacklistRequest,
	listMediaItems as listMediaItemsRequest,
	listMediaRequests as listMediaRequestsRequest,
	loadMediaDiscoverSection as loadMediaDiscoverSectionRequest,
	loadMediaDiscoverSections as loadMediaDiscoverSectionsRequest,
	loadSettings as loadSettingsRequest,
	login as loginRequest,
	logout as logoutRequest,
	matchLibraryScanItem as matchLibraryScanItemRequest,
	mediaTypeForLibraryKind,
	saveCustomFormat as saveCustomFormatRequest,
	saveDownloadClient as saveDownloadClientRequest,
	saveIndexer as saveIndexerRequest,
	saveLibraryFolder as saveLibraryFolderRequest,
	saveMediaProfile as saveMediaProfileRequest,
	saveMetadataProvider as saveMetadataProviderRequest,
	savePathMapping as savePathMappingRequest,
	rescanMediaItemFiles as rescanMediaItemFilesRequest,
	scanLibraryFolder as scanLibraryFolderRequest,
	saveTag as saveTagRequest,
	saveUser as saveUserRequest,
	searchMedia as searchMediaRequest,
	searchMediaReleases as searchMediaReleasesRequest,
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest
} from '$lib/settings/api';
import {
	customFormatFormFromFormat,
	downloadClientFormFromClient,
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLibraryFolderForm,
	emptyMediaProfileForm,
	emptyPathMappingForm,
	emptyUserForm,
	indexerFormFromIndexer,
	mediaProfileFormFromProfile,
	userFormFromUser
} from '$lib/settings/forms';
import type { MediaActionSelection } from '$lib/components/app/mediaActionTypes';
import { mediaMetadataDetail } from '$lib/components/app/mediaDetail';
import type { LibraryScanImportRow } from '$lib/components/settings/libraryScanImport';
import type {
	AppView,
	CustomFormat,
	CustomFormatForm as CustomFormatFormValue,
	DiscoverBlacklistItem,
	DownloadActivity,
	DownloadActivityStatus,
	DownloadClient,
	DownloadClientForm as DownloadClientFormValue,
	HomeSection,
	Indexer,
	IndexerForm as IndexerFormValue,
	IntegrationTestResults,
	LibraryFolder,
	LibraryFolderForm as LibraryFolderFormValue,
	LibraryMediaKind,
	LibraryScan,
	ManagedUser,
	MediaAdvancedSearchRequest,
	MediaCollection,
	MediaDiscoverSection,
	MediaItem,
	MediaProfile,
	MediaProfileForm as MediaProfileFormValue,
	MediaMetadataDetails,
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchGroup,
	MediaSearchResult,
	MediaType,
	MetadataCacheResponse,
	MetadataProvider,
	MetadataProviderForm as MetadataProviderFormValue,
	MetadataProviderType,
	PathMapping,
	PathMappingForm,
	ReleaseCandidate,
	ReleaseSearchResults,
	SettingsSection,
	SystemSection,
	Tag,
	TagForm,
	UserForm as UserFormValue,
	UserSummary
} from '$lib/settings/types';

interface ServerEventEnvelope<T = unknown> {
	id: string;
	type: string;
	time: string;
	data: T;
}

export interface AppShellOptions {
	initialView?: AppView;
	initialHomeSection?: HomeSection;
	initialSettingsSection?: SettingsSection;
	initialSystemSection?: SystemSection;
	initialSelectedMediaItemId?: string;
	initialSelectedRequestId?: string;
	initialAdvancedQuery?: string;
	initialMetadataProvider?: string;
	initialMetadataType?: string;
	initialMetadataExternalId?: string;
	initialCollectionProvider?: string;
	initialCollectionId?: string;
	initialDiscoverSectionId?: string;
	initialRelatedSectionKind?: RelatedSectionKind;
	initialPeopleSectionKind?: PeopleSectionKind;
}

export type RelatedSectionKind = 'recommendations' | 'similar';
export type PeopleSectionKind = 'cast' | 'crew';

export function createAppShellController(options: AppShellOptions) {
	let authenticated = $state(false);
	let loading = $state(true);
	let savingDownloadClient = $state(false);
	let savingIndexer = $state(false);
	let savingMetadataProviderId = $state<string | undefined>();
	let savingLibraryFolder = $state(false);
	let savingPathMapping = $state(false);
	let deletingPathMappingId = $state<string | undefined>();
	let savingMediaProfile = $state(false);
	let deletingMediaProfileId = $state<string | undefined>();
	let savingCustomFormat = $state(false);
	let deletingCustomFormatId = $state<string | undefined>();
	let savingTag = $state(false);
	let deletingTagId = $state<string | undefined>();
	let savingUser = $state(false);
	let message = $state('');
	let errorMessage = $state('');
	let username = $state('admin');
	let password = $state('admin');
	let downloadClients = $state<DownloadClient[]>([]);
	let indexers = $state<Indexer[]>([]);
	let metadataProviders = $state<MetadataProvider[]>([]);
	let metadataCache = $state<MetadataCacheResponse>(emptyMetadataCache());
	let libraryFolders = $state<LibraryFolder[]>([]);
	let pathMappings = $state<PathMapping[]>([]);
	let mediaProfiles = $state<MediaProfile[]>([]);
	let customFormats = $state<CustomFormat[]>([]);
	let users = $state<ManagedUser[]>([]);
	let tags = $state<Tag[]>([]);
	let currentUser = $state<UserSummary | undefined>();
	let mediaItems = $state<MediaItem[]>([]);
	let mediaRequests = $state<MediaRequest[]>([]);
	let discoverSections = $state<MediaDiscoverSection[]>([]);
	let discoverSection = $state<MediaDiscoverSection | undefined>();
	let discoverBlacklist = $state<DiscoverBlacklistItem[]>([]);
	let metadataDetail = $state<MediaMetadataDetails | undefined>();
	let mediaCollection = $state<MediaCollection | undefined>();
	let autocompleteGroups = $state<MediaSearchGroup[]>([]);
	let advancedSearchGroups = $state<MediaSearchGroup[]>([]);
	let releaseResults = $state<ReleaseSearchResults>({});
	let activities = $state<DownloadActivity[]>([]);
	let downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	let indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	let libraryFolderForm = $state<LibraryFolderFormValue>(emptyLibraryFolderForm());
	let pathMappingForm = $state<PathMappingForm>(emptyPathMappingForm());
	let mediaProfileForm = $state<MediaProfileFormValue>(emptyMediaProfileForm());
	let customFormatForm = $state<CustomFormatFormValue>(emptyCustomFormatForm());
	let tagForm = $state<TagForm>(emptyTagForm());
	let userForm = $state<UserFormValue>(emptyUserForm());
	let testingIndexerId = $state<string | undefined>();
	let testingMetadataProviderId = $state<string | undefined>();
	let loadingMetadataCache = $state(false);
	let clearingMetadataCache = $state(false);
	let metadataCachePattern = $state('');
	let loadingDiscover = $state(false);
	let loadingDiscoverSection = $state(false);
	let loadingMoreDiscoverSection = $state(false);
	let discoverSectionPage = $state(1);
	let discoverSectionHasMore = $state(true);
	let loadingBlacklist = $state(false);
	let loadingMetadataDetail = $state(false);
	let loadingMediaCollection = $state(false);
	let loadingAutocomplete = $state(false);
	let searchingAdvanced = $state(false);
	let addingKey = $state<string | undefined>();
	let blacklistingKey = $state<string | undefined>();
	let removingBlacklistId = $state<string | undefined>();
	let savingMediaAction = $state(false);
	let activeMediaCandidate = $state<MediaSearchResult | undefined>();
	let mediaDeleteCandidate = $state<MediaItem | undefined>();
	let approvingRequestId = $state<string | undefined>();
	let searchingItemId = $state<string | undefined>();
	let scanningMediaItemId = $state<string | undefined>();
	let grabbingKey = $state<string | undefined>();
	let deletingMediaItemId = $state<string | undefined>();
	let cancellingActivityId = $state<string | undefined>();
	let deletingActivityId = $state<string | undefined>();
	let loadingActivity = $state(false);
	let scanningLibraryFolderId = $state<string | undefined>();
	let libraryScansByFolder = $state<Record<string, LibraryScan>>({});
	let openLibraryFolderId = $state<string | undefined>();
	let indexerTests = $state<IntegrationTestResults>({});
	let metadataProviderTests = $state<IntegrationTestResults>({});
	let activeView = $state<AppView>(options.initialView ?? 'home');
	let activeHomeSection = $state<HomeSection>(options.initialHomeSection ?? 'discover');
	let activeSettingsSection = $state<SettingsSection>(options.initialSettingsSection ?? 'library');
	let activeSystemSection = $state<SystemSection>(options.initialSystemSection ?? 'status');
	let activeDiscoverSectionId = $state<string | undefined>(options.initialDiscoverSectionId);
	let activeRelatedSectionKind = $state<RelatedSectionKind>(
		options.initialRelatedSectionKind ?? 'recommendations'
	);
	let activePeopleSectionKind = $state<PeopleSectionKind>(
		options.initialPeopleSectionKind ?? 'cast'
	);
	let selectedMediaItemId = $state<string | undefined>(options.initialSelectedMediaItemId);
	let selectedRequestId = $state<string | undefined>(options.initialSelectedRequestId);
	let searchQuery = $state(options.initialAdvancedQuery ?? '');
	let eventSource: EventSource | undefined;
	let mediaPeopleDetail = $derived(
		metadataDetail ??
			(selectedMediaItemId
				? mediaItems
						.filter((item) => item.type === (activeHomeSection === 'movies' ? 'movie' : 'series'))
						.find((item) => item.id === selectedMediaItemId)
				: undefined)
	);
	let mediaPeopleMetadataDetail = $derived(
		mediaPeopleDetail && 'id' in mediaPeopleDetail
			? mediaMetadataDetail(mediaPeopleDetail)
			: mediaPeopleDetail
	);
	let relatedMediaSection = $derived(
		relatedSectionFromDetail(metadataDetail, activeRelatedSectionKind)
	);
	let isAdmin = $derived(currentUser?.role === 'admin');
	let activePrimarySection = $derived(
		activeView === 'settings'
			? 'settings'
			: activeView === 'system'
				? 'system'
				: activeView === 'discover-section'
					? 'discover'
					: activeHomeSection === 'movies' ||
						  activeHomeSection === 'series' ||
						  activeHomeSection === 'wanted'
						? 'library'
						: activeHomeSection
	);
	let activeSubmenuSection = $derived(
		activeView === 'system'
			? activeSystemSection
			: activeView === 'discover-section'
				? activeDiscoverSectionId
				: activePrimarySection === 'library'
					? activeHomeSection
					: activePrimarySection === 'discover'
						? activeHomeSection
						: activeSettingsSection
	);
	let primaryItems = $derived(
		isAdmin
			? [...basePrimaryItems, settingsPrimaryItem, systemPrimaryItem]
			: basePrimaryItems.filter((item) => item.value !== 'blacklist')
	);

	async function initialise() {
		loading = true;
		errorMessage = '';

		const session = await currentSessionRequest();
		authenticated = Boolean(session?.authenticated);
		currentUser = session?.user;
		if (authenticated) {
			if (currentUser?.role === 'admin') {
				await loadSettings();
				await loadDiscoverBlacklist();
			} else if (activeView === 'settings' || activeView === 'system') {
				activeView = 'home';
				activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			} else if (activeHomeSection === 'blacklist') {
				activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (
				activeView === 'metadata-detail' ||
				activeView === 'media-people' ||
				activeView === 'related-section'
			) {
				await loadMetadataDetail();
			} else if (activeView === 'media-collection') {
				await loadMediaCollection();
			} else if (activeView === 'discover-section') {
				await loadDiscoverSection();
			}
			connectEvents();
		}

		loading = false;
	}

	function connectEvents() {
		if (!authenticated || eventSource) {
			return;
		}
		const source = new EventSource('/api/events', { withCredentials: true });
		eventSource = source;
		source.addEventListener('activity.download.updated', (event) => {
			const activity = parseEventData<DownloadActivity>(event);
			if (!activity) {
				return;
			}
			upsertActivity(activity);
			updateMediaStatusFromActivity(activity);
			if (activity.status === 'completed') {
				void loadMediaItems();
			}
		});
		source.onerror = () => {
			if (!authenticated) {
				disconnectEvents();
			}
		};
	}

	function disconnectEvents() {
		eventSource?.close();
		eventSource = undefined;
	}

	async function login(event: SubmitEvent) {
		event.preventDefault();
		clearNotice();

		try {
			const session = await loginRequest(username, password);
			authenticated = true;
			currentUser = session.user;
			if (currentUser?.role === 'admin') {
				await loadSettings();
			} else if (activeView === 'settings' || activeView === 'system') {
				activeView = 'home';
				activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (
				activeView === 'metadata-detail' ||
				activeView === 'media-people' ||
				activeView === 'related-section'
			) {
				await loadMetadataDetail();
			} else if (activeView === 'media-collection') {
				await loadMediaCollection();
			} else if (activeView === 'discover-section') {
				await loadDiscoverSection();
			}
			connectEvents();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Login failed');
		}
	}

	async function logout() {
		clearNotice();

		try {
			await logoutRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not log out');
		} finally {
			disconnectEvents();
			authenticated = false;
			currentUser = undefined;
			activeView = 'home';
			activeHomeSection = 'discover';
			downloadClients = [];
			indexers = [];
			metadataProviders = [];
			mediaProfiles = [];
			customFormats = [];
			users = [];
			tags = [];
			mediaItems = [];
			mediaRequests = [];
			discoverSections = [];
			discoverSection = undefined;
			discoverSectionPage = 1;
			discoverSectionHasMore = true;
			metadataDetail = undefined;
			mediaCollection = undefined;
			autocompleteGroups = [];
			advancedSearchGroups = [];
			releaseResults = {};
			activities = [];
			libraryFolders = [];
			pathMappings = [];
			libraryScansByFolder = {};
			openLibraryFolderId = undefined;
			downloadForm = emptyDownloadClientForm();
			indexerForm = emptyIndexerForm();
			libraryFolderForm = emptyLibraryFolderForm();
			pathMappingForm = emptyPathMappingForm();
			mediaProfileForm = emptyMediaProfileForm();
			customFormatForm = emptyCustomFormatForm();
			tagForm = emptyTagForm();
			userForm = emptyUserForm();
		}
	}

	function upsertActivity(activity: DownloadActivity) {
		activities = [activity, ...activities.filter((item) => item.id !== activity.id)];
	}

	function updateMediaStatusFromActivity(activity: DownloadActivity) {
		const status = mediaStatusFromActivity(activity.status);
		if (!status) {
			return;
		}
		mediaItems = mediaItems.map((item) =>
			item.id === activity.mediaItemId ? { ...item, status } : item
		);
	}

	function mediaStatusFromActivity(status: DownloadActivityStatus) {
		if (status === 'completed') {
			return 'downloaded';
		}
		if (status === 'queued' || status === 'grabbed' || status === 'downloading') {
			return 'downloading';
		}
		return undefined;
	}

	function parseEventData<T>(event: Event) {
		const message = event as MessageEvent<string>;
		try {
			return (JSON.parse(message.data) as ServerEventEnvelope<T>).data;
		} catch {
			return undefined;
		}
	}

	function showProfile() {
		clearNotice();
		message = 'Profile settings are not implemented yet';
	}

	async function loadSettings() {
		try {
			const settings = await loadSettingsRequest();
			downloadClients = settings.downloadClients;
			indexers = settings.indexers;
			metadataProviders = settings.metadataProviders;
			metadataCache = settings.metadataCache;
			libraryFolders = settings.libraryFolders;
			pathMappings = settings.pathMappings;
			mediaProfiles = settings.mediaProfiles;
			customFormats = settings.customFormats;
			users = settings.users;
			tags = settings.tags;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load settings');
		}
	}

	async function loadLibrary() {
		await Promise.all([loadMediaItems(), loadMediaRequests(), loadDownloadActivity()]);
	}

	async function loadDiscoverSections() {
		loadingDiscover = true;
		try {
			discoverSections = await loadMediaDiscoverSectionsRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load discover sections');
		} finally {
			loadingDiscover = false;
		}
	}

	async function loadDiscoverSection() {
		if (!activeDiscoverSectionId) {
			return;
		}
		loadingDiscoverSection = true;
		discoverSectionPage = 1;
		discoverSectionHasMore = true;
		try {
			discoverSection = await loadMediaDiscoverSectionRequest(activeDiscoverSectionId, 1);
			discoverSectionHasMore = (discoverSection.results ?? []).length > 0;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load discover section');
		} finally {
			loadingDiscoverSection = false;
		}
	}

	async function loadMoreDiscoverSection() {
		if (
			!activeDiscoverSectionId ||
			loadingDiscoverSection ||
			loadingMoreDiscoverSection ||
			!discoverSectionHasMore
		) {
			return;
		}
		loadingMoreDiscoverSection = true;
		const nextPage = discoverSectionPage + 1;
		try {
			const nextSection = await loadMediaDiscoverSectionRequest(activeDiscoverSectionId, nextPage);
			const existingKeys = (discoverSection?.results ?? []).map(discoverResultKey);
			const nextResults = (nextSection.results ?? []).filter((result) => {
				const key = discoverResultKey(result);
				if (existingKeys.includes(key)) {
					return false;
				}
				existingKeys.push(key);
				return true;
			});
			discoverSectionPage = nextPage;
			discoverSectionHasMore = nextResults.length > 0;
			discoverSection = {
				...nextSection,
				results: [...(discoverSection?.results ?? []), ...nextResults]
			};
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load more discover results');
		} finally {
			loadingMoreDiscoverSection = false;
		}
	}

	async function loadDiscoverBlacklist() {
		if (!isAdmin) {
			return;
		}
		loadingBlacklist = true;
		try {
			discoverBlacklist = await listDiscoverBlacklistRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load discover blacklist');
		} finally {
			loadingBlacklist = false;
		}
	}

	async function blacklistDiscoverMedia(candidate: MediaSearchResult) {
		if (!isAdmin) {
			return;
		}
		blacklistingKey = discoverResultKey(candidate);
		try {
			const item = await addDiscoverBlacklistItemRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath
			});
			discoverBlacklist = [
				item,
				...discoverBlacklist.filter((entry) => !sameDiscoverBlacklistItem(entry, item))
			];
			discoverSections = filterDiscoverSections(discoverSections);
			if (discoverSection) {
				discoverSection = filterDiscoverSection(discoverSection);
			}
			message = `${candidate.title} hidden from discover`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not add media to discover blacklist');
		} finally {
			blacklistingKey = undefined;
		}
	}

	async function removeDiscoverBlacklistItem(item: DiscoverBlacklistItem) {
		if (!isAdmin) {
			return;
		}
		removingBlacklistId = item.id;
		try {
			await deleteDiscoverBlacklistItemRequest(item.id);
			discoverBlacklist = discoverBlacklist.filter((entry) => entry.id !== item.id);
			message = `${item.title} removed from blacklist`;
			await loadDiscoverSections();
			if (activeView === 'discover-section') {
				await loadDiscoverSection();
			}
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not remove media from discover blacklist');
		} finally {
			removingBlacklistId = undefined;
		}
	}

	async function loadMetadataDetail() {
		if (
			!options.initialMetadataProvider ||
			!options.initialMetadataType ||
			!options.initialMetadataExternalId
		) {
			return;
		}
		loadingMetadataDetail = true;
		try {
			metadataDetail = await getMediaMetadataDetailsRequest(
				options.initialMetadataProvider as MetadataProviderType,
				options.initialMetadataType as MediaType,
				options.initialMetadataExternalId
			);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load media details');
		} finally {
			loadingMetadataDetail = false;
		}
	}

	async function loadMediaCollection() {
		if (!options.initialCollectionProvider || !options.initialCollectionId) {
			return;
		}
		loadingMediaCollection = true;
		try {
			mediaCollection = await getMediaCollectionRequest(
				options.initialCollectionProvider as MetadataProviderType,
				options.initialCollectionId
			);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load media collection');
		} finally {
			loadingMediaCollection = false;
		}
	}

	async function loadMediaItems() {
		try {
			mediaItems = await listMediaItemsRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load media items');
		}
	}

	async function loadMediaRequests() {
		try {
			mediaRequests = await listMediaRequestsRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load media requests');
		}
	}

	async function loadDownloadActivity() {
		loadingActivity = true;
		try {
			activities = await listDownloadActivityRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load download activity');
		} finally {
			loadingActivity = false;
		}
	}

	async function autocompleteMedia(query: string) {
		const trimmed = query.trim();
		if (trimmed.length < 2) {
			autocompleteGroups = [];
			return;
		}
		loadingAutocomplete = true;
		autocompleteGroups = [];
		try {
			const groups = await autocompleteMediaRequest(trimmed, 'library');
			if (searchQuery.trim() !== trimmed) {
				return;
			}
			autocompleteGroups = groups;
		} catch {
			autocompleteGroups = [];
		} finally {
			loadingAutocomplete = false;
		}
	}

	async function advancedSearch(request: MediaAdvancedSearchRequest) {
		searchingAdvanced = true;
		clearNotice();

		try {
			advancedSearchGroups = await advancedSearchMediaRequest(request);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not search media');
		} finally {
			searchingAdvanced = false;
		}
	}

	function selectAutocompleteResult(result: MediaSearchResult) {
		if (result.id) {
			void goto(
				resolve(result.type === 'movie' ? '/movies/[id]' : '/series/[id]', { id: result.id })
			);
			return;
		}
		if (result.externalProvider && result.externalId) {
			void goto(
				resolve('/media/[provider]/[type]/[externalId]', {
					provider: result.externalProvider,
					type: result.type,
					externalId: result.externalId
				})
			);
			return;
		}
		searchQuery = result.title;
		void goto(resolve(`/search/advanced?q=${encodeURIComponent(result.title)}`));
	}

	function openAdvancedSearch(query: string) {
		searchQuery = query;
		void goto(resolve(`/search/advanced?q=${encodeURIComponent(query)}`));
	}

	function addMedia(candidate: MediaSearchResult) {
		activeMediaCandidate = candidate;
		clearNotice();
	}

	function closeMediaAction() {
		if (!savingMediaAction) {
			activeMediaCandidate = undefined;
		}
	}

	async function confirmMediaAction(selection: MediaActionSelection) {
		const candidate = activeMediaCandidate;
		if (!candidate) {
			return;
		}
		addingKey = candidateKey(candidate);
		savingMediaAction = true;
		clearNotice();

		try {
			if (isAdmin) {
				if (!selection.qualityProfileId || !selection.libraryFolderId) {
					throw new Error('Quality profile and library folder are required');
				}
				await createMediaItemRequest({
					title: candidate.title,
					type: candidate.type,
					year: candidate.year,
					monitored: selection.monitorMode !== 'none',
					monitorMode: selection.monitorMode,
					seriesType: candidate.type === 'series' ? selection.seriesType : undefined,
					minimumAvailability: selection.minimumAvailability,
					startSearch: selection.startSearch,
					externalProvider: candidate.externalProvider,
					externalId: candidate.externalId,
					overview: candidate.overview,
					posterPath: candidate.posterPath,
					qualityProfileId: selection.qualityProfileId,
					libraryFolderId: selection.libraryFolderId,
					tags: selection.tags
				});
				await loadMediaItems();
				await loadSettings();
				message =
					selection.monitorMode === 'none'
						? 'Media item added to library'
						: selection.monitorMode === 'collection'
							? 'Media collection added to monitored'
							: 'Media item added to monitored';
				activeHomeSection = candidate.type === 'movie' ? 'movies' : 'series';
				activeMediaCandidate = undefined;
				void goto(resolve(candidate.type === 'movie' ? '/movies' : '/series'));
				return;
			}

			const request = await createMediaRequestRequest({
				title: candidate.title,
				type: candidate.type,
				monitorMode: selection.monitorMode,
				seriesType: candidate.type === 'series' ? selection.seriesType : undefined,
				minimumAvailability: selection.minimumAvailability,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath,
				tags: selection.tags
			});
			mediaRequests = [request, ...mediaRequests.filter((item) => item.id !== request.id)];
			message = 'Media request created';
			activeHomeSection = 'requests';
			activeMediaCandidate = undefined;
			void goto(resolve('/requests'));
		} catch (error) {
			errorMessage = errorMessageFrom(
				error,
				isAdmin ? 'Could not add media item' : 'Could not create media request'
			);
		} finally {
			addingKey = undefined;
			savingMediaAction = false;
		}
	}

	async function findReleases(item: MediaItem) {
		searchingItemId = item.id;
		clearNotice();

		try {
			const job = await enqueueMediaReleaseSearchRequest(item.id);
			releaseResults = {
				...releaseResults,
				[item.id]: { loaded: false, releases: [], errors: [`${job.message} (#${job.jobId})`] }
			};
			message = job.message;
			window.setTimeout(() => void loadReleaseResults(item.id), 1200);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not enqueue release search');
		} finally {
			searchingItemId = undefined;
		}
	}

	async function autoSearchMedia(item: MediaItem) {
		searchingItemId = item.id;
		clearNotice();
		try {
			const job = await enqueueMediaAutomaticSearchRequest(item.id);
			message = `${job.message} (#${job.jobId})`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not enqueue automatic search');
		} finally {
			searchingItemId = undefined;
		}
	}

	async function rescanMediaFiles(item: MediaItem) {
		scanningMediaItemId = item.id;
		clearNotice();

		try {
			const updated = await rescanMediaItemFilesRequest(item.id);
			mediaItems = [updated, ...mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)];
			message = `File scan completed: ${updated.filePaths.length} media, ${updated.metadataFilePaths.length} metadata`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not rescan media folder');
		} finally {
			scanningMediaItemId = undefined;
		}
	}

	async function deleteMediaFile(item: MediaItem, path: string) {
		clearNotice();
		try {
			const updated = await deleteMediaItemFileRequest(item.id, path);
			mediaItems = [updated, ...mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)];
			message = 'Media file deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete media file');
		}
	}

	function deleteMediaItem(item: MediaItem) {
		clearNotice();
		if (mediaItemFileCount(item) > 0) {
			mediaDeleteCandidate = item;
			return;
		}
		void removeMediaItem(item, false);
	}

	function closeMediaDelete() {
		if (!deletingMediaItemId) {
			mediaDeleteCandidate = undefined;
		}
	}

	async function confirmMediaDelete(keepFiles: boolean) {
		const item = mediaDeleteCandidate;
		if (!item) {
			return;
		}
		await removeMediaItem(item, keepFiles);
	}

	async function removeMediaItem(item: MediaItem, keepFiles: boolean) {
		deletingMediaItemId = item.id;
		clearNotice();

		try {
			await deleteMediaItemRequest(item.id, { keepFiles });
			mediaItems = mediaItems.filter((mediaItem) => mediaItem.id !== item.id);
			releaseResults = omitResult(releaseResults, item.id);
			activities = activities.filter((activity) => activity.mediaItemId !== item.id);
			mediaDeleteCandidate = undefined;
			message = keepFiles ? 'Media item removed; files kept' : 'Media item and files removed';
			if (selectedMediaItemId === item.id) {
				selectedMediaItemId = undefined;
				void goto(resolve(item.type === 'movie' ? '/movies' : '/series'));
			}
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not remove media item');
		} finally {
			deletingMediaItemId = undefined;
		}
	}

	function mediaItemFileCount(item: MediaItem) {
		return (item.filePaths?.length ?? 0) + (item.metadataFilePaths?.length ?? 0);
	}

	async function approveMediaRequest(request: MediaRequest, approval: MediaRequestApproveRequest) {
		approvingRequestId = request.id;
		clearNotice();

		try {
			const result = await approveMediaRequestRequest(request.id, approval);
			mediaRequests = mediaRequests.map((item) =>
				item.id === result.request.id ? result.request : item
			);
			mediaItems = [
				result.mediaItem,
				...mediaItems.filter((item) => item.id !== result.mediaItem.id)
			];
			message = 'Media request approved and added to monitored';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not approve media request');
		} finally {
			approvingRequestId = undefined;
		}
	}

	async function loadReleaseResults(id: string) {
		try {
			const results = await searchMediaReleasesRequest(id);
			releaseResults = {
				...releaseResults,
				[id]: { loaded: true, releases: results.releases, errors: results.errors }
			};
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load release results');
		}
	}

	async function grabRelease(item: MediaItem, release: ReleaseCandidate) {
		grabbingKey = `${item.id}:${release.id}`;
		clearNotice();

		try {
			const result = await grabMediaReleaseRequest(item.id, release);
			activities = [
				result.activity,
				...activities.filter((activity) => activity.id !== result.activity.id)
			];
			updateMediaStatusFromActivity(result.activity);
			message = `${result.message} (#${result.jobId})`;
			activeHomeSection = 'activity';
			void goto(resolve('/activity'));
			window.setTimeout(() => void loadDownloadActivity(), 1200);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not enqueue download');
		} finally {
			grabbingKey = undefined;
		}
	}

	async function cancelActivity(activity: DownloadActivity) {
		cancellingActivityId = activity.id;
		clearNotice();

		try {
			const cancelled = await cancelDownloadActivityRequest(activity.id);
			upsertActivity(cancelled);
			await loadMediaItems();
			message = 'Download activity cancelled';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not cancel download activity');
		} finally {
			cancellingActivityId = undefined;
		}
	}

	async function deleteActivity(activity: DownloadActivity) {
		deletingActivityId = activity.id;
		clearNotice();

		try {
			await deleteDownloadActivityRequest(activity.id);
			activities = activities.filter((item) => item.id !== activity.id);
			message = 'Download activity deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete download activity');
		} finally {
			deletingActivityId = undefined;
		}
	}

	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		savingDownloadClient = true;
		clearNotice();

		try {
			await saveDownloadClientRequest(downloadForm);
			downloadForm = emptyDownloadClientForm();
			message = 'Download client saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save download client');
		} finally {
			savingDownloadClient = false;
		}
	}

	async function saveIndexer(event: SubmitEvent) {
		event.preventDefault();
		savingIndexer = true;
		clearNotice();

		try {
			await saveIndexerRequest(indexerForm);
			indexerForm = emptyIndexerForm();
			message = 'Indexer saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save indexer');
		} finally {
			savingIndexer = false;
		}
	}

	async function saveMetadataProvider(form: MetadataProviderFormValue) {
		savingMetadataProviderId = form.id;
		clearNotice();

		try {
			await saveMetadataProviderRequest(form);
			message = 'Metadata provider saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save metadata provider');
		} finally {
			savingMetadataProviderId = undefined;
		}
	}

	async function saveLibraryFolder(event: SubmitEvent) {
		event.preventDefault();
		savingLibraryFolder = true;
		clearNotice();

		try {
			const result = await saveLibraryFolderRequest(libraryFolderForm);
			libraryFolderForm = emptyLibraryFolderForm();
			libraryFolders = [
				result.folder,
				...libraryFolders.filter((folder) => folder.id !== result.folder.id)
			];
			libraryScansByFolder = {
				...libraryScansByFolder,
				[result.folder.id]: result.scan
			};
			openLibraryFolderId = result.folder.id;
			message = `Library scan completed: ${result.scan.manualCount} pending`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not add library folder');
		} finally {
			savingLibraryFolder = false;
		}
	}

	async function savePathMapping(event: SubmitEvent) {
		event.preventDefault();
		savingPathMapping = true;
		clearNotice();

		try {
			const mapping = await savePathMappingRequest(pathMappingForm);
			pathMappingForm = emptyPathMappingForm();
			pathMappings = [mapping, ...pathMappings.filter((item) => item.id !== mapping.id)];
			message = 'Path mapping saved';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save path mapping');
		} finally {
			savingPathMapping = false;
		}
	}

	async function saveUser(event: SubmitEvent) {
		event.preventDefault();
		savingUser = true;
		clearNotice();

		try {
			await saveUserRequest(userForm);
			userForm = emptyUserForm();
			message = 'User saved';
			await loadSettings();
			if (currentUser && users.some((user) => user.id === currentUser?.id)) {
				const updatedUser = users.find((user) => user.id === currentUser?.id);
				if (updatedUser) {
					currentUser = {
						id: updatedUser.id,
						username: updatedUser.username,
						role: updatedUser.role
					};
				}
			}
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save user');
		} finally {
			savingUser = false;
		}
	}

	async function saveTag(event: SubmitEvent) {
		event.preventDefault();
		savingTag = true;
		clearNotice();

		try {
			await saveTagRequest(tagForm);
			tagForm = emptyTagForm();
			message = 'Tag saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save tag');
		} finally {
			savingTag = false;
		}
	}

	async function saveMediaProfile(event: SubmitEvent) {
		event.preventDefault();
		savingMediaProfile = true;
		clearNotice();

		try {
			await saveMediaProfileRequest(mediaProfileForm);
			mediaProfileForm = emptyMediaProfileForm();
			message = 'Profile saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save profile');
		} finally {
			savingMediaProfile = false;
		}
	}

	async function saveCustomFormat(event: SubmitEvent) {
		event.preventDefault();
		savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(customFormatForm);
			customFormatForm = emptyCustomFormatForm();
			message = 'Custom format saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save custom format');
		} finally {
			savingCustomFormat = false;
		}
	}

	async function importCustomFormat(format: CustomFormatFormValue) {
		savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(format);
			message = 'Custom format imported';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not import custom format');
			throw error;
		} finally {
			savingCustomFormat = false;
		}
	}

	async function deleteDownloadClient(id: string) {
		clearNotice();

		try {
			await deleteDownloadClientRequest(id);
			if (downloadForm.id === id) {
				downloadForm = emptyDownloadClientForm();
			}
			message = 'Download client deleted';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete download client');
		}
	}

	async function deleteIndexer(id: string) {
		clearNotice();

		try {
			await deleteIndexerRequest(id);
			if (indexerForm.id === id) {
				indexerForm = emptyIndexerForm();
			}
			indexerTests = omitResult(indexerTests, id);
			message = 'Indexer deleted';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete indexer');
		}
	}

	async function deleteLibraryFolder(id: string) {
		clearNotice();

		try {
			await deleteLibraryFolderRequest(id);
			libraryFolders = libraryFolders.filter((folder) => folder.id !== id);
			const remainingScans = { ...libraryScansByFolder };
			delete remainingScans[id];
			libraryScansByFolder = remainingScans;
			if (openLibraryFolderId === id) {
				openLibraryFolderId = undefined;
			}
			message = 'Library folder deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete library folder');
		}
	}

	async function scanLibraryFolder(id: string) {
		scanningLibraryFolderId = id;
		clearNotice();

		try {
			const scan = await scanLibraryFolderRequest(id);
			libraryScansByFolder = { ...libraryScansByFolder, [scan.folderId]: scan };
			openLibraryFolderId = scan.folderId;
			message = `Library scan completed: ${scan.manualCount} pending`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not scan library folder');
		} finally {
			scanningLibraryFolderId = undefined;
		}
	}

	async function deletePathMapping(id: string) {
		deletingPathMappingId = id;
		clearNotice();

		try {
			await deletePathMappingRequest(id);
			pathMappings = pathMappings.filter((mapping) => mapping.id !== id);
			message = 'Path mapping deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete path mapping');
		} finally {
			deletingPathMappingId = undefined;
		}
	}

	async function deleteUser(id: string) {
		clearNotice();

		try {
			await deleteUserRequest(id);
			if (userForm.id === id) {
				userForm = emptyUserForm();
			}
			users = users.filter((user) => user.id !== id);
			message = 'User deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete user');
		}
	}

	async function deleteTag(id: string) {
		deletingTagId = id;
		clearNotice();

		try {
			await deleteTagRequest(id);
			if (tagForm.id === id) {
				tagForm = emptyTagForm();
			}
			tags = tags.filter((tag) => tag.id !== id);
			message = 'Tag deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete tag');
		} finally {
			deletingTagId = undefined;
		}
	}

	async function deleteMediaProfile(id: string) {
		deletingMediaProfileId = id;
		clearNotice();

		try {
			await deleteMediaProfileRequest(id);
			if (mediaProfileForm.id === id) {
				mediaProfileForm = emptyMediaProfileForm();
			}
			mediaProfiles = mediaProfiles.filter((profile) => profile.id !== id);
			message = 'Profile deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete profile');
		} finally {
			deletingMediaProfileId = undefined;
		}
	}

	async function deleteCustomFormat(id: string) {
		deletingCustomFormatId = id;
		clearNotice();

		try {
			await deleteCustomFormatRequest(id);
			if (customFormatForm.id === id) {
				customFormatForm = emptyCustomFormatForm();
			}
			customFormats = customFormats.filter((format) => format.id !== id);
			message = 'Custom format deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete custom format');
		} finally {
			deletingCustomFormatId = undefined;
		}
	}

	async function searchLibraryMatch(kind: LibraryMediaKind, query: string) {
		return await searchMediaRequest({
			type: mediaTypeForLibraryKind(kind),
			query: query.trim()
		});
	}

	async function importLibraryScanRows(scan: LibraryScan, rows: LibraryScanImportRow[]) {
		clearNotice();

		try {
			const results: Awaited<ReturnType<typeof matchLibraryScanItemRequest>>[] = [];
			for (const row of rows) {
				results.push(await matchLibraryScanItemRequest(scan.id, row.item.id, row.request));
			}
			const importedMediaIds = results.map((result) => result.mediaItem.id);
			mediaItems = [
				...results.map((result) => result.mediaItem),
				...mediaItems.filter((item) => !importedMediaIds.includes(item.id))
			];
			libraryScansByFolder = {
				...libraryScansByFolder,
				[scan.folderId]: {
					...scan,
					manualCount: Math.max(0, scan.manualCount - results.length),
					items: scan.items.map(
						(item) => results.find((result) => result.item.id === item.id)?.item ?? item
					)
				}
			};
			message = `Imported ${results.length} media item${results.length === 1 ? '' : 's'}`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not import library items');
		}
	}

	async function testDownloadClientConfig(form: DownloadClientFormValue) {
		clearNotice();

		try {
			return await testDownloadClientConfigRequest(form);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test download client');
			throw error;
		}
	}

	async function testIndexer(id: string) {
		clearNotice();
		testingIndexerId = id;

		try {
			const result = await testIndexerRequest(id);
			indexerTests = { ...indexerTests, [id]: result };
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test indexer');
		} finally {
			testingIndexerId = undefined;
		}
	}

	async function testMetadataProvider(id: string) {
		clearNotice();
		testingMetadataProviderId = id;

		try {
			const result = await testMetadataProviderRequest(id);
			metadataProviderTests = { ...metadataProviderTests, [id]: result };
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test metadata provider');
		} finally {
			testingMetadataProviderId = undefined;
		}
	}

	async function refreshMetadataCache() {
		loadingMetadataCache = true;
		clearNotice();

		try {
			metadataCache = await getMetadataCacheRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load metadata cache');
		} finally {
			loadingMetadataCache = false;
		}
	}

	async function clearMetadataCache() {
		clearingMetadataCache = true;
		clearNotice();

		try {
			const deletedCount = await clearMetadataCacheRequest();
			metadataCache = await getMetadataCacheRequest();
			message = `Metadata cache reset: ${deletedCount} entries deleted`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not reset metadata cache');
		} finally {
			clearingMetadataCache = false;
		}
	}

	async function clearMetadataCachePattern(event: SubmitEvent) {
		event.preventDefault();
		const pattern = metadataCachePattern.trim();
		if (!pattern) {
			return;
		}
		clearingMetadataCache = true;
		clearNotice();

		try {
			const deletedCount = await clearMetadataCacheByPatternRequest(pattern);
			metadataCachePattern = '';
			metadataCache = await getMetadataCacheRequest();
			message = `Metadata cache reset: ${deletedCount} matching entries deleted`;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not reset matching metadata cache entries');
		} finally {
			clearingMetadataCache = false;
		}
	}

	function clearNotice() {
		errorMessage = '';
		message = '';
	}

	function omitResult<TValue>(results: Record<string, TValue | undefined>, id: string) {
		const { [id]: _removed, ...remaining } = results;
		return remaining;
	}

	function candidateKey(candidate: MediaSearchResult) {
		return `${candidate.type}:${candidate.title}:${candidate.year ?? ''}`;
	}

	function discoverResultKey(candidate: MediaSearchResult) {
		return `${candidate.type}:${candidate.externalProvider ?? ''}:${candidate.externalId ?? ''}:${candidate.title}:${candidate.year ?? ''}`;
	}

	function filterDiscoverSections(sections: MediaDiscoverSection[]) {
		return sections.map(filterDiscoverSection);
	}

	function filterDiscoverSection(section: MediaDiscoverSection): MediaDiscoverSection {
		return {
			...section,
			results: section.results.filter((result) => !isDiscoverBlacklisted(result))
		};
	}

	function relatedSectionFromDetail(
		detail: MediaMetadataDetails | undefined,
		kind: RelatedSectionKind
	): MediaDiscoverSection | undefined {
		if (!detail) {
			return undefined;
		}
		const results = kind === 'recommendations' ? detail.recommendations : detail.similar;
		return filterDiscoverSection({
			id: kind,
			title:
				kind === 'recommendations'
					? 'Recommendations'
					: detail.type === 'movie'
						? 'Similar Movies'
						: 'Similar Series',
			providerName: detail.externalProvider?.toUpperCase() ?? 'Metadata',
			mediaType: detail.type,
			results: results ?? []
		});
	}

	function isDiscoverBlacklisted(result: MediaSearchResult) {
		return discoverBlacklist.some((item) => sameDiscoverBlacklistItem(item, result));
	}

	function sameDiscoverBlacklistItem(
		item: DiscoverBlacklistItem,
		result: DiscoverBlacklistItem | MediaSearchResult
	) {
		const itemExternalKey = discoverExternalKey(item);
		const resultExternalKey = discoverExternalKey(result);
		if (itemExternalKey && resultExternalKey && itemExternalKey === resultExternalKey) {
			return true;
		}
		return discoverTitleKey(item) === discoverTitleKey(result);
	}

	function discoverExternalKey(item: DiscoverBlacklistItem | MediaSearchResult) {
		if (!item.externalProvider || !item.externalId) {
			return '';
		}
		return `${item.type}:${item.externalProvider}:${item.externalId}`.trim().toLowerCase();
	}

	function discoverTitleKey(item: DiscoverBlacklistItem | MediaSearchResult) {
		return `${item.type}:${item.title.trim().toLowerCase()}:${item.year ?? ''}`;
	}

	function selectHomeSection(section: HomeSection) {
		if (section === 'blacklist' && !isAdmin) {
			return;
		}
		activeView = 'home';
		activeHomeSection = section;
		void goto(resolve(`/${section}`));
	}

	function selectSettingsSection(section: string) {
		if (!isAdmin) {
			return;
		}
		activeSettingsSection = section as SettingsSection;
		void goto(resolve(settingsSectionHref(activeSettingsSection)));
	}

	function selectSystemSection(section: string) {
		if (!isAdmin) {
			return;
		}
		activeSystemSection = section as SystemSection;
		void goto(resolve(systemSectionHref(activeSystemSection)));
	}

	function selectSubmenuSection(section: string) {
		if (activeView === 'system') {
			selectSystemSection(section);
			return;
		}
		if (activePrimarySection === 'library') {
			selectHomeSection(section as HomeSection);
			return;
		}
		if (activePrimarySection === 'discover') {
			if (section === 'discover') {
				selectHomeSection('discover');
				return;
			}
			activeView = 'discover-section';
			activeHomeSection = 'discover';
			activeDiscoverSectionId = section;
			discoverSection = undefined;
			discoverSectionPage = 1;
			discoverSectionHasMore = true;
			void goto(resolve('/discover/[sectionId]', { sectionId: section }));
			void loadDiscoverSection();
			return;
		}
		selectSettingsSection(section);
	}

	function selectPrimarySection(section: string) {
		if (section === 'library') {
			selectHomeSection('movies');
			return;
		}
		if (section === 'settings') {
			if (!isAdmin) {
				return;
			}
			activeView = 'settings';
			activeSettingsSection = 'library';
			void goto(resolve('/settings/library'));
			return;
		}
		if (section === 'system') {
			if (!isAdmin) {
				return;
			}
			activeView = 'system';
			activeSystemSection = 'status';
			void goto(resolve('/system/status'));
			return;
		}
		selectHomeSection(section as HomeSection);
	}

	function editDownloadClient(client: DownloadClient) {
		downloadForm = downloadClientFormFromClient(client);
		activeSettingsSection = 'download-clients';
		void goto(resolve('/settings/download-clients'));
	}

	function editIndexer(indexer: Indexer) {
		indexerForm = indexerFormFromIndexer(indexer);
		activeSettingsSection = 'indexers';
		void goto(resolve('/settings/indexers'));
	}

	function editUser(user: ManagedUser) {
		userForm = userFormFromUser(user);
		activeSettingsSection = 'users';
		void goto(resolve('/settings/users'));
	}

	function editTag(tag: Tag) {
		tagForm = { id: tag.id, name: tag.name };
		activeSettingsSection = 'tags';
		void goto(resolve('/settings/tags'));
	}

	function editMediaProfile(profile: MediaProfile) {
		mediaProfileForm = mediaProfileFormFromProfile(profile);
		activeSettingsSection = 'profiles';
		void goto(resolve('/settings/profiles'));
	}

	function editCustomFormat(format: CustomFormat) {
		customFormatForm = customFormatFormFromFormat(format);
		activeSettingsSection = 'custom-formats';
		void goto(resolve('/settings/custom-formats'));
	}

	function errorMessageFrom(error: unknown, fallback: string) {
		return error instanceof Error ? error.message : fallback;
	}

	function emptyTagForm(): TagForm {
		return { name: '' };
	}

	return {
		get authenticated() {
			return authenticated;
		},
		get loading() {
			return loading;
		},
		get savingDownloadClient() {
			return savingDownloadClient;
		},
		get savingIndexer() {
			return savingIndexer;
		},
		get savingMetadataProviderId() {
			return savingMetadataProviderId;
		},
		get savingLibraryFolder() {
			return savingLibraryFolder;
		},
		get savingPathMapping() {
			return savingPathMapping;
		},
		get deletingPathMappingId() {
			return deletingPathMappingId;
		},
		get savingMediaProfile() {
			return savingMediaProfile;
		},
		get deletingMediaProfileId() {
			return deletingMediaProfileId;
		},
		get savingCustomFormat() {
			return savingCustomFormat;
		},
		get deletingCustomFormatId() {
			return deletingCustomFormatId;
		},
		get savingTag() {
			return savingTag;
		},
		get deletingTagId() {
			return deletingTagId;
		},
		get savingUser() {
			return savingUser;
		},
		get message() {
			return message;
		},
		get errorMessage() {
			return errorMessage;
		},
		clearNotice,
		get username() {
			return username;
		},
		set username(value: string) {
			username = value;
		},
		get password() {
			return password;
		},
		set password(value: string) {
			password = value;
		},
		get downloadClients() {
			return downloadClients;
		},
		get indexers() {
			return indexers;
		},
		get metadataProviders() {
			return metadataProviders;
		},
		get metadataCache() {
			return metadataCache;
		},
		get libraryFolders() {
			return libraryFolders;
		},
		get pathMappings() {
			return pathMappings;
		},
		get mediaProfiles() {
			return mediaProfiles;
		},
		get customFormats() {
			return customFormats;
		},
		get users() {
			return users;
		},
		get tags() {
			return tags;
		},
		get currentUser() {
			return currentUser;
		},
		get mediaItems() {
			return mediaItems;
		},
		get mediaRequests() {
			return mediaRequests;
		},
		get discoverSections() {
			return discoverSections;
		},
		get discoverSection() {
			return discoverSection;
		},
		get relatedMediaSection() {
			return relatedMediaSection;
		},
		get discoverBlacklist() {
			return discoverBlacklist;
		},
		get metadataDetail() {
			return metadataDetail;
		},
		get mediaPeopleDetail() {
			return mediaPeopleMetadataDetail;
		},
		get activePeopleSectionKind() {
			return activePeopleSectionKind;
		},
		get mediaCollection() {
			return mediaCollection;
		},
		get autocompleteGroups() {
			return autocompleteGroups;
		},
		get advancedSearchGroups() {
			return advancedSearchGroups;
		},
		get releaseResults() {
			return releaseResults;
		},
		get activities() {
			return activities;
		},
		get downloadForm() {
			return downloadForm;
		},
		set downloadForm(value: DownloadClientFormValue) {
			downloadForm = value;
		},
		get indexerForm() {
			return indexerForm;
		},
		set indexerForm(value: IndexerFormValue) {
			indexerForm = value;
		},
		get libraryFolderForm() {
			return libraryFolderForm;
		},
		set libraryFolderForm(value: LibraryFolderFormValue) {
			libraryFolderForm = value;
		},
		get pathMappingForm() {
			return pathMappingForm;
		},
		set pathMappingForm(value: PathMappingForm) {
			pathMappingForm = value;
		},
		get mediaProfileForm() {
			return mediaProfileForm;
		},
		set mediaProfileForm(value: MediaProfileFormValue) {
			mediaProfileForm = value;
		},
		get customFormatForm() {
			return customFormatForm;
		},
		set customFormatForm(value: CustomFormatFormValue) {
			customFormatForm = value;
		},
		get tagForm() {
			return tagForm;
		},
		set tagForm(value: TagForm) {
			tagForm = value;
		},
		get userForm() {
			return userForm;
		},
		set userForm(value: UserFormValue) {
			userForm = value;
		},
		get testingIndexerId() {
			return testingIndexerId;
		},
		get testingMetadataProviderId() {
			return testingMetadataProviderId;
		},
		get loadingMetadataCache() {
			return loadingMetadataCache;
		},
		get clearingMetadataCache() {
			return clearingMetadataCache;
		},
		get metadataCachePattern() {
			return metadataCachePattern;
		},
		set metadataCachePattern(value: string) {
			metadataCachePattern = value;
		},
		get loadingDiscover() {
			return loadingDiscover;
		},
		get loadingDiscoverSection() {
			return loadingDiscoverSection;
		},
		get loadingMoreDiscoverSection() {
			return loadingMoreDiscoverSection;
		},
		get discoverSectionHasMore() {
			return discoverSectionHasMore;
		},
		get loadingBlacklist() {
			return loadingBlacklist;
		},
		get loadingMetadataDetail() {
			return loadingMetadataDetail;
		},
		get loadingMediaCollection() {
			return loadingMediaCollection;
		},
		get loadingAutocomplete() {
			return loadingAutocomplete;
		},
		get searchingAdvanced() {
			return searchingAdvanced;
		},
		get addingKey() {
			return addingKey;
		},
		get blacklistingKey() {
			return blacklistingKey;
		},
		get removingBlacklistId() {
			return removingBlacklistId;
		},
		get savingMediaAction() {
			return savingMediaAction;
		},
		get activeMediaCandidate() {
			return activeMediaCandidate;
		},
		get mediaDeleteCandidate() {
			return mediaDeleteCandidate;
		},
		get approvingRequestId() {
			return approvingRequestId;
		},
		get searchingItemId() {
			return searchingItemId;
		},
		get scanningMediaItemId() {
			return scanningMediaItemId;
		},
		get grabbingKey() {
			return grabbingKey;
		},
		get deletingMediaItemId() {
			return deletingMediaItemId;
		},
		get cancellingActivityId() {
			return cancellingActivityId;
		},
		get deletingActivityId() {
			return deletingActivityId;
		},
		get loadingActivity() {
			return loadingActivity;
		},
		get scanningLibraryFolderId() {
			return scanningLibraryFolderId;
		},
		get libraryScansByFolder() {
			return libraryScansByFolder;
		},
		get openLibraryFolderId() {
			return openLibraryFolderId;
		},
		get indexerTests() {
			return indexerTests;
		},
		get metadataProviderTests() {
			return metadataProviderTests;
		},
		get activeView() {
			return activeView;
		},
		get activeHomeSection() {
			return activeHomeSection;
		},
		get activeSettingsSection() {
			return activeSettingsSection;
		},
		get activeSystemSection() {
			return activeSystemSection;
		},
		get selectedMediaItemId() {
			return selectedMediaItemId;
		},
		get selectedRequestId() {
			return selectedRequestId;
		},
		get searchQuery() {
			return searchQuery;
		},
		set searchQuery(value: string) {
			searchQuery = value;
		},
		get isAdmin() {
			return isAdmin;
		},
		get activePrimarySection() {
			return activePrimarySection;
		},
		get activeSubmenuSection() {
			return activeSubmenuSection;
		},
		get primaryItems() {
			return primaryItems;
		},
		initialise,
		login,
		logout,
		showProfile,
		connectEvents,
		disconnectEvents,
		autocompleteMedia,
		selectAutocompleteResult,
		openAdvancedSearch,
		advancedSearch,
		addMedia,
		loadMoreDiscoverSection,
		blacklistDiscoverMedia,
		removeDiscoverBlacklistItem,
		closeMediaAction,
		confirmMediaAction,
		closeMediaDelete,
		confirmMediaDelete,
		approveMediaRequest,
		findReleases,
		autoSearchMedia,
		rescanMediaFiles,
		deleteMediaFile,
		deleteMediaItem,
		grabRelease,
		cancelActivity,
		deleteActivity,
		loadDownloadActivity,
		saveDownloadClient,
		testDownloadClientConfig,
		saveIndexer,
		saveMetadataProvider,
		refreshMetadataCache,
		clearMetadataCache,
		clearMetadataCachePattern,
		saveLibraryFolder,
		savePathMapping,
		saveMediaProfile,
		saveCustomFormat,
		importCustomFormat,
		saveTag,
		saveUser,
		editDownloadClient,
		editIndexer,
		editMediaProfile,
		editCustomFormat,
		editTag,
		editUser,
		deleteDownloadClient,
		deleteIndexer,
		deleteLibraryFolder,
		scanLibraryFolder,
		deletePathMapping,
		deleteMediaProfile,
		deleteCustomFormat,
		deleteTag,
		deleteUser,
		testIndexer,
		testMetadataProvider,
		searchLibraryMatch,
		importLibraryScanRows,
		selectPrimarySection,
		selectSettingsSection,
		selectSystemSection,
		selectSubmenuSection,
		cancelDownloadClient: () => (downloadForm = emptyDownloadClientForm()),
		cancelIndexer: () => (indexerForm = emptyIndexerForm()),
		cancelMediaProfile: () => (mediaProfileForm = emptyMediaProfileForm()),
		cancelCustomFormat: () => (customFormatForm = emptyCustomFormatForm()),
		cancelTag: () => (tagForm = emptyTagForm()),
		cancelUser: () => (userForm = emptyUserForm())
	};
}
