<script lang="ts">
	import AppShell from '$lib/components/app/AppShell.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<AppShell
	initialView="metadata-detail"
	initialMetadataProvider={params.provider}
	initialMetadataType={params.type}
	initialMetadataExternalId={params.externalId}
/>
