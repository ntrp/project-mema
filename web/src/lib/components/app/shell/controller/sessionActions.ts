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
import { loadAppRouteData, type RouteDataDeps } from './routeData';
import { defaultRouteState } from './routeState';
import type { AppShellState } from './state.svelte';

interface SessionDeps {
	clearNotice: () => void;
	events: EventConnectionDeps;
	routeData: RouteDataDeps;
}

export function createSessionActions(state: AppShellState, deps: SessionDeps) {
	const clearNotice = deps.clearNotice;
	const routeData = deps.routeData;
	const eventDeps = deps.events;
	async function initialise() {
		state.loading = true;
		state.errorMessage = '';

		const session = await currentSessionRequest();
		state.authenticated = Boolean(session?.authenticated);
		state.currentUser = session?.user;
		if (state.authenticated) {
			if (
				state.currentUser?.role !== 'admin' &&
				(state.activeView === 'settings' || state.activeView === 'system')
			) {
				redirectToDiscover();
			} else if (state.activeHomeSection === 'blacklist') {
				redirectToDiscover();
			}
			await loadAppRouteData(state.route, state.currentUser?.role === 'admin', routeData);
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
			if (
				state.currentUser?.role !== 'admin' &&
				(state.activeView === 'settings' || state.activeView === 'system')
			) {
				redirectToDiscover();
			}
			await loadAppRouteData(state.route, state.currentUser?.role === 'admin', routeData);
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
			state.profile = undefined;
			state.profileErrorMessage = '';
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
			state.personDetail = undefined;
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

	function redirectToDiscover() {
		state.route = defaultRouteState();
		state.activeView = 'home';
		state.activeHomeSection = 'discover';
		void goto(resolve('/discover'));
	}

	return { initialise, connectEvents, disconnectEvents, login, logout };
}
