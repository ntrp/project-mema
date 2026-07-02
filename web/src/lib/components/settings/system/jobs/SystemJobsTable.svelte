<script lang="ts">
	import BanIcon from '@lucide/svelte/icons/ban';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { SystemJob } from '$lib/settings/types';
	import { createRowPulse } from '../cache/rowPulse.svelte';

	interface Props {
		jobs: SystemJob[];
		abortingId?: number;
		onAbort: (_job: SystemJob) => void;
	}

	let { jobs, abortingId, onAbort }: Props = $props();
	const rowPulse = createRowPulse();
	const rowKeys = $derived(jobs.map((job) => `${job.id}:${job.status}`));

	$effect(() => rowPulse.update(rowKeys));

	function canAbort(job: SystemJob) {
		return !['completed', 'cancelled', 'discarded'].includes(job.status);
	}

	function statusClass(status: string) {
		if (status === 'running') return 'border-sky-500/50 bg-sky-500/10 text-sky-300';
		if (status === 'completed') return 'border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
		if (status === 'cancelled' || status === 'discarded') {
			return 'border-destructive/50 bg-destructive/10 text-destructive';
		}
		return 'border-amber-500/50 bg-amber-500/10 text-amber-300';
	}
</script>

<div class="min-h-0 overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-full table-auto border-collapse">
		<Table.Header class="sticky top-0 bg-card">
			<Table.Row>
				<Table.Head class="w-px">Status</Table.Head>
				<Table.Head class="w-px">ID</Table.Head>
				<Table.Head>Kind</Table.Head>
				<Table.Head class="w-px">Queue</Table.Head>
				<Table.Head class="w-px">Attempt</Table.Head>
				<Table.Head>Info</Table.Head>
				<Table.Head class="w-px">Scheduled</Table.Head>
				<Table.Head class="w-px text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each jobs as job (job.id)}
				{@const rowKey = `${job.id}:${job.status}`}
				<Table.Row class={rowPulse.classFor(rowKey)}>
					<Table.Cell class="w-px">
						<Badge variant="outline" class={statusClass(job.status)}>{job.status}</Badge>
					</Table.Cell>
					<Table.Cell class="w-px font-mono text-xs">{job.id}</Table.Cell>
					<Table.Cell class="max-w-60">
						<strong class="block truncate">{job.kind}</strong>
					</Table.Cell>
					<Table.Cell class="w-px">{job.queue}</Table.Cell>
					<Table.Cell class="w-px">{job.attempt}/{job.maxAttempts}</Table.Cell>
					<Table.Cell class="max-w-96">
						<span class="block truncate text-sm text-muted-foreground">{job.infoMessage}</span>
					</Table.Cell>
					<Table.Cell class="w-px">{formatDateTime(job.scheduledAt)}</Table.Cell>
					<Table.Cell class="w-px text-right">
						{#if canAbort(job)}
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label="Abort job"
											disabled={abortingId === job.id}
											onclick={() => onAbort(job)}
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
					<Table.Cell colspan={8} class="text-muted-foreground"
						>No jobs match the filters.</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
