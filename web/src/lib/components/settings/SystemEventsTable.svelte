<script lang="ts">
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTimeWithSeconds } from '$lib/settings/dateFormat';
	import type { SystemEvent } from '$lib/settings/types';
	import SystemEventSeverityIcon from './SystemEventSeverityIcon.svelte';
	import SystemLogAttributesButton from './SystemLogAttributesButton.svelte';

	interface Props {
		events: SystemEvent[];
		loading: boolean;
		hasMore: boolean;
		loadingMore: boolean;
		deletingId?: string;
		onDelete: (id: string) => void;
		onLoadMore: () => void;
	}

	let { events, loading, hasMore, loadingMore, deletingId, onDelete, onLoadMore }: Props = $props();

	function hasData(event: SystemEvent) {
		return Object.keys(event.data ?? {}).length > 0;
	}

	function errorText(event: SystemEvent) {
		const error = event.data?.error;
		if (typeof error === 'string') {
			return error;
		}
		const message = event.data?.message;
		return event.severity === 'error' && typeof message === 'string' ? message : '';
	}

	function handleScroll(event: Event) {
		const target = event.currentTarget as unknown as {
			scrollHeight: number;
			scrollTop: number;
			clientHeight: number;
		};
		const remaining = target.scrollHeight - target.scrollTop - target.clientHeight;
		if (remaining < 160 && hasMore && !loadingMore && !loading) {
			onLoadMore();
		}
	}

	function eventRowClass(severity: SystemEvent['severity']) {
		if (severity === 'error') return 'bg-destructive/5';
		if (severity === 'warning') return 'bg-secondary/40';
		return '';
	}
</script>

<div class="max-h-[min(62vh,680px)] overflow-auto" onscroll={handleScroll}>
	<Table.Root class="table-fixed">
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-[10.5rem] whitespace-nowrap">Time</Table.Head>
				<Table.Head class="w-20 whitespace-nowrap">Severity</Table.Head>
				<Table.Head class="w-28 whitespace-nowrap">Category</Table.Head>
				<Table.Head>Message</Table.Head>
				<Table.Head>Error</Table.Head>
				<Table.Head class="w-20"></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each events as event (event.id)}
				<Table.Row class={eventRowClass(event.severity)}>
					<Table.Cell class="whitespace-nowrap py-1.5"
						>{formatDateTimeWithSeconds(event.createdAt)}</Table.Cell
					>
					<Table.Cell class="py-1.5"
						><SystemEventSeverityIcon severity={event.severity} /></Table.Cell
					>
					<Table.Cell class="whitespace-nowrap py-1.5">{event.category}</Table.Cell>
					<Table.Cell class="break-words py-1.5">{event.message}</Table.Cell>
					<Table.Cell class="break-words py-1.5">{errorText(event)}</Table.Cell>
					<Table.Cell class="py-1.5">
						<div class="flex items-center justify-end gap-1">
							{#if hasData(event)}
								<SystemLogAttributesButton attributes={event.data} />
							{/if}
							<SettingsRowActionButton
								label={`Delete event ${event.message}`}
								icon="delete"
								variant="destructive"
								disabled={deletingId === event.id}
								onclick={() => onDelete(event.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={6}>
						{loading ? 'Loading events' : 'No events match the selected severity.'}
					</Table.Cell>
				</Table.Row>
			{/each}
			{#if loadingMore}
				<Table.Row>
					<Table.Cell colspan={6}>Loading more events</Table.Cell>
				</Table.Row>
			{/if}
		</Table.Body>
	</Table.Root>
</div>
