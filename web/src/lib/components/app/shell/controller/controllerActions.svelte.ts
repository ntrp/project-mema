import type { createActivityCache } from '$lib/features/activity/cache';
import type { createLibraryCache } from '$lib/features/library/cache';
import type { createReleaseCache } from '$lib/features/releases/cache';
import type { createDiscoverBlacklistCache } from '$lib/features/discovery/blacklist/cache';
import type { createDiscoverContentCache } from '$lib/features/discovery/content/cache';
import type { createSearchCache } from '$lib/features/search/cache';
import type { MediaAdvancedSearchRequest } from '$lib/features/search/api';
import { createDiscoveryActions } from './discoveryActions';
import { createEventActions, mediaStatusFromActivity } from './events';
import { createFormCancelActions } from './formCancelActions';
import { createLoadActions } from './loaders';
import { createMediaActions } from './mediaActions';
import { createMediaComponentActions } from './mediaComponentActions';
import { createMediaFulfillmentActions } from './mediaFulfillmentActions';
import { createMediaMetadataActions } from './mediaMetadataActions';
import { createMediaSubtitleActions } from './mediaSubtitleActions';
import { createNavigationActions } from './navigationActions';
import { createNoticeActions } from './noticeActions';
import { createProfileActions } from './profileActions';
import { createReleaseActions } from './releaseActions';
import { createRouteActions } from './routeActions';
import { createSearchActions } from './searchActions';
import { createSessionActions } from './sessionActions';
import { createSettingsDeleteActions } from './settingsDeleteActions';
import { createSettingsEditActions } from './settingsEditActions';
import { createSettingsSaveActions } from './settingsSaveActions';
import { createSettingsTestCacheActions } from './settingsTestCacheActions';
import type { createControllerResourceRuntime } from './controllerResourceRuntime.svelte';
import type { AppShellState } from './state.svelte';

interface Deps {
	activity: ReturnType<typeof createActivityCache>;
	library: ReturnType<typeof createLibraryCache>;
	releases: ReturnType<typeof createReleaseCache>;
	blacklist: ReturnType<typeof createDiscoverBlacklistCache>;
	discovery: ReturnType<typeof createDiscoverContentCache>;
	search: ReturnType<typeof createSearchCache>;
	resources: ReturnType<typeof createControllerResourceRuntime>;
	setAutocomplete: (_query: string) => void;
	setAdvanced: (_request: MediaAdvancedSearchRequest) => void;
}

export function createControllerActions(state: AppShellState, deps: Deps) {
	const notices = createNoticeActions(state);
	const profile = createProfileActions(state, notices);
	const events = createEventActions(state);
	const loads = createLoadActions(state);
	const updateStatus = (activity: import('$lib/settings/types').DownloadActivity) => {
		const status = mediaStatusFromActivity(activity.status);
		if (status)
			deps.library.mapItems((item) =>
				item.id === activity.mediaItemId ? { ...item, status } : item
			);
	};
	const refreshSettings = async () => {
		await Promise.all([
			loads.loadSettingsSection(state.activeSettingsSection),
			deps.resources.catalogCache.refresh()
		]);
	};
	const discovery = createDiscoveryActions(state, deps.blacklist, deps.discovery);
	const search = createSearchActions(state, {
		...notices,
		setAutocompleteQuery: deps.setAutocomplete,
		setAdvancedRequest: deps.setAdvanced
	});
	const media = createMediaActions(state, {
		...notices,
		loadMediaItems: deps.library.refreshItems,
		removeActivityForMedia: deps.activity.removeForMedia,
		removeReleaseResults: deps.releases.remove,
		mediaItems: deps.library.items,
		upsertMediaItem: deps.library.upsertItem,
		mapMediaItems: deps.library.mapItems,
		removeMediaItem: deps.library.removeItem,
		upsertMediaRequest: deps.library.upsertRequest,
		mapMediaRequests: deps.library.mapRequests
	});
	const mediaComponents = createMediaComponentActions(state, {
		...notices,
		loadMediaItems: deps.library.refreshItems
	});
	const mediaFulfillment = createMediaFulfillmentActions(state, {
		...notices,
		loadMediaItems: deps.library.refreshItems
	});
	const mediaMetadata = createMediaMetadataActions(state, {
		...notices,
		upsertMediaItem: deps.library.upsertItem
	});
	const mediaSubtitles = createMediaSubtitleActions(state, {
		...notices,
		upsertMediaItem: deps.library.upsertItem
	});
	const releaseActions = createReleaseActions(state, {
		...notices,
		upsertActivity: deps.activity.upsert,
		refreshActivity: deps.activity.refresh,
		updateMediaStatusFromActivity: updateStatus,
		setReleaseResult: deps.releases.set,
		loadReleaseResult: deps.releases.load
	});
	const catalog = deps.resources.catalogCache;
	const scans = deps.resources.scans;
	const settingsSave = createSettingsSaveActions(state, {
		...notices,
		loadSettings: refreshSettings,
		loadMediaItems: deps.library.refreshItems,
		mediaItems: deps.library.items,
		users: () => deps.resources.queries.users.data ?? [],
		upsertLibraryFolder: catalog.upsertLibraryFolder,
		upsertPathMapping: catalog.upsertPathMapping,
		upsertLibraryScan: scans.upsert
	});
	const settingsDelete = createSettingsDeleteActions(state, {
		...notices,
		loadSettings: refreshSettings,
		refreshMediaItems: deps.library.refreshItems,
		removeLanguage: catalog.removeLanguage,
		removeTag: catalog.removeTag,
		removeUser: catalog.removeUser,
		removeDownloadClient: catalog.removeDownloadClient,
		removeIndexer: catalog.removeIndexer,
		removeSubtitleProvider: catalog.removeSubtitleProvider,
		removeLibraryFolder: catalog.removeLibraryFolder,
		removePathMapping: catalog.removePathMapping,
		removeMediaProfile: catalog.removeMediaProfile,
		removeCustomFormat: catalog.removeCustomFormat,
		upsertLibraryScan: scans.upsert,
		removeLibraryScan: scans.remove
	});
	const routeData = {
		loadSettingsSection: loads.loadSettingsSection,
		loadSystemSettings: loads.loadSystemSettings,
		loadMediaActionSettings: loads.loadMediaActionSettings,
		loadProfile: profile.loadProfile
	};
	const session = createSessionActions(state, {
		...notices,
		events: {
			loadMediaItems: deps.library.refreshItems,
			upsertActivity: deps.activity.upsert,
			updateMediaStatusFromActivity: updateStatus,
			appendIndexerSearchHistory: events.appendIndexerSearchHistory,
			upsertIndexerSearchCache: events.upsertIndexerSearchCache,
			upsertMetadataCache: events.upsertMetadataCache,
			appendMetadataSearchHistory: events.appendMetadataSearchHistory,
			updateFulfillmentJobExecution: mediaFulfillment.updateFulfillmentJobExecution,
			parseEventData: events.parseEventData
		},
		clearActivityCache: deps.activity.clear,
		clearLibraryCache: deps.library.clear,
		clearReleaseCache: deps.releases.clear,
		clearDiscoverBlacklistCache: deps.blacklist.clear,
		clearDiscoverContentCache: deps.discovery.clear,
		clearSearchCache: deps.search.clear,
		clearSettingsCatalogCache: catalog.clear,
		clearServerResourceCache: deps.resources.server.clear,
		clearLibraryScanCache: scans.clear,
		routeData
	});
	return [
		notices,
		profile,
		session,
		discovery,
		search,
		media,
		mediaComponents,
		mediaFulfillment,
		mediaMetadata,
		mediaSubtitles,
		releaseActions,
		loads,
		settingsSave,
		settingsDelete,
		createSettingsTestCacheActions(state, { ...notices, loadSettings: refreshSettings }),
		createSettingsEditActions(state),
		createRouteActions(state, { routeData }),
		createNavigationActions(state),
		createFormCancelActions(state)
	];
}
