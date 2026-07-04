import {
	cancelDownloadActivity as cancelDownloadActivityRequest,
	clearReleaseBlocklist as clearReleaseBlocklistRequest,
	deleteReleaseBlocklistItem as deleteReleaseBlocklistItemRequest,
	deleteDownloadActivity as deleteDownloadActivityRequest
} from '$lib/settings/api';
import type { DownloadActivity, ReleaseBlocklistItem } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface ActivityDeps {
	clearNotice: () => void;
	loadMediaItems: () => Promise<void>;
	upsertActivity: (activity: DownloadActivity) => void;
}

export function createActivityActions(state: AppShellState, deps: ActivityDeps) {
	const clearNotice = deps.clearNotice;
	const loadMediaItems = deps.loadMediaItems;
	const upsertActivity = deps.upsertActivity;
	async function cancelActivity(activity: DownloadActivity) {
		state.cancellingActivityId = activity.id;
		clearNotice();

		try {
			const cancelled = await cancelDownloadActivityRequest(activity.id);
			upsertActivity(cancelled);
			await loadMediaItems();
			state.message = 'Download activity cancelled';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not cancel download activity');
		} finally {
			state.cancellingActivityId = undefined;
		}
	}

	async function deleteActivity(activity: DownloadActivity) {
		state.deletingActivityId = activity.id;
		clearNotice();

		try {
			await deleteDownloadActivityRequest(activity.id);
			state.activities = state.activities.filter((item) => item.id !== activity.id);
			state.message = 'Download activity deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete download activity');
		} finally {
			state.deletingActivityId = undefined;
		}
	}

	async function deleteReleaseBlocklistItem(item: ReleaseBlocklistItem) {
		state.deletingReleaseBlocklistId = item.id;
		clearNotice();

		try {
			await deleteReleaseBlocklistItemRequest(item.id);
			state.releaseBlocklist = state.releaseBlocklist.filter((block) => block.id !== item.id);
			state.message = 'Release blocklist entry removed';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not remove release blocklist entry');
		} finally {
			state.deletingReleaseBlocklistId = undefined;
		}
	}

	async function clearReleaseBlocklist() {
		state.clearingReleaseBlocklist = true;
		clearNotice();

		try {
			await clearReleaseBlocklistRequest();
			state.releaseBlocklist = [];
			state.message = 'Release blocklist cleared';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not clear release blocklist');
		} finally {
			state.clearingReleaseBlocklist = false;
		}
	}

	return { cancelActivity, deleteActivity, deleteReleaseBlocklistItem, clearReleaseBlocklist };
}
