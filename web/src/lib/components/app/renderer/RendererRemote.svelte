<script lang="ts">
	import PauseIcon from '@lucide/svelte/icons/pause';
	import { onMount } from 'svelte';
	import PlayIcon from '@lucide/svelte/icons/play';
	import SquareIcon from '@lucide/svelte/icons/square';
	import Volume2Icon from '@lucide/svelte/icons/volume-2';
	import VolumeXIcon from '@lucide/svelte/icons/volume-x';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { Renderer, RendererStatus, RendererTransportRequest } from './types';
	import {
		rendererStatus,
		rendererTransport,
		seekRenderer,
		setRendererMute,
		setRendererVolume
	} from './api';

	interface Props {
		renderer: Renderer;
	}

	let { renderer }: Props = $props();
	let status = $state.raw<RendererStatus>();
	let error = $state('');
	let busy = $state(false);
	let refreshing = false;

	async function refresh() {
		if (refreshing) return;
		refreshing = true;
		try {
			status = await rendererStatus(renderer.id);
			error = '';
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Renderer is offline.';
		} finally {
			refreshing = false;
		}
	}

	onMount(() => {
		void refresh();
		const timer = globalThis.setInterval(() => void refresh(), 2000);
		return () => globalThis.clearInterval(timer);
	});

	async function command(run: () => Promise<void>) {
		busy = true;
		error = '';
		try {
			await run();
			await refresh();
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Renderer action failed.';
		} finally {
			busy = false;
		}
	}

	function transport(action: RendererTransportRequest['action']) {
		void command(() => rendererTransport(renderer.id, action));
	}

	function sliderValue(event: Event) {
		return Number((event.currentTarget as HTMLInputElement).value);
	}
</script>

<div class="grid gap-4 rounded-md border border-border p-3" aria-label={`${renderer.name} remote`}>
	<div class="flex items-center justify-between gap-3">
		<div>
			<p class="font-semibold">{renderer.name}</p>
			<p class="text-xs text-muted-foreground">
				{status?.transportState?.replaceAll('_', ' ') ?? 'Connecting…'}
			</p>
		</div>
		<div class="flex gap-1">
			{#each [{ action: 'play', label: 'Play', icon: PlayIcon, supported: renderer.supports.play }, { action: 'pause', label: 'Pause', icon: PauseIcon, supported: renderer.supports.pause }, { action: 'stop', label: 'Stop', icon: SquareIcon, supported: renderer.supports.stop }] as control (control.action)}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								{...props}
								type="button"
								size="icon-sm"
								variant="outline"
								aria-label={control.label}
								disabled={busy || !control.supported}
								onclick={() => transport(control.action as RendererTransportRequest['action'])}
							>
								<control.icon aria-hidden="true" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content
						>{control.supported
							? control.label
							: `${control.label} is not supported`}</Tooltip.Content
					>
				</Tooltip.Root>
			{/each}
		</div>
	</div>

	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<label {...props} class="grid gap-1 text-xs font-medium">
					<span>Position {status?.positionSeconds ?? 0}s / {status?.durationSeconds ?? 0}s</span>
					<input
						type="range"
						min="0"
						max={Math.max(status?.durationSeconds ?? 0, 1)}
						value={status?.positionSeconds ?? 0}
						disabled={busy || !renderer.supports.seek || !status?.durationSeconds}
						onchange={(event) => void command(() => seekRenderer(renderer.id, sliderValue(event)))}
					/>
				</label>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>
			{renderer.supports.seek ? 'Seek playback' : 'Seek is not supported'}
		</Tooltip.Content>
	</Tooltip.Root>

	<div class="flex items-end gap-2">
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<label {...props} class="grid grow gap-1 text-xs font-medium">
						<span>Volume {status?.volume ?? 0}%</span>
						<input
							type="range"
							min="0"
							max="100"
							value={status?.volume ?? 0}
							disabled={busy || !renderer.supports.volume}
							onchange={(event) =>
								void command(() => setRendererVolume(renderer.id, sliderValue(event)))}
						/>
					</label>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>
				{renderer.supports.volume ? 'Set volume' : 'Volume is not supported'}
			</Tooltip.Content>
		</Tooltip.Root>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						size="icon-sm"
						variant="outline"
						aria-label={status?.muted ? 'Unmute' : 'Mute'}
						disabled={busy || !renderer.supports.mute}
						onclick={() => void command(() => setRendererMute(renderer.id, !status?.muted))}
					>
						{#if status?.muted}<VolumeXIcon aria-hidden="true" />{:else}<Volume2Icon
								aria-hidden="true"
							/>{/if}
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>
				{renderer.supports.mute ? 'Toggle mute' : 'Mute is not supported'}
			</Tooltip.Content>
		</Tooltip.Root>
	</div>

	{#if error}
		<p class="text-sm text-destructive" role="alert">
			{error} Try again after checking the device.
		</p>
	{/if}
</div>
