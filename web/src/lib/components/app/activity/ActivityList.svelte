<script lang="ts">
	import ActivityManualImportModal from './ActivityManualImportModal.svelte';
	import ActivityBlocklistTable from './ActivityBlocklistTable.svelte';
	import ActivityQueueTable from './queue/ActivityQueueTable.svelte';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		activitySectionDescription,
		activitySectionHeading,
		visibleInActivitySection
	} from './activitySections';
	import type {
		ActivitySection,
		DownloadActivity,
		ManualImportRequest,
		ReleaseBlocklistItem
	} from '$lib/settings/types';

	interface Props {
		section?: ActivitySection;
		activities: DownloadActivity[];
		releaseBlocklist?: ReleaseBlocklistItem[];
		loading: boolean;
		canManage: boolean;
		cancellingId?: string;
		retryingId?: string;
		deletingId?: string;
		deletingBlocklistId?: string;
		clearingReleaseBlocklist?: boolean;
		onRefresh: () => void;
		onCancel: (_activity: DownloadActivity) => void;
		onRetry: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
		onDeleteReleaseBlocklistItem?: (_item: ReleaseBlocklistItem) => void;
		onClearReleaseBlocklist?: () => void;
		onManualImport?: (_id: string, _request: ManualImportRequest) => Promise<void>;
	}

	let {
		section = 'queue',
		activities,
		releaseBlocklist = [],
		loading,
		canManage,
		cancellingId,
		retryingId,
		deletingId,
		deletingBlocklistId,
		clearingReleaseBlocklist = false,
		onRefresh,
		onCancel,
		onRetry,
		onDelete,
		onDeleteReleaseBlocklistItem = () => {},
		onClearReleaseBlocklist = () => {},
		onManualImport = async () => {}
	}: Props = $props();
	let manualImportActivity = $state<DownloadActivity | undefined>();
	let importingId = $state<string | undefined>();
	let importError = $state<string | undefined>();
	const visibleActivities = $derived(
		activities.filter((activity) => visibleInActivitySection(activity, section))
	);
	const heading = $derived(activitySectionHeading(section));
	const sectionDescription = $derived(activitySectionDescription(section));

	function openManualImport(activity: DownloadActivity) {
		manualImportActivity = activity;
		importError = undefined;
	}

	async function submitManualImport(request: ManualImportRequest) {
		if (!manualImportActivity) return;
		importingId = manualImportActivity.id;
		importError = undefined;
		try {
			await onManualImport(manualImportActivity.id, request);
			manualImportActivity = undefined;
			await onRefresh();
		} catch (error) {
			importError = error instanceof Error ? error.message : 'Manual import failed';
		} finally {
			importingId = undefined;
		}
	}
</script>

<PageHeading eyebrow="Activity" title={heading.title} titleId="home-title">
	{#snippet actions()}
		<Button type="button" variant="outline" disabled={loading} onclick={onRefresh}>
			{loading ? 'Refreshing' : 'Refresh'}
		</Button>
	{/snippet}
</PageHeading>

{#if section === 'blocklist'}
	{#if releaseBlocklist.length > 0}
		<ActivityBlocklistTable
			items={releaseBlocklist}
			{canManage}
			deletingId={deletingBlocklistId}
			clearing={clearingReleaseBlocklist}
			onDelete={onDeleteReleaseBlocklistItem}
			onClear={onClearReleaseBlocklist}
		/>
	{:else}
		<EmptyState
			class="my-[18px] grid min-h-60 w-full place-items-center content-center gap-[18px] text-center"
			title={heading.empty}
			description={sectionDescription}
		/>
	{/if}
{:else if visibleActivities.length > 0}
	<ActivityQueueTable
		activities={visibleActivities}
		{canManage}
		{cancellingId}
		{retryingId}
		{deletingId}
		onManualImport={openManualImport}
		{onCancel}
		{onRetry}
		{onDelete}
	/>
{:else}
	<EmptyState
		class="my-[18px] grid min-h-60 w-full place-items-center content-center gap-[18px] text-center"
		title={heading.empty}
		description={sectionDescription}
	/>
{/if}

{#if manualImportActivity}
	<ActivityManualImportModal
		activity={manualImportActivity}
		importing={importingId === manualImportActivity.id}
		error={importError}
		onImport={submitManualImport}
		onClose={() => (manualImportActivity = undefined)}
	/>
{/if}
