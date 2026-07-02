import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	currentSession as currentSessionRequest,
	login as loginRequest,
	logout as logoutRequest
} from '$lib/settings/api';
import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLibraryFolderForm,
	emptyMediaProfileForm,
	emptyPathMappingForm,
	emptyUserForm
} from '$lib/settings/forms';
import { connectAppEvents, disconnectAppEvents, type EventConnectionDeps } from './eventConnection';
import { emptyTagForm, errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SessionDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
	loadDiscoverBlacklist: () => Promise<void>;
	loadLibrary: () => Promise<void>;
	loadDiscoverSections: () => Promise<void>;
	loadMetadataDetail: () => Promise<void>;
	loadMediaCollection: () => Promise<void>;
	loadDiscoverSection: () => Promise<void>;
	events: EventConnectionDeps;
}

export function createSessionActions(state: AppShellState, deps: SessionDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	const loadDiscoverBlacklist = deps.loadDiscoverBlacklist;
	const loadLibrary = deps.loadLibrary;
	const loadDiscoverSections = deps.loadDiscoverSections;
	const loadMetadataDetail = deps.loadMetadataDetail;
	const loadMediaCollection = deps.loadMediaCollection;
	const loadDiscoverSection = deps.loadDiscoverSection;
	const eventDeps = deps.events;
	async function initialise() {
		state.loading = true;
		state.errorMessage = '';

		const session = await currentSessionRequest();
		state.authenticated = Boolean(session?.authenticated);
		state.currentUser = session?.user;
		if (state.authenticated) {
			if (state.currentUser?.role === 'admin') {
				await loadSettings();
				await loadDiscoverBlacklist();
			} else if (state.activeView === 'settings' || state.activeView === 'system') {
				state.activeView = 'home';
				state.activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			} else if (state.activeHomeSection === 'blacklist') {
				state.activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (
				state.activeView === 'metadata-detail' ||
				state.activeView === 'media-people' ||
				state.activeView === 'related-section'
			) {
				await loadMetadataDetail();
			} else if (state.activeView === 'media-collection') {
				await loadMediaCollection();
			} else if (state.activeView === 'discover-section') {
				await loadDiscoverSection();
			}
			connectEvents();
		}

		state.loading = false;
	}

	function connectEvents() {
		connectAppEvents(state, eventDeps);
	}

	function disconnectEvents() {
		disconnectAppEvents(state);
	}

	async function login(event: SubmitEvent) {
		event.preventDefault();
		clearNotice();

		try {
			const session = await loginRequest(state.username, state.password);
			state.authenticated = true;
			state.currentUser = session.user;
			if (state.currentUser?.role === 'admin') {
				await loadSettings();
			} else if (state.activeView === 'settings' || state.activeView === 'system') {
				state.activeView = 'home';
				state.activeHomeSection = 'discover';
				void goto(resolve('/discover'));
			}
			await loadLibrary();
			await loadDiscoverSections();
			if (
				state.activeView === 'metadata-detail' ||
				state.activeView === 'media-people' ||
				state.activeView === 'related-section'
			) {
				await loadMetadataDetail();
			} else if (state.activeView === 'media-collection') {
				await loadMediaCollection();
			} else if (state.activeView === 'discover-section') {
				await loadDiscoverSection();
			}
			connectEvents();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Login failed');
		}
	}

	async function logout() {
		clearNotice();

		try {
			await logoutRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not log out');
		} finally {
			disconnectEvents();
			state.authenticated = false;
			state.currentUser = undefined;
			state.activeView = 'home';
			state.activeHomeSection = 'discover';
			state.downloadClients = [];
			state.indexers = [];
			state.metadataProviders = [];
			state.mediaProfiles = [];
			state.customFormats = [];
			state.users = [];
			state.tags = [];
			state.languages = [];
			state.mediaItems = [];
			state.mediaRequests = [];
			state.discoverSections = [];
			state.discoverSection = undefined;
			state.discoverSectionPage = 1;
			state.discoverSectionHasMore = true;
			state.metadataDetail = undefined;
			state.mediaCollection = undefined;
			state.autocompleteGroups = [];
			state.advancedSearchGroups = [];
			state.releaseResults = {};
			state.activities = [];
			state.libraryFolders = [];
			state.pathMappings = [];
			state.libraryScansByFolder = {};
			state.openLibraryFolderId = undefined;
			state.downloadForm = emptyDownloadClientForm();
			state.indexerForm = emptyIndexerForm();
			state.libraryFolderForm = emptyLibraryFolderForm();
			state.pathMappingForm = emptyPathMappingForm();
			state.mediaProfileForm = emptyMediaProfileForm();
			state.customFormatForm = emptyCustomFormatForm();
			state.tagForm = emptyTagForm();
			state.userForm = emptyUserForm();
		}
	}

	return { initialise, connectEvents, disconnectEvents, login, logout };
}
