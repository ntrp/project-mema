<script lang="ts">
	import ActivityList from '$lib/components/app/activity/ActivityList.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import type { ActivitySection } from '$lib/settings/types';
	import { onDestroy } from 'svelte';
	import { createActivityMutations } from './mutations.svelte';
	import { createDownloadActivityQuery, createReleaseBlocklistQuery } from './queries.svelte';
	import { connectActivityQueryEvents } from './realtime.svelte';

	interface Props {
		section?: ActivitySection;
	}

	let { section = 'queue' }: Props = $props();
	const app = getAppShellContext();
	const downloads = createDownloadActivityQuery(() => section !== 'blocklist');
	const blocklist = createReleaseBlocklistQuery(() => section === 'blocklist');
	const mutations = createActivityMutations();
	onDestroy(connectActivityQueryEvents());
	const errorText = (error: unknown, fallback: string) =>
		error instanceof Error ? error.message : fallback;

	async function run(action: () => Promise<unknown>, success: string, fallback: string) {
		app.clearNotice();
		try {
			await action();
			app.message = success;
		} catch (error) {
			app.errorMessage = errorText(error, fallback);
			throw error;
		}
	}
</script>

<ActivityList
	{section}
	activities={downloads.data ?? []}
	releaseBlocklist={blocklist.data ?? []}
	loading={section === 'blocklist' ? blocklist.isFetching : downloads.isFetching}
	canManage={app.isAdmin}
	cancellingId={mutations.cancel.isPending ? mutations.cancel.variables : undefined}
	deletingId={mutations.deleteDownload.isPending ? mutations.deleteDownload.variables : undefined}
	deletingBlocklistId={mutations.deleteBlocklistItem.isPending
		? mutations.deleteBlocklistItem.variables
		: undefined}
	clearingReleaseBlocklist={mutations.clearBlocklist.isPending}
	onRefresh={() => (section === 'blocklist' ? blocklist.refetch() : downloads.refetch())}
	onCancel={(activity) =>
		void run(
			() => mutations.cancel.mutateAsync(activity.id),
			'Download activity cancelled',
			'Could not cancel download activity'
		)}
	onDelete={(activity) =>
		void run(
			() => mutations.deleteDownload.mutateAsync(activity.id),
			'Download activity deleted',
			'Could not delete download activity'
		)}
	onDeleteReleaseBlocklistItem={(item) =>
		void run(
			() => mutations.deleteBlocklistItem.mutateAsync(item.id),
			'Release blocklist entry removed',
			'Could not remove release blocklist entry'
		)}
	onClearReleaseBlocklist={() =>
		void run(
			() => mutations.clearBlocklist.mutateAsync(),
			'Release blocklist cleared',
			'Could not clear release blocklist'
		)}
	onManualImport={(id, request) =>
		run(
			() => mutations.manualImport.mutateAsync({ id, request }),
			'Download manually imported',
			'Manual import failed'
		).then(() => {})}
/>
