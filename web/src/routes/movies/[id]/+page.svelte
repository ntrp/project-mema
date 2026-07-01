<script lang="ts">
	import HomeRoute from '$lib/features/home/HomeRoute.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<HomeRoute section="movies" selectedMediaItemId={params.id} />
