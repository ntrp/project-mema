<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import FileOutputIcon from '@lucide/svelte/icons/file-output';
	import FileVideoIcon from '@lucide/svelte/icons/file-video';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import MusicIcon from '@lucide/svelte/icons/music';
	import PackageIcon from '@lucide/svelte/icons/package';
	import WandIcon from '@lucide/svelte/icons/wand-sparkles';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
	import { mediaFulfillmentActionKey } from '$lib/settings/mediaFulfillmentActionKey';
	import type { MediaFulfillmentActionRequest } from '$lib/settings/types';

	interface Props {
		row: MediaFileDetailRow;
		canManage: boolean;
		pendingFulfillmentActionKeys?: string[];
		onFulfillmentAction: (_request: MediaFulfillmentActionRequest) => void | Promise<void>;
	}

	let { row, canManage, pendingFulfillmentActionKeys = [], onFulfillmentAction }: Props = $props();
	let pendingAction = $state<string | undefined>();

	type Action = {
		key: string;
		label: string;
		request: MediaFulfillmentActionRequest;
		icon: string;
	};

	const actions = $derived(fulfillmentActions(row));

	function fulfillmentActions(track: MediaFileDetailRow): Action[] {
		if (track.chapterSummary || track.type === 'chapter' || track.unwanted) return [];
		if (track.type === 'video') return videoActions(track);
		if (track.type === 'audio') return audioActions(track);
		if (track.type === 'subtitle') return subtitleActions(track);
		return [];
	}

	function videoActions(track: MediaFileDetailRow): Action[] {
		if (track.visualState !== 'partial' && track.visualState !== 'pending_operation') return [];
		if (!track.trackId) return [];
		if (!hasSupportedVideoTranscodeMismatch(track)) return [];
		return [action('Transcode video', 'video_transcode', 'video', track.languageId, 'video')];
	}

	function audioActions(track: MediaFileDetailRow): Action[] {
		if (track.missing) return [];
		if (track.visualState !== 'partial' && track.visualState !== 'pending_operation') return [];
		if (!track.trackId) return [];
		if (!hasTargetCodecMismatch(track)) return [];
		return [action('Transcode audio', 'audio_transcode', 'audio', track.languageId, 'binary')];
	}

	function subtitleActions(track: MediaFileDetailRow): Action[] {
		if (track.missing) return [];
		if (matchesOperation(track, 'embed')) {
			return [action('Embed subtitle', 'subtitle_embed', 'subtitle', track.languageId, 'captions')];
		}
		if (matchesOperation(track, 'extract')) {
			return [
				action('Extract subtitle', 'subtitle_extraction', 'subtitle', track.languageId, 'output')
			];
		}
		if (track.visualState === 'partial' || matchesOperation(track, 'convert')) {
			return [
				action('Convert subtitle', 'subtitle_conversion', 'subtitle', track.languageId, 'wand')
			];
		}
		return [];
	}

	function action(
		label: string,
		operation: MediaFulfillmentActionRequest['operation'],
		targetType: MediaFulfillmentActionRequest['targetType'],
		languageId: string | undefined,
		icon: string
	): Action {
		return {
			key: `${operation}:${row.trackId ?? row.otherFileId ?? row.key}`,
			label,
			icon,
			request: {
				operation,
				filePath: row.filePath,
				targetType,
				languageId,
				trackId: row.trackId,
				otherFileId: row.otherFileId
			}
		};
	}

	function actionKey(request: MediaFulfillmentActionRequest) {
		return mediaFulfillmentActionKey(request);
	}

	function matchesOperation(track: MediaFileDetailRow, value: string) {
		return (track.operationLabel ?? '').toLowerCase().includes(value);
	}

	function hasTargetCodecMismatch(track: MediaFileDetailRow) {
		return (track.details ?? []).some((detail) => detail.toLowerCase().includes('codec'));
	}

	function hasSupportedVideoTranscodeMismatch(track: MediaFileDetailRow) {
		return (track.details ?? []).some((detail) => {
			const normalized = detail.toLowerCase();
			return normalized.includes('codec') || normalized.includes('pixel format');
		});
	}

	function isPending(item: Action) {
		const key = actionKey(item.request);
		return pendingAction === key || pendingFulfillmentActionKeys.includes(key);
	}

	async function clickAction(event: MouseEvent, item: Action) {
		event.stopPropagation();
		if (isPending(item)) {
			await goto(resolve('/system/jobs'));
			return;
		}
		pendingAction = actionKey(item.request);
		try {
			await onFulfillmentAction(item.request);
		} finally {
			if (pendingAction === actionKey(item.request)) {
				pendingAction = undefined;
			}
		}
	}
</script>

{#if actions.length > 0}
	<span class="inline-flex justify-end gap-1">
		{#each actions as item (item.key)}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						{@const pending = isPending(item)}
						<Button
							{...props}
							type="button"
							size="icon-sm"
							variant="outline"
							disabled={!canManage || (!!pendingAction && !pending)}
							aria-label={item.label}
							aria-busy={pending}
							onclick={(event) => clickAction(event, item)}
							onkeydown={(event) => event.stopPropagation()}
						>
							{#if pending}
								<LoaderCircleIcon class="animate-spin" aria-hidden="true" />
							{:else if item.icon === 'video'}
								<FileVideoIcon aria-hidden="true" />
							{:else if item.icon === 'music'}
								<MusicIcon aria-hidden="true" />
							{:else if item.icon === 'captions'}
								<CaptionsIcon aria-hidden="true" />
							{:else if item.icon === 'output'}
								<FileOutputIcon aria-hidden="true" />
							{:else if item.icon === 'package'}
								<PackageIcon aria-hidden="true" />
							{:else}
								<WandIcon aria-hidden="true" />
							{/if}
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>{item.label}</Tooltip.Content>
			</Tooltip.Root>
		{/each}
	</span>
{/if}
