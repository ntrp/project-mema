<script lang="ts">
	import ActivityList from '$lib/components/app/activity/ActivityList.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import type { ActivitySection } from '$lib/settings/types';

	interface Props {
		section?: ActivitySection;
	}

	let { section = 'queue' }: Props = $props();
	const app = getAppShellContext();
</script>

<ActivityList
	{section}
	activities={app.activities}
	releaseBlocklist={app.releaseBlocklist}
	loading={app.loadingActivity}
	canManage={app.isAdmin}
	cancellingId={app.cancellingActivityId}
	deletingId={app.deletingActivityId}
	deletingBlocklistId={app.deletingReleaseBlocklistId}
	clearingReleaseBlocklist={app.clearingReleaseBlocklist}
	onRefresh={section === 'blocklist' ? app.loadReleaseBlocklist : app.loadDownloadActivity}
	onCancel={app.cancelActivity}
	onDelete={app.deleteActivity}
	onDeleteReleaseBlocklistItem={app.deleteReleaseBlocklistItem}
	onClearReleaseBlocklist={app.clearReleaseBlocklist}
/>
