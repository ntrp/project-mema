<script lang="ts">
	import MediaCollectionArea from '$lib/components/app/media/collection/MediaCollectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import { createMediaCollectionQuery } from './queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const collection = createMediaCollectionQuery(
		() => app.route?.collectionProvider,
		() => app.route?.collectionId
	);
</script>

<MediaCollectionArea
	collection={collection.data}
	mediaItems={library.data ?? []}
	loading={collection.isFetching}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
