<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { abortSystemJob, listSystemJobs } from '$lib/settings/api';
	import type { SystemJob } from '$lib/settings/types';
	import { parseSystemEvent } from '../events/systemEventStream';
	import SystemJobAbortDialog from './SystemJobAbortDialog.svelte';
	import SystemJobsFilters from './SystemJobsFilters.svelte';
	import SystemJobsTable from './SystemJobsTable.svelte';

	const reloadDelayMs = 180;
	const defaultStatuses = ['available', 'scheduled', 'retryable', 'running'];
	const knownQueues = ['media_search', 'downloads'];
	const knownKinds = [
		'media.release_search',
		'media.auto_search_download',
		'media.rss_sync',
		'media.grab_release',
		'download.activity_sync'
	];

	let allJobs = $state<SystemJob[]>([]);
	let loading = $state(false);
	let errorMessage = $state('');
	let selectedStatuses = $state([...defaultStatuses]);
	let selectedQueues = $state<string[]>([]);
	let selectedKinds = $state<string[]>([]);
	let query = $state('');
	let abortingId = $state<number | undefined>();
	let abortCandidate = $state<SystemJob | undefined>();
	let mounted = false;

	const jobs = $derived(allJobs.filter(matchesFilters));
	const filterKey = $derived(
		JSON.stringify({
			statuses: selectedStatuses,
			queues: selectedQueues,
			kinds: selectedKinds,
			query
		})
	);
	const queueOptions = $derived(
		optionList(
			knownQueues,
			allJobs.map((job) => job.queue)
		)
	);
	const kindOptions = $derived(
		optionList(
			knownKinds,
			allJobs.map((job) => job.kind)
		)
	);

	onMount(() => {
		mounted = true;
		void loadJobs();
		const source = new EventSource('/api/events', { withCredentials: true });
		source.addEventListener('system.job.updated', (event) => {
			const job = parseSystemEvent<SystemJob>(event);
			if (job) {
				applyJobUpdate(job);
			}
		});
		source.addEventListener('error', () => {
			errorMessage = errorMessage || 'Job event stream disconnected';
		});
		source.addEventListener('open', () => {
			if (errorMessage === 'Job event stream disconnected') {
				errorMessage = '';
			}
		});
		return () => source.close();
	});

	$effect(() => {
		filterKey;
		if (!mounted) return;
		const timeout = window.setTimeout(() => void loadJobs(), reloadDelayMs);
		return () => window.clearTimeout(timeout);
	});

	async function loadJobs() {
		if (loading) return;
		loading = true;
		errorMessage = '';
		try {
			allJobs = await listSystemJobs({
				status: selectedStatuses.length > 0 ? selectedStatuses : undefined,
				query: query.trim() || undefined,
				limit: 200
			});
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load jobs';
		} finally {
			loading = false;
		}
	}

	function applyJobUpdate(job: SystemJob) {
		if (!matchesServerFilters(job)) {
			allJobs = allJobs.filter((current) => current.id !== job.id);
			return;
		}
		allJobs = [job, ...allJobs.filter((current) => current.id !== job.id)].slice(0, 200);
	}

	function matchesFilters(job: SystemJob) {
		if (!matchesServerFilters(job)) return false;
		if (selectedQueues.length > 0 && !selectedQueues.includes(job.queue)) return false;
		if (selectedKinds.length > 0 && !selectedKinds.includes(job.kind)) return false;
		return true;
	}

	function matchesServerFilters(job: SystemJob) {
		if (selectedStatuses.length > 0 && !selectedStatuses.includes(job.status)) return false;
		const trimmedQuery = query.trim().toLowerCase();
		if (!trimmedQuery) return true;
		return [job.kind, job.queue, job.args, job.errors, job.infoMessage]
			.join(' ')
			.toLowerCase()
			.includes(trimmedQuery);
	}

	function optionList(defaults: string[], values: string[]) {
		return Array.from(new Set([...defaults, ...values].filter(Boolean)))
			.sort((left, right) => left.localeCompare(right))
			.map((value) => ({ value, label: value }));
	}

	async function abortJob() {
		if (!abortCandidate) return;
		abortingId = abortCandidate.id;
		errorMessage = '';
		try {
			const updated = await abortSystemJob(abortCandidate.id);
			applyJobUpdate(updated);
			abortCandidate = undefined;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not abort job';
		} finally {
			abortingId = undefined;
		}
	}
</script>

<Card.Root aria-labelledby="system-jobs-title">
	<Card.Header>
		<div>
			<Card.Description class="flex items-center gap-2">
				<span class="relative flex size-2.5">
					<span
						class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"
					></span>
					<span class="relative inline-flex size-2.5 rounded-full bg-emerald-500"></span>
				</span>
				<span>Live</span>
			</Card.Description>
			<Card.Title id="system-jobs-title">Jobs</Card.Title>
		</div>
		<Card.Action>
			<Button type="button" variant="outline" disabled={loading} onclick={() => void loadJobs()}>
				{loading ? 'Refreshing' : 'Refresh'}
			</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content
		class="grid max-h-[calc(100vh-13rem)] min-h-[32rem] grid-rows-[auto_auto_minmax(0,1fr)] gap-4"
	>
		<SystemJobsFilters
			bind:selectedStatuses
			bind:selectedQueues
			bind:selectedKinds
			bind:query
			{queueOptions}
			{kindOptions}
		/>
		{#if errorMessage}
			<p
				class="m-0 rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive"
			>
				{errorMessage}
			</p>
		{/if}
		<SystemJobsTable {jobs} {abortingId} onAbort={(job) => (abortCandidate = job)} />
	</Card.Content>
</Card.Root>

<SystemJobAbortDialog
	job={abortCandidate}
	onClose={() => (abortCandidate = undefined)}
	onAbort={abortJob}
/>
