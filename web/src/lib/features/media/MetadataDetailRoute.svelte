<script lang="ts">
	import MetadataDetailArea from '$lib/components/app/media/metadata/MetadataDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<MetadataDetailArea
	detail={app.metadataDetail}
	loading={app.loadingMetadataDetail}
	mediaItems={library.data ?? []}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
