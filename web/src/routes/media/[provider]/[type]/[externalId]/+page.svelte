<script lang="ts">
	import AppShell from '$lib/components/app/shell/AppShell.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

{#key `${params.provider}:${params.type}:${params.externalId}`}
	<AppShell
		initialView="metadata-detail"
		initialMetadataProvider={params.provider}
		initialMetadataType={params.type}
		initialMetadataExternalId={params.externalId}
	/>
{/key}
