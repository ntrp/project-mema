<script lang="ts">
	import MediaCollectionArea from '$lib/components/app/media/collection/MediaCollectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<MediaCollectionArea
	collection={app.mediaCollection}
	mediaItems={library.data ?? []}
	loading={app.loadingMediaCollection}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
