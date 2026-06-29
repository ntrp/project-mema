<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import type {
		MediaAdvancedSearchRequest,
		MediaSearchGroup,
		MediaSearchResult,
		MediaType,
		MetadataProvider
	} from '$lib/settings/types';

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
	let type = $state<MediaType | 'any'>('any');
	let year = $state('');
	let selectedProviderIds = $state<string[]>([]);

	const enabledProviders = $derived(metadataProviders.filter((provider) => provider.enabled));
	const resultCount = $derived(groups.reduce((count, group) => count + group.results.length, 0));

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
			type: type === 'any' ? undefined : type,
			year: Number.isFinite(parsedYear) ? parsedYear : undefined,
			providerIds: selectedProviderIds.length > 0 ? selectedProviderIds : undefined,
			limit: 30
		};
		await onSearch(request);
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.id ?? ''}:${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function groupDomId(group: MediaSearchGroup) {
		return `advanced-${group.sourceType}-${group.sourceName.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`;
	}

	function candidateKey(candidate: MediaSearchResult) {
		return `${candidate.type}:${candidate.title}:${candidate.year ?? ''}`;
	}

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w185${path}`;
	}
</script>

<section class="workspace-main advanced-search" aria-labelledby="advanced-search-title">
	<div class="page-heading">
		<p>Search</p>
		<h1 id="advanced-search-title">Advanced media search</h1>
	</div>

	<form class="advanced-search-form panel" onsubmit={submit}>
		<label class="wide">
			<span>Title</span>
			<input bind:value={query} placeholder="Movie or series title" autocomplete="off" />
		</label>
		<label>
			<span>Type</span>
			<select bind:value={type}>
				<option value="any">Any</option>
				<option value="movie">Movie</option>
				<option value="series">Series</option>
			</select>
		</label>
		<label>
			<span>Year</span>
			<input bind:value={year} inputmode="numeric" placeholder="Optional" />
		</label>
		<fieldset class="provider-picker wide">
			<legend>Metadata providers</legend>
			{#if enabledProviders.length > 0}
				<div class="provider-options">
					{#each enabledProviders as provider (provider.id)}
						<label>
							<input
								type="checkbox"
								checked={selectedProviderIds.includes(provider.id)}
								onchange={() => toggleProvider(provider.id)}
							/>
							<span>{provider.name}</span>
						</label>
					{/each}
				</div>
			{:else}
				<p>No enabled metadata providers are configured.</p>
			{/if}
		</fieldset>
		<div class="form-actions wide">
			<button type="submit" disabled={searching || query.trim().length === 0}>
				{searching ? 'Searching' : 'Search'}
			</button>
			<span class="muted">{resultCount} results</span>
		</div>
	</form>

	<div class="advanced-results" aria-label="Advanced search results">
		{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
			{#if group.results.length > 0}
				{@const headingId = groupDomId(group)}
				<section class="search-result-group" aria-labelledby={headingId}>
					<div class="section-heading">
						<h2 id={headingId}>{group.sourceName}</h2>
						<span>{group.sourceType}</span>
					</div>
					<div class="wide-card-list">
						{#each group.results as result (resultKey(result))}
							<article class="wide-media-card">
								<div class="wide-poster">
									{#if posterUrl(result.posterPath)}
										<img src={posterUrl(result.posterPath)} alt="" loading="lazy" />
									{:else}
										<div class="poster-placeholder">{result.type}</div>
									{/if}
								</div>
								<div class="wide-media-body">
									<div>
										<h3>
											{#if result.id}
												<a
													href={result.type === 'movie'
														? resolve('/movies/[id]', { id: result.id })
														: resolve('/series/[id]', { id: result.id })}>{result.title}</a
												>
											{:else if result.externalProvider && result.externalId}
												<a
													href={resolve('/media/[provider]/[type]/[externalId]', {
														provider: result.externalProvider,
														type: result.type,
														externalId: result.externalId
													})}>{result.title}</a
												>
											{:else}
												{result.title}
											{/if}
										</h3>
										<p>{result.type}{result.year ? ` · ${result.year}` : ''}</p>
									</div>
									{#if result.overview}
										<p>{result.overview}</p>
									{/if}
								</div>
								<div class="wide-media-actions">
									{#if group.sourceType === 'library'}
										<span class="status-pill">In library</span>
									{:else}
										<button
											type="button"
											disabled={addingKey === candidateKey(result)}
											onclick={() => onAdd(result)}
										>
											{addingKey === candidateKey(result) ? 'Working' : actionLabel}
										</button>
									{/if}
								</div>
							</article>
						{/each}
					</div>
				</section>
			{/if}
		{/each}
	</div>
</section>
