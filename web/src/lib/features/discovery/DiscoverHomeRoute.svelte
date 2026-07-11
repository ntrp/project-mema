<script lang="ts">
	import MediaSearchPanel from '$lib/components/app/discovery/MediaSearchPanel.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<MediaSearchPanel
	sections={app.discoverSections}
	mediaItems={library.data ?? []}
	loading={app.loadingDiscover}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	canManage={app.isAdmin}
	blacklist={app.discoverBlacklist}
	onAdd={app.addMedia}
	onBlacklist={app.blacklistDiscoverMedia}
/>
