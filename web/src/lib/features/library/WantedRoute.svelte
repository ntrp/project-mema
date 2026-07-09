<script lang="ts">
	import WantedMediaTable from '$lib/components/app/home/WantedMediaTable.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<WantedMediaTable
	items={app.mediaItems}
	languages={app.languages}
	searchingItemId={app.searchingItemId}
	grabbingKey={app.grabbingKey}
	canManage={app.isAdmin}
	onFindReleases={app.findReleases}
	onGrabRelease={app.grabRelease}
/>
