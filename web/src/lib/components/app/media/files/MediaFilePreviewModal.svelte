<script lang="ts">
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import MusicIcon from '@lucide/svelte/icons/music';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import {
		chapterSeconds,
		formatPlaybackTime,
		mediaFilePreviewUrl,
		mediaFileVlcUrl,
		metadataAudioTrackOptions,
		playlistDownloadName
	} from '$lib/components/app/media/files/mediaFilePlayback';

	interface Props {
		mediaItemId: string;
		row: MediaFileRow;
		onClose: () => void;
	}
	interface ChapterOption {
		key: string;
		title: string;
		time: string;
		seconds: number;
	}

	let { mediaItemId, row, onClose }: Props = $props();
	let open = $state(true);
	let videoElement = $state<globalThis.HTMLVideoElement>();
	let selectedAudioTrack = $state('');
	let playbackError = $state(false);
	const fileName = $derived(row.relativePath.split(/[\\/]/).filter(Boolean).pop() ?? 'media');
	const vlcUrl = $derived(row.path ? mediaFileVlcUrl(mediaItemId, row.path) : '');
	const metadataAudioTracks = $derived(metadataAudioTrackOptions(row));
	const activeAudioTrackKey = $derived(selectedAudioTrack || metadataAudioTracks[0]?.key || '');
	const selectedAudio = $derived(
		metadataAudioTracks.find((track) => track.key === activeAudioTrackKey) ?? metadataAudioTracks[0]
	);
	const previewUrl = $derived(
		row.path ? mediaFilePreviewUrl(mediaItemId, row.path, selectedAudio?.streamIndex) : ''
	);
	const canSwitchAudio = $derived(metadataAudioTracks.length > 1);
	const selectedAudioLabel = $derived(
		selectedAudio?.label ?? metadataAudioTracks[0]?.label ?? 'Audio track'
	);
	const chapters = $derived(
		row.chapters
			.map((chapter, index): ChapterOption | undefined => {
				const seconds = chapterSeconds(chapter.startTime);
				if (seconds === undefined) return undefined;
				return {
					key: `chapter-${chapter.index}-${index}`,
					title: chapter.title?.trim() || `Chapter ${chapter.index + 1 || index + 1}`,
					time: formatPlaybackTime(seconds),
					seconds
				};
			})
			.filter((chapter): chapter is ChapterOption => Boolean(chapter))
	);

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) onClose();
	}

	function selectAudioTrack(key: string) {
		selectedAudioTrack = key;
		playbackError = false;
	}

	function seekChapter(seconds: number) {
		if (!videoElement) return;
		videoElement.currentTime = seconds;
		void videoElement.play().catch(() => {});
	}

	function resetPlaybackError() {
		playbackError = false;
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(1180px,calc(100vw-32px))] gap-4 p-4 sm:max-w-none">
		<Dialog.Header class="pr-10">
			<div class="flex flex-wrap items-start justify-between gap-3">
				<div class="grid min-w-0 gap-1">
					<Dialog.Title class="break-anywhere text-xl">{fileName}</Dialog.Title>
					<Dialog.Description class="break-anywhere">{row.relativePath}</Dialog.Description>
				</div>
				{#if row.path}
					<Button
						href={vlcUrl}
						download={playlistDownloadName(fileName)}
						variant="outline"
						size="sm"
						class="shrink-0"
					>
						<DownloadIcon aria-hidden="true" />
						VLC playlist
					</Button>
				{/if}
			</div>
		</Dialog.Header>
		<div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_300px]">
			<div class="grid content-start gap-2">
				<div class="aspect-video overflow-hidden rounded-md bg-black">
					{#if previewUrl}
						{#key previewUrl}
							<!-- svelte-ignore a11y_media_has_caption -->
							<video
								bind:this={videoElement}
								class="block size-full"
								controls
								preload="metadata"
								src={previewUrl}
								onloadeddata={resetPlaybackError}
								onerror={() => (playbackError = true)}
							></video>
						{/key}
					{/if}
				</div>
				{#if playbackError}
					<p
						class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
					>
						Preview could not start. Use the VLC playlist for the original stream.
					</p>
				{/if}
			</div>
			<div class="grid content-start gap-3">
				<section class="grid gap-2 rounded-md border border-border p-3">
					<h3 class="m-0 inline-flex items-center gap-2 text-sm font-semibold">
						<MusicIcon aria-hidden="true" />
						Audio
					</h3>
					{#if canSwitchAudio}
						<Select.Root type="single" value={activeAudioTrackKey} onValueChange={selectAudioTrack}>
							<Select.Trigger class="w-full">{selectedAudioLabel}</Select.Trigger>
							<Select.Content>
								{#each metadataAudioTracks as track (track.key)}
									<Select.Item value={track.key} label={track.label} />
								{/each}
							</Select.Content>
						</Select.Root>
					{:else if metadataAudioTracks.length > 0}
						<ul class="m-0 grid gap-1 p-0 text-sm text-muted-foreground">
							{#each metadataAudioTracks as track (track.key)}
								<li class="list-none rounded-md bg-muted/40 px-2 py-1.5">{track.label}</li>
							{/each}
						</ul>
					{:else}
						<p class="m-0 text-sm text-muted-foreground">No audio tracks found.</p>
					{/if}
				</section>
				<section class="grid gap-2 rounded-md border border-border p-3">
					<h3 class="m-0 inline-flex items-center gap-2 text-sm font-semibold">
						<ClapperboardIcon aria-hidden="true" />
						Chapters
					</h3>
					<div class="grid max-h-80 gap-1 overflow-y-auto pr-1">
						{#each chapters as chapter (chapter.key)}
							<Button
								type="button"
								variant="ghost"
								size="sm"
								class="h-auto justify-start gap-2 px-2 py-1.5 text-left"
								onclick={() => seekChapter(chapter.seconds)}
							>
								<span class="w-14 shrink-0 font-mono text-xs text-muted-foreground">
									{chapter.time}
								</span>
								<span class="min-w-0 truncate">{chapter.title}</span>
							</Button>
						{:else}
							<p class="m-0 text-sm text-muted-foreground">No chapters found.</p>
						{/each}
					</div>
				</section>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
