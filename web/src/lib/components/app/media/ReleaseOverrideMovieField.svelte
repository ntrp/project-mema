<script lang="ts">
	import { searchMedia } from '$lib/settings/api';
	import type { MediaSearchResult } from '$lib/settings/types';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';

	interface Props {
		value: string;
	}

	let { value = $bindable() }: Props = $props();
	let results = $state<MediaSearchResult[]>([]);
	let loading = $state(false);
	let open = $state(false);
	let selectedIndex = $state(-1);
	let requestNumber = 0;

	const trimmed = $derived(value.trim());
	const selected = $derived(selectedIndex >= 0 ? results[selectedIndex] : undefined);

	function handleInput() {
		open = true;
		selectedIndex = -1;
		void searchMovies(trimmed);
	}

	async function searchMovies(query: string) {
		const current = ++requestNumber;
		if (query.length < 2) {
			results = [];
			loading = false;
			return;
		}
		loading = true;
		try {
			const found = await searchMedia({ query, type: 'movie' });
			if (current === requestNumber) {
				results = found.slice(0, 6);
			}
		} finally {
			if (current === requestNumber) {
				loading = false;
			}
		}
	}

	function choose(result: MediaSearchResult) {
		value = result.title;
		open = false;
		selectedIndex = -1;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
			if (results.length === 0) return;
			event.preventDefault();
			open = true;
			selectedIndex =
				event.key === 'ArrowDown'
					? Math.min(selectedIndex + 1, results.length - 1)
					: Math.max(selectedIndex - 1, 0);
		}
		if (event.key === 'Enter' && selected) {
			event.preventDefault();
			choose(selected);
		}
		if (event.key === 'Escape') {
			open = false;
		}
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}
</script>

<div class="relative grid gap-1.5">
	<Label for="override-movie">Movie</Label>
	<Input
		id="override-movie"
		bind:value
		autocomplete="off"
		oninput={handleInput}
		onfocus={() => {
			open = true;
			if (trimmed.length >= 2 && results.length === 0) void searchMovies(trimmed);
		}}
		onkeydown={handleKeydown}
		onblur={() => {
			window.setTimeout(() => (open = false), 120);
		}}
	/>
	{#if open && trimmed.length >= 2}
		<div
			class="absolute inset-x-0 top-[calc(100%+6px)] z-30 max-h-72 overflow-auto rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-md"
			role="listbox"
			aria-label="Movie matches"
		>
			{#if results.length > 0}
				{#each results as result, index (resultKey(result))}
					<Button
						type="button"
						variant="ghost"
						role="option"
						aria-selected={index === selectedIndex}
						class={cn(
							'grid min-h-10 w-full justify-start rounded-sm px-2 py-1.5 text-left',
							index === selectedIndex && 'bg-accent text-accent-foreground'
						)}
						onpointerdown={(event) => event.preventDefault()}
						onclick={() => choose(result)}
					>
						<span class="truncate">{result.title}</span>
						{#if result.year}
							<small class="text-xs text-muted-foreground">{result.year}</small>
						{/if}
					</Button>
				{/each}
			{:else}
				<div class="px-2 py-1 text-xs font-bold text-muted-foreground uppercase">
					{loading ? 'Searching' : 'No matches'}
				</div>
			{/if}
		</div>
	{/if}
</div>
