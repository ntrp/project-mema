<script lang="ts">
	import { resolve } from '$app/paths';
	import { providerDisplayName, providerPageUrl } from '$lib/settings/providerLinks';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		groups: MediaSearchGroup[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { groups, addingKey, actionLabel, onAdd }: Props = $props();

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

	function externalUrl(result: MediaSearchResult) {
		return providerPageUrl(result.externalProvider, result.type, result.externalId);
	}

	function externalLabel(result: MediaSearchResult) {
		return providerDisplayName(result.externalProvider);
	}
</script>

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
								{#if externalUrl(result)}
									<!-- eslint-disable svelte/no-navigation-without-resolve -->
									<a
										class="external-link"
										href={externalUrl(result)}
										target="_blank"
										rel="noreferrer"
										aria-label={`Open ${externalLabel(result)} page in a new tab`}
									>
										<span class="app-icon" aria-hidden="true">open_in_new</span>
										<span>{externalLabel(result)}</span>
									</a>
									<!-- eslint-enable svelte/no-navigation-without-resolve -->
								{/if}
								{#if group.sourceType === 'library'}
									<span class="status-pill">In library</span>
								{:else}
									<button
										type="button"
										class="add-action-button"
										disabled={addingKey === candidateKey(result)}
										onclick={() => onAdd(result)}
									>
										<span class="app-icon" aria-hidden="true">add</span>
										<span>{addingKey === candidateKey(result) ? 'Working' : actionLabel}</span>
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
