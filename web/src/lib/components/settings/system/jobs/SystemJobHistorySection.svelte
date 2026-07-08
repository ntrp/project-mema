<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { SystemJobExecution } from '$lib/settings/types';
	import SystemJobHistoryTable from './SystemJobHistoryTable.svelte';
	import SystemJobsFilters from './SystemJobsFilters.svelte';

	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		executions: SystemJobExecution[];
		loading: boolean;
		hasMore: boolean;
		loadingMore: boolean;
		savingRetention: boolean;
		retentionDays: number;
		selectedStatuses: string[];
		selectedQueues: string[];
		selectedKinds: string[];
		query: string;
		queueOptions: Option[];
		kindOptions: Option[];
		abortingId?: number;
		loadingLogsId?: number;
		onSaveRetention: () => void;
		onAbort: (id: number) => void;
		onLogs: (execution: SystemJobExecution) => void;
		onLoadMore: () => void;
	}

	let {
		executions,
		loading,
		hasMore,
		loadingMore,
		savingRetention,
		retentionDays = $bindable(),
		selectedStatuses = $bindable(),
		selectedQueues = $bindable(),
		selectedKinds = $bindable(),
		query = $bindable(),
		queueOptions,
		kindOptions,
		abortingId,
		loadingLogsId,
		onSaveRetention,
		onAbort,
		onLogs,
		onLoadMore
	}: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<div>
			<Card.Description>Filtered execution history</Card.Description>
			<Card.Title>Execution History</Card.Title>
		</div>
		<Card.Action>
			<div class="flex items-center gap-2">
				<Input
					class="w-24"
					type="number"
					min="1"
					max="365"
					bind:value={retentionDays}
					aria-label="Retention days"
				/>
				<Button variant="outline" disabled={savingRetention} onclick={onSaveRetention}
					>{savingRetention ? 'Saving' : 'Save'}</Button
				>
			</div>
		</Card.Action>
	</Card.Header>
	<Card.Content
		class="grid max-h-[calc(100vh-13rem)] min-h-[32rem] grid-rows-[auto_minmax(0,1fr)] gap-4"
	>
		<SystemJobsFilters
			bind:selectedStatuses
			bind:selectedQueues
			bind:selectedKinds
			bind:query
			{queueOptions}
			{kindOptions}
		/>
		<SystemJobHistoryTable
			{executions}
			{loading}
			{hasMore}
			{loadingMore}
			{abortingId}
			{loadingLogsId}
			{onAbort}
			{onLogs}
			{onLoadMore}
		/>
	</Card.Content>
</Card.Root>
