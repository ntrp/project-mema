<script lang="ts">
	import { providerPageUrl } from '$lib/settings/providerLinks';
	import type { MediaMetadataDetails, MediaMetadataFact } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		facts: MediaMetadataFact[];
	}

	let { detail, facts }: Props = $props();

	type InfoRow = {
		label: string;
		value: string | string[];
	};

	const factMap = $derived(new Map(facts.map((fact) => [fact.label, fact.value])));
	const score = $derived(detail.voteAverage ? Math.round(detail.voteAverage * 10) : undefined);
	const rows = $derived(infoRows(detail, factMap));
	const externalLinks = $derived(links(detail, factMap));

	function infoRows(details: MediaMetadataDetails, lookup: Map<string, string>): InfoRow[] {
		return [
			row('Status', details.status),
			row(
				details.type === 'movie' ? 'Release Dates' : 'First Aired',
				releaseDates(details, lookup)
			),
			row('Revenue', lookup.get('Revenue')),
			row('Budget', lookup.get('Budget')),
			row('Original Language', languageName(details.originalLanguage)),
			row('Production Countries', lookup.get('Production Countries')),
			row('Studios', lookup.get('Studios') ?? lookup.get('Networks'))
		].filter((item): item is InfoRow => Boolean(item));
	}

	function row(label: string, value: string | string[] | undefined): InfoRow | undefined {
		if (!value || (Array.isArray(value) && value.length === 0)) {
			return undefined;
		}
		return { label, value };
	}

	function releaseDates(details: MediaMetadataDetails, lookup: Map<string, string>) {
		const values = [
			lookup.get('Theatrical Release Date') ?? details.releaseDate,
			lookup.get('Digital Release Date'),
			lookup.get('Physical Release Date')
		].filter((value): value is string => Boolean(value));
		return [...new Set(values)];
	}

	function languageName(code?: string) {
		if (!code) return undefined;
		try {
			return new Intl.DisplayNames(undefined, { type: 'language' }).of(code) ?? code.toUpperCase();
		} catch {
			return code.toUpperCase();
		}
	}

	function links(details: MediaMetadataDetails, lookup: Map<string, string>) {
		const items = [];
		const providerUrl = providerPageUrl(details.externalProvider, details.type, details.externalId);
		if (providerUrl) {
			items.push({ label: details.externalProvider.toUpperCase(), href: providerUrl });
		}
		const imdbId = lookup.get('IMDb ID');
		if (imdbId) {
			items.push({ label: 'IMDb', href: `https://www.imdb.com/title/${imdbId}/` });
		}
		const wikidataId = lookup.get('Wikidata ID');
		if (wikidataId) {
			items.push({ label: 'Wikidata', href: `https://www.wikidata.org/wiki/${wikidataId}` });
		}
		const query = encodeURIComponent(details.title);
		items.push(
			{ label: 'Rotten Tomatoes', href: `https://www.rottentomatoes.com/search?search=${query}` },
			{ label: 'Letterboxd', href: `https://letterboxd.com/search/${query}/` }
		);
		return items;
	}
</script>

<aside class="metadata-overview-card" aria-label="Media facts">
	{#if score}
		<div class="metadata-score-row">
			<span class="metadata-score">
				<span class="tmdb-mark" aria-label="TMDB">TMDb</span>
				{score}%
			</span>
		</div>
	{/if}

	<div class="metadata-overview-card-rows">
		{#each rows as row (`${row.label}:${row.value}`)}
			<div>
				<strong>{row.label}</strong>
				{#if Array.isArray(row.value)}
					<span class="metadata-date-list">
						{#each row.value as value (value)}
							<span>{value}</span>
						{/each}
					</span>
				{:else}
					<span>{row.value}</span>
				{/if}
			</div>
		{/each}
	</div>

	{#if externalLinks.length > 0}
		<div class="metadata-source-row" aria-label="Metadata sources">
			{#each externalLinks as link (link.href)}
				<!-- eslint-disable svelte/no-navigation-without-resolve -->
				<a href={link.href} target="_blank" rel="noreferrer">{link.label}</a>
				<!-- eslint-enable svelte/no-navigation-without-resolve -->
			{/each}
		</div>
	{/if}
</aside>
