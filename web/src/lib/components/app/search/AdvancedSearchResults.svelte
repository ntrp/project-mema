<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import { resolve } from '$app/paths';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import { providerDisplayName, providerPageUrl } from '$lib/settings/providerLinks';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';
	import MediaAddButton from '$lib/components/app/media/shared/MediaAddButton.svelte';
	import PosterPlaceholder from '$lib/components/app/media/posters/PosterPlaceholder.svelte';

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

<div class="grid gap-[22px]" aria-label="Advanced search results">
	{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
		{#if group.results.length > 0}
			{@const headingId = groupDomId(group)}
			<section aria-labelledby={headingId}>
				<SectionHeading title={group.sourceName} titleId={headingId}>
					{#snippet actions()}
						<span>{group.sourceType}</span>
					{/snippet}
				</SectionHeading>
				<div class="grid gap-2.5">
					{#each group.results as result (resultKey(result))}
						<article
							class="grid items-center gap-3.5 rounded-md border border-border bg-muted p-2.5 md:grid-cols-[82px_minmax(0,1fr)_auto]"
						>
							<div class="aspect-[2/3] overflow-hidden rounded-md bg-card">
								{#if posterUrl(result.posterPath)}
									<img
										class="block h-full w-full object-cover"
										src={posterUrl(result.posterPath)}
										alt=""
										loading="lazy"
									/>
								{:else}
									<PosterPlaceholder label={result.type} class="h-full min-h-0" />
								{/if}
							</div>
							<div class="grid min-w-0 gap-2">
								<div>
									<h3 class="m-0 text-base leading-tight">
										{#if result.id}
											<a
												class="text-foreground no-underline hover:text-primary-hover"
												href={result.type === 'movie'
													? resolve('/movies/[id]', { id: result.id })
													: resolve('/series/[id]', { id: result.id })}>{result.title}</a
											>
										{:else if result.externalProvider && result.externalId}
											<a
												class="text-foreground no-underline hover:text-primary-hover"
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
									<p class="m-0 text-sm text-muted-foreground">
										{result.type}{result.year ? ` · ${result.year}` : ''}
									</p>
								</div>
								{#if result.overview}
									<p class="line-clamp-2 m-0 text-sm text-muted-foreground">{result.overview}</p>
								{/if}
							</div>
							<div class="flex items-center justify-end gap-2.5">
								{#if externalUrl(result)}
									<Button
										variant="outline"
										size="sm"
										href={externalUrl(result)}
										target="_blank"
										rel="noreferrer"
										aria-label={`Open ${externalLabel(result)} page in a new tab`}
									>
										<ExternalLinkIcon aria-hidden="true" />
										<span>{externalLabel(result)}</span>
									</Button>
								{/if}
								{#if group.sourceType === 'library'}
									<StatusPill tone="success">In library</StatusPill>
								{:else}
									<MediaAddButton
										{result}
										adding={addingKey === candidateKey(result)}
										label={actionLabel}
										{onAdd}
									/>
								{/if}
							</div>
						</article>
					{/each}
				</div>
			</section>
		{/if}
	{/each}
</div>
