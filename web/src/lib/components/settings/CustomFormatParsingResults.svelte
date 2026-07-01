<script lang="ts">
	import CustomFormatSpecChip from './CustomFormatSpecChip.svelte';
	import type { CustomFormatParsingResponse } from '$lib/settings/types';

	interface Props {
		result: CustomFormatParsingResponse;
	}

	type Row = [string, string | number | undefined];

	let { result }: Props = $props();

	const scoreRows = $derived<Row[]>([
		['Matched profile', result.matchedProfile?.name],
		['Calculated score', result.calculatedScore]
	]);
	const releaseRows = $derived<Row[]>([
		['Release title', result.release.releaseTitle],
		['Movie title', result.release.movieTitle],
		['Year', result.release.year],
		['Edition', result.release.edition],
		['Release type', result.details.releaseType],
		['Release group', result.release.releaseGroup],
		['Release hash', result.release.releaseHash]
	]);
	const qualityRows = $derived<Row[]>([
		['Quality', result.quality.quality],
		['Source', result.quality.source],
		['Resolution', result.quality.resolution],
		['Video codec', result.quality.videoCodec],
		['Audio codec', result.quality.audioCodec],
		['Audio channels', result.quality.audioChannels],
		['Version', result.quality.version],
		['Proper', yesNo(result.quality.proper)],
		['Repack', yesNo(result.quality.repack)],
		['Real', yesNo(result.quality.real)]
	]);

	function valueOrDash(value: string | number | undefined) {
		return value === undefined || value === '' ? '-' : String(value);
	}

	function yesNo(value: boolean) {
		return value ? 'Yes' : '-';
	}
</script>

{#snippet fieldList(rows: Row[])}
	<dl class="m-0 grid gap-2 sm:grid-cols-2">
		{#each rows as row (row[0])}
			<div class="grid grid-cols-[minmax(130px,0.4fr)_minmax(0,1fr)] gap-3 max-sm:grid-cols-1">
				<dt class="text-right font-extrabold text-muted-foreground max-sm:text-left">{row[0]}</dt>
				<dd class="m-0 min-w-0 text-foreground [overflow-wrap:anywhere]">
					{valueOrDash(row[1])}
				</dd>
			</div>
		{/each}
	</dl>
{/snippet}

<div class="grid gap-4.5">
	<section class="grid min-w-0 gap-2.5 mt-6">
		<h3 class="m-0 border-b border-border pb-2 text-2xl leading-tight text-foreground">Release</h3>
		{@render fieldList(releaseRows)}
	</section>

	<section class="grid min-w-0 gap-2.5">
		<h3 class="m-0 border-b border-border pb-2 text-2xl leading-tight text-foreground">Quality</h3>
		{@render fieldList(qualityRows)}
	</section>

	<section class="grid min-w-0 gap-2.5">
		<h3 class="m-0 border-b border-border pb-2 text-2xl leading-tight text-foreground">
			Languages
		</h3>
		<div class="flex flex-wrap gap-1.5">
			{#each result.languages as language (language)}
				<span
					class="rounded-md bg-primary px-2 py-1 text-xs font-extrabold text-primary-foreground"
				>
					{language}
				</span>
			{:else}
				<span
					class="rounded-md bg-primary px-2 py-1 text-xs font-extrabold text-primary-foreground"
				>
					-
				</span>
			{/each}
		</div>
	</section>

	<section class="grid min-w-0 gap-2.5">
		<h3
			class="m-0 flex items-center justify-between gap-3 border-b border-border pb-2 text-2xl leading-tight text-foreground"
		>
			<span>Matched custom formats</span>
			<span class="rounded-md bg-muted px-2 py-1 text-xs font-extrabold text-muted-foreground">
				{result.matchedCustomFormats.length}
			</span>
		</h3>
		{@render fieldList(scoreRows)}
		<div class="grid grid-cols-[repeat(auto-fit,minmax(min(100%,260px),1fr))] gap-3">
			{#each result.matchedCustomFormats as format (format.id)}
				<article class="grid min-h-0 gap-3 rounded-md border border-border bg-card p-3">
					<div class="flex items-start justify-between gap-3">
						<h4 class="m-0 text-xl text-muted-foreground">{format.name}</h4>
						<span
							class="min-w-[42px] rounded-md bg-primary/10 px-2 py-1.5 text-center text-sm leading-none font-bold text-primary"
						>
							{format.score}
						</span>
					</div>
					<div class="flex flex-wrap content-start gap-1.5">
						{#each format.matchedSpecs as spec (spec.id)}
							<CustomFormatSpecChip {spec} />
						{/each}
					</div>
				</article>
			{:else}
				<p class="m-0 text-sm leading-6 text-muted-foreground">No custom formats matched</p>
			{/each}
		</div>
	</section>
</div>
