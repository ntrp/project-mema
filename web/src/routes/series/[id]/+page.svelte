<script lang="ts">
	import MediaItemDetailRoute from '$lib/features/library/MediaItemDetailRoute.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<MediaItemDetailRoute mediaType="serie" id={params.id} />
