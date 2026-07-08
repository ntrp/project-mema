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
	import {
		canAbortStatus,
		formatInterval,
		scheduleCategoryLabel,
		statusClass
	} from './systemJobDisplay';
	interface Props {
		schedules: SystemJobSchedule[];
		updatingId?: string;
		updatingIntervalId?: string;
		runningId?: string;
		abortingId?: number;
		onPause: (schedule: SystemJobSchedule) => void;
		onResume: (schedule: SystemJobSchedule) => void;
		onRun: (schedule: SystemJobSchedule) => void;
		onSaveInterval: (schedule: SystemJobSchedule, intervalSeconds: number) => void;
		onAbort: (id: number) => void;
	}
	let {
		schedules,
		updatingId,
		updatingIntervalId,
		runningId,
		abortingId,
		onPause,
		onResume,
		onRun,
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
	const hasActiveRun = (schedule: SystemJobSchedule) =>
		!!schedule.activeRiverJobId && canAbortStatus(schedule.activeStatus);
	const progressValue = (schedule: SystemJobSchedule) =>
		schedule.activeStatus ? schedule.activeProgressPercent : 0;
</script>
{#snippet actionButton(label: string, disabled: boolean, onclick: () => void)}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant={label === 'Abort active run' ? 'destructive' : 'outline'}
					size="icon-sm"
					aria-label={label}
					{disabled}
					{onclick}
				>
					{#if label === 'Run schedule now'}
						<PlayIcon aria-hidden="true" />
					{:else if label === 'Abort active run'}
						<BanIcon aria-hidden="true" />
					{:else if label === 'Enable automatic job'}
						<PlayIcon aria-hidden="true" />
					{:else if label === 'Save interval'}
						<SaveIcon aria-hidden="true" />
					{:else}
						<PauseIcon aria-hidden="true" />
					{/if}
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{label}</Tooltip.Content>
	</Tooltip.Root>
{/snippet}
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
								<Badge variant="outline">{scheduleCategoryLabel(schedule.category)}</Badge>
								{#if schedule.historyPolicy === 'routine'}<Badge variant="outline">routine</Badge>{/if}
								{#if schedule.manualActionAvailable}<Badge variant="secondary">manual</Badge>{/if}
							</div>
							<span class="block truncate text-xs text-muted-foreground">{schedule.description || schedule.kind}</span>
					</Table.Cell>
					<Table.Cell class="w-px">
						<Badge
							variant="outline"
							class={statusClass(schedule.enabled ? schedule.activeStatus : 'cancelled')}
							>{schedule.enabled ? schedule.activeStatus || 'idle' : 'disabled'}</Badge
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
								{@render actionButton(
									'Save interval',
									updatingIntervalId === schedule.id ||
										intervalDraft(schedule) === schedule.intervalSeconds,
									() => onSaveInterval(schedule, intervalDraft(schedule))
								)}
							</div>
						{:else}
							{formatInterval(schedule.intervalSeconds)}
						{/if}
					</Table.Cell>
					<Table.Cell>
						<SystemJobProgress
							value={progressValue(schedule)}
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
							{@render actionButton(
								'Run schedule now',
								runningId === schedule.id || hasActiveRun(schedule),
								() => onRun(schedule)
							)}
							{@render actionButton(
								schedule.enabled ? 'Disable automatic job' : 'Enable automatic job',
								updatingId === schedule.id,
								() => (schedule.paused ? onResume(schedule) : onPause(schedule))
							)}
							{#if hasActiveRun(schedule)}
								{@render actionButton(
									'Abort active run',
									abortingId === schedule.activeRiverJobId,
									() => onAbort(schedule.activeRiverJobId!)
								)}
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
