<script lang="ts">
	import { tick } from 'svelte';
	import type { ReleaseSearchLogEntry } from '$lib/components/app/media/release-search/releaseSearchLog';

	interface Props {
		messages: ReleaseSearchLogEntry[];
	}

	let { messages }: Props = $props();

	let expanded = $state(false);
	let viewport = $state<HTMLDivElement>();
	const visibleMessage = $derived(messages.at(-1) ?? messages[0]);

	$effect(() => {
		messages;
		expanded;
		void scrollToBottom();
	});

	async function scrollToBottom() {
		await tick();
		if (viewport) {
			viewport.scrollTop = viewport.scrollHeight;
		}
	}

	function lineText(entry?: ReleaseSearchLogEntry, index?: number) {
		if (!entry) {
			return 'Press search to start';
		}
		const prefix = entry.timestamp ? `[${entry.timestamp}] ` : '';
		if (!entry.resultMessage) {
			return `${prefix}${entry.message}`;
		}
		const branch = hasLaterResult(index) ? '├─ ' : '└─ ';
		const duration = entry.durationMs !== undefined ? ` --(${entry.durationMs}ms)--> ` : ' --> ';
		return `${prefix}${branch}${entry.message}${duration}${entry.resultMessage}`;
	}

	function hasLaterResult(index?: number) {
		if (index === undefined) {
			return false;
		}
		return messages.slice(index + 1).some((message) => message.resultMessage);
	}
</script>

<button
	type="button"
	class="w-full rounded-md bg-black px-3 py-2.5 text-left font-mono text-xs font-medium text-white"
	aria-expanded={expanded}
	onclick={() => (expanded = !expanded)}
>
	{#if expanded}
		<div bind:this={viewport} class="max-h-40 overflow-y-auto">
			<div class="grid min-h-40 content-end gap-1">
				{#each messages as message, index (message.id)}
					<p class="m-0 flex items-center gap-2">
						<span class="min-w-0 truncate">{lineText(message, index)}</span>
						{#if message.cacheHit}
							<span
								class="shrink-0 rounded-sm bg-emerald-500 px-1.5 py-0.5 text-[10px] leading-none font-bold text-black uppercase"
							>
								cache hit
							</span>
						{/if}
					</p>
				{/each}
			</div>
		</div>
	{:else}
		<p class="m-0 flex items-center gap-2">
			<span class="min-w-0 truncate">{lineText(visibleMessage)}</span>
			{#if visibleMessage?.cacheHit}
				<span
					class="shrink-0 rounded-sm bg-emerald-500 px-1.5 py-0.5 text-[10px] leading-none font-bold text-black uppercase"
				>
					cache hit
				</span>
			{/if}
		</p>
	{/if}
</button>
