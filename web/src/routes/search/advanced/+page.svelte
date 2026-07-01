<script lang="ts">
	import AdvancedSearchRoute from '$lib/features/search/AdvancedSearchRoute.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<AdvancedSearchRoute initialQuery={data.q} />
