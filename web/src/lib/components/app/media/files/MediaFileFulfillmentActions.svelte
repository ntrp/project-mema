<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import FileOutputIcon from '@lucide/svelte/icons/file-output';
	import FileVideoIcon from '@lucide/svelte/icons/file-video';
	import MusicIcon from '@lucide/svelte/icons/music';
	import PackageIcon from '@lucide/svelte/icons/package';
	import WandIcon from '@lucide/svelte/icons/wand-sparkles';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFulfillmentActionRequest } from '$lib/settings/types';

	interface Props {
		row: MediaFileDetailRow;
		canManage: boolean;
		onFulfillmentAction: (_request: MediaFulfillmentActionRequest) => void | Promise<void>;
	}

	let { row, canManage, onFulfillmentAction }: Props = $props();

	type Action = { label: string; request: MediaFulfillmentActionRequest; icon: string };

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
		if (matchesOperation(track, 'remux')) {
			return [action('Remux container', 'container_remux', 'video', track.languageId, 'package')];
		}
		return [action('Transcode video', 'video_transcode', 'video', track.languageId, 'video')];
	}

	function audioActions(track: MediaFileDetailRow): Action[] {
		if (track.missing) {
			return [action('Source audio', 'audio_sourcing', 'audio', track.languageId, 'music')];
		}
		if (track.visualState !== 'partial' && track.visualState !== 'pending_operation') return [];
		return [action('Transcode audio', 'audio_transcode', 'audio', track.languageId, 'music')];
	}

	function subtitleActions(track: MediaFileDetailRow): Action[] {
		if (track.missing) {
			return [
				action('Download subtitle', 'subtitle_download', 'subtitle', track.languageId, 'download')
			];
		}
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
			label,
			icon,
			request: {
				operation,
				targetType,
				languageId,
				trackId: row.trackId,
				otherFileId: row.otherFileId
			}
		};
	}

	function matchesOperation(track: MediaFileDetailRow, value: string) {
		return (track.operationLabel ?? '').toLowerCase().includes(value);
	}

	function clickAction(event: MouseEvent, request: MediaFulfillmentActionRequest) {
		event.stopPropagation();
		void onFulfillmentAction(request);
	}
</script>

{#if actions.length > 0}
	<span class="inline-flex justify-end gap-1">
		{#each actions as item (item.label)}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							size="icon-sm"
							variant="outline"
							disabled={!canManage}
							aria-label={item.label}
							onclick={(event) => clickAction(event, item.request)}
							onkeydown={(event) => event.stopPropagation()}
						>
							{#if item.icon === 'video'}
								<FileVideoIcon aria-hidden="true" />
							{:else if item.icon === 'music'}
								<MusicIcon aria-hidden="true" />
							{:else if item.icon === 'download'}
								<DownloadIcon aria-hidden="true" />
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
