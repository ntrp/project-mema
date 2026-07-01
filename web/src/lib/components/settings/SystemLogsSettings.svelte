<script lang="ts">
	/* global EventSource, HTMLSelectElement, MessageEvent, HTMLDivElement */
	import { onMount, tick } from 'svelte';

	import { getSystemLogLevel, updateSystemLogLevel } from '$lib/settings/api';
	import { formatTimeWithSeconds } from '$lib/settings/dateFormat';
	import type { SystemLogEntry, SystemLogLevel } from '$lib/settings/types';
	import SystemLogAttributesButton from './SystemLogAttributesButton.svelte';

	type StreamEnvelope<T> = {
		data: T;
	};

	const levels: SystemLogLevel[] = ['debug', 'info', 'warn', 'error'];
	const maxEntries = 500;

	interface Props {
		onConnectionChange?: (connected: boolean) => void;
	}

	let { onConnectionChange }: Props = $props();

	let entries = $state<SystemLogEntry[]>([]);
	let level = $state<SystemLogLevel>('info');
	let loading = $state(true);
	let saving = $state(false);
	let errorMessage = $state('');
	let followLogs = $state(true);
	let logViewport = $state<HTMLDivElement>();
	let lastScrollTop = 0;

	onMount(() => {
		void loadLevel();

		const source = new EventSource('/api/system/logs', { withCredentials: true });
		source.addEventListener('open', () => {
			onConnectionChange?.(true);
		});
		source.addEventListener('error', () => {
			onConnectionChange?.(false);
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

		return () => {
			onConnectionChange?.(false);
			source.close();
		};
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
		if (followLogs) {
			void scrollToBottom();
		}
	}

	async function scrollToBottom() {
		await tick();
		if (logViewport) {
			logViewport.scrollTop = logViewport.scrollHeight;
			lastScrollTop = logViewport.scrollTop;
		}
	}

	function clearLogs() {
		entries = [];
	}

	function enableFollow() {
		followLogs = true;
		void scrollToBottom();
	}

	function handleLogScroll(event: Event) {
		const target = event.currentTarget as HTMLDivElement;
		if (target.scrollTop < lastScrollTop - 4) {
			followLogs = false;
		}
		lastScrollTop = target.scrollTop;
	}

	function attributeText(entry: SystemLogEntry) {
		if (!entry.attributes || Object.keys(entry.attributes).length === 0) {
			return '';
		}
		return JSON.stringify(entry.attributes);
	}
</script>

<section class="panel log-settings-panel" aria-label="Logs">
	<div class="section-heading">
		<div class="log-controls">
			<button type="button" class="secondary compact-action" onclick={clearLogs}>Clear logs</button>
			<button
				type="button"
				class="secondary compact-action"
				class:active-follow={followLogs}
				aria-pressed={followLogs}
				onclick={enableFollow}
			>
				Follow logs
			</button>
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

	<div
		class="log-viewer"
		data-log-viewer
		bind:this={logViewport}
		onscroll={handleLogScroll}
		aria-live="polite"
		aria-label="Application logs"
	>
		{#each entries as entry (entry.id)}
			<div class={`log-row level-${entry.level}`}>
				<time datetime={entry.time}>{formatTimeWithSeconds(entry.time)}</time>
				<span>{entry.level.toUpperCase()}</span>
				<p>{entry.message}</p>
				{#if entry.attributes && attributeText(entry)}
					<SystemLogAttributesButton attributes={entry.attributes} />
				{/if}
			</div>
		{:else}
			<div class="log-empty">Waiting for log entries</div>
		{/each}
	</div>
</section>
