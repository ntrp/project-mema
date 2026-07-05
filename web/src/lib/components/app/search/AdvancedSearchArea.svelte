<script lang="ts">
	import { onMount } from 'svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import type {
		MediaAdvancedSearchRequest,
		MediaSearchGroup,
		MediaSearchResult,
		MetadataProvider
	} from '$lib/settings/types';
	import AdvancedSearchResults from './AdvancedSearchResults.svelte';

	type SearchTarget = 'all' | 'movie' | 'series' | 'people';

	interface Props {
		initialQuery: string;
		metadataProviders: MetadataProvider[];
		groups: MediaSearchGroup[];
		searching: boolean;
		addingKey?: string;
		actionLabel: string;
		onSearch: (_request: MediaAdvancedSearchRequest) => void | Promise<void>;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let {
		initialQuery,
		metadataProviders,
		groups,
		searching,
		addingKey,
		actionLabel,
		onSearch,
		onAdd
	}: Props = $props();

	let query = $state('');
	let target = $state<SearchTarget>('all');
	let year = $state('');
	let selectedProviderIds = $state<string[]>([]);

	const targetOptions: { value: SearchTarget; label: string }[] = [
		{ value: 'all', label: 'All' },
		{ value: 'movie', label: 'Movies' },
		{ value: 'series', label: 'Series' },
		{ value: 'people', label: 'People' }
	];
	const enabledProviders = $derived(metadataProviders.filter((provider) => provider.enabled));
	const resultCount = $derived(
		groups.reduce((count, group) => count + group.results.length + (group.people?.length ?? 0), 0)
	);
	const includeMedia = $derived(target !== 'people');
	const includePeople = $derived(target === 'all' || target === 'people');

	onMount(() => {
		query = initialQuery;
		selectedProviderIds = enabledProviders.map((provider) => provider.id);
		if (query.trim().length > 0) {
			void submitSearch();
		}
	});

	function toggleProvider(id: string) {
		selectedProviderIds = selectedProviderIds.includes(id)
			? selectedProviderIds.filter((providerId) => providerId !== id)
			: [...selectedProviderIds, id];
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		void submitSearch();
	}

	async function submitSearch() {
		const parsedYear = Number.parseInt(year, 10);
		const request: MediaAdvancedSearchRequest = {
			query: query.trim(),
			type: target === 'movie' ? 'movie' : target === 'series' ? 'serie' : undefined,
			includeMedia,
			includePeople,
			year: includeMedia && Number.isFinite(parsedYear) ? parsedYear : undefined,
			providerIds: selectedProviderIds.length > 0 ? selectedProviderIds : undefined,
			limit: 30
		};
		await onSearch(request);
	}
</script>

<section class="advanced-search grid min-w-0 gap-[18px]" aria-labelledby="advanced-search-title">
	<PageHeading eyebrow="Search" title="Advanced search" titleId="advanced-search-title" />

	<form
		class="grid gap-4 rounded-md border border-border bg-card p-5 md:grid-cols-[minmax(0,1fr)_180px_150px]"
		onsubmit={submit}
	>
		<label class="grid gap-1.5 md:col-span-3">
			<span class="text-sm font-bold text-muted-foreground">Title</span>
			<Input bind:value={query} placeholder="Title or person name" autocomplete="off" />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Search</span>
			<SettingsSelect
				value={target}
				options={targetOptions}
				onValueChange={(value) => (target = value as SearchTarget)}
			/>
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Year</span>
			<Input
				bind:value={year}
				inputmode="numeric"
				placeholder="Optional"
				disabled={!includeMedia}
			/>
		</label>
		<fieldset class="grid gap-2.5 rounded-md border border-border p-3 md:col-span-3">
			<legend class="px-1 text-sm font-extrabold text-muted-foreground">Metadata providers</legend>
			{#if enabledProviders.length > 0}
				<div class="flex flex-wrap gap-2.5">
					{#each enabledProviders as provider (provider.id)}
						<label
							class="flex items-center gap-2 rounded-md border border-border bg-background px-2.5 py-1.5"
						>
							<Checkbox
								checked={selectedProviderIds.includes(provider.id)}
								onCheckedChange={() => toggleProvider(provider.id)}
							/>
							<span>{provider.name}</span>
						</label>
					{/each}
				</div>
			{:else}
				<p class="m-0 text-sm text-muted-foreground">
					No enabled metadata providers are configured.
				</p>
			{/if}
		</fieldset>
		<div class="flex items-center gap-3 md:col-span-3">
			<Button type="submit" disabled={searching || query.trim().length === 0}>
				{searching ? 'Searching' : 'Search'}
			</Button>
			<span class="text-sm text-muted-foreground">{resultCount} results</span>
		</div>
	</form>

	<AdvancedSearchResults {groups} {addingKey} {actionLabel} {onAdd} />
</section>
