<script lang="ts">
	import { onMount } from 'svelte';
	import { clearSystemEvents, deleteSystemEvent, listSystemEvents } from './api';
	import type { SystemEvent } from '$lib/settings/types';
	import LivePulseDot from '$lib/components/shared/LivePulseDot.svelte';
	import * as Card from '$lib/components/ui/card';
	import ClearSystemEventsModal from './ClearSystemEventsModal.svelte';
	import SystemEventsControls from './SystemEventsControls.svelte';
	import SystemEventsTable from './SystemEventsTable.svelte';
	import { subscribeSystemEvents } from './systemEventSubscription';

	const eventPageLimit = 100;
	type SeverityFilter = 'info' | 'warning' | 'error';

	interface Props {
		onConnectionChange?: (connected: boolean) => void;
	}

	let { onConnectionChange }: Props = $props();

	let events = $state<SystemEvent[]>([]);
	let loading = $state(true);
	let loadingMore = $state(false);
	let hasMore = $state(false);
	let severityFilter = $state<SeverityFilter>('info');
	let clearing = $state(false);
	let clearModalOpen = $state(false);
	let deletingId = $state<string | undefined>();
	let errorMessage = $state('');
	let message = $state('');

	const visibleEvents = $derived(events.filter((event) => severityVisible(event, severityFilter)));

	onMount(() => {
		void load();
		const close = subscribeSystemEvents({
			onOpen: () => onConnectionChange?.(true),
			onError: () => onConnectionChange?.(false),
			onCreated: (event) => (events = [event, ...events.filter((item) => item.id !== event.id)]),
			onDeleted: (id) => (events = events.filter((item) => item.id !== id)),
			onCleared: () => {
				events = [];
				hasMore = false;
			}
		});
		return () => {
			onConnectionChange?.(false);
			close();
		};
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			const nextEvents = await listSystemEvents({ limit: eventPageLimit });
			events = nextEvents.events;
			hasMore = nextEvents.hasMore;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load events';
		} finally {
			loading = false;
		}
	}

	async function deleteEvent(id: string) {
		deletingId = id;
		errorMessage = '';
		try {
			await deleteSystemEvent(id);
			events = events.filter((event) => event.id !== id);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not delete event';
		} finally {
			deletingId = undefined;
		}
	}

	async function clearEvents() {
		clearing = true;
		errorMessage = '';
		message = '';
		try {
			await clearSystemEvents();
			events = [];
			hasMore = false;
			clearModalOpen = false;
			message = 'Events cleared';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not clear events';
		} finally {
			clearing = false;
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
			const response = await listSystemEvents({ before, limit: eventPageLimit });
			const existing = new Set(events.map((event) => event.id));
			events = [...events, ...response.events.filter((event) => !existing.has(event.id))];
			hasMore = response.hasMore;
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
