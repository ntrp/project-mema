/* global $derived, $state, EventSource, MessageEvent */
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';

import {
	basePrimaryItems,
	settingsPrimaryItem,
	settingsSectionHref
} from '$lib/components/app/appNavigation';
import {
	advancedSearchMedia as advancedSearchMediaRequest,
	approveMediaRequest as approveMediaRequestRequest,
	autocompleteMedia as autocompleteMediaRequest,
	cancelDownloadActivity as cancelDownloadActivityRequest,
	clearMetadataCache as clearMetadataCacheRequest,
	clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
	createMediaItem as createMediaItemRequest,
	createMediaRequest as createMediaRequestRequest,
	currentSession as currentSessionRequest,
	deleteDownloadClient as deleteDownloadClientRequest,
	deleteIndexer as deleteIndexerRequest,
	deleteLibraryFolder as deleteLibraryFolderRequest,
	deleteMediaItem as deleteMediaItemRequest,
	deleteMediaProfile as deleteMediaProfileRequest,
	deletePathMapping as deletePathMappingRequest,
	deleteTag as deleteTagRequest,
	deleteUser as deleteUserRequest,
	emptyMetadataCache,
	enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
	getLibraryScan as getLibraryScanRequest,
	getMediaMetadataDetails as getMediaMetadataDetailsRequest,
	getMetadataCache as getMetadataCacheRequest,
	grabMediaRelease as grabMediaReleaseRequest,
	listDownloadActivity as listDownloadActivityRequest,
	listMediaItems as listMediaItemsRequest,
	listMediaRequests as listMediaRequestsRequest,
	loadMediaDiscoverSections as loadMediaDiscoverSectionsRequest,
	loadSettings as loadSettingsRequest,
	login as loginRequest,
	logout as logoutRequest,
	matchLibraryScanItem as matchLibraryScanItemRequest,
	mediaTypeForLibraryKind,
	saveDownloadClient as saveDownloadClientRequest,
	saveIndexer as saveIndexerRequest,
	saveLibraryFolder as saveLibraryFolderRequest,
	saveMediaProfile as saveMediaProfileRequest,
	saveMetadataProvider as saveMetadataProviderRequest,
	savePathMapping as savePathMappingRequest,
	saveTag as saveTagRequest,
	saveUser as saveUserRequest,
	searchMedia as searchMediaRequest,
	searchMediaReleases as searchMediaReleasesRequest,
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest
} from '$lib/settings/api';
import {
	downloadClientFormFromClient,
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
import type {
	AppView,
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
	LibraryScanItem,
	LibraryScanItemMatchRequest,
	ManagedUser,
	MediaAdvancedSearchRequest,
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
	initialSelectedMediaItemId?: string;
	initialSelectedRequestId?: string;
	initialLibraryScanId?: string;
	initialAdvancedQuery?: string;
	initialMetadataProvider?: string;
	initialMetadataType?: string;
	initialMetadataExternalId?: string;
}

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
	let users = $state<ManagedUser[]>([]);
	let tags = $state<Tag[]>([]);
	let currentUser = $state<UserSummary | undefined>();
	let mediaItems = $state<MediaItem[]>([]);
	let mediaRequests = $state<MediaRequest[]>([]);
	let discoverSections = $state<MediaDiscoverSection[]>([]);
	let metadataDetail = $state<MediaMetadataDetails | undefined>();
	let autocompleteGroups = $state<MediaSearchGroup[]>([]);
	let advancedSearchGroups = $state<MediaSearchGroup[]>([]);
	let releaseResults = $state<ReleaseSearchResults>({});
	let activities = $state<DownloadActivity[]>([]);
	let downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	let indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	let libraryFolderForm = $state<LibraryFolderFormValue>(emptyLibraryFolderForm());
	let pathMappingForm = $state<PathMappingForm>(emptyPathMappingForm());
	let mediaProfileForm = $state<MediaProfileFormValue>(emptyMediaProfileForm());
	let tagForm = $state<TagForm>(emptyTagForm());
	let userForm = $state<UserFormValue>(emptyUserForm());
	let testingIndexerId = $state<string | undefined>();
	let testingMetadataProviderId = $state<string | undefined>();
	let loadingMetadataCache = $state(false);
	let clearingMetadataCache = $state(false);
	let metadataCachePattern = $state('');
	let loadingDiscover = $state(false);
	let loadingMetadataDetail = $state(false);
	let loadingAutocomplete = $state(false);
	let searchingAdvanced = $state(false);
	let addingKey = $state<string | undefined>();
	let savingMediaAction = $state(false);
	let activeMediaCandidate = $state<MediaSearchResult | undefined>();
	let approvingRequestId = $state<string | undefined>();
	let searchingItemId = $state<string | undefined>();
	let grabbingKey = $state<string | undefined>();
	let deletingMediaItemId = $state<string | undefined>();
	let cancellingActivityId = $state<string | undefined>();
	let loadingActivity = $state(false);
	let loadingLibraryScan = $state(false);
	let indexerTests = $state<IntegrationTestResults>({});
	let metadataProviderTests = $state<IntegrationTestResults>({});
	let activeView = $state<AppView>(options.initialView ?? 'home');
	let activeHomeSection = $state<HomeSection>(options.initialHomeSection ?? 'discover');
	let activeSettingsSection = $state<SettingsSection>(options.initialSettingsSection ?? 'library');
	let selectedMediaItemId = $state<string | undefined>(options.initialSelectedMediaItemId);
	let selectedRequestId = $state<string | undefined>(options.initialSelectedRequestId);
	let activeLibraryScanId = $state<string | undefined>(options.initialLibraryScanId);
	let activeLibraryScan = $state<LibraryScan | undefined>();
	let searchQuery = $state(options.initialAdvancedQuery ?? '');
	let isAdmin = $derived(currentUser?.role === 'admin');
	let activePrimarySection = $derived(activeView === 'settings' ? 'settings' : activeHomeSection);
	let primaryItems = $derived(
		isAdmin ? [...basePrimaryItems, settingsPrimaryItem] : basePrimaryItems
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
			} else if (activeView === 'settings') {
				activeView = 'home';
				activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (activeView === 'metadata-detail') {
				await loadMetadataDetail();
			}
		}

		loading = false;
	}

	function connectEvents() {
		const source = new EventSource('/api/events', { withCredentials: true });
		source.addEventListener('activity.download.updated', (event) => {
			const activity = parseEventData<DownloadActivity>(event);
			if (!activity) {
				return;
			}
			upsertActivity(activity);
			updateMediaStatusFromActivity(activity);
		});
		source.onerror = () => {
			if (!authenticated) {
				source.close();
			}
		};
		return () => source.close();
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
			} else if (activeView === 'settings') {
				activeView = 'home';
				activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (activeView === 'metadata-detail') {
				await loadMetadataDetail();
			}
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
			authenticated = false;
			currentUser = undefined;
			activeView = 'home';
			activeHomeSection = 'discover';
			downloadClients = [];
			indexers = [];
			metadataProviders = [];
			mediaProfiles = [];
			users = [];
			tags = [];
			mediaItems = [];
			mediaRequests = [];
			discoverSections = [];
			metadataDetail = undefined;
			autocompleteGroups = [];
			advancedSearchGroups = [];
			releaseResults = {};
			activities = [];
			libraryFolders = [];
			pathMappings = [];
			activeLibraryScan = undefined;
			activeLibraryScanId = undefined;
			downloadForm = emptyDownloadClientForm();
			indexerForm = emptyIndexerForm();
			libraryFolderForm = emptyLibraryFolderForm();
			pathMappingForm = emptyPathMappingForm();
			mediaProfileForm = emptyMediaProfileForm();
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
			users = settings.users;
			tags = settings.tags;
			if (activeLibraryScanId) {
				await loadLibraryScan(activeLibraryScanId);
			}
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

	async function confirmMediaAction(
		qualityProfileId?: string,
		libraryFolderId?: string,
		selectedTags: string[] = []
	) {
		const candidate = activeMediaCandidate;
		if (!candidate) {
			return;
		}
		addingKey = candidateKey(candidate);
		savingMediaAction = true;
		clearNotice();

		try {
			if (isAdmin) {
				if (!qualityProfileId || !libraryFolderId) {
					throw new Error('Quality profile and library folder are required');
				}
				const item = await createMediaItemRequest({
					title: candidate.title,
					type: candidate.type,
					year: candidate.year,
					monitored: true,
					externalProvider: candidate.externalProvider,
					externalId: candidate.externalId,
					overview: candidate.overview,
					posterPath: candidate.posterPath,
					qualityProfileId,
					libraryFolderId,
					tags: selectedTags
				});
				mediaItems = [item, ...mediaItems.filter((mediaItem) => mediaItem.id !== item.id)];
				await loadSettings();
				message = 'Media item added to monitored';
				activeHomeSection = candidate.type === 'movie' ? 'movies' : 'series';
				activeMediaCandidate = undefined;
				void goto(resolve(candidate.type === 'movie' ? '/movies' : '/series'));
				return;
			}

			const request = await createMediaRequestRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath,
				tags: selectedTags
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

	async function deleteMediaItem(item: MediaItem) {
		deletingMediaItemId = item.id;
		clearNotice();

		try {
			await deleteMediaItemRequest(item.id);
			mediaItems = mediaItems.filter((mediaItem) => mediaItem.id !== item.id);
			releaseResults = omitResult(releaseResults, item.id);
			activities = activities.filter((activity) => activity.mediaItemId !== item.id);
			message = 'Media item removed';
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
			activeLibraryScan = result.scan;
			activeLibraryScanId = result.scan.id;
			message = `Library scan completed: ${result.scan.autoMatchedCount} auto-added, ${result.scan.manualCount} pending`;
			await loadMediaItems();
			void goto(resolve(`/settings/library/scans/${result.scan.id}`));
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
			if (activeLibraryScan?.folderId === id) {
				activeLibraryScan = undefined;
				activeLibraryScanId = undefined;
				void goto(resolve('/settings/library'));
			}
			message = 'Library folder deleted';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete library folder');
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

	async function loadLibraryScan(id: string) {
		loadingLibraryScan = true;
		clearNotice();

		try {
			activeLibraryScan = await getLibraryScanRequest(id);
			activeLibraryScanId = id;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load library scan');
		} finally {
			loadingLibraryScan = false;
		}
	}

	async function searchLibraryMatch(kind: LibraryMediaKind, query: string) {
		return await searchMediaRequest({
			type: mediaTypeForLibraryKind(kind),
			query: query.trim()
		});
	}

	async function matchLibraryScanItem(item: LibraryScanItem, request: LibraryScanItemMatchRequest) {
		if (!activeLibraryScanId) {
			return;
		}
		clearNotice();

		try {
			const result = await matchLibraryScanItemRequest(activeLibraryScanId, item.id, request);
			activeLibraryScan = {
				...activeLibraryScan!,
				manualCount: Math.max(0, activeLibraryScan!.manualCount - 1),
				items: activeLibraryScan!.items.map((scanItem) =>
					scanItem.id === item.id ? result.item : scanItem
				)
			};
			mediaItems = [
				result.mediaItem,
				...mediaItems.filter((mediaItem) => mediaItem.id !== result.mediaItem.id)
			];
			message = 'Library item added to monitored';
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not match library item');
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

	function selectHomeSection(section: HomeSection) {
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

	function selectPrimarySection(section: string) {
		if (section === 'settings') {
			if (!isAdmin) {
				return;
			}
			activeView = 'settings';
			activeSettingsSection = 'library';
			void goto(resolve('/settings/library'));
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
		get metadataDetail() {
			return metadataDetail;
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
		get loadingMetadataDetail() {
			return loadingMetadataDetail;
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
		get savingMediaAction() {
			return savingMediaAction;
		},
		get activeMediaCandidate() {
			return activeMediaCandidate;
		},
		get approvingRequestId() {
			return approvingRequestId;
		},
		get searchingItemId() {
			return searchingItemId;
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
		get loadingActivity() {
			return loadingActivity;
		},
		get loadingLibraryScan() {
			return loadingLibraryScan;
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
		get selectedMediaItemId() {
			return selectedMediaItemId;
		},
		get selectedRequestId() {
			return selectedRequestId;
		},
		get activeLibraryScan() {
			return activeLibraryScan;
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
		get primaryItems() {
			return primaryItems;
		},
		initialise,
		login,
		logout,
		showProfile,
		connectEvents,
		autocompleteMedia,
		selectAutocompleteResult,
		openAdvancedSearch,
		advancedSearch,
		addMedia,
		closeMediaAction,
		confirmMediaAction,
		approveMediaRequest,
		findReleases,
		deleteMediaItem,
		grabRelease,
		cancelActivity,
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
		saveTag,
		saveUser,
		editDownloadClient,
		editIndexer,
		editMediaProfile,
		editTag,
		editUser,
		deleteDownloadClient,
		deleteIndexer,
		deleteLibraryFolder,
		deletePathMapping,
		deleteMediaProfile,
		deleteTag,
		deleteUser,
		testIndexer,
		testMetadataProvider,
		searchLibraryMatch,
		matchLibraryScanItem,
		selectPrimarySection,
		selectSettingsSection,
		cancelDownloadClient: () => (downloadForm = emptyDownloadClientForm()),
		cancelIndexer: () => (indexerForm = emptyIndexerForm()),
		cancelMediaProfile: () => (mediaProfileForm = emptyMediaProfileForm()),
		cancelTag: () => (tagForm = emptyTagForm()),
		cancelUser: () => (userForm = emptyUserForm())
	};
}
