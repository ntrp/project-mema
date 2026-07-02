<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MatchInfo } from './releaseCandidateDisplay';
	import { signedScore } from './releaseCandidateDisplay';

	interface Props {
		match: MatchInfo;
	}

	let { match }: Props = $props();
	let open = $state(false);

	const contributorGroups = $derived(
		[
			{ label: 'Custom formats', values: match.customFormatContributors ?? [] },
			{ label: 'Languages', values: match.languageContributors ?? [] }
		].filter((group) => group.values.length > 0)
	);
</script>

<Tooltip.Root bind:open>
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
	{#if open}
		<Tooltip.Content class="max-w-80">
			<div class="grid gap-1 text-left">
				{#if contributorGroups.length > 0}
					{#each contributorGroups as group (group.label)}
						<span class="pt-1 font-bold first:pt-0">{group.label}</span>
						{#each group.values as contributor (`${group.label}:${contributor.label}:${contributor.score}`)}
							<div class="grid grid-cols-[1fr_auto] gap-4">
								<span>{contributor.label}</span>
								<span class="font-mono tabular-nums">{signedScore(contributor.score)}</span>
							</div>
						{/each}
					{/each}
				{:else}
					<span>No matched custom formats or scored languages.</span>
				{/if}
			</div>
		</Tooltip.Content>
	{/if}
</Tooltip.Root>
