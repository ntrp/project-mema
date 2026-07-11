<script lang="ts">
	import DiscoverSectionArea from '$lib/components/app/discovery/DiscoverSectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<DiscoverSectionArea
	section={app.discoverSection}
	mediaItems={library.data ?? []}
	loading={app.loadingDiscoverSection}
	loadingMore={app.loadingMoreDiscoverSection}
	hasMore={app.discoverSectionHasMore}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	canManage={app.isAdmin}
	blacklist={app.discoverBlacklist}
	onAdd={app.addMedia}
	onBlacklist={app.blacklistDiscoverMedia}
	onLoadMore={app.loadMoreDiscoverSection}
/>
