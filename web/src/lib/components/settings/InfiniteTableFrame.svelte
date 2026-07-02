<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		hasMore: boolean;
		loading?: boolean;
		onLoadMore: () => void | Promise<void>;
		children: Snippet;
	}

	let { hasMore, loading = false, onLoadMore, children }: Props = $props();
	let frame: HTMLDivElement;
	let loadingMore = $state(false);

	async function handleScroll() {
		if (!frame || !hasMore || loading || loadingMore) return;
		const remaining = frame.scrollHeight - frame.scrollTop - frame.clientHeight;
		if (remaining > 80) return;

		loadingMore = true;
		try {
			await onLoadMore();
		} finally {
			loadingMore = false;
		}
	}
</script>

<div
	bind:this={frame}
	onscroll={handleScroll}
	class="min-h-0 overflow-auto rounded-md border border-border"
>
	{@render children()}
	{#if hasMore}
		<p class="m-0 border-t border-border p-2 text-center text-xs text-muted-foreground">
			{loading || loadingMore ? 'Loading more...' : 'Scroll for more'}
		</p>
	{/if}
</div>
