<script lang="ts">
	import { onMount } from 'svelte';
	import type { SystemEvent } from '$lib/settings/types';
	import LivePulseDot from '$lib/components/shared/LivePulseDot.svelte';
	import * as Card from '$lib/components/ui/card';
	import ClearSystemEventsModal from './ClearSystemEventsModal.svelte';
	import SystemEventsControls from './SystemEventsControls.svelte';
	import SystemEventsTable from './SystemEventsTable.svelte';
	import { subscribeSystemEvents } from './systemEventSubscription';
	import { createSystemEventsResource } from '$lib/features/settings/resources/systemEvents.svelte';

	type SeverityFilter = 'info' | 'warning' | 'error';

	interface Props {
		onConnectionChange?: (connected: boolean) => void;
	}

	let { onConnectionChange }: Props = $props();

	const resource = createSystemEventsResource();
	const events = $derived(resource.query.data?.events ?? []);
	const hasMore = $derived(resource.query.data?.hasMore ?? false);
	const loading = $derived(resource.query.isPending || resource.query.isFetching);
	let loadingMore = $state(false);
	let severityFilter = $state<SeverityFilter>('info');
	const clearing = $derived(resource.clear.isPending);
	let clearModalOpen = $state(false);
	const deletingId = $derived(resource.remove.variables);
	let errorMessage = $state('');
	let message = $state('');

	const visibleEvents = $derived(events.filter((event) => severityVisible(event, severityFilter)));

	onMount(() => {
		const close = subscribeSystemEvents({
			onOpen: () => onConnectionChange?.(true),
			onError: () => onConnectionChange?.(false),
			onCreated: resource.created,
			onDeleted: resource.deleted,
			onCleared: resource.cleared
		});
		return () => {
			onConnectionChange?.(false);
			close();
		};
	});

	async function deleteEvent(id: string) {
		errorMessage = '';
		try {
			await resource.remove.mutateAsync(id);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not delete event';
		} finally {
			/* mutation owns pending state */
		}
	}

	async function clearEvents() {
		errorMessage = '';
		message = '';
		try {
			await resource.clear.mutateAsync();
			clearModalOpen = false;
			message = 'Events cleared';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not clear events';
		} finally {
			/* mutation owns pending state */
		}
	}

	async function loadMoreEvents() {
		const before = events.at(-1)?.createdAt;
		if (!before || loading || loadingMore || !hasMore) {
			return;
		}
		loadingMore = true;
		errorMessage = '';
		try {
			await resource.loadMore(before);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load more events';
		} finally {
			loadingMore = false;
		}
	}

	function severityVisible(event: SystemEvent, filter: SeverityFilter) {
		if (filter === 'info') {
			return true;
		}
		if (filter === 'warning') {
			return event.severity === 'warning' || event.severity === 'error';
		}
		return event.severity === 'error';
	}
</script>

<Card.Root aria-labelledby="system-events-title">
	<Card.Header>
		<div>
			<Card.Description class="flex items-center gap-2">
				<LivePulseDot />
				<span>Live</span>
			</Card.Description>
			<Card.Title id="system-events-title">Events</Card.Title>
		</div>
		<Card.Action>
			<SystemEventsControls
				{severityFilter}
				{loading}
				{clearing}
				eventsEmpty={events.length === 0}
				onSeverityChange={(severity) => (severityFilter = severity)}
				onClear={() => (clearModalOpen = true)}
			/>
		</Card.Action>
	</Card.Header>

	<Card.Content class="grid gap-4">
		{#if errorMessage}
			<p class="m-0 font-bold text-destructive">{errorMessage}</p>
		{/if}
		{#if message}
			<p class="m-0 text-sm leading-6 text-muted-foreground">{message}</p>
		{/if}

		<SystemEventsTable
			events={visibleEvents}
			{loading}
			{hasMore}
			{loadingMore}
			{deletingId}
			onDelete={(id) => void deleteEvent(id)}
			onLoadMore={() => void loadMoreEvents()}
		/>
	</Card.Content>
</Card.Root>

{#if clearModalOpen}
	<ClearSystemEventsModal
		{clearing}
		onCancel={() => (clearModalOpen = false)}
		onConfirm={() => void clearEvents()}
	/>
{/if}
