<script lang="ts">
	import { resolve } from '$app/paths';
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import DiscIcon from '@lucide/svelte/icons/disc-3';
	import MonitorPlayIcon from '@lucide/svelte/icons/monitor-play';
	import PauseIcon from '@lucide/svelte/icons/pause';
	import PlayIcon from '@lucide/svelte/icons/play';
	import SquareIcon from '@lucide/svelte/icons/square';
	import { formatDate } from '$lib/settings/dateFormat';
	import { displayLanguage } from '$lib/settings/languageDisplay';
	import type { MediaMetadataDetails, MediaMetadataFact } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		facts: MediaMetadataFact[];
	}

	let { detail, facts }: Props = $props();

	type InfoRow = {
		label: string;
		value: string | string[] | ReleaseDateItem[];
	};

	type ReleaseDateItem = {
		kind: 'cinema' | 'digital' | 'physical';
		label: string;
		date: string;
	};

	const factMap = $derived(new Map(facts.map((fact) => [fact.label, fact.value])));
	const score = $derived(detail.voteAverage ? Math.round(detail.voteAverage * 10) : undefined);
	const rows = $derived(infoRows(detail, factMap));

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

	function row(
		label: string,
		value: string | string[] | ReleaseDateItem[] | undefined
	): InfoRow | undefined {
		if (!value || (Array.isArray(value) && value.length === 0)) {
			return undefined;
		}
		return { label, value };
	}

	function releaseDates(
		details: MediaMetadataDetails,
		lookup: Map<string, string>
	): string[] | ReleaseDateItem[] {
		if (details.type === 'series') {
			return details.firstAirDate ? [formatDate(details.firstAirDate)] : [];
		}
		return releaseDateItems([
			{
				kind: 'cinema',
				label: 'Cinema',
				date: lookup.get('Theatrical Release Date') ?? details.releaseDate
			},
			{ kind: 'digital', label: 'Digital', date: lookup.get('Digital Release Date') },
			{ kind: 'physical', label: 'Physical', date: lookup.get('Physical Release Date') }
		]);
	}

	function releaseDateItems(items: (Omit<ReleaseDateItem, 'date'> & { date?: string })[]) {
		return items
			.filter((item): item is Omit<ReleaseDateItem, 'date'> & { date: string } =>
				Boolean(item.date)
			)
			.map((item) => ({ ...item, date: formatDate(item.date) }));
	}

	function releaseDateIcon(kind: ReleaseDateItem['kind']) {
		switch (kind) {
			case 'cinema':
				return ClapperboardIcon;
			case 'physical':
				return DiscIcon;
			default:
				return MonitorPlayIcon;
		}
	}

	function seriesStatusIcon(value: string) {
		switch (value.trim().toLowerCase()) {
			case 'continuing':
			case 'returning series':
				return PlayIcon;
			case 'ended':
				return SquareIcon;
			case 'on hold':
			case 'on_hold':
				return PauseIcon;
			default:
				return undefined;
		}
	}

	function isReleaseDateItems(value: string[] | ReleaseDateItem[]): value is ReleaseDateItem[] {
		return value.length > 0 && typeof value[0] !== 'string';
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
		return code ? displayLanguage(code) : undefined;
	}

	function discoverHref(label: string, value: string) {
		if (label === 'Studios') {
			return `${resolve('/discover/movies')}?studios=${encodeURIComponent(value)}`;
		}
		if (label === 'Original Language' && detail.originalLanguage) {
			return `${resolve('/discover/movies')}?originalLanguages=${encodeURIComponent(detail.originalLanguage)}`;
		}
		return undefined;
	}
</script>

<aside
	class="overflow-hidden rounded-md border border-border bg-card text-sm"
	aria-label="Media facts"
>
	{#if score}
		<div class="flex flex-wrap items-center justify-center gap-2.5 px-4 py-3.5">
			<span class="inline-flex items-center gap-3 text-sm font-black text-foreground">
				<span
					class="inline-flex min-h-5.5 items-center rounded-md bg-primary px-2.25 py-0.75 text-xs leading-none font-black text-primary-foreground"
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
				class="grid grid-cols-[max-content_minmax(0,1fr)] gap-3 border-b border-border px-4 py-3.25"
			>
				<strong class="wrap-anywhere text-sm whitespace-nowrap text-foreground">
					{row.label}
				</strong>
				{#if Array.isArray(row.value)}
					<span class="grid justify-items-end gap-1 text-right text-muted-foreground">
						{#if isReleaseDateItems(row.value)}
							{#each row.value as value (`${value.kind}:${value.date}`)}
								{@const Icon = releaseDateIcon(value.kind)}
								<span class="inline-flex items-center justify-end gap-1.5">
									<Icon aria-hidden="true" class="size-3.5 text-foreground" />
									<span class="sr-only">{value.label}</span>
									<span>{value.date}</span>
								</span>
							{/each}
						{:else}
							{#each row.value as value (value)}
								{@const href = discoverHref(row.label, value)}
								{#if href}
									<a class="text-muted-foreground hover:text-foreground" {href}>{value}</a>
								{:else}
									<span>{value}</span>
								{/if}
							{/each}
						{/if}
					</span>
				{:else}
					{@const href = discoverHref(row.label, row.value)}
					<span
						class="wrap-anywhere inline-flex items-center justify-end gap-1.5 text-right text-muted-foreground"
					>
						{#if row.label === 'Status' && detail.type === 'series' && seriesStatusIcon(row.value)}
							{@const Icon = seriesStatusIcon(row.value)}
							<Icon aria-hidden="true" class="size-3.5 text-foreground" />
						{/if}
						{#if href}
							<a class="text-muted-foreground hover:text-foreground" {href}>{row.value}</a>
						{:else}
							<span>{row.value}</span>
						{/if}
					</span>
				{/if}
			</div>
		{/each}
	</div>
</aside>
