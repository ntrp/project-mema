<script lang="ts">
	import MediaCollectionArea from '$lib/components/app/media/collection/MediaCollectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<MediaCollectionArea
	collection={app.mediaCollection}
	mediaItems={app.mediaItems}
	loading={app.loadingMediaCollection}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
