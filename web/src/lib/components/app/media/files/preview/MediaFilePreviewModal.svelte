<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import InfoIcon from '@lucide/svelte/icons/info';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import MediaFileInfoPanel from '$lib/components/app/media/files/preview/MediaFileInfoPanel.svelte';
	import MediaFileVideoPlayer from '$lib/components/app/media/files/preview/MediaFileVideoPlayer.svelte';
	import {
		mediaFilePreviewInfoUrl,
		mediaFilePreviewUrl,
		mediaFileTextTracks,
		mediaFileVlcUrl,
		metadataAudioTrackOptions
	} from '$lib/components/app/media/files/preview/mediaFilePlayback';
	import {
		type MediaFilePlaybackStats,
		type MediaFilePreviewInfo
	} from '$lib/components/app/media/files/preview/mediaFilePreviewInfo';

	interface Props {
		mediaItemId: string;
		mediaTitle: string;
		row: MediaFileRow;
		onClose: () => void;
	}

	let { mediaItemId, mediaTitle, row, onClose }: Props = $props();
	let open = $state(true);
	let selectedAudioTrack = $state('');
	let playbackError = $state(false);
	let infoOpen = $state(false);
	let infoLoading = $state(false);
	let infoError = $state('');
	let previewInfo = $state<MediaFilePreviewInfo>();
	let previewStartTime = $state(0);
	let playbackStats = $state<MediaFilePlaybackStats>({
		playing: false,
		variableBitRate: false
	});
	const fileName = $derived(row.relativePath.split(/[\\/]/).filter(Boolean).pop() ?? 'media');
	const vlcPlaylistName = $derived(playlistName(fileName));
	const vlcUrl = $derived(row.path ? mediaFileVlcUrl(mediaItemId, row.path) : '');
	const metadataAudioTracks = $derived(metadataAudioTrackOptions(row));
	const activeAudioTrackKey = $derived(selectedAudioTrack || metadataAudioTracks[0]?.key || '');
	const selectedAudio = $derived(
		metadataAudioTracks.find((track) => track.key === activeAudioTrackKey) ?? metadataAudioTracks[0]
	);
	const restartOnSeek = $derived(
		previewInfo?.streamingMode === 'remux' || previewInfo?.streamingMode === 'transcode'
	);
	const previewUrl = $derived(
		row.path
			? mediaFilePreviewUrl(mediaItemId, row.path, selectedAudio?.streamIndex, previewStartTime)
			: ''
	);
	const textTracks = $derived(mediaFileTextTracks(mediaItemId, row));

	$effect(() => {
		if (!row.path) {
			previewInfo = undefined;
			previewStartTime = 0;
			return;
		}
		const audioTrackIndex = selectedAudio?.streamIndex;
		const controller = new globalThis.AbortController();
		infoLoading = true;
		infoError = '';
		void globalThis
			.fetch(mediaFilePreviewInfoUrl(mediaItemId, row.path, audioTrackIndex), {
				headers: { Accept: 'application/json' },
				signal: controller.signal
			})
			.then((response) => {
				if (!response.ok) throw new Error('Preview info request failed');
				return response.json() as Promise<MediaFilePreviewInfo>;
			})
			.then((nextInfo) => {
				previewInfo = nextInfo;
				infoLoading = false;
			})
			.catch((error: unknown) => {
				if (error instanceof globalThis.DOMException && error.name === 'AbortError') return;
				infoError = 'Could not load media info.';
				infoLoading = false;
			});
		return () => controller.abort();
	});

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) onClose();
	}

	function selectAudioTrack(key: string) {
		selectedAudioTrack = key;
		previewStartTime = 0;
		playbackError = false;
	}

	function resetPlaybackError() {
		playbackError = false;
	}

	function updatePlaybackStats(stats: MediaFilePlaybackStats) {
		playbackStats = stats;
	}

	function restartPreviewAt(timeSeconds: number) {
		const target = boundedStartTime(timeSeconds, previewInfo?.durationSeconds);
		if (Math.abs(target - previewStartTime) < 0.5) return;
		previewStartTime = target;
		playbackError = false;
	}

	function playlistName(name: string) {
		const base = name.replace(/\.[^.]+$/, '').trim() || 'media-stream';
		return `${base}.m3u`;
	}

	function boundedStartTime(value: number, duration?: number) {
		if (!Number.isFinite(value) || value <= 0) return 0;
		if (!duration || !Number.isFinite(duration)) return value;
		return Math.max(0, Math.min(value, Math.max(0, duration - 1)));
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(1180px,calc(100vw-32px))] gap-4 p-4 sm:max-w-none">
		<Dialog.Header class="pr-10">
			<div class="flex flex-wrap items-start justify-between gap-3">
				<div class="grid min-w-0 gap-1">
					<Dialog.Title class="break-anywhere text-xl">{mediaTitle}</Dialog.Title>
					<Dialog.Description class="break-anywhere">{fileName}</Dialog.Description>
				</div>
				<div class="flex shrink-0 flex-wrap items-center gap-2">
					<Button
						type="button"
						variant={infoOpen ? 'secondary' : 'outline'}
						size="sm"
						aria-expanded={infoOpen}
						onclick={() => (infoOpen = !infoOpen)}
					>
						<InfoIcon aria-hidden="true" />
						Media info
					</Button>
					{#if row.path}
						<Button href={vlcUrl} download={vlcPlaylistName} variant="outline" size="sm">
							<ExternalLinkIcon aria-hidden="true" />
							Play in VLC
						</Button>
					{/if}
				</div>
			</div>
		</Dialog.Header>
		<div class="grid content-start gap-2">
			<div class="relative aspect-video overflow-hidden rounded-md bg-black">
				{#if previewUrl}
					{#key previewUrl}
						<MediaFileVideoPlayer
							src={previewUrl}
							durationSeconds={previewInfo?.durationSeconds}
							{textTracks}
							audioTracks={metadataAudioTracks}
							{activeAudioTrackKey}
							{restartOnSeek}
							sourceStartTime={previewStartTime}
							onAudioTrackChange={selectAudioTrack}
							onSeekRequest={restartPreviewAt}
							onPlaybackStatsChange={updatePlaybackStats}
							onLoaded={resetPlaybackError}
							onError={() => (playbackError = true)}
						/>
					{/key}
				{/if}
				{#if infoOpen}
					<div class="absolute top-3 left-3 z-20 max-h-[calc(100%-1.5rem)] overflow-auto">
						<MediaFileInfoPanel
							info={previewInfo}
							{playbackStats}
							loading={infoLoading}
							error={infoError}
						/>
					</div>
				{/if}
			</div>
			{#if playbackError}
				<p
					class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
				>
					Preview could not start. Use VLC to stream the original file.
				</p>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
