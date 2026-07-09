<script lang="ts">
	import MediaSearchPanel from '$lib/components/app/discovery/MediaSearchPanel.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<MediaSearchPanel
	sections={app.discoverSections}
	mediaItems={app.mediaItems}
	loading={app.loadingDiscover}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	canManage={app.isAdmin}
	blacklist={app.discoverBlacklist}
	onAdd={app.addMedia}
	onBlacklist={app.blacklistDiscoverMedia}
/>
