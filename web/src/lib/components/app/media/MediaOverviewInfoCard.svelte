<script lang="ts">
	import { providerPageUrl } from '$lib/settings/providerLinks';
	import { formatDate } from '$lib/settings/dateFormat';
	import { Button } from '$lib/components/ui/button';
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
	const externalLinks = $derived(links(detail));

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
			row('Production Countries', factList(lookup.get('Production Countries'))),
			row('Studios', factList(lookup.get('Studios') ?? lookup.get('Networks')))
		].filter((item): item is InfoRow => Boolean(item));
	}

	function row(label: string, value: string | string[] | undefined): InfoRow | undefined {
		if (!value || (Array.isArray(value) && value.length === 0)) {
			return undefined;
		}
		return { label, value };
	}

	function releaseDates(details: MediaMetadataDetails, lookup: Map<string, string>) {
		if (details.type === 'series') {
			return details.firstAirDate ? [formatDate(details.firstAirDate)] : [];
		}
		const values = [
			lookup.get('Theatrical Release Date') ?? details.releaseDate,
			lookup.get('Digital Release Date'),
			lookup.get('Physical Release Date')
		].filter((value): value is string => Boolean(value));
		return [...new Set(values)].map(formatDate);
	}

	function factList(value?: string) {
		if (!value) {
			return [];
		}
		const separator = value.includes('\n') ? /\n+/ : /,\s*/;
		return value
			.split(separator)
			.map((item) => item.trim())
			.filter(Boolean);
	}

	function languageName(code?: string) {
		if (!code) return undefined;
		try {
			return new Intl.DisplayNames(undefined, { type: 'language' }).of(code) ?? code.toUpperCase();
		} catch {
			return code.toUpperCase();
		}
	}

	function links(details: MediaMetadataDetails) {
		const items = [];
		const providerUrl = providerPageUrl(details.externalProvider, details.type, details.externalId);
		if (providerUrl) {
			items.push({ label: details.externalProvider.toUpperCase(), href: providerUrl });
		}
		return items;
	}
</script>

<aside class="overflow-hidden rounded-md border border-border bg-card" aria-label="Media facts">
	{#if score}
		<div class="flex flex-wrap items-center justify-center gap-2.5 px-4 py-3.5">
			<span class="inline-flex items-center gap-3 text-lg font-black text-foreground">
				<span
					class="inline-flex min-h-[22px] items-center rounded-md bg-primary px-[9px] py-[3px] text-xs leading-none font-black text-primary-foreground"
					aria-label="TMDB"
				>
					TMDb
				</span>
				{score}%
			</span>
		</div>
	{/if}

	<div class="border-t border-border">
		{#each rows as row (`${row.label}:${row.value}`)}
			<div
				class="grid grid-cols-[max-content_minmax(0,1fr)] gap-3 border-b border-border px-4 py-[13px]"
			>
				<strong class="[overflow-wrap:anywhere] text-sm whitespace-nowrap text-foreground">
					{row.label}
				</strong>
				{#if Array.isArray(row.value)}
					<span class="grid justify-items-end gap-1 text-right text-muted-foreground">
						{#each row.value as value (value)}
							<span>{value}</span>
						{/each}
					</span>
				{:else}
					<span class="[overflow-wrap:anywhere] text-right text-muted-foreground">{row.value}</span>
				{/if}
			</div>
		{/each}
	</div>

	{#if externalLinks.length > 0}
		<div
			class="flex flex-wrap items-center justify-center gap-2.5 px-4 py-3.5"
			aria-label="Metadata sources"
		>
			{#each externalLinks as link (link.href)}
				<Button variant="outline" size="xs" href={link.href} target="_blank" rel="noreferrer">
					{link.label}
				</Button>
			{/each}
		</div>
	{/if}
</aside>
