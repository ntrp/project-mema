<script lang="ts">
	import ActivityManualImportModal from './ActivityManualImportModal.svelte';
	import ActivityActions from './ActivityActions.svelte';
	import ActivityProgressCell from './ActivityProgressCell.svelte';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import { activityDisplay } from './activityDisplay';
	import { manualImportDownloadActivity } from '$lib/settings/api';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

	interface Props {
		activities: DownloadActivity[];
		loading: boolean;
		canManage: boolean;
		cancellingId?: string;
		deletingId?: string;
		onRefresh: () => void;
		onCancel: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
	}

	let {
		activities,
		loading,
		canManage,
		cancellingId,
		deletingId,
		onRefresh,
		onCancel,
		onDelete
	}: Props = $props();
	let manualImportActivity = $state<DownloadActivity | undefined>();
	let importingId = $state<string | undefined>();
	let importError = $state<string | undefined>();

	function openManualImport(activity: DownloadActivity) {
		manualImportActivity = activity;
		importError = undefined;
	}

	async function submitManualImport(request: ManualImportRequest) {
		if (!manualImportActivity) return;
		importingId = manualImportActivity.id;
		importError = undefined;
		try {
			await manualImportDownloadActivity(manualImportActivity.id, request);
			manualImportActivity = undefined;
			await onRefresh();
		} catch (error) {
			importError = error instanceof Error ? error.message : 'Manual import failed';
		} finally {
			importingId = undefined;
		}
	}
</script>

<PageHeading eyebrow="Activity" title="Downloads and imports" titleId="home-title">
	{#snippet actions()}
		<Button type="button" variant="outline" disabled={loading} onclick={onRefresh}>
			{loading ? 'Refreshing' : 'Refresh'}
		</Button>
	{/snippet}
</PageHeading>

{#if activities.length > 0}
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
				{#each activities as activity (activity.id)}
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
	>
		<p class="m-0 text-lg font-black text-foreground">No download activity yet</p>
	</EmptyState>
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
