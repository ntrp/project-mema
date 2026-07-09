<script lang="ts">
	import DiscoverBlacklistArea from '$lib/components/app/discovery/DiscoverBlacklistArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<DiscoverBlacklistArea
	items={app.discoverBlacklist}
	loading={app.loadingBlacklist}
	removingId={app.removingBlacklistId}
	onRemove={app.removeDiscoverBlacklistItem}
/>
