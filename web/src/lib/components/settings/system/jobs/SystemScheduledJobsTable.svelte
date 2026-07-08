<script lang="ts">
	import BanIcon from '@lucide/svelte/icons/ban';
	import PauseIcon from '@lucide/svelte/icons/pause';
	import PlayIcon from '@lucide/svelte/icons/play';
	import SaveIcon from '@lucide/svelte/icons/save';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { SystemJobSchedule } from '$lib/settings/types';
	import { createRowPulse } from '../cache/rowPulse.svelte';
	import SystemJobProgress from './SystemJobProgress.svelte';
	import { canAbortStatus, formatInterval, statusClass } from './systemJobDisplay';

	interface Props {
		schedules: SystemJobSchedule[];
		updatingId?: string;
		updatingIntervalId?: string;
		abortingId?: number;
		onPause: (schedule: SystemJobSchedule) => void;
		onResume: (schedule: SystemJobSchedule) => void;
		onSaveInterval: (schedule: SystemJobSchedule, intervalSeconds: number) => void;
		onAbort: (id: number) => void;
	}

	let {
		schedules,
		updatingId,
		updatingIntervalId,
		abortingId,
		onPause,
		onResume,
		onSaveInterval,
		onAbort
	}: Props = $props();

	const rowPulse = createRowPulse();
	const rowKeys = $derived(schedules.map((schedule) => schedule.id));
	let intervalDrafts = $state<Record<string, number>>({});

	$effect(() => rowPulse.update(rowKeys));

	function intervalDraft(schedule: SystemJobSchedule) {
		return intervalDrafts[schedule.id] ?? schedule.intervalSeconds;
	}

	function updateIntervalDraft(schedule: SystemJobSchedule, value: number) {
		intervalDrafts = { ...intervalDrafts, [schedule.id]: value };
	}
</script>

<div class="overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-full table-auto border-collapse">
		<Table.Header class="bg-card">
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head class="w-px">Status</Table.Head>
				<Table.Head class="w-px">Interval</Table.Head>
				<Table.Head>Progress</Table.Head>
				<Table.Head class="w-px">Next</Table.Head>
				<Table.Head class="w-px">Last</Table.Head>
				<Table.Head class="w-px text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each schedules as schedule (schedule.id)}
				<Table.Row class={rowPulse.classFor(schedule.id)}>
					<Table.Cell>
						<div class="flex min-w-0 items-center gap-2">
							<strong class="truncate">{schedule.name}</strong>
							{#if schedule.historyPolicy === 'routine'}
								<Badge variant="outline">routine</Badge>
							{/if}
						</div>
						<span class="block truncate text-xs text-muted-foreground">{schedule.kind}</span>
					</Table.Cell>
					<Table.Cell class="w-px">
						<Badge
							variant="outline"
							class={statusClass(schedule.paused ? 'cancelled' : schedule.activeStatus)}
							>{schedule.paused ? 'paused' : schedule.activeStatus || 'idle'}</Badge
						>
					</Table.Cell>
					<Table.Cell class="w-px">
						{#if schedule.intervalConfigurable}
							<div class="flex items-center gap-1">
								<Input
									class="w-20"
									type="number"
									min="15"
									step="15"
									value={intervalDraft(schedule)}
									aria-label={`Interval seconds for ${schedule.name}`}
									oninput={(event) =>
										updateIntervalDraft(schedule, event.currentTarget.valueAsNumber)}
								/>
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<Button
												{...props}
												type="button"
												variant="outline"
												size="icon-sm"
												aria-label="Save interval"
												disabled={updatingIntervalId === schedule.id ||
													intervalDraft(schedule) === schedule.intervalSeconds}
												onclick={() => onSaveInterval(schedule, intervalDraft(schedule))}
											>
												<SaveIcon aria-hidden="true" />
											</Button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>Save interval</Tooltip.Content>
								</Tooltip.Root>
							</div>
						{:else}
							{formatInterval(schedule.intervalSeconds)}
						{/if}
					</Table.Cell>
					<Table.Cell>
						<SystemJobProgress
							value={schedule.activeProgressPercent}
							label={schedule.activeProgressLabel ||
								schedule.activeInfoMessage ||
								schedule.activeStatus ||
								'Idle'}
						/>
					</Table.Cell>
					<Table.Cell class="w-px whitespace-nowrap">
						{schedule.nextRunAt ? formatDateTime(schedule.nextRunAt) : 'Paused'}
					</Table.Cell>
					<Table.Cell class="w-px whitespace-nowrap">
						{schedule.lastCreatedAt ? formatDateTime(schedule.lastCreatedAt) : 'Never'}
					</Table.Cell>
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
											aria-label={schedule.paused ? 'Resume schedule' : 'Pause schedule'}
											disabled={updatingId === schedule.id}
											onclick={() => (schedule.paused ? onResume(schedule) : onPause(schedule))}
										>
											{#if schedule.paused}
												<PlayIcon aria-hidden="true" />
											{:else}
												<PauseIcon aria-hidden="true" />
											{/if}
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content
									>{schedule.paused ? 'Resume schedule' : 'Pause schedule'}</Tooltip.Content
								>
							</Tooltip.Root>
							{#if schedule.activeRiverJobId && canAbortStatus(schedule.activeStatus)}
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<Button
												{...props}
												type="button"
												variant="destructive"
												size="icon-sm"
												aria-label="Abort active run"
												disabled={abortingId === schedule.activeRiverJobId}
												onclick={() => onAbort(schedule.activeRiverJobId!)}
											>
												<BanIcon aria-hidden="true" />
											</Button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>Abort active run</Tooltip.Content>
								</Tooltip.Root>
							{/if}
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={7} class="text-muted-foreground"
						>No fixed scheduled jobs are registered.</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
