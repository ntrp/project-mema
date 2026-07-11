import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
	grabMediaRelease as grabMediaReleaseRequest,
	searchMediaReleases as searchMediaReleasesRequest
} from '$lib/settings/api';
import type {
	DownloadActivity,
	MediaItem,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

const RELEASE_SEARCH_POLL_MS = 1000;
const RELEASE_SEARCH_MAX_POLLS = 120;

interface ReleaseDeps {
	clearNotice: () => void;
	upsertActivity: (_activity: DownloadActivity) => void;
	refreshActivity: () => Promise<void>;
	updateMediaStatusFromActivity: (activity: DownloadActivity) => void;
}

export function createReleaseActions(state: AppShellState, deps: ReleaseDeps) {
	const clearNotice = deps.clearNotice;
	const updateMediaStatusFromActivity = deps.updateMediaStatusFromActivity;

	async function findReleases(item: MediaItem, query?: string) {
		state.searchingItemId = item.id;
		clearNotice();

		try {
			const job = await enqueueMediaReleaseSearchRequest(item.id, query);
			const queuedMessage = `${job.message} (#${job.jobId})`;
			state.releaseResults = {
				...state.releaseResults,
				[item.id]: { loaded: false, releases: [], errors: [queuedMessage] }
			};
			state.message = job.message;
			await pollReleaseResults(item.id, queuedMessage);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue release search');
		} finally {
			state.searchingItemId = undefined;
		}
	}

	async function pollReleaseResults(id: string, queuedMessage: string) {
		for (let attempt = 0; attempt < RELEASE_SEARCH_MAX_POLLS; attempt += 1) {
			await sleep(RELEASE_SEARCH_POLL_MS);
			const loaded = await loadReleaseResults(id, false);
			if (loaded) return;
			state.releaseResults = {
				...state.releaseResults,
				[id]: { loaded: false, releases: [], errors: [queuedMessage] }
			};
		}
		state.releaseResults = {
			...state.releaseResults,
			[id]: { loaded: true, releases: [], errors: ['Release search is still running.'] }
		};
	}

	async function loadReleaseResults(id: string, markEmptyLoaded = true) {
		try {
			const results = await searchMediaReleasesRequest(id);
			const complete = results.releases.length > 0 || results.errors.length > 0;
			if (!complete && !markEmptyLoaded) {
				return false;
			}
			state.releaseResults = {
				...state.releaseResults,
				[id]: { loaded: true, releases: results.releases, errors: results.errors }
			};
			return true;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load release results');
			return true;
		}
	}

	async function grabRelease(
		item: MediaItem,
		release: ReleaseCandidate,
		overrideMatch = false,
		overrideDetails?: ReleaseOverrideDetails
	) {
		state.grabbingKey = `${item.id}:${release.id}`;
		clearNotice();

		try {
			const result = await grabMediaReleaseRequest(
				item.id,
				release,
				overrideMatch,
				overrideDetails
			);
			deps.upsertActivity(result.activity);
			updateMediaStatusFromActivity(result.activity);
			state.message = `${result.message} (#${result.jobId})`;
			state.activeHomeSection = 'activity';
			void goto(resolve('/activity'));
			window.setTimeout(() => void deps.refreshActivity(), 1200);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue download');
		} finally {
			state.grabbingKey = undefined;
		}
	}

	return { findReleases, grabRelease };
}

function sleep(ms: number) {
	return new Promise((resolve) => window.setTimeout(resolve, ms));
}
