<script lang="ts">
	import MediaProfileEditorRoute from '$lib/features/settings/MediaProfileEditorRoute.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<MediaProfileEditorRoute profileId={params.id} />
