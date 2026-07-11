<script lang="ts">
	import BellIcon from '@lucide/svelte/icons/bell';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { subscribeToAppEvent } from '$lib/app/realtime/appEventSource';
	import { createEventNotificationsResource } from '$lib/features/settings/resources/eventNotifications.svelte';
	import { formatCompactDateTime } from '$lib/settings/dateFormat';
	import type { SystemEvent } from '$lib/settings/types';
	import SystemEventSeverityIcon from '$lib/components/settings/system/events/SystemEventSeverityIcon.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	const maxEvents = 20;

	let open = $state(false);
	const resource = createEventNotificationsResource(() => open);
	const events = $derived(
		(resource.query.data?.events ?? []).filter(
			(event) => event.severity === 'warning' || event.severity === 'error'
		)
	);
	const loaded = $derived(!resource.query.isPending && !resource.query.isFetching);

	const visibleEvents = $derived(events.slice(0, maxEvents));
	const errorTotal = $derived(events.filter((event) => event.severity === 'error').length);

	onMount(() => {
		const unsubscribeCreated = subscribeToAppEvent<SystemEvent>(
			'system.event.created',
			({ data: nextEvent }) => {
				if (!nextEvent || (nextEvent.severity !== 'warning' && nextEvent.severity !== 'error')) {
					return;
				}
				resource.created(nextEvent);
			}
		);
		const unsubscribeDeleted = subscribeToAppEvent<{ id: string }>(
			'system.event.deleted',
			({ data: deleted }) => {
				if (deleted?.id) {
					resource.deleted(deleted.id);
				}
			}
		);
		const unsubscribeCleared = subscribeToAppEvent('system.events.cleared', () => {
			resource.cleared();
		});
		return () => {
			unsubscribeCreated();
			unsubscribeDeleted();
			unsubscribeCleared();
		};
	});

	function handleOpenChange(nextOpen: boolean) {
		if (nextOpen) void resource.query.refetch();
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
