<script lang="ts">
	import CastIcon from '@lucide/svelte/icons/cast';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { Renderer } from './types';
	import { listRenderers, playOnRenderer } from './api';
	import RendererRemote from './RendererRemote.svelte';

	interface Props {
		mediaItemId: string;
		path: string;
	}

	let { mediaItemId, path }: Props = $props();
	let open = $state(false);
	let renderers = $state.raw<Renderer[]>([]);
	let active = $state.raw<Renderer>();
	let loading = $state(false);
	let playingId = $state('');
	let error = $state('');

	async function loadRenderers() {
		loading = true;
		error = '';
		try {
			renderers = await listRenderers();
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not discover renderers.';
		} finally {
			loading = false;
		}
	}

	function handleOpenChange(next: boolean) {
		open = next;
		if (next) void loadRenderers();
	}

	async function play(renderer: Renderer) {
		playingId = renderer.id;
		error = '';
		try {
			await playOnRenderer(renderer.id, mediaItemId, path);
			active = renderer;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not start playback.';
		} finally {
			playingId = '';
		}
	}
</script>

<Dialog.Root {open} onOpenChange={handleOpenChange}>
	<Dialog.Trigger>
		{#snippet child({ props })}
			<Button {...props} type="button" variant="outline" size="sm">
				<CastIcon aria-hidden="true" />
				Play on device
			</Button>
		{/snippet}
	</Dialog.Trigger>
	<Dialog.Content class="w-[min(520px,calc(100vw-32px))] sm:max-w-none">
		<Dialog.Header>
			<Dialog.Title>Play on device</Dialog.Title>
			<Dialog.Description
				>Choose a renderer discovered on your configured DLNA network.</Dialog.Description
			>
		</Dialog.Header>

		{#if active}
			<RendererRemote renderer={active} />
		{:else}
			<div class="grid gap-2">
				<div class="flex items-center justify-between">
					<p class="text-sm font-medium">Available devices</p>
					<Button
						type="button"
						variant="ghost"
						size="icon-sm"
						aria-label="Refresh devices"
						disabled={loading}
						onclick={() => void loadRenderers()}
					>
						<RefreshCwIcon class={loading ? 'animate-spin' : undefined} aria-hidden="true" />
					</Button>
				</div>
				{#each renderers as renderer (renderer.id)}
					<button
						type="button"
						class="flex items-center justify-between rounded-md border border-border px-3 py-2 text-left hover:bg-muted disabled:opacity-60"
						disabled={playingId !== '' || !renderer.supports.play}
						onclick={() => void play(renderer)}
					>
						<span>
							<span class="block font-medium">{renderer.name}</span>
							<span class="block text-xs text-muted-foreground">
								{renderer.manufacturer}
								{renderer.modelName}
							</span>
						</span>
						<span class="text-xs">{playingId === renderer.id ? 'Starting…' : 'Play'}</span>
					</button>
				{:else}
					{#if !loading}
						<p class="rounded-md border border-dashed p-4 text-sm text-muted-foreground">
							No renderers are online. Check the device and DLNA interfaces, then refresh.
						</p>
					{/if}
				{/each}
			</div>
		{/if}

		{#if error}
			<p class="text-sm text-destructive" role="alert">{error} Check the device and try again.</p>
		{/if}
	</Dialog.Content>
</Dialog.Root>
