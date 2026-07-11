<script lang="ts">
	import { tick } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { basename } from './dlnaDecisionTrace';
	import { Input } from '$lib/components/ui/input';
	import { listMediaItems } from '$lib/features/library/api';
	import { autocompleteMedia } from '$lib/settings/domains/searchMetadata';
	import type { MediaItem, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		mediaPath: string;
		onMediaPath: (_value: string) => void;
		onSelectedMediaId?: (_value: string) => void;
	}

	type MediaSuggestion = {
		id?: string;
		title: string;
		path: string;
	};

	let { mediaPath, onMediaPath, onSelectedMediaId = () => {} }: Props = $props();

	let open = $state(false);
	let query = $state('');
	let searchInput = $state<HTMLInputElement | null>(null);

	const normalizedQuery = $derived(query.trim());
	const canSearch = $derived(normalizedQuery.length >= 2);
	const mediaItems = createQuery(() => ({
		queryKey: ['settings', 'dlna', 'trace-media-items'],
		queryFn: listMediaItems,
		enabled: canSearch && !mediaPath
	}));
	const mediaMatches = createQuery(() => ({
		queryKey: ['settings', 'dlna', 'trace-media', normalizedQuery],
		queryFn: () => autocompleteMedia(normalizedQuery, 'library'),
		enabled: canSearch && !mediaPath
	}));
	const loading = $derived(mediaItems.isFetching || mediaMatches.isFetching);
	const selectedFileName = $derived(mediaPath ? basename(mediaPath) : '');
	const searchedSuggestions = $derived.by(() =>
		buildSearchSuggestions(mediaMatches.data ?? [], mediaItems.data?.items ?? [])
	);
	const suggestions = $derived(mediaPath ? [] : searchedSuggestions);

	function buildSearchSuggestions(groups: { results: MediaSearchResult[] }[], items: MediaItem[]) {
		const itemsById = new Map(items.map((item) => [item.id, item]));
		const suggestions = groups.flatMap((group) =>
			group.results.flatMap((result) => {
				if (!result.id) return [];
				const item = itemsById.get(result.id);
				return item?.filePaths[0]
					? [
							{
								id: result.id,
								title: result.title,
								path: item.filePaths[0]
							} satisfies MediaSuggestion
						]
					: [];
			})
		);
		return [
			...new Map(suggestions.map((suggestion) => [suggestion.path, suggestion])).values()
		].slice(0, 12);
	}

	function openDropdown() {
		if (mediaPath || !canSearch) return;
		open = true;
	}

	function handleFocusOut(event: globalThis.FocusEvent) {
		const currentTarget = event.currentTarget as HTMLElement | null;
		if (!currentTarget?.contains(event.relatedTarget as globalThis.Node | null)) {
			open = false;
		}
	}

	function chooseMedia(suggestion: MediaSuggestion) {
		onSelectedMediaId(suggestion.id ?? '');
		onMediaPath(suggestion.path);
		query = '';
		open = false;
	}

	async function clearMedia() {
		onSelectedMediaId('');
		onMediaPath('');
		query = '';
		open = false;
		await tick();
		searchInput?.focus();
	}

	async function syncQueryOpenState() {
		await tick();
		open = canSearch;
	}
</script>

<div class="relative min-w-0 w-full" onfocusout={handleFocusOut}>
	{#if mediaPath}
		<button
			type="button"
			class="grid h-9 w-full gap-0.5 rounded-md border border-border bg-background px-3 py-1.5 text-left text-sm shadow-xs transition-colors hover:bg-muted/60"
			onclick={clearMedia}
			aria-label="Change selected media file"
		>
			<span class="truncate font-medium">{selectedFileName}</span>
		</button>
	{:else}
		<div class="relative">
			<Input
				bind:ref={searchInput}
				bind:value={query}
				class="pr-10"
				placeholder="Type at least 2 characters to search media"
				autocomplete="off"
				onfocus={openDropdown}
				onclick={openDropdown}
				oninput={syncQueryOpenState}
			/>
			{#if open}
				<div
					class="absolute left-0 top-[calc(100%+6px)] z-30 grid w-full gap-1 rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-md"
				>
					{#if loading}
						<p class="m-0 px-2 py-2 text-sm text-muted-foreground">Searching media…</p>
					{:else if suggestions.length > 0}
						{#each suggestions as suggestion (suggestion.path)}
							<button
								type="button"
								class="grid gap-0.5 rounded-sm px-2 py-2 text-left hover:bg-accent hover:text-accent-foreground"
								onmousedown={(event) => event.preventDefault()}
								onclick={() => chooseMedia(suggestion)}
							>
								<span class="truncate text-sm font-medium">{suggestion.title}</span>
								<span class="truncate text-xs text-muted-foreground">{suggestion.path}</span>
							</button>
						{/each}
					{:else}
						<p class="m-0 px-2 py-2 text-sm text-muted-foreground">
							{canSearch ? 'No media matches.' : 'Type at least 2 characters to search media.'}
						</p>
					{/if}
				</div>
			{/if}
		</div>
	{/if}
</div>
