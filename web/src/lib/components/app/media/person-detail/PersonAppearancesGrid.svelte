<script lang="ts">
	import MediaPosterCard from '$lib/components/app/media/posters/MediaPosterCard.svelte';
	import type { MediaItem, MediaSearchResult, PersonAppearance } from '$lib/settings/types';
	import { appearanceResult, resultKey } from './personDetail';

	interface Props {
		appearances: PersonAppearance[];
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { appearances, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();

	const libraryKeys = $derived(new Set(mediaItems.map((item) => resultKey(item))));

	function inLibrary(result: MediaSearchResult) {
		return libraryKeys.has(resultKey(result));
	}
</script>

<div
	class="grid grid-cols-[repeat(auto-fill,minmax(164px,1fr))] gap-x-5 gap-y-6 max-sm:grid-cols-2"
>
	{#each appearances as appearance (`${appearance.type}:${appearance.externalProvider}:${appearance.externalId}`)}
		{@const result = appearanceResult(appearance)}
		<div class="grid min-w-0 gap-2">
			<MediaPosterCard
				{result}
				adding={addingKey === resultKey(result)}
				{actionLabel}
				inLibrary={inLibrary(result)}
				{onAdd}
			/>
			<p class="m-0 min-h-5 truncate text-center text-sm font-medium text-muted-foreground">
				{#if appearance.role}
					as {appearance.role}
				{:else}
					&nbsp;
				{/if}
			</p>
		</div>
	{/each}
</div>
