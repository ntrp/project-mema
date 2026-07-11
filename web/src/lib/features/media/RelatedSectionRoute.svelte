<script lang="ts">
	import DiscoverSectionArea from '$lib/components/app/discovery/DiscoverSectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	function noop() {}

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<DiscoverSectionArea
	section={app.relatedMediaSection}
	mediaItems={library.data ?? []}
	loading={app.loadingMetadataDetail}
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
