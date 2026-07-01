<script lang="ts">
	import type { Snippet } from 'svelte';
	import { imageUrl } from './mediaDetail';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { MediaMetadataEpisode } from '$lib/settings/types';

	interface Props {
		episode: MediaMetadataEpisode;
		title: string;
		beforeTitle?: Snippet;
		children?: Snippet;
	}

	let { episode, title, beforeTitle, children }: Props = $props();

	const stillUrl = $derived(imageUrl(episode.stillPath, 'w300'));
</script>

<article
	class="grid min-h-[118px] items-start gap-5.5 border-t border-border py-4.5 first:border-t-0 md:grid-cols-[minmax(0,1fr)_244px] max-[980px]:md:grid-cols-[minmax(0,1fr)_200px] max-sm:grid-cols-1"
>
	<div class="grid min-w-0 gap-2.5">
		<h3 class="m-0 flex flex-wrap items-center gap-2 text-xl text-foreground">
			{@render beforeTitle?.()}
			<span class="min-w-0">{title}</span>
			{#if episode.airDate}
				<span class="rounded-md bg-muted px-2 py-0.5 text-xs font-black text-muted-foreground">
					{formatDate(episode.airDate)}
				</span>
			{/if}
		</h3>
		<p class="m-0 text-sm leading-6 text-muted-foreground">
			{episode.overview ?? 'No episode overview available.'}
		</p>
	</div>
	{#if stillUrl}
		<img
			class="aspect-video w-full rounded-md object-cover md:w-[244px] max-[980px]:md:w-[200px] max-sm:max-w-80"
			src={stillUrl}
			alt=""
			loading="lazy"
		/>
	{/if}
	{@render children?.()}
</article>
