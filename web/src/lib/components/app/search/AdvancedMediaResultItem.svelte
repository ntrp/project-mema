<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import MediaAddButton from '$lib/components/app/media/shared/MediaAddButton.svelte';
	import PosterPlaceholder from '$lib/components/app/media/posters/PosterPlaceholder.svelte';
	import { isUnreleasedMedia } from '$lib/components/app/media/shared/mediaRelease';
	import { providerDisplayName } from '$lib/settings/providerLinks';
	import type { MediaSearchResult } from '$lib/settings/types';
	import {
		externalMediaUrl,
		imageUrl,
		mediaCandidateKey,
		mediaHref
	} from './advancedSearchResults';

	interface Props {
		result: MediaSearchResult;
		inLibrary: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { result, inLibrary, addingKey, actionLabel, onAdd }: Props = $props();
	const href = $derived(mediaHref(result));
	const externalUrl = $derived(externalMediaUrl(result));
	const externalLabel = $derived(providerDisplayName(result.externalProvider));
	const isUnreleased = $derived(isUnreleasedMedia(result));
</script>

<article
	class={`${isUnreleased ? 'border-yellow-400 ' : 'border-border '}relative grid items-center gap-3.5 rounded-md border bg-muted p-2.5 transition-colors hover:bg-accent/20 md:grid-cols-[82px_minmax(0,1fr)_auto]`}
>
	{#if href}
		<a class="absolute inset-0 z-0 rounded-md" {href} aria-label={`Open ${result.title}`}></a>
	{/if}
	<div
		class={`${href ? 'pointer-events-none ' : ''}${isUnreleased ? 'border-yellow-400 ' : 'border-border '}relative z-10 aspect-[2/3] overflow-hidden rounded-md border bg-card`}
	>
		{#if imageUrl(result.posterPath)}
			<img
				class="block h-full w-full object-cover"
				src={imageUrl(result.posterPath)}
				alt=""
				loading="lazy"
			/>
		{:else}
			<PosterPlaceholder label={result.title} class="h-full min-h-0" />
		{/if}
	</div>
	<div class={`${href ? 'pointer-events-none ' : ''}relative z-10 grid min-w-0 gap-2`}>
		<div>
			<h3 class="m-0 text-base leading-tight">
				{result.title}
			</h3>
			<p class="m-0 text-sm text-muted-foreground">
				{result.type}{result.year ? ` · ${result.year}` : ''}
			</p>
		</div>
		{#if result.overview}
			<p class="line-clamp-2 m-0 text-sm text-muted-foreground">{result.overview}</p>
		{/if}
	</div>
	<div class="relative z-20 flex items-center justify-end gap-2.5">
		{#if externalUrl}
			<Button
				variant="outline"
				size="sm"
				href={externalUrl}
				target="_blank"
				rel="noreferrer"
				aria-label={`Open ${externalLabel} page in a new tab`}
			>
				<ExternalLinkIcon aria-hidden="true" />
				<span>{externalLabel}</span>
			</Button>
		{/if}
		{#if inLibrary}
			<StatusPill tone="success">In library</StatusPill>
		{:else}
			<MediaAddButton
				{result}
				adding={addingKey === mediaCandidateKey(result)}
				label={actionLabel}
				{onAdd}
			/>
		{/if}
	</div>
</article>
