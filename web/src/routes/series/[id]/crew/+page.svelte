<script lang="ts">
	import AppShell from '$lib/components/app/shell/AppShell.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<AppShell
	initialView="media-people"
	initialHomeSection="series"
	initialSelectedMediaItemId={params.id}
	initialPeopleSectionKind="crew"
/>
