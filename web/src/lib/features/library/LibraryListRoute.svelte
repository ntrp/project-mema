<script lang="ts">
	import MediaItemList from '$lib/components/app/home/MediaItemList.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import type { MediaItem, MediaType } from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
	}

	let { mediaType }: Props = $props();
	const app = getAppShellContext();
	const items = $derived(app.mediaItems.filter((item: MediaItem) => item.type === mediaType));
</script>

<MediaItemList {mediaType} {items} />
