<script lang="ts">
	import DiscoverSectionArea from '$lib/components/app/discovery/DiscoverSectionArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	function noop() {}

	const app = getAppShellContext();
</script>

<DiscoverSectionArea
	section={app.relatedMediaSection}
	mediaItems={app.mediaItems}
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
