<script lang="ts">
	import { page } from '$app/state';
	import AppShell from '$lib/components/app/AppShell.svelte';
</script>

<AppShell
	initialView="advanced-search"
	initialAdvancedQuery={page.url.searchParams.get('q') ?? ''}
/>
