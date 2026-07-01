<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import AppShell from '$lib/components/app/shell/AppShell.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
	<link
		rel="stylesheet"
		href="https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&display=swap"
	/>
</svelte:head>

<AppShell>
	{@render children()}
</AppShell>
