<script lang="ts">
	import WantedMediaTable from '$lib/components/app/home/WantedMediaTable.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from './queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<WantedMediaTable
	items={library.data ?? []}
	languages={app.languages}
	searchingItemId={app.searchingItemId}
	grabbingKey={app.grabbingKey}
	canManage={app.isAdmin}
	onFindReleases={app.findReleases}
	onGrabRelease={app.grabRelease}
/>
