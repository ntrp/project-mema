<script lang="ts">
	/* global EventSource, HTMLSelectElement, MessageEvent, HTMLDivElement */
	import { onMount, tick } from 'svelte';

	import { getSystemLogLevel, updateSystemLogLevel } from '$lib/settings/api';
	import type { SystemLogEntry, SystemLogLevel } from '$lib/settings/types';

	type StreamEnvelope<T> = {
		data: T;
	};

	const levels: SystemLogLevel[] = ['debug', 'info', 'warn', 'error'];
	const maxEntries = 500;

	let entries = $state<SystemLogEntry[]>([]);
	let level = $state<SystemLogLevel>('info');
	let loading = $state(true);
	let saving = $state(false);
	let connected = $state(false);
	let streamMessage = $state('Connecting');
	let errorMessage = $state('');

	onMount(() => {
		void loadLevel();

		const source = new EventSource('/api/system/logs', { withCredentials: true });
		source.addEventListener('open', () => {
			connected = true;
			streamMessage = 'Live';
		});
		source.addEventListener('error', () => {
			connected = false;
			streamMessage = 'Reconnecting';
		});
		source.addEventListener('system.log', (event) =>
			appendEntry(parseEvent<SystemLogEntry>(event))
		);
		source.addEventListener('system.log.level', (event) => {
			const payload = parseEvent<{ level?: SystemLogLevel }>(event);
			if (payload?.level) {
				level = payload.level;
			}
		});

		return () => source.close();
	});

	async function loadLevel() {
		loading = true;
		errorMessage = '';
		try {
			const response = await getSystemLogLevel();
			level = response.level;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load log level';
		} finally {
			loading = false;
		}
	}

	async function changeLevel(event: Event) {
		const select = event.currentTarget as HTMLSelectElement;
		const nextLevel = select.value as SystemLogLevel;
		saving = true;
		errorMessage = '';
		try {
			const response = await updateSystemLogLevel(nextLevel);
			level = response.level;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not update log level';
		} finally {
			saving = false;
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

	function appendEntry(entry?: SystemLogEntry) {
		if (!entry) {
			return;
		}
		entries = [...entries, entry].slice(-maxEntries);
		void scrollToBottom();
	}

	async function scrollToBottom() {
		await tick();
		const logViewport = document.querySelector<HTMLDivElement>('[data-log-viewer]');
		if (logViewport) {
			logViewport.scrollTop = logViewport.scrollHeight;
		}
	}

	function formatTime(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		}).format(new Date(value));
	}

	function attributeText(entry: SystemLogEntry) {
		if (!entry.attributes || Object.keys(entry.attributes).length === 0) {
			return '';
		}
		return JSON.stringify(entry.attributes);
	}
</script>

<section class="panel log-settings-panel" aria-labelledby="system-logs-title">
	<div class="section-heading">
		<div>
			<p>System</p>
			<h2 id="system-logs-title">Logs</h2>
		</div>
		<div class="log-controls">
			<span class:connected class="log-stream-state">{streamMessage}</span>
			<label>
				<span>Verbosity</span>
				<select value={level} disabled={loading || saving} onchange={changeLevel}>
					{#each levels as option (option)}
						<option value={option}>{option.toUpperCase()}</option>
					{/each}
				</select>
			</label>
		</div>
	</div>

	{#if errorMessage}
		<p class="inline-error">{errorMessage}</p>
	{/if}

	<div class="log-viewer" data-log-viewer aria-live="polite" aria-label="Application logs">
		{#each entries as entry (entry.id)}
			<div class={`log-row level-${entry.level}`}>
				<time datetime={entry.time}>{formatTime(entry.time)}</time>
				<span>{entry.level.toUpperCase()}</span>
				<p>{entry.message}</p>
				{#if attributeText(entry)}
					<code>{attributeText(entry)}</code>
				{/if}
			</div>
		{:else}
			<div class="log-empty">Waiting for log entries</div>
		{/each}
	</div>
</section>
