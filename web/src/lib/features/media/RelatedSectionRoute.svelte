<script lang="ts">
	import DiscoverSectionArea from '$lib/components/app/discovery/DiscoverSectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import { relatedSectionFromDetail } from '$lib/components/app/shell/controller/discoverFilters';
	import { createMetadataDetailQuery } from './queries.svelte';

	function noop() {}

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const metadata = createMetadataDetailQuery({
		provider: () => app.route?.metadataProvider,
		type: () => app.route?.metadataType,
		id: () => app.route?.metadataExternalId
	});
	const section = $derived(
		relatedSectionFromDetail(metadata.data, app.activeRelatedSectionKind, app.discoverBlacklist)
	);
</script>

<DiscoverSectionArea
	{section}
	mediaItems={library.data ?? []}
	loading={metadata.isFetching}
	loadingMore={false}
	hasMore={false}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	canManage={app.isAdmin}
	blacklist={app.discoverBlacklist}
	onAdd={app.addMedia}
	onBlacklist={app.blacklistDiscoverMedia}
	onLoadMore={noop}
/>
