<script lang="ts">
	import BellIcon from '@lucide/svelte/icons/bell';
	import { resolve } from '$app/paths';
	import { listSystemEvents } from '$lib/settings/api';
	import { formatCompactDateTime } from '$lib/settings/dateFormat';
	import type { SystemEvent } from '$lib/settings/types';
	import SystemEventSeverityIcon from '$lib/components/settings/system/events/SystemEventSeverityIcon.svelte';
	import { parseSystemEvent } from '$lib/components/settings/system/events/systemEventStream';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	const maxEvents = 20;

	let open = $state(false);
	let events = $state<SystemEvent[]>([]);
	let errorTotal = $state(0);
	let loaded = $state(false);
	let source: EventSource | undefined;

	const visibleEvents = $derived(events.slice(0, maxEvents));

	$effect(() => {
		if (!open) {
			closeEvents();
			return;
		}
		void load();
		connectEvents();
		return closeEvents;
	});

	function connectEvents() {
		if (source) return;
		source = new EventSource('/api/events', { withCredentials: true });
		source.addEventListener('system.event.created', (event) => {
			const nextEvent = parseSystemEvent<SystemEvent>(event);
			if (!nextEvent || (nextEvent.severity !== 'warning' && nextEvent.severity !== 'error')) {
				return;
			}
			events = [nextEvent, ...events.filter((item) => item.id !== nextEvent.id)];
			errorTotal = events.filter((item) => item.severity === 'error').length;
		});
		source.addEventListener('system.event.deleted', (event) => {
			const deleted = parseSystemEvent<{ id: string }>(event);
			if (deleted?.id) {
				events = events.filter((item) => item.id !== deleted.id);
				errorTotal = events.filter((item) => item.severity === 'error').length;
			}
		});
		source.addEventListener('system.events.cleared', () => {
			events = [];
			errorTotal = 0;
		});
	}

	function closeEvents() {
		source?.close();
		source = undefined;
	}

	async function load() {
		loaded = false;
		try {
			const response = await listSystemEvents();
			events = response.events.filter(
				(event) => event.severity === 'warning' || event.severity === 'error'
			);
			errorTotal = events.filter((event) => event.severity === 'error').length;
		} catch {
			events = [];
			errorTotal = 0;
		} finally {
			loaded = true;
		}
	}
</script>

<DropdownMenu.Root bind:open>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button
				type="button"
				variant="outline"
				size="icon"
				class="relative"
				aria-label="Warning and error events"
				{...props}
			>
				<BellIcon aria-hidden="true" />
				{#if errorTotal > 0}
					<span
						class="absolute -top-1 -right-1 h-4 min-w-4 rounded-full border border-background bg-destructive px-1 text-[10px] leading-3.5 font-black text-destructive-foreground"
					>
						{errorTotal}
					</span>
				{/if}
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content
		align="end"
		class="grid max-h-[min(460px,calc(100vh-72px))] w-[min(380px,calc(100vw-28px))] gap-1 overflow-auto"
	>
		{#if open}
			{#if visibleEvents.length > 0}
				{#each visibleEvents as event (event.id)}
					<a
						class="grid gap-1 rounded-sm p-2 text-popover-foreground no-underline hover:bg-accent hover:text-accent-foreground"
						href={resolve('/system/events')}
						role="menuitem"
						onclick={() => (open = false)}
					>
						<SystemEventSeverityIcon severity={event.severity} />
						<strong>{event.message}</strong>
						<small class="text-xs text-muted-foreground">
							{event.category} - {formatCompactDateTime(event.createdAt)}
						</small>
					</a>
				{/each}
			{:else}
				<p class="m-0 p-2 text-sm text-muted-foreground">
					{loaded ? 'No warning or error events.' : 'Loading events'}
				</p>
			{/if}
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>
