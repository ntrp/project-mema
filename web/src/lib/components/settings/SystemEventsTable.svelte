<script lang="ts">
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
</script>

<div class="table-wrap events-table-wrap" onscroll={handleScroll}>
	<table class="data-table events-table">
		<colgroup>
			<col class="event-time-column" />
			<col class="event-severity-column" />
			<col class="event-category-column" />
			<col class="event-message-column" />
			<col class="event-error-column" />
			<col class="event-actions-column" />
		</colgroup>
		<thead>
			<tr>
				<th>Time</th>
				<th>Severity</th>
				<th>Category</th>
				<th>Message</th>
				<th>Error</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each events as event (event.id)}
				<tr class={`event-row event-${event.severity}`}>
					<td>{formatDateTimeWithSeconds(event.createdAt)}</td>
					<td><SystemEventSeverityIcon severity={event.severity} /></td>
					<td>{event.category}</td>
					<td>{event.message}</td>
					<td>{errorText(event)}</td>
					<td class="row-actions">
						{#if hasData(event)}
							<SystemLogAttributesButton attributes={event.data} />
						{/if}
						<button
							type="button"
							class="danger icon-button"
							aria-label={`Delete event ${event.message}`}
							disabled={deletingId === event.id}
							onclick={() => onDelete(event.id)}
						>
							<span class="app-icon" aria-hidden="true">delete</span>
						</button>
					</td>
				</tr>
			{:else}
				<tr>
					<td colspan="6"
						>{loading ? 'Loading events' : 'No events match the selected severity.'}</td
					>
				</tr>
			{/each}
			{#if loadingMore}
				<tr>
					<td colspan="6">Loading more events</td>
				</tr>
			{/if}
		</tbody>
	</table>
</div>
