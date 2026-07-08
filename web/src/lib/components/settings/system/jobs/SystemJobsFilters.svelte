<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import SystemJobsMultiSelect from './SystemJobsMultiSelect.svelte';

	const statuses = [
		'running',
		'available',
		'scheduled',
		'retryable',
		'pending',
		'completed',
		'cancelled',
		'discarded'
	];

	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		selectedStatuses: string[];
		selectedQueues: string[];
		selectedKinds: string[];
		queueOptions: Option[];
		kindOptions: Option[];
		query: string;
		includeRoutine: boolean;
	}

	let {
		selectedStatuses = $bindable(),
		selectedQueues = $bindable(),
		selectedKinds = $bindable(),
		queueOptions,
		kindOptions,
		query = $bindable(),
		includeRoutine = $bindable()
	}: Props = $props();

	const statusOptions = statuses.map((status) => ({ value: status, label: status }));
</script>

<div class="grid gap-4 rounded-md border border-border p-4">
	<div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)]">
		<SystemJobsMultiSelect
			id="job-status-filter"
			label="Status"
			bind:values={selectedStatuses}
			options={statusOptions}
			placeholder="All statuses"
		/>
		<SystemJobsMultiSelect
			id="job-queue-filter"
			label="Queue"
			bind:values={selectedQueues}
			options={queueOptions}
			placeholder="All queues"
		/>
		<SystemJobsMultiSelect
			id="job-kind-filter"
			label="Kind"
			bind:values={selectedKinds}
			options={kindOptions}
			placeholder="All kinds"
		/>
		<div class="grid gap-1.5">
			<Label for="job-query-filter">Search</Label>
			<Input id="job-query-filter" bind:value={query} placeholder="Kind, queue, args, errors" />
		</div>
	</div>
	<div class="flex min-h-9 items-center gap-2">
		<Checkbox id="job-include-routine" bind:checked={includeRoutine} />
		<Label for="job-include-routine">Include routine runs</Label>
	</div>
</div>
