<script lang="ts">
	import { resolve } from '$app/paths';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import type { MediaRequest } from '$lib/settings/types';
	import PosterPlaceholder from '../media/PosterPlaceholder.svelte';

	interface Props {
		request: MediaRequest;
	}

	let { request }: Props = $props();

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w185${path}`;
	}

	function statusTone(status: MediaRequest['status']) {
		if (status === 'approved') {
			return 'success';
		}
		if (status === 'pending') {
			return 'pending';
		}
		return 'muted';
	}
</script>

<a
	class="grid items-stretch gap-3.5 rounded-md border border-border bg-muted p-2.5 text-foreground no-underline hover:border-primary focus-visible:border-primary focus-visible:outline-none md:grid-cols-[82px_minmax(0,1fr)_auto]"
	href={resolve('/requests/[id]', { id: request.id })}
>
	<div class="aspect-[2/3] overflow-hidden rounded-md bg-card">
		{#if posterUrl(request.posterPath)}
			<img
				class="block h-full w-full object-cover"
				src={posterUrl(request.posterPath)}
				alt=""
				loading="lazy"
			/>
		{:else}
			<PosterPlaceholder label={request.type} class="h-full min-h-0" />
		{/if}
	</div>
	<div class="grid min-w-0 gap-2">
		<h3 class="m-0 text-base leading-tight">{request.title}</h3>
		<p class="m-0 text-sm text-muted-foreground">
			{request.type}{request.year ? ` · ${request.year}` : ''} · Requested by {request.requestedByUsername}
		</p>
		{#if request.overview}
			<p class="line-clamp-2 m-0 text-sm text-muted-foreground">{request.overview}</p>
		{/if}
		{#if request.tags?.length}
			<div class="flex flex-wrap gap-1.5" aria-label="Tags">
				{#each request.tags.slice(0, 3) as tag (tag)}
					<StatusPill class="max-w-24 truncate">{tag}</StatusPill>
				{/each}
			</div>
		{/if}
	</div>
	<StatusPill tone={statusTone(request.status)}>{request.status}</StatusPill>
</a>
