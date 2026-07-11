<script lang="ts">
	import MediaItemList from '$lib/components/app/home/MediaItemList.svelte';
	import type { MediaItem, MediaType } from '$lib/settings/types';
	import { createMediaItemsQuery } from './queries.svelte';

	interface Props {
		mediaType: MediaType;
	}

	let { mediaType }: Props = $props();
	const library = createMediaItemsQuery();
	const items = $derived((library.data ?? []).filter((item: MediaItem) => item.type === mediaType));
</script>

<MediaItemList {mediaType} {items} />
