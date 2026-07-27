<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import { createQuery } from '@tanstack/svelte-query';
	import InfoIcon from '@lucide/svelte/icons/info';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import MediaFileInfoPanel from '$lib/components/app/media/files/preview/MediaFileInfoPanel.svelte';
	import PlayOnDevice from '$lib/components/app/renderer/PlayOnDevice.svelte';
	import MediaFilePlaybackError from '$lib/components/app/media/files/preview/MediaFilePlaybackError.svelte';
	import MediaFileVideoPlayer from '$lib/components/app/media/files/preview/MediaFileVideoPlayer.svelte';
	import {
		mediaFilePreviewClientProfile,
		mediaFilePreviewInfoUrl,
		mediaFilePreviewSourceType,
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
	let playbackError = $state('');
	let infoOpen = $state(false);
	let playbackStats = $state<MediaFilePlaybackStats>({
		playing: false,
		variableBitRate: false
	});
	const fileName = $derived(row.relativePath.split(/[\\/]/).filter(Boolean).pop() ?? 'media');
	const clientProfile = mediaFilePreviewClientProfile();
	const vlcUrl = $derived(row.fileId ? mediaFileVlcUrl(mediaItemId, row.fileId) : '');
	const metadataAudioTracks = $derived(metadataAudioTrackOptions(row));
	const activeAudioTrackKey = $derived(selectedAudioTrack || metadataAudioTracks[0]?.key || '');
	const selectedAudio = $derived(
		metadataAudioTracks.find((track) => track.key === activeAudioTrackKey) ?? metadataAudioTracks[0]
	);
	const previewInfoQuery = createQuery(() => ({
		queryKey: [
			'media-file-preview-info',
			mediaItemId,
			row.path,
			selectedAudio?.streamIndex,
			clientProfile
		],
		enabled: Boolean(row.path),
		queryFn: async ({ signal }) => {
			const response = await globalThis.fetch(
				mediaFilePreviewInfoUrl(
					mediaItemId,
					row.path ?? '',
					selectedAudio?.streamIndex,
					clientProfile
				),
				{ headers: { Accept: 'application/json' }, signal }
			);
			if (!response.ok) throw new Error('Preview info request failed');
			return response.json() as Promise<MediaFilePreviewInfo>;
		}
	}));
	const previewInfo = $derived(previewInfoQuery.data);
	const infoLoading = $derived(previewInfoQuery.isFetching);
	const infoError = $derived(previewInfoQuery.isError ? 'Could not load media info.' : '');
	const previewUrl = $derived(
		row.path
			? mediaFilePreviewUrl(mediaItemId, row.path, selectedAudio?.streamIndex, clientProfile)
			: ''
	);
	const previewSourceType = $derived(mediaFilePreviewSourceType(previewInfo?.deliveryProtocol));
	const previewPlayerKey = $derived(
		`${previewUrl}|${previewSourceType}|${previewInfo?.durationSeconds ?? ''}`
	);
	const textTracks = $derived(mediaFileTextTracks(mediaItemId, row));

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) onClose();
	}

	function selectAudioTrack(key: string) {
		selectedAudioTrack = key;
		playbackError = '';
	}

	function resetPlaybackError() {
		playbackError = '';
	}

	function updatePlaybackStats(stats: MediaFilePlaybackStats) {
		playbackStats = stats;
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
						<PlayOnDevice {mediaItemId} path={row.path} />
					{/if}
					{#if vlcUrl}
						<Button href={vlcUrl} variant="outline" size="sm">
							<ExternalLinkIcon aria-hidden="true" />
							Play in VLC
						</Button>
					{/if}
				</div>
			</div>
		</Dialog.Header>
		<div class="grid content-start gap-2">
			<div class="relative aspect-video overflow-hidden rounded-md bg-black">
				{#if previewUrl && previewInfo}
					{#key previewPlayerKey}
						<MediaFileVideoPlayer
							src={previewUrl}
							sourceType={previewSourceType}
							durationSeconds={previewInfo?.durationSeconds}
							{textTracks}
							audioTracks={metadataAudioTracks}
							{activeAudioTrackKey}
							restartOnSeek={false}
							sourceStartTime={0}
							onAudioTrackChange={selectAudioTrack}
							onSeekRequest={() => undefined}
							onPlaybackStatsChange={updatePlaybackStats}
							onLoaded={resetPlaybackError}
							onError={(message) => (playbackError = message)}
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
				<MediaFilePlaybackError message={playbackError} />
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
