import { enqueueMediaFulfillmentAction } from '$lib/features/releases/api';
import { mediaFulfillmentActionKey } from '$lib/settings/mediaFulfillmentActionKey';
import type {
	MediaFulfillmentActionRequest,
	MediaItem,
	SystemJobExecution
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface MediaFulfillmentDeps {
	clearNotice: () => void;
	loadMediaItems: () => Promise<void>;
}

export function createMediaFulfillmentActions(state: AppShellState, deps: MediaFulfillmentDeps) {
	async function enqueueMediaFulfillment(item: MediaItem, request: MediaFulfillmentActionRequest) {
		deps.clearNotice();
		const key = mediaFulfillmentActionKey(request);
		state.pendingFulfillmentActions = { ...state.pendingFulfillmentActions, [key]: 0 };
		try {
			const job = await enqueueMediaFulfillmentAction(item.id, request);
			if (key in state.pendingFulfillmentActions) {
				state.pendingFulfillmentActions = { ...state.pendingFulfillmentActions, [key]: job.jobId };
			}
			state.message = `${job.message} (#${job.jobId})`;
			await deps.loadMediaItems();
			return job;
		} catch (error) {
			const next = { ...state.pendingFulfillmentActions };
			delete next[key];
			state.pendingFulfillmentActions = next;
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue fulfillment job');
		}
	}

	function updateFulfillmentJobExecution(execution: SystemJobExecution) {
		const key = mediaFulfillmentActionKeyFromExecution(execution);
		if (!key) return;
		if (!finalJobStatus(execution.status)) {
			state.pendingFulfillmentActions = {
				...state.pendingFulfillmentActions,
				[key]: execution.riverJobId
			};
			return;
		}
		const next = { ...state.pendingFulfillmentActions };
		delete next[key];
		state.pendingFulfillmentActions = next;
	}

	return { enqueueMediaFulfillment, updateFulfillmentJobExecution };
}

function mediaFulfillmentActionKeyFromExecution(execution: SystemJobExecution) {
	if (!execution.kind.startsWith('media.fulfillment.')) return '';
	const args = parseJobArgs(execution.args);
	const operation = fulfillmentOperationFromKind(execution.kind);
	if (!operation) return '';
	return mediaFulfillmentActionKey({
		operation,
		filePath: textArg(args.file_path),
		targetType: textArg(args.target_type) as MediaFulfillmentActionRequest['targetType'],
		languageId: textArg(args.language_id),
		trackId: textArg(args.track_id),
		otherFileId: textArg(args.other_file_id),
		externalSubtitleId: textArg(args.external_subtitle_id)
	});
}

function fulfillmentOperationFromKind(
	kind: string
): MediaFulfillmentActionRequest['operation'] | '' {
	switch (kind.replace('media.fulfillment.', '')) {
		case 'video_transcode':
			return 'video_transcode';
		case 'audio_transcode':
			return 'audio_transcode';
		case 'container_remux':
			return 'container_remux';
		case 'subtitle_embed':
			return 'subtitle_embed';
		case 'subtitle_extract':
			return 'subtitle_extraction';
		case 'subtitle_convert':
			return 'subtitle_conversion';
		default:
			return '';
	}
}

function parseJobArgs(value: string) {
	try {
		return JSON.parse(value) as Record<string, unknown>;
	} catch {
		return {};
	}
}

function textArg(value: unknown) {
	return typeof value === 'string' && value.trim() ? value : undefined;
}

function finalJobStatus(value: string) {
	return value === 'completed' || value === 'cancelled' || value === 'discarded';
}
