<script lang="ts">
	import AppShell from '$lib/components/app/shell/AppShell.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<AppShell
	initialView="related-section"
	initialMetadataProvider={params.provider}
	initialMetadataType={params.type}
	initialMetadataExternalId={params.externalId}
	initialRelatedSectionKind="similar"
/>
