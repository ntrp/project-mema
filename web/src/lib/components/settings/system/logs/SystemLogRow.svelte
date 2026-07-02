<script lang="ts">
	import { formatDateTimeWithSeconds } from '$lib/settings/dateFormat';
	import type { SystemLogEntry, SystemLogLevel } from '$lib/settings/types';
	import SystemLogAttributesButton from './SystemLogAttributesButton.svelte';

	interface Props {
		entry: SystemLogEntry;
	}

	let { entry }: Props = $props();

	function attributeText(entry: SystemLogEntry) {
		if (!entry.attributes || Object.keys(entry.attributes).length === 0) {
			return '';
		}
		return JSON.stringify(entry.attributes);
	}

	function logLevelClass(level: SystemLogLevel) {
		switch (level) {
			case 'debug':
				return 'text-primary';
			case 'info':
				return 'text-primary';
			case 'warn':
				return 'text-secondary-foreground';
			case 'error':
				return 'text-destructive';
		}
	}
</script>

<div
	class="grid grid-cols-[max-content_44px_minmax(0,1fr)_20px] items-center gap-4 border-b border-border/30 px-1 py-px text-xs leading-tight text-foreground"
>
	<time
		class="self-center font-extrabold whitespace-nowrap text-muted-foreground"
		datetime={entry.time}
	>
		{formatDateTimeWithSeconds(entry.time)}
	</time>
	<span class={`self-center font-extrabold ${logLevelClass(entry.level)}`}>
		{entry.level.toUpperCase()}
	</span>
	<p class="m-0 self-center break-words">{entry.message}</p>
	{#if entry.attributes && attributeText(entry)}
		<SystemLogAttributesButton attributes={entry.attributes} />
	{/if}
</div>
