<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';

	import AppNav from '$lib/components/app/AppNav.svelte';
	import AdvancedSearchArea from '$lib/components/app/AdvancedSearchArea.svelte';
	import HomeArea from '$lib/components/app/HomeArea.svelte';
	import MediaActionModal from '$lib/components/app/MediaActionModal.svelte';
	import MetadataDetailArea from '$lib/components/app/MetadataDetailArea.svelte';
	import SettingsArea from '$lib/components/app/SettingsArea.svelte';
	import SidebarMenu from '$lib/components/app/SidebarMenu.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import {
		approveMediaRequest as approveMediaRequestRequest,
		clearMetadataCache as clearMetadataCacheRequest,
		clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
		createMediaItem as createMediaItemRequest,
		createMediaRequest as createMediaRequestRequest,
		currentSession as currentSessionRequest,
		deleteDownloadClient as deleteDownloadClientRequest,
		deleteIndexer as deleteIndexerRequest,
		deleteLibraryFolder as deleteLibraryFolderRequest,
		deleteMediaItem as deleteMediaItemRequest,
		deleteTag as deleteTagRequest,
		deleteUser as deleteUserRequest,
		enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
		advancedSearchMedia as advancedSearchMediaRequest,
		autocompleteMedia as autocompleteMediaRequest,
		getLibraryScan as getLibraryScanRequest,
		grabMediaRelease as grabMediaReleaseRequest,
		getMetadataCache as getMetadataCacheRequest,
		getMediaMetadataDetails as getMediaMetadataDetailsRequest,
		listDownloadActivity as listDownloadActivityRequest,
		listMediaRequests as listMediaRequestsRequest,
		loadMediaDiscoverSections as loadMediaDiscoverSectionsRequest,
		listMediaItems as listMediaItemsRequest,
		loadSettings as loadSettingsRequest,
		login as loginRequest,
		logout as logoutRequest,
		matchLibraryScanItem as matchLibraryScanItemRequest,
		mediaTypeForLibraryKind,
		emptyMetadataCache,
		saveDownloadClient as saveDownloadClientRequest,
		saveIndexer as saveIndexerRequest,
		saveLibraryFolder as saveLibraryFolderRequest,
		saveMetadataProvider as saveMetadataProviderRequest,
		saveTag as saveTagRequest,
		saveUser as saveUserRequest,
		searchMedia as searchMediaRequest,
		searchMediaReleases as searchMediaReleasesRequest,
		testDownloadClient as testDownloadClientRequest,
		testIndexer as testIndexerRequest,
		testMetadataProvider as testMetadataProviderRequest
	} from '$lib/settings/api';
	import {
		downloadClientFormFromClient,
		emptyDownloadClientForm,
		emptyIndexerForm,
		emptyLibraryFolderForm,
		emptyUserForm,
		indexerFormFromIndexer,
		userFormFromUser
	} from '$lib/settings/forms';
	import { qualityProfiles } from '$lib/settings/qualityProfiles';
	import '$lib/settings/styles.css';
	import type {
		AppView,
		DownloadActivity,
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
		ReleaseCandidate,
		ReleaseSearchResults,
		SettingsSection,
		Tag,
		TagForm,
		UserForm as UserFormValue,
		UserSummary
	} from '$lib/settings/types';

	interface Props {
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

	let {
		initialView = 'home',
		initialHomeSection = 'discover',
		initialSettingsSection = 'library',
		initialSelectedMediaItemId,
		initialSelectedRequestId,
		initialLibraryScanId,
		initialAdvancedQuery = '',
		initialMetadataProvider,
		initialMetadataType,
		initialMetadataExternalId
	}: Props = $props();
	const routeDefaults = (() => ({
		view: initialView,
		homeSection: initialHomeSection,
		settingsSection: initialSettingsSection,
		selectedMediaItemId: initialSelectedMediaItemId,
		selectedRequestId: initialSelectedRequestId,
		libraryScanId: initialLibraryScanId
	}))();

	let authenticated = $state(false);
	let loading = $state(true);
	let savingDownloadClient = $state(false);
	let savingIndexer = $state(false);
	let savingMetadataProviderId = $state<string | undefined>();
	let savingLibraryFolder = $state(false);
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
	let tagForm = $state<TagForm>(emptyTagForm());
	let userForm = $state<UserFormValue>(emptyUserForm());
	let testingDownloadClientId = $state<string | undefined>();
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
	let loadingActivity = $state(false);
	let loadingLibraryScan = $state(false);
	let downloadClientTests = $state<IntegrationTestResults>({});
	let indexerTests = $state<IntegrationTestResults>({});
	let metadataProviderTests = $state<IntegrationTestResults>({});
	let activeView = $state<AppView>(routeDefaults.view);
	let activeHomeSection = $state<HomeSection>(routeDefaults.homeSection);
	let activeSettingsSection = $state<SettingsSection>(routeDefaults.settingsSection);
	let selectedMediaItemId = $state<string | undefined>(routeDefaults.selectedMediaItemId);
	let selectedRequestId = $state<string | undefined>(routeDefaults.selectedRequestId);
	let activeLibraryScanId = $state<string | undefined>(routeDefaults.libraryScanId);
	let activeLibraryScan = $state<LibraryScan | undefined>();
	let searchQuery = $state('');
	let isAdmin = $derived(currentUser?.role === 'admin');
	let activePrimarySection = $derived(activeView === 'settings' ? 'settings' : activeHomeSection);
	type PrimaryItem = {
		value: HomeSection | 'settings';
		label: string;
		icon: 'discover' | 'movies' | 'series' | 'activity' | 'settings';
		href: '/discover' | '/requests' | '/movies' | '/series' | '/activity' | '/settings/library';
		children?: readonly {
			value: SettingsSection;
			label: string;
			href:
				| '/settings/library'
				| '/settings/download-clients'
				| '/settings/indexers'
				| '/settings/metadata'
				| '/settings/tags'
				| '/settings/users';
		}[];
	};
	const settingsItems = [
		{ value: 'library', label: 'Library', href: '/settings/library' },
		{ value: 'download-clients', label: 'Download clients', href: '/settings/download-clients' },
		{ value: 'indexers', label: 'Indexers', href: '/settings/indexers' },
		{ value: 'metadata', label: 'Metadata', href: '/settings/metadata' },
		{ value: 'tags', label: 'Tags', href: '/settings/tags' },
		{ value: 'users', label: 'Users', href: '/settings/users' }
	] satisfies PrimaryItem['children'];
	const basePrimaryItems = [
		{ value: 'discover', label: 'Discover', icon: 'discover', href: '/discover' },
		{ value: 'requests', label: 'Requests', icon: 'activity', href: '/requests' },
		{ value: 'movies', label: 'Movies', icon: 'movies', href: '/movies' },
		{ value: 'series', label: 'Series', icon: 'series', href: '/series' },
		{ value: 'activity', label: 'Activity', icon: 'activity', href: '/activity' }
	] satisfies PrimaryItem[];
	const settingsPrimaryItem = {
		value: 'settings',
		label: 'Settings',
		icon: 'settings',
		href: '/settings/library',
		children: settingsItems
	} satisfies PrimaryItem;
	let primaryItems = $derived(
		isAdmin ? [...basePrimaryItems, settingsPrimaryItem] : basePrimaryItems
	);

	onMount(() => {
		searchQuery = initialAdvancedQuery;
		void initialise();
	});

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
			activeLibraryScan = undefined;
			activeLibraryScanId = undefined;
			downloadForm = emptyDownloadClientForm();
			indexerForm = emptyIndexerForm();
			libraryFolderForm = emptyLibraryFolderForm();
			tagForm = emptyTagForm();
			userForm = emptyUserForm();
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
		if (!initialMetadataProvider || !initialMetadataType || !initialMetadataExternalId) {
			return;
		}
		loadingMetadataDetail = true;
		try {
			metadataDetail = await getMediaMetadataDetailsRequest(
				initialMetadataProvider as MetadataProviderType,
				initialMetadataType as MediaType,
				initialMetadataExternalId
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
			if (result.type === 'movie') {
				void goto(resolve('/movies/[id]', { id: result.id }));
			} else {
				void goto(resolve('/series/[id]', { id: result.id }));
			}
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
		if (savingMediaAction) {
			return;
		}
		activeMediaCandidate = undefined;
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

	async function deleteDownloadClient(id: string) {
		clearNotice();

		try {
			await deleteDownloadClientRequest(id);
			if (downloadForm.id === id) {
				downloadForm = emptyDownloadClientForm();
			}
			downloadClientTests = omitResult(downloadClientTests, id);
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

	async function testDownloadClient(id: string) {
		clearNotice();
		testingDownloadClientId = id;

		try {
			const result = await testDownloadClientRequest(id);
			downloadClientTests = { ...downloadClientTests, [id]: result };
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test download client');
		} finally {
			testingDownloadClientId = undefined;
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

	function selectSettingsSection(section: SettingsSection) {
		if (!isAdmin) {
			return;
		}
		activeSettingsSection = section;
		void goto(resolve(`/settings/${section}`));
	}

	function selectPrimarySection(section: HomeSection | 'settings') {
		if (section === 'settings') {
			if (!isAdmin) {
				return;
			}
			activeView = 'settings';
			activeSettingsSection = 'library';
			void goto(resolve('/settings/library'));
			return;
		}
		selectHomeSection(section);
	}

	function errorMessageFrom(error: unknown, fallback: string) {
		return error instanceof Error ? error.message : fallback;
	}

	function emptyTagForm(): TagForm {
		return { name: '' };
	}
</script>

<svelte:head>
	<title>mema</title>
	<meta
		name="description"
		content="Self-hosted media manager for movies, series, books, music, clients, and indexers"
	/>
</svelte:head>

{#if loading}
	<main class="shell">
		<section class="panel">
			<p class="muted">Loading app</p>
		</section>
	</main>
{:else if !authenticated}
	<main class="shell login-shell">
		<div class="login-brand">
			<span class="brand-mark large" aria-hidden="true">M</span>
			<h1>mema</h1>
		</div>
		<NoticeStack {message} {errorMessage} />
		<AuthPanel bind:username bind:password onLogin={login} />
	</main>
{:else}
	<div class="app-frame">
		<SidebarMenu
			title="mema"
			items={primaryItems}
			active={activePrimarySection}
			activeSubmenu={activeSettingsSection}
			onSelect={(section) => selectPrimarySection(section as HomeSection | 'settings')}
			onSubmenuSelect={(section) => selectSettingsSection(section as SettingsSection)}
		/>
		<div class="app-main">
			<AppNav
				bind:searchQuery
				groups={autocompleteGroups}
				loading={loadingAutocomplete}
				onSearch={autocompleteMedia}
				onSelect={selectAutocompleteResult}
				onAdvancedSearch={openAdvancedSearch}
				onProfile={showProfile}
				onLogout={logout}
			/>
			<main class="app-content">
				<NoticeStack {message} {errorMessage} />
				{#if activeView === 'settings' && isAdmin}
					<SettingsArea
						bind:downloadForm
						bind:indexerForm
						bind:libraryFolderForm
						bind:tagForm
						bind:userForm
						activeSection={activeSettingsSection}
						{downloadClients}
						{indexers}
						{metadataProviders}
						{metadataCache}
						{libraryFolders}
						{users}
						{tags}
						{currentUser}
						{activeLibraryScan}
						{savingDownloadClient}
						{savingIndexer}
						{savingMetadataProviderId}
						{loadingMetadataCache}
						{clearingMetadataCache}
						{savingLibraryFolder}
						{savingTag}
						{deletingTagId}
						{savingUser}
						bind:metadataCachePattern
						{loadingLibraryScan}
						{testingDownloadClientId}
						{testingIndexerId}
						{testingMetadataProviderId}
						{downloadClientTests}
						{indexerTests}
						{metadataProviderTests}
						onSaveDownloadClient={saveDownloadClient}
						onSaveIndexer={saveIndexer}
						onSaveMetadataProvider={saveMetadataProvider}
						onRefreshMetadataCache={refreshMetadataCache}
						onClearMetadataCache={clearMetadataCache}
						onClearMetadataCachePattern={clearMetadataCachePattern}
						onSaveLibraryFolder={saveLibraryFolder}
						onSaveTag={saveTag}
						onSaveUser={saveUser}
						onCancelDownloadClient={() => (downloadForm = emptyDownloadClientForm())}
						onCancelIndexer={() => (indexerForm = emptyIndexerForm())}
						onCancelTag={() => (tagForm = emptyTagForm())}
						onCancelUser={() => (userForm = emptyUserForm())}
						onEditDownloadClient={(client) => {
							downloadForm = downloadClientFormFromClient(client);
							activeSettingsSection = 'download-clients';
							void goto(resolve('/settings/download-clients'));
						}}
						onEditIndexer={(indexer) => {
							indexerForm = indexerFormFromIndexer(indexer);
							activeSettingsSection = 'indexers';
							void goto(resolve('/settings/indexers'));
						}}
						onEditUser={(user) => {
							userForm = userFormFromUser(user);
							activeSettingsSection = 'users';
							void goto(resolve('/settings/users'));
						}}
						onEditTag={(tag) => {
							tagForm = { id: tag.id, name: tag.name };
							activeSettingsSection = 'tags';
							void goto(resolve('/settings/tags'));
						}}
						onDeleteDownloadClient={deleteDownloadClient}
						onDeleteIndexer={deleteIndexer}
						onDeleteLibraryFolder={deleteLibraryFolder}
						onDeleteTag={deleteTag}
						onDeleteUser={deleteUser}
						onTestDownloadClient={testDownloadClient}
						onTestIndexer={testIndexer}
						onTestMetadataProvider={testMetadataProvider}
						onSearchLibraryMatch={searchLibraryMatch}
						onMatchLibraryScanItem={matchLibraryScanItem}
					/>
				{:else if activeView === 'advanced-search'}
					<AdvancedSearchArea
						initialQuery={initialAdvancedQuery}
						{metadataProviders}
						groups={advancedSearchGroups}
						searching={searchingAdvanced}
						{addingKey}
						actionLabel={isAdmin ? 'Add' : 'Request'}
						onSearch={advancedSearch}
						onAdd={addMedia}
					/>
				{:else if activeView === 'metadata-detail'}
					<MetadataDetailArea
						detail={metadataDetail}
						loading={loadingMetadataDetail}
						{addingKey}
						actionLabel={isAdmin ? 'Add' : 'Request'}
						onAdd={addMedia}
					/>
				{:else}
					<HomeArea
						activeSection={activeHomeSection}
						{selectedMediaItemId}
						{selectedRequestId}
						{mediaItems}
						{mediaRequests}
						{discoverSections}
						{libraryFolders}
						{qualityProfiles}
						{releaseResults}
						{activities}
						{loadingDiscover}
						{addingKey}
						{approvingRequestId}
						{searchingItemId}
						{grabbingKey}
						{deletingMediaItemId}
						{loadingActivity}
						canManage={isAdmin}
						onAddMedia={addMedia}
						onApproveMediaRequest={approveMediaRequest}
						onFindReleases={findReleases}
						onDeleteMedia={deleteMediaItem}
						onGrabRelease={grabRelease}
						onRefreshActivity={loadDownloadActivity}
					/>
				{/if}
			</main>
		</div>
	</div>
	{#if activeMediaCandidate}
		<MediaActionModal
			candidate={activeMediaCandidate}
			{isAdmin}
			{libraryFolders}
			{qualityProfiles}
			{tags}
			saving={savingMediaAction}
			onClose={closeMediaAction}
			onConfirm={confirmMediaAction}
		/>
	{/if}
{/if}
