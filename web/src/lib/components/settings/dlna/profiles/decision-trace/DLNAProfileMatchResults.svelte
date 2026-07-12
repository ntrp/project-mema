<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import type { DLNAProfileMatchCandidate, DLNAProfileMatchView } from './dlnaDecisionTrace';

	interface Props {
		match?: DLNAProfileMatchView;
	}

	let { match }: Props = $props();

	function isOutcomeCandidate(candidate: DLNAProfileMatchCandidate) {
		return candidate.profileId === match?.profileId;
	}

	function candidateStatus(candidate: DLNAProfileMatchCandidate) {
		if (candidate.selected) return 'Selected automatically';
		if (isOutcomeCandidate(candidate)) return 'Used by selection';
		return candidate.qualified ? 'Qualified' : 'Did not qualify';
	}

	function fieldLabel(field: string) {
		return field === 'any' ? 'Any request signal (legacy token)' : field;
	}
</script>

<section class="grid gap-3" aria-label="DLNA profile match">
	<h3 class="m-0 text-sm font-semibold">Profile match</h3>
	{#if match}
		<div class="grid gap-3 rounded-md border border-border bg-card p-4">
			<div class="flex flex-wrap items-start justify-between gap-3">
				<div class="grid gap-1">
					<span class="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
						>Matched profile</span
					>
					<strong class="text-base">{match.profileName}</strong>
					<span class="text-xs text-muted-foreground">{match.profileId}</span>
				</div>
				<Badge variant="default">{match.selectionMethod}</Badge>
			</div>
			<div class="grid gap-3 text-sm sm:grid-cols-3">
				<div>
					<span class="text-muted-foreground">Score</span>
					<p class="m-0 font-medium">{match.score}</p>
				</div>
				<div>
					<span class="text-muted-foreground">Strongest rule</span>
					<p class="m-0 font-medium">{match.winningRule || 'Not applicable'}</p>
				</div>
				<div>
					<span class="text-muted-foreground">Reason</span>
					<p class="m-0 font-medium">{match.matchReason || 'No automatic match'}</p>
				</div>
			</div>
			{#if match.fallbackPath}
				<p class="m-0 text-xs text-muted-foreground">Fallback path: {match.fallbackPath}</p>
			{/if}
		</div>

		<div class="grid gap-2">
			<div>
				<h4 class="m-0 text-sm font-semibold">Candidate profiles</h4>
				<p class="m-0 text-xs text-muted-foreground">
					Profiles must reach their minimum score. Qualified profiles are ranked by priority, then
					score, then profile ID.
				</p>
			</div>
			{#each match.candidates as candidate (candidate.profileId)}
				<details class="rounded-md border border-border bg-card">
					<summary class="cursor-pointer list-none px-4 py-3">
						<div class="flex flex-wrap items-center justify-between gap-3">
							<div class="grid gap-0.5">
								<span class="text-sm font-semibold">{candidate.profileName}</span>
								<span class="text-xs text-muted-foreground">{candidate.profileId}</span>
							</div>
							<div class="flex flex-wrap items-center gap-2">
								<Badge variant="outline">Score {candidate.score}/{candidate.minimumScore}</Badge>
								<Badge variant="secondary">Priority {candidate.priority}</Badge>
								<Badge
									variant={candidate.selected || isOutcomeCandidate(candidate)
										? 'default'
										: candidate.qualified
											? 'secondary'
											: 'destructive'}>{candidateStatus(candidate)}</Badge
								>
							</div>
						</div>
					</summary>
					<div class="border-t border-border px-4 py-3">
						{#if candidate.ruleTrace.length > 0}
							<ul class="m-0 grid list-none gap-3 p-0">
								{#each candidate.ruleTrace as rule, index (`${candidate.profileId}-${index}`)}
									<li class="grid gap-1 text-sm">
										<div class="flex flex-wrap items-center justify-between gap-2">
											<span class="font-medium">{fieldLabel(rule.field)}: {rule.rule}</span>
											<Badge variant={rule.result === 'pass' ? 'default' : 'outline'}>
												{rule.result === 'pass' ? `+${rule.score}` : 'No match'}
											</Badge>
										</div>
										<span class="break-all text-xs text-muted-foreground">{rule.value || '—'}</span>
									</li>
								{/each}
							</ul>
						{:else}
							<p class="m-0 text-sm text-muted-foreground">No matching rules configured.</p>
						{/if}
					</div>
				</details>
			{/each}
		</div>

		{#if match.headersSummary.length > 0}
			<details class="rounded-md border border-border bg-muted/30 px-4 py-3">
				<summary class="cursor-pointer text-sm font-medium">Request signals used</summary>
				<ul class="mb-0 mt-2 grid gap-1 pl-5 text-xs text-muted-foreground">
					{#each match.headersSummary as header (header)}
						<li class="break-all">{header}</li>
					{/each}
				</ul>
			</details>
		{/if}
	{:else}
		<p class="m-0 rounded-md border border-border p-4 text-sm text-muted-foreground">
			Run trace to compare profile candidates.
		</p>
	{/if}
</section>
