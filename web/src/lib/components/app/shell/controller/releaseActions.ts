import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
	grabMediaRelease as grabMediaReleaseRequest,
	searchMediaReleases as searchMediaReleasesRequest
} from '$lib/settings/api';
import type { DownloadActivity, MediaItem, ReleaseCandidate } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface ReleaseDeps {
	clearNotice: () => void;
	loadDownloadActivity: () => Promise<void>;
	updateMediaStatusFromActivity: (activity: DownloadActivity) => void;
}

export function createReleaseActions(state: AppShellState, deps: ReleaseDeps) {
	const clearNotice = deps.clearNotice;
	const loadDownloadActivity = deps.loadDownloadActivity;
	const updateMediaStatusFromActivity = deps.updateMediaStatusFromActivity;

	async function findReleases(item: MediaItem) {
		state.searchingItemId = item.id;
		clearNotice();

		try {
			const job = await enqueueMediaReleaseSearchRequest(item.id);
			state.releaseResults = {
				...state.releaseResults,
				[item.id]: { loaded: false, releases: [], errors: [`${job.message} (#${job.jobId})`] }
			};
			state.message = job.message;
			window.setTimeout(() => void loadReleaseResults(item.id), 1200);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue release search');
		} finally {
			state.searchingItemId = undefined;
		}
	}

	async function loadReleaseResults(id: string) {
		try {
			const results = await searchMediaReleasesRequest(id);
			state.releaseResults = {
				...state.releaseResults,
				[id]: { loaded: true, releases: results.releases, errors: results.errors }
			};
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load release results');
		}
	}

	async function grabRelease(item: MediaItem, release: ReleaseCandidate) {
		state.grabbingKey = `${item.id}:${release.id}`;
		clearNotice();

		try {
			const result = await grabMediaReleaseRequest(item.id, release);
			state.activities = [
				result.activity,
				...state.activities.filter((activity) => activity.id !== result.activity.id)
			];
			updateMediaStatusFromActivity(result.activity);
			state.message = `${result.message} (#${result.jobId})`;
			state.activeHomeSection = 'activity';
			void goto(resolve('/activity'));
			window.setTimeout(() => void loadDownloadActivity(), 1200);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue download');
		} finally {
			state.grabbingKey = undefined;
		}
	}

	return { findReleases, grabRelease };
}
