<script lang="ts">
	import { onMount } from 'svelte';
	import { Badge } from '$lib/components/ui/badge';
	import {
		categoryBadges,
		privacyBadgeClass,
		protocolBadgeClass
	} from '$lib/components/settings/indexers/indexerCatalogPresentation';
	import { flattenCategories, privacyLabel } from './indexerCatalogFilters';
	import IndexerCatalogSupportIcons from './IndexerCatalogSupportIcons.svelte';
	import type { IndexerCatalogEntry } from '$lib/settings/types';

	interface Props {
		entries: IndexerCatalogEntry[];
		hasMore?: boolean;
		onEndReached?: () => void;
		onSelect: (_entry: IndexerCatalogEntry) => void;
	}

	let { entries, hasMore = false, onEndReached, onSelect }: Props = $props();
	let scroller: HTMLDivElement | undefined = $state();
	let sentinel: HTMLTableRowElement | undefined = $state();
	let pendingMore = $state(false);
	let lastEntryCount = $state(0);

	$effect(() => {
		if (entries.length !== lastEntryCount) {
			pendingMore = false;
			lastEntryCount = entries.length;
		}
	});

	function handleRowKeydown(event: KeyboardEvent, entry: IndexerCatalogEntry) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onSelect(entry);
		}
	}

	function requestMore() {
		if (!hasMore || pendingMore) return;
		pendingMore = true;
		onEndReached?.();
	}

	function handleScroll(event: Event) {
		const target = event.currentTarget as HTMLDivElement;
		if (target.scrollTop + target.clientHeight >= target.scrollHeight - 160) {
			requestMore();
		}
	}

	onMount(() => {
		if (!scroller || !sentinel) return;
		const observer = new IntersectionObserver(
			(observed) => {
				if (observed.some((entry) => entry.isIntersecting)) {
					requestMore();
				}
			},
			{ root: scroller, rootMargin: '160px' }
		);
		observer.observe(sentinel);
		return () => observer.disconnect();
	});
</script>

<div
	bind:this={scroller}
	data-testid="indexer-catalog-table"
	class="max-h-[min(560px,calc(100vh-300px))] overflow-auto rounded-md border border-border"
	onscroll={handleScroll}
>
	<table class="w-full min-w-215 table-auto border-collapse text-sm">
		<thead class="sticky top-0 z-10 bg-card text-left text-xs font-extrabold text-muted-foreground">
			<tr class="border-b border-border">
				<th class="w-px px-3 py-2">Protocol</th>
				<th class="w-px px-3 py-2">Name</th>
				<th class="px-3 py-2">Description</th>
				<th class="w-px px-3 py-2">Privacy</th>
				<th class="w-px px-3 py-2">Language</th>
				<th class="w-px px-3 py-2">Supports</th>
				<th class="w-px px-3 py-2">Categories</th>
			</tr>
		</thead>
		<tbody>
			{#each entries as entry (entry.definitionId)}
				<tr
					class="cursor-pointer border-b border-border last:border-0 hover:bg-muted/60 focus-visible:bg-muted/60 focus-visible:outline-none"
					tabindex="0"
					role="button"
					onclick={() => onSelect(entry)}
					onkeydown={(event) => handleRowKeydown(event, entry)}
				>
					<td class="w-px px-3 py-2">
						<Badge variant="outline" class={protocolBadgeClass(entry.protocol)}
							>{entry.protocol}</Badge
						>
					</td>
					<td class="w-px px-3 py-2">
						<div class="whitespace-nowrap font-bold text-foreground">{entry.name}</div>
					</td>
					<td class="px-3 py-2">
						<div class="line-clamp-2 text-xs text-muted-foreground">{entry.description}</div>
					</td>
					<td class="w-px px-3 py-2">
						<Badge variant="outline" class={'uppercase ' + privacyBadgeClass(entry.privacy)}
							>{privacyLabel(entry.privacy)}</Badge
						>
					</td>
					<td class="w-px px-3 py-2 whitespace-nowrap">{entry.language}</td>
					<td class="w-px px-3 py-2">
						<IndexerCatalogSupportIcons {entry} />
					</td>
					<td class="w-px px-3 py-2">
						<div class="flex max-w-72 flex-nowrap gap-1">
							{#each categoryBadges(entry, flattenCategories) as category (category.id)}
								<Badge variant="outline" class="border-sky-500/50 bg-sky-500/10 text-sky-300">
									{category.name}
								</Badge>
							{/each}
						</div>
					</td>
				</tr>
			{:else}
				<tr>
					<td class="px-3 py-8 text-center text-muted-foreground" colspan="7">
						No catalog indexers match the selected filters.
					</td>
				</tr>
			{/each}
			{#if hasMore}
				<tr bind:this={sentinel} aria-hidden="true">
					<td class="h-2 p-0" colspan="7"></td>
				</tr>
			{/if}
		</tbody>
	</table>
</div>
