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
import { createMediaComponentActions } from './mediaComponentActions';
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
	const notices = createNoticeActions(state);
	const profile = createProfileActions(state, notices);
	const events = createEventActions(state);
	const loads = createLoadActions(state);
	const discovery = createDiscoveryActions(state);
	const search = createSearchActions(state, notices);
	const media = createMediaActions(state, {
		...notices,
		loadMediaItems: loads.loadMediaItems,
		loadSettings: loads.loadSettings
	});
	const mediaComponents = createMediaComponentActions(state, {
		...notices,
		loadMediaItems: loads.loadMediaItems
	});
	const mediaMetadata = createMediaMetadataActions(state, notices);
	const mediaSubtitles = createMediaSubtitleActions(state, notices);
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
		loadSettings: loads.loadSettings,
		loadMediaItems: loads.loadMediaItems
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
		loadMetadataDetail: loads.loadMetadataDetail,
		loadPersonDetail: loads.loadPersonDetail,
		loadProfile: profile.loadProfile
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
		loadPersonDetail: loads.loadPersonDetail,
		loadMediaCollection: loads.loadMediaCollection,
		loadDiscoverSection: discovery.loadDiscoverSection,
		loadProfile: profile.loadProfile,
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
		mediaMetadata,
		mediaSubtitles,
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
