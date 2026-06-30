<script lang="ts">
	/* global EventSource, MessageEvent */
	import { onMount } from 'svelte';
	import {
		deleteSystemEvent,
		getSystemEventSettings,
		listSystemEvents,
		updateSystemEventSettings
	} from '$lib/settings/api';
	import type { SystemEvent } from '$lib/settings/types';
	import SystemLogAttributesButton from './SystemLogAttributesButton.svelte';

	type StreamEnvelope<T> = {
		data: T;
	};

	const maxEvents = 300;

	let events = $state<SystemEvent[]>([]);
	let loading = $state(true);
	let retentionDays = $state(30);
	let saving = $state(false);
	let deletingId = $state<string | undefined>();
	let connected = $state(false);
	let errorMessage = $state('');
	let message = $state('');

	onMount(() => {
		void load();
		const source = new EventSource('/api/events', { withCredentials: true });
		source.addEventListener('open', () => {
			connected = true;
		});
		source.addEventListener('error', () => {
			connected = false;
		});
		source.addEventListener('system.event.created', (event) => {
			const nextEvent = parseEvent<SystemEvent>(event);
			if (nextEvent) {
				events = [nextEvent, ...events.filter((item) => item.id !== nextEvent.id)].slice(
					0,
					maxEvents
				);
			}
		});
		source.addEventListener('system.event.deleted', (event) => {
			const deleted = parseEvent<{ id: string }>(event);
			if (deleted?.id) {
				events = events.filter((item) => item.id !== deleted.id);
			}
		});
		return () => source.close();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			const [nextEvents, nextSettings] = await Promise.all([
				listSystemEvents(),
				getSystemEventSettings()
			]);
			events = nextEvents;
			retentionDays = nextSettings.retentionDays;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load events';
		} finally {
			loading = false;
		}
	}

	async function saveSettings(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		errorMessage = '';
		message = '';
		try {
			retentionDays = (await updateSystemEventSettings({ retentionDays })).retentionDays;
			events = await listSystemEvents();
			message = 'Event settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save event settings';
		} finally {
			saving = false;
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

	function parseEvent<T>(event: Event) {
		try {
			const message = event as MessageEvent<string>;
			const envelope = JSON.parse(message.data) as StreamEnvelope<T>;
			return envelope.data;
		} catch {
			return undefined;
		}
	}

	function formatTime(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			month: 'short',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		}).format(new Date(value));
	}

	function hasData(event: SystemEvent) {
		return Object.keys(event.data ?? {}).length > 0;
	}
</script>

<section class="panel log-settings-panel" aria-label="Events">
	<div class="section-heading">
		<div class="log-controls">
			<span class:connected class="log-stream-state">{connected ? 'Live' : 'Reconnecting'}</span>
			<button type="button" class="secondary compact-action" disabled={loading} onclick={load}>
				Refresh
			</button>
		</div>
	</div>

	{#if errorMessage}
		<p class="inline-error">{errorMessage}</p>
	{/if}
	{#if message}
		<p class="muted">{message}</p>
	{/if}

	<form class="settings-form compact-form event-settings-form" onsubmit={saveSettings}>
		<label>
			<span>Retention days</span>
			<input type="number" min="1" max="365" bind:value={retentionDays} />
		</label>
		<div class="form-actions">
			<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</button>
		</div>
	</form>

	<div class="table-wrap">
		<table class="data-table events-table">
			<thead>
				<tr>
					<th>Time</th>
					<th>Severity</th>
					<th>Category</th>
					<th>Message</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each events as event (event.id)}
					<tr class={`event-row event-${event.severity}`}>
						<td>{formatTime(event.createdAt)}</td>
						<td><span class="status-pill">{event.severity}</span></td>
						<td>{event.category}</td>
						<td>{event.message}</td>
						<td class="row-actions">
							{#if hasData(event)}
								<SystemLogAttributesButton attributes={event.data} />
							{/if}
							<button
								type="button"
								class="danger icon-button"
								aria-label={`Delete event ${event.message}`}
								disabled={deletingId === event.id}
								onclick={() => deleteEvent(event.id)}
							>
								<span class="app-icon" aria-hidden="true">delete</span>
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="5">{loading ? 'Loading events' : 'No events recorded.'}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
