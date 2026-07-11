import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptyUserForm
} from '$lib/settings/forms';
import { useQueryClient } from '@tanstack/svelte-query';
import { createActivityCache } from '$lib/features/activity/cache';
import { createLibraryCache } from '$lib/features/library/cache';
import { emptyTagForm } from './helpers';
import { createDiscoveryActions } from './discoveryActions';
import { createEventActions, mediaStatusFromActivity } from './events';
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
import { AppShellState } from './state.svelte';
import { defaultRouteState, type AppRouteState } from './routeState';

export type { PeopleSectionKind, RelatedSectionKind } from './types';
export type { AppRouteState } from './routeState';

export function createAppShellController(route: AppRouteState = defaultRouteState()) {
	const state = new AppShellState(route);
	const activityCache = createActivityCache(useQueryClient());
	const libraryCache = createLibraryCache(useQueryClient());
	const updateMediaStatusFromActivity = (
		activity: import('$lib/settings/types').DownloadActivity
	) => {
		const status = mediaStatusFromActivity(activity.status);
		if (status)
			libraryCache.mapItems((item) =>
				item.id === activity.mediaItemId ? { ...item, status } : item
			);
	};
	const notices = createNoticeActions(state);
	const profile = createProfileActions(state, notices);
	const events = createEventActions(state);
	const loads = createLoadActions(state);
	const discovery = createDiscoveryActions(state);
	const search = createSearchActions(state, notices);
	const media = createMediaActions(state, {
		...notices,
		loadMediaItems: libraryCache.refreshItems,
		loadSettings: loads.loadSettings,
		removeActivityForMedia: activityCache.removeForMedia,
		mediaItems: libraryCache.items,
		upsertMediaItem: libraryCache.upsertItem,
		mapMediaItems: libraryCache.mapItems,
		removeMediaItem: libraryCache.removeItem,
		upsertMediaRequest: libraryCache.upsertRequest,
		mapMediaRequests: libraryCache.mapRequests
	});
	const mediaComponents = createMediaComponentActions(state, {
		...notices,
		loadMediaItems: libraryCache.refreshItems
	});
	const mediaFulfillment = createMediaFulfillmentActions(state, {
		...notices,
		loadMediaItems: libraryCache.refreshItems
	});
	const mediaMetadata = createMediaMetadataActions(state, {
		...notices,
		upsertMediaItem: libraryCache.upsertItem
	});
	const mediaSubtitles = createMediaSubtitleActions(state, {
		...notices,
		upsertMediaItem: libraryCache.upsertItem
	});
	const releases = createReleaseActions(state, {
		...notices,
		upsertActivity: activityCache.upsert,
		refreshActivity: activityCache.refresh,
		updateMediaStatusFromActivity
	});
	const settingsSave = createSettingsSaveActions(state, {
		...notices,
		loadSettings: loads.loadSettings,
		loadMediaItems: libraryCache.refreshItems,
		mediaItems: libraryCache.items
	});
	const settingsDelete = createSettingsDeleteActions(state, {
		...notices,
		loadSettings: loads.loadSettings,
		refreshMediaItems: libraryCache.refreshItems
	});
	const settingsTestCache = createSettingsTestCacheActions(state, {
		...notices,
		loadSettings: loads.loadSettings
	});
	const settingsEdit = createSettingsEditActions(state);
	const routeActions = createRouteActions(state, {
		routeData: {
			loadSettings: loads.loadSettings,
			loadDiscoverBlacklist: discovery.loadDiscoverBlacklist,
			loadDiscoverSections: discovery.loadDiscoverSections,
			loadDiscoverSection: discovery.loadDiscoverSection,
			loadMetadataDetail: loads.loadMetadataDetail,
			loadPersonDetail: loads.loadPersonDetail,
			loadMediaCollection: loads.loadMediaCollection,
			loadProfile: profile.loadProfile
		}
	});
	const navigation = createNavigationActions(state, {
		loadDiscoverSection: discovery.loadDiscoverSection
	});
	const session = createSessionActions(state, {
		...notices,
		events: {
			loadMediaItems: libraryCache.refreshItems,
			upsertActivity: activityCache.upsert,
			updateMediaStatusFromActivity,
			appendIndexerSearchHistory: events.appendIndexerSearchHistory,
			upsertIndexerSearchCache: events.upsertIndexerSearchCache,
			upsertMetadataCache: events.upsertMetadataCache,
			appendMetadataSearchHistory: events.appendMetadataSearchHistory,
			updateFulfillmentJobExecution: mediaFulfillment.updateFulfillmentJobExecution,
			parseEventData: events.parseEventData
		},
		clearActivityCache: activityCache.clear,
		clearLibraryCache: libraryCache.clear,
		routeData: {
			loadSettings: loads.loadSettings,
			loadDiscoverBlacklist: discovery.loadDiscoverBlacklist,
			loadDiscoverSections: discovery.loadDiscoverSections,
			loadDiscoverSection: discovery.loadDiscoverSection,
			loadMetadataDetail: loads.loadMetadataDetail,
			loadPersonDetail: loads.loadPersonDetail,
			loadMediaCollection: loads.loadMediaCollection,
			loadProfile: profile.loadProfile
		}
	});
	function cancelDownloadClient() {
		state.downloadForm = emptyDownloadClientForm();
	}

	function cancelIndexer() {
		state.indexerForm = emptyIndexerForm();
	}

	function cancelMediaProfile() {
		state.mediaProfileForm = emptyMediaProfileForm();
	}

	function cancelCustomFormat() {
		state.customFormatForm = emptyCustomFormatForm();
	}

	function cancelTag() {
		state.tagForm = emptyTagForm();
	}

	function cancelLanguage() {
		state.languageForm = emptyLanguageForm();
	}

	function cancelUser() {
		state.userForm = emptyUserForm();
	}

	return Object.assign(
		state,
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
		releases,
		loads,
		settingsSave,
		settingsDelete,
		settingsTestCache,
		settingsEdit,
		routeActions,
		navigation,
		{
			cancelDownloadClient,
			cancelIndexer,
			cancelMediaProfile,
			cancelCustomFormat,
			cancelTag,
			cancelLanguage,
			cancelUser
		}
	);
}
