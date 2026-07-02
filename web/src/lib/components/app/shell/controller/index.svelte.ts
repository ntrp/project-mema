import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptyUserForm
} from '$lib/settings/forms';
import { emptyTagForm } from './helpers';
import { createActivityActions } from './activityActions';
import { createDiscoveryActions } from './discoveryActions';
import { createEventActions } from './events';
import { createLoadActions } from './loaders';
import { createMediaActions } from './mediaActions';
import { createMediaMetadataActions } from './mediaMetadataActions';
import { createNavigationActions } from './navigationActions';
import { createNoticeActions } from './noticeActions';
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
	const notices = createNoticeActions(state);
	const events = createEventActions(state);
	const loads = createLoadActions(state);
	const discovery = createDiscoveryActions(state);
	const search = createSearchActions(state, notices);
	const media = createMediaActions(state, {
		...notices,
		loadMediaItems: loads.loadMediaItems,
		loadSettings: loads.loadSettings
	});
	const mediaMetadata = createMediaMetadataActions(state, notices);
	const releases = createReleaseActions(state, {
		...notices,
		loadDownloadActivity: loads.loadDownloadActivity,
		updateMediaStatusFromActivity: events.updateMediaStatusFromActivity
	});
	const activity = createActivityActions(state, {
		...notices,
		loadMediaItems: loads.loadMediaItems,
		upsertActivity: events.upsertActivity
	});
	const settingsSave = createSettingsSaveActions(state, {
		...notices,
		loadSettings: loads.loadSettings
	});
	const settingsDelete = createSettingsDeleteActions(state, {
		...notices,
		loadSettings: loads.loadSettings
	});
	const settingsTestCache = createSettingsTestCacheActions(state, {
		...notices,
		loadSettings: loads.loadSettings
	});
	const settingsEdit = createSettingsEditActions(state);
	const routeActions = createRouteActions(state, {
		loadDiscoverSection: discovery.loadDiscoverSection,
		loadMediaCollection: loads.loadMediaCollection,
		loadMetadataDetail: loads.loadMetadataDetail
	});
	const navigation = createNavigationActions(state, {
		loadDiscoverSection: discovery.loadDiscoverSection
	});
	const session = createSessionActions(state, {
		...notices,
		loadSettings: loads.loadSettings,
		loadDiscoverBlacklist: discovery.loadDiscoverBlacklist,
		loadLibrary: loads.loadLibrary,
		loadDiscoverSections: discovery.loadDiscoverSections,
		loadMetadataDetail: loads.loadMetadataDetail,
		loadMediaCollection: loads.loadMediaCollection,
		loadDiscoverSection: discovery.loadDiscoverSection,
		events: {
			loadMediaItems: loads.loadMediaItems,
			upsertActivity: events.upsertActivity,
			updateMediaStatusFromActivity: events.updateMediaStatusFromActivity,
			appendIndexerSearchHistory: events.appendIndexerSearchHistory,
			upsertIndexerSearchCache: events.upsertIndexerSearchCache,
			upsertMetadataCache: events.upsertMetadataCache,
			appendMetadataSearchHistory: events.appendMetadataSearchHistory,
			parseEventData: events.parseEventData
		}
	});

	return Object.assign(
		state,
		notices,
		session,
		discovery,
		search,
		media,
		mediaMetadata,
		releases,
		activity,
		loads,
		settingsSave,
		settingsDelete,
		settingsTestCache,
		settingsEdit,
		routeActions,
		navigation,
		{
			cancelDownloadClient: () => (state.downloadForm = emptyDownloadClientForm()),
			cancelIndexer: () => (state.indexerForm = emptyIndexerForm()),
			cancelMediaProfile: () => (state.mediaProfileForm = emptyMediaProfileForm()),
			cancelCustomFormat: () => (state.customFormatForm = emptyCustomFormatForm()),
			cancelTag: () => (state.tagForm = emptyTagForm()),
			cancelLanguage: () => (state.languageForm = emptyLanguageForm()),
			cancelUser: () => (state.userForm = emptyUserForm())
		}
	);
}
