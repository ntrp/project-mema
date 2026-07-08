<script lang="ts">
	import BanIcon from '@lucide/svelte/icons/ban';
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
		jobs: SystemJobExecution[];
		abortingId?: number;
		onAbort: (id: number) => void;
	}

	let { jobs, abortingId, onAbort }: Props = $props();

	const rowPulse = createRowPulse();
	const rowKeys = $derived(jobs.map((job) => String(job.riverJobId)));

	$effect(() => rowPulse.update(rowKeys));
</script>

<div class="overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-full table-auto border-collapse">
		<Table.Header class="bg-card">
			<Table.Row>
				<Table.Head class="w-px">Status</Table.Head>
				<Table.Head class="w-px">ID</Table.Head>
				<Table.Head>Kind</Table.Head>
				<Table.Head class="w-px">Queue</Table.Head>
				<Table.Head>Progress</Table.Head>
				<Table.Head class="w-px">Scheduled</Table.Head>
				<Table.Head class="w-px text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each jobs as job (job.riverJobId)}
				<Table.Row class={rowPulse.classFor(String(job.riverJobId))}>
					<Table.Cell class="w-px">
						<Badge variant="outline" class={statusClass(job.status)}>{job.status}</Badge>
					</Table.Cell>
					<Table.Cell class="w-px font-mono text-xs">{job.riverJobId}</Table.Cell>
					<Table.Cell class="max-w-72">
						<strong class="block truncate">{job.kind}</strong>
					</Table.Cell>
					<Table.Cell class="w-px">{job.queue}</Table.Cell>
					<Table.Cell>
						<SystemJobProgress
							value={job.progressPercent}
							label={executionMessage(job.progressLabel, job.infoMessage, job.status)}
						/>
					</Table.Cell>
					<Table.Cell class="w-px whitespace-nowrap">{formatDateTime(job.scheduledAt)}</Table.Cell>
					<Table.Cell class="w-px text-right">
						{#if canAbortStatus(job.status)}
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label="Abort job"
											disabled={abortingId === job.riverJobId}
											onclick={() => onAbort(job.riverJobId)}
										>
											<BanIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Abort job</Tooltip.Content>
							</Tooltip.Root>
						{/if}
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={7} class="text-muted-foreground"
						>No one-shot jobs are running or scheduled.</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
