<script lang="ts">
	import MetadataDetailArea from '$lib/components/app/media/MetadataDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<MetadataDetailArea
	detail={app.metadataDetail}
	loading={app.loadingMetadataDetail}
	mediaItems={app.mediaItems}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
