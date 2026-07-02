<script lang="ts">
	import EditMediaProfileRoute from '$lib/features/settings/routes/profiles/EditMediaProfileRoute.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<EditMediaProfileRoute profileId={params.id} />
