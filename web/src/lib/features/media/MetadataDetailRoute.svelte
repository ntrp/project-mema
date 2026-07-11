<script lang="ts">
	import MetadataDetailArea from '$lib/components/app/media/metadata/MetadataDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import { createMetadataDetailQuery } from './queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const detail = createMetadataDetailQuery({
		provider: () => app.route?.metadataProvider,
		type: () => app.route?.metadataType,
		id: () => app.route?.metadataExternalId
	});
</script>

<MetadataDetailArea
	detail={detail.data}
	loading={detail.isFetching}
	mediaItems={library.data ?? []}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
