<script lang="ts">
	import ActivityManualImportModal from './ActivityManualImportModal.svelte';
	import ActivityBlocklistTable from './ActivityBlocklistTable.svelte';
	import ActivityActions from './ActivityActions.svelte';
	import ActivityProgressCell from './ActivityProgressCell.svelte';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import { activityDisplay } from './activityDisplay';
	import { activitySectionHeading, visibleInActivitySection } from './activitySections';
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
		deletingId?: string;
		deletingBlocklistId?: string;
		clearingReleaseBlocklist?: boolean;
		onRefresh: () => void;
		onCancel: (_activity: DownloadActivity) => void;
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
		deletingId,
		deletingBlocklistId,
		clearingReleaseBlocklist = false,
		onRefresh,
		onCancel,
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
	const sectionDescription = $derived(
		section === 'history'
			? 'Completed and cancelled background activity appears here so you can review finished media work.'
			: section === 'blocklist'
				? 'Automatic blocks for broken or temporarily unavailable releases will appear here and expire according to the indexer search setting.'
				: 'Queued downloads and in-progress grabs appear here while Mema fetches monitored media.'
	);

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
	<Card class="overflow-hidden p-0">
		<Table.Root class="[&_td]:whitespace-nowrap [&_th]:whitespace-nowrap">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-[42px]"><span class="sr-only">Select</span></Table.Head>
					<Table.Head class="min-w-55 whitespace-normal">Media</Table.Head>
					<Table.Head>Year</Table.Head>
					<Table.Head>Languages</Table.Head>
					<Table.Head>Quality</Table.Head>
					<Table.Head>Formats</Table.Head>
					<Table.Head>Time left</Table.Head>
					<Table.Head class="min-w-55 whitespace-normal">Progress</Table.Head>
					<Table.Head class="text-right">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each visibleActivities as activity (activity.id)}
					{@const display = activityDisplay(activity)}
					<Table.Row>
						<Table.Cell class="w-[42px]">
							<Checkbox aria-label={`Select ${activity.releaseTitle}`} />
						</Table.Cell>
						<Table.Cell class="grid min-w-55 gap-1 whitespace-normal">
							<strong>{activity.mediaTitle}</strong>
							<small class="text-xs text-muted-foreground"
								>{activity.downloadClientName} · {activity.indexerName}</small
							>
							{#if activity.error}
								<small class="max-w-55 text-xs text-muted-foreground">{activity.error}</small>
							{/if}
						</Table.Cell>
						<Table.Cell>{display.year}</Table.Cell>
						<Table.Cell>{display.languages.length ? display.languages.join(', ') : '-'}</Table.Cell>
						<Table.Cell>{display.quality}</Table.Cell>
						<Table.Cell>
							{#if display.formats.length}
								<div class="flex min-w-28 flex-wrap gap-1">
									{#each display.formats as format (format)}
										<span
											class="rounded-md border border-primary/30 bg-primary/10 px-1.5 py-0.5 text-[11px] font-extrabold text-primary"
											>{format}</span
										>
									{/each}
								</div>
							{:else}
								-
							{/if}
						</Table.Cell>
						<Table.Cell>{display.timeLeft}</Table.Cell>
						<Table.Cell>
							<ActivityProgressCell {activity} />
						</Table.Cell>
						<Table.Cell class="min-w-22">
							<ActivityActions
								{activity}
								{canManage}
								{cancellingId}
								{deletingId}
								onManualImport={openManualImport}
								{onCancel}
								{onDelete}
							/>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card>
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
