<script lang="ts">
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
	<dl>
		{#each rows as row (row[0])}
			<div>
				<dt>{row[0]}</dt>
				<dd>{valueOrDash(row[1])}</dd>
			</div>
		{/each}
	</dl>
{/snippet}

<div class="test-parsing-results">
	<section class="test-parsing-section">
		<h3>Release</h3>
		{@render fieldList(releaseRows)}
	</section>

	<section class="test-parsing-section">
		<h3>Quality</h3>
		{@render fieldList(qualityRows)}
	</section>

	<section class="test-parsing-section">
		<h3>Languages</h3>
		<div class="test-parsing-tags">
			{#each result.languages as language (language)}
				<span>{language}</span>
			{:else}
				<span>-</span>
			{/each}
		</div>
	</section>

	<section class="test-parsing-section">
		<h3 class="test-parsing-section-title">
			<span>Matched custom formats</span>
			<span>{result.matchedCustomFormats.length}</span>
		</h3>
		{@render fieldList(scoreRows)}
		<div class="test-parsing-match-list">
			{#each result.matchedCustomFormats as format (format.id)}
				<article class="custom-format-card test-parsing-format-card">
					<div class="custom-format-card-header">
						<h4>{format.name}</h4>
						<span class="custom-format-score">{format.score}</span>
					</div>
					<div class="custom-format-tags">
						{#each format.matchedSpecs as spec (spec.id)}
							<span title={`${spec.type}: ${spec.value}`}>{spec.name}</span>
						{/each}
					</div>
				</article>
			{:else}
				<p class="empty">No custom formats matched</p>
			{/each}
		</div>
	</section>
</div>
