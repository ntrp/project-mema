<script lang="ts">
	import BanIcon from '@lucide/svelte/icons/ban';
	import ScrollTextIcon from '@lucide/svelte/icons/scroll-text';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { SystemJobExecution } from '$lib/settings/types';
	import { createRowPulse } from '../cache/rowPulse.svelte';
	import SystemJobProgress from './SystemJobProgress.svelte';
	import { canAbortStatus, executionMessage, statusClass } from './systemJobDisplay';

	interface Props {
		executions: SystemJobExecution[];
		loading: boolean;
		hasMore: boolean;
		loadingMore: boolean;
		abortingId?: number;
		loadingLogsId?: number;
		onAbort: (id: number) => void;
		onLogs: (execution: SystemJobExecution) => void;
		onLoadMore: () => void;
	}

	let {
		executions,
		loading,
		hasMore,
		loadingMore,
		abortingId,
		loadingLogsId,
		onAbort,
		onLogs,
		onLoadMore
	}: Props = $props();

	const rowPulse = createRowPulse();
	const rowKeys = $derived(executions.map((execution) => String(execution.riverJobId)));

	$effect(() => rowPulse.update(rowKeys));

	function handleScroll(event: Event) {
		const target = event.currentTarget as HTMLDivElement;
		const remaining = target.scrollHeight - target.scrollTop - target.clientHeight;
		if (remaining < 160 && hasMore && !loadingMore && !loading) {
			onLoadMore();
		}
	}
</script>

<div class="min-h-0 overflow-auto rounded-md border border-border" onscroll={handleScroll}>
	<Table.Root class="min-w-full table-auto border-collapse">
		<Table.Header class="sticky top-0 bg-card">
			<Table.Row>
				<Table.Head class="w-px">Status</Table.Head>
				<Table.Head class="w-px">ID</Table.Head>
				<Table.Head>Kind</Table.Head>
				<Table.Head class="w-px">Queue</Table.Head>
				<Table.Head>Progress</Table.Head>
				<Table.Head class="w-px">Attempt</Table.Head>
				<Table.Head class="w-px">Updated</Table.Head>
				<Table.Head class="w-px text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each executions as execution (execution.riverJobId)}
				<Table.Row class={rowPulse.classFor(String(execution.riverJobId))}>
					<Table.Cell class="w-px">
						<Badge variant="outline" class={statusClass(execution.status)}>{execution.status}</Badge
						>
					</Table.Cell>
					<Table.Cell class="w-px font-mono text-xs">{execution.riverJobId}</Table.Cell>
					<Table.Cell class="max-w-72">
						<strong class="block truncate">{execution.kind}</strong>
						<span class="block truncate text-xs text-muted-foreground"
							>{execution.scheduleId || execution.classification}</span
						>
					</Table.Cell>
					<Table.Cell class="w-px">{execution.queue}</Table.Cell>
					<Table.Cell>
						<SystemJobProgress
							value={execution.progressPercent}
							label={executionMessage(
								execution.progressLabel,
								execution.infoMessage,
								execution.status
							)}
						/>
					</Table.Cell>
					<Table.Cell class="w-px">{execution.attempt}/{execution.maxAttempts}</Table.Cell>
					<Table.Cell class="w-px whitespace-nowrap"
						>{formatDateTime(execution.updatedAt)}</Table.Cell
					>
					<Table.Cell class="w-px text-right">
						<div class="flex justify-end gap-1">
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="outline"
											size="icon-sm"
											aria-label="Show execution logs"
											disabled={loadingLogsId === execution.riverJobId}
											onclick={() => onLogs(execution)}
										>
											<ScrollTextIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Show execution logs</Tooltip.Content>
							</Tooltip.Root>
							{#if canAbortStatus(execution.status)}
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<Button
												{...props}
												type="button"
												variant="destructive"
												size="icon-sm"
												aria-label="Abort job"
												disabled={abortingId === execution.riverJobId}
												onclick={() => onAbort(execution.riverJobId)}
											>
												<BanIcon aria-hidden="true" />
											</Button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>Abort job</Tooltip.Content>
								</Tooltip.Root>
							{/if}
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={8} class="text-muted-foreground">
						{loading ? 'Loading execution history' : 'No executions match the filters.'}
					</Table.Cell>
				</Table.Row>
			{/each}
			{#if loadingMore}
				<Table.Row><Table.Cell colspan={8}>Loading more executions</Table.Cell></Table.Row>
			{/if}
		</Table.Body>
	</Table.Root>
</div>
