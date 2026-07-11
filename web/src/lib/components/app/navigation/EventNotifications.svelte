<script lang="ts">
	import BellIcon from '@lucide/svelte/icons/bell';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { subscribeToAppEvent } from '$lib/app/realtime/appEventSource';
	import { listSystemEvents } from '$lib/components/settings/system/events/api';
	import { formatCompactDateTime } from '$lib/settings/dateFormat';
	import type { SystemEvent } from '$lib/settings/types';
	import SystemEventSeverityIcon from '$lib/components/settings/system/events/SystemEventSeverityIcon.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	const maxEvents = 20;

	let open = $state(false);
	let events = $state<SystemEvent[]>([]);
	let loaded = $state(false);

	const visibleEvents = $derived(events.slice(0, maxEvents));
	const errorTotal = $derived(events.filter((event) => event.severity === 'error').length);

	onMount(() => {
		const unsubscribeCreated = subscribeToAppEvent<SystemEvent>(
			'system.event.created',
			({ data: nextEvent }) => {
				if (!nextEvent || (nextEvent.severity !== 'warning' && nextEvent.severity !== 'error')) {
					return;
				}
				events = [nextEvent, ...events.filter((item) => item.id !== nextEvent.id)];
			}
		);
		const unsubscribeDeleted = subscribeToAppEvent<{ id: string }>(
			'system.event.deleted',
			({ data: deleted }) => {
				if (deleted?.id) {
					events = events.filter((item) => item.id !== deleted.id);
				}
			}
		);
		const unsubscribeCleared = subscribeToAppEvent('system.events.cleared', () => {
			events = [];
		});
		return () => {
			unsubscribeCreated();
			unsubscribeDeleted();
			unsubscribeCleared();
		};
	});

	function handleOpenChange(nextOpen: boolean) {
		if (nextOpen) void load();
	}

	async function load() {
		loaded = false;
		try {
			const response = await listSystemEvents();
			events = response.events.filter(
				(event) => event.severity === 'warning' || event.severity === 'error'
			);
		} catch {
			events = [];
		} finally {
			loaded = true;
		}
	}
</script>

<DropdownMenu.Root bind:open onOpenChange={handleOpenChange}>
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
