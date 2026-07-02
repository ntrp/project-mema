<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MatchInfo } from './releaseCandidateDisplay';
	import { signedScore } from './releaseCandidateDisplay';

	interface Props {
		match: MatchInfo;
	}

	let { match }: Props = $props();
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span
				{...props}
				class="inline-flex min-w-10 cursor-help justify-end font-mono tabular-nums"
				aria-label="Show score contributors"
			>
				{signedScore(match.score)}
			</span>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content class="max-w-80">
		<div class="grid gap-1 text-left">
			<span class="font-bold">Score contributors</span>
			{#each match.scoreContributors as contributor (`${contributor.label}:${contributor.score}`)}
				<div class="grid grid-cols-[1fr_auto] gap-4">
					<span>{contributor.label}</span>
					<span class="font-mono tabular-nums">{signedScore(contributor.score)}</span>
				</div>
			{/each}
		</div>
	</Tooltip.Content>
</Tooltip.Root>
