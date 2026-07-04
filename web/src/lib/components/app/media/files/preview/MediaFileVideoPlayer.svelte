<script lang="ts">
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import 'video.js/dist/video-js.css';
	import type Player from 'video.js/dist/types/player';
	import type {
		AudioTrackOption,
		MediaFileTextTrack
	} from '$lib/components/app/media/files/preview/mediaFilePlayback';
	import { addAudioTracks } from '$lib/components/app/media/files/preview/mediaFileVideoAudioTracks';
	import { addSourceTimeline } from '$lib/components/app/media/files/preview/mediaFileVideoSeeking';
	import {
		emptyPlaybackStats,
		watchPlaybackStats
	} from '$lib/components/app/media/files/preview/mediaFilePlaybackStats';
	import type { MediaFilePlaybackStats } from '$lib/components/app/media/files/preview/mediaFilePreviewInfo';

	interface Props {
		src: string;
		durationSeconds?: number;
		textTracks: MediaFileTextTrack[];
		audioTracks: AudioTrackOption[];
		activeAudioTrackKey: string;
		restartOnSeek: boolean;
		sourceStartTime: number;
		onAudioTrackChange: (_key: string) => void;
		onSeekRequest: (_timeSeconds: number) => void;
		onPlaybackStatsChange: (_stats: MediaFilePlaybackStats) => void;
		onLoaded: () => void;
		onError: () => void;
	}

	let {
		src,
		durationSeconds,
		textTracks,
		audioTracks,
		activeAudioTrackKey,
		restartOnSeek,
		sourceStartTime,
		onAudioTrackChange,
		onSeekRequest,
		onPlaybackStatsChange,
		onLoaded,
		onError
	}: Props = $props();
	let videoElement = $state<globalThis.HTMLVideoElement>();
	let player = $state<Player>();
	let starting = $state(true);

	$effect(() => {
		if (!videoElement || !src) return;
		const source = src;
		const currentAudioTracks = audioTracks;
		const currentAudioTrackKey = activeAudioTrackKey;
		let disposed = false;
		let instance: Player | undefined;
		let removeAudioTrackListener: (() => void) | undefined;
		let removePlaybackStatsWatcher: (() => void) | undefined;
		let removeSeekHandler: (() => void) | undefined;
		onPlaybackStatsChange(emptyPlaybackStats());
		starting = true;

		void import('video.js').then(({ default: videojs }) => {
			if (disposed || !videoElement) return;
			const currentInstance = videojs(videoElement, {
				autoplay: 'play',
				controls: true,
				fill: true,
				inactivityTimeout: 0,
				preload: 'auto',
				responsive: true,
				sources: [{ src: source, type: 'video/mp4' }]
			});
			instance = currentInstance;
			player = currentInstance;
			removeAudioTrackListener = addAudioTracks(
				videojs,
				currentInstance,
				currentAudioTracks,
				currentAudioTrackKey,
				onAudioTrackChange
			);
			removeSeekHandler = addSourceTimeline(
				currentInstance,
				restartOnSeek,
				sourceStartTime,
				durationSeconds,
				onSeekRequest
			);
			removePlaybackStatsWatcher = watchPlaybackStats(
				currentInstance,
				videoElement,
				source,
				onPlaybackStatsChange
			);
			currentInstance.on('loadeddata', handleLoaded);
			currentInstance.on('error', onError);
			currentInstance.on('canplay', hideStarting);
			currentInstance.on('playing', hideStarting);
			currentInstance.ready(() => startPlayback(currentInstance));
		});

		return () => {
			disposed = true;
			if (instance) {
				instance.off('loadeddata', handleLoaded);
				instance.off('error', onError);
				instance.off('canplay', hideStarting);
				instance.off('playing', hideStarting);
				removeAudioTrackListener?.();
				removePlaybackStatsWatcher?.();
				removeSeekHandler?.();
				instance.dispose();
				if (player === instance) player = undefined;
			}
		};
	});

	$effect(() => {
		const instance = player;
		const sourceDuration = durationSeconds;
		if (!instance || !validDuration(sourceDuration)) return;
		const applyDuration = () => {
			const currentDuration = instance.duration();
			if (!validDuration(currentDuration) || Math.abs(currentDuration - sourceDuration) > 0.5) {
				instance.duration(sourceDuration);
			}
		};
		applyDuration();
		instance.on('loadedmetadata', applyDuration);
		instance.on('durationchange', applyDuration);
		return () => {
			instance.off('loadedmetadata', applyDuration);
			instance.off('durationchange', applyDuration);
		};
	});

	function validDuration(value: number | undefined): value is number {
		return typeof value === 'number' && Number.isFinite(value) && value > 0;
	}

	function startPlayback(instance: Player) {
		void instance.play()?.catch(() => undefined);
	}

	function handleLoaded() {
		hideStarting();
		onLoaded();
	}

	function hideStarting() {
		starting = false;
	}
</script>

<div class="relative size-full">
	<video
		bind:this={videoElement}
		class="video-js vjs-big-play-centered size-full"
		autoplay
		controls
		playsinline
	>
		{#each textTracks as track (track.key)}
			<track
				kind={track.kind}
				label={track.label}
				src={track.src}
				srclang={track.srclang}
				default={track.default}
			/>
		{/each}
	</video>
	{#if starting}
		<div
			class="pointer-events-none absolute inset-0 z-10 grid place-items-center bg-black/20 text-white"
			aria-live="polite"
			aria-label="Loading preview"
		>
			<LoaderCircleIcon class="size-10 animate-spin drop-shadow" aria-hidden="true" />
		</div>
	{/if}
</div>

<style>
	:global(.video-js) {
		width: 100%;
		height: 100%;
	}

	:global(.video-js .vjs-menu-button-popup .vjs-menu .vjs-menu-content) {
		max-height: 18rem;
	}
</style>
