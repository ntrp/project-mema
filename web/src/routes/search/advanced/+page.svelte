<script lang="ts">
	import { page } from '$app/state';
	import AppShell from '$lib/components/app/shell/AppShell.svelte';
</script>

<AppShell
	initialView="advanced-search"
	initialAdvancedQuery={page.url.searchParams.get('q') ?? ''}
/>
