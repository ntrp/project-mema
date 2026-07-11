import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	approveMediaRequest as approveMediaRequestRequest,
	createMediaItem as createMediaItemRequest,
	createMediaRequest as createMediaRequestRequest
} from '$lib/features/library/commands';
import type { MediaActionSelection } from '$lib/components/app/media/actions/mediaActionTypes';
import type {
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchResult
} from '$lib/settings/types';
import { candidateKey, errorMessageFrom } from '../helpers';
import type { AppShellState } from '../state.svelte';
import type { MediaDeps } from './types';

export function createMediaCreateActions(state: AppShellState, deps: MediaDeps) {
	function addMedia(candidate: MediaSearchResult) {
		state.activeMediaCandidate = candidate;
		deps.clearNotice();
	}

	function closeMediaAction() {
		if (!state.savingMediaAction) state.activeMediaCandidate = undefined;
	}

	async function confirmMediaAction(selection: MediaActionSelection) {
		const candidate = state.activeMediaCandidate;
		if (!candidate) return;
		state.addingKey = candidateKey(candidate);
		state.savingMediaAction = true;
		deps.clearNotice();
		try {
			if (state.isAdmin) {
				await addLibraryItem(candidate, selection);
				return;
			}
			const request = await createMediaRequestRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath
			});
			deps.upsertMediaRequest(request);
			state.message = 'Media request created';
			state.activeHomeSection = 'requests';
			state.activeMediaCandidate = undefined;
			void goto(resolve('/requests'));
		} catch (error) {
			state.errorMessage = errorMessageFrom(
				error,
				state.isAdmin ? 'Could not add media item' : 'Could not create media request'
			);
		} finally {
			state.addingKey = undefined;
			state.savingMediaAction = false;
		}
	}

	async function addLibraryItem(candidate: MediaSearchResult, selection: MediaActionSelection) {
		if (!selection.qualityProfileId || !selection.libraryFolderId) {
			throw new Error('Quality profile and library folder are required');
		}
		await createMediaItemRequest({
			title: candidate.title,
			type: candidate.type,
			year: candidate.year,
			monitored: selection.monitorMode !== 'none',
			monitorMode: selection.monitorMode,
			seriesType: candidate.type === 'serie' ? selection.seriesType : undefined,
			minimumAvailability: selection.minimumAvailability,
			startSearch: selection.startSearch,
			externalProvider: candidate.externalProvider,
			externalId: candidate.externalId,
			overview: candidate.overview,
			posterPath: candidate.posterPath,
			qualityProfileId: selection.qualityProfileId,
			libraryFolderId: selection.libraryFolderId,
			tags: selection.tags
		});
		await deps.loadMediaItems();
		state.message =
			selection.monitorMode === 'none'
				? 'Media item added to library'
				: selection.monitorMode === 'collection'
					? 'Media collection added to monitored'
					: 'Media item added to monitored';
		state.activeHomeSection = candidate.type === 'movie' ? 'movies' : 'series';
		state.activeMediaCandidate = undefined;
		void goto(resolve(candidate.type === 'movie' ? '/movies' : '/series'));
	}

	async function approveMediaRequest(request: MediaRequest, approval: MediaRequestApproveRequest) {
		state.approvingRequestId = request.id;
		deps.clearNotice();
		try {
			const result = await approveMediaRequestRequest(request.id, approval);
			deps.mapMediaRequests((item) => (item.id === result.request.id ? result.request : item));
			deps.upsertMediaItem(result.mediaItem);
			state.message = 'Media request approved and added to monitored';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not approve media request');
		} finally {
			state.approvingRequestId = undefined;
		}
	}

	return { addMedia, closeMediaAction, confirmMediaAction, approveMediaRequest };
}
