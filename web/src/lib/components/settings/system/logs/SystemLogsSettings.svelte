<script lang="ts">
	import { onMount, tick } from 'svelte';

	import { Card } from '$lib/components/ui/card';
	import { createLogLevelResource } from '$lib/features/settings/resources/systemGeneral.svelte';
	import type { SystemLogEntry, SystemLogLevel } from '$lib/settings/types';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import SystemLogRow from './SystemLogRow.svelte';
	import SystemLogsToolbar from './SystemLogsToolbar.svelte';

	type StreamEnvelope<T> = {
		data: T;
	};

	const levels: SystemLogLevel[] = ['debug', 'info', 'warn', 'error'];
	const levelOptions = levels.map((value) => ({ value, label: value.toUpperCase() }));
	const maxEntries = 500;

	interface Props {
		onConnectionChange?: (connected: boolean) => void;
	}

	let { onConnectionChange }: Props = $props();

	let entries = $state<SystemLogEntry[]>([]);
	let level = $state<SystemLogLevel>('info');
	const levelResource = createLogLevelResource();
	const loading = $derived(levelResource.query.isPending || levelResource.query.isFetching);
	const saving = $derived(levelResource.save.isPending);
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
		errorMessage = '';
		try {
			const { data: response } = await levelResource.query.refetch();
			if (!response) return;
			level = response.level;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load log level';
		} finally {
			/* query owns loading state */
		}
	}

	async function changeLevel(value: string) {
		const nextLevel = value as SystemLogLevel;
		errorMessage = '';
		try {
			const response = await levelResource.save.mutateAsync(nextLevel);
			level = response.level;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not update log level';
		} finally {
			/* mutation owns saving state */
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

	function trackViewport(node: HTMLDivElement) {
		logViewport = node;
		return { destroy: () => (logViewport = undefined) };
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
</script>

<Card class="gap-4 p-5" aria-label="Logs">
	<SectionHeading>
		{#snippet actions()}
			<SystemLogsToolbar
				{level}
				{levelOptions}
				{loading}
				{saving}
				{followLogs}
				onClearLogs={clearLogs}
				onEnableFollow={enableFollow}
				onLevelChange={changeLevel}
			/>
		{/snippet}
	</SectionHeading>

	{#if errorMessage}
		<p class="m-0 font-bold text-destructive">{errorMessage}</p>
	{/if}

	<div
		class="h-[min(58vh,620px)] overflow-auto bg-background px-1 py-px font-mono"
		data-log-viewer
		{@attach trackViewport}
		onscroll={handleLogScroll}
		aria-live="polite"
		aria-label="Application logs"
	>
		<div
			class="sticky top-0 z-10 grid grid-cols-[max-content_44px_minmax(0,1fr)_20px] items-center gap-4 border-b border-border bg-background px-1 py-1.5 text-xs leading-tight font-extrabold text-muted-foreground uppercase"
			aria-hidden="true"
		>
			<span>Time</span>
			<span>Level</span>
			<span>Message</span>
			<span></span>
		</div>
		{#each entries as entry (entry.id)}
			<SystemLogRow {entry} />
		{:else}
			<p class="m-0 p-5 text-center text-muted-foreground">Waiting for log entries</p>
		{/each}
	</div>
</Card>
