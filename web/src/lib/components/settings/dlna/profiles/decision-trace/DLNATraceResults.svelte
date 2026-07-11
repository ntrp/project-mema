<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import {
		filterTraceSteps,
		type DLNATraceSummaryItem,
		type DLNATraceStep
	} from './dlnaDecisionTrace';

	interface Props {
		steps: DLNATraceStep[];
		summary: DLNATraceSummaryItem[];
		hideFailedSteps?: boolean;
	}

	let { steps, summary, hideFailedSteps = true }: Props = $props();
	const visibleSteps = $derived(filterTraceSteps(steps, hideFailedSteps));
</script>

<section class="grid gap-3" aria-label="DLNA decision trace results">
	<div
		class="h-[clamp(16rem,calc(100vh-22rem),24rem)] overflow-y-auto rounded-md border border-border bg-card"
	>
		{#if visibleSteps.length > 0}
			<ul class="grid divide-y divide-border">
				{#each visibleSteps as step, index (step.id)}
					<li class="grid gap-2 px-4 py-3">
						<div class="flex flex-wrap items-center justify-between gap-2">
							<div class="flex flex-wrap items-center gap-2">
								<Badge variant="secondary">{step.stage}</Badge>
								<span class="text-sm font-semibold">Step {index + 1}: {step.field}</span>
							</div>
							<div class="flex items-center gap-2">
								{#if step.score !== undefined}
									<Badge variant="outline">Score {step.score}</Badge>
								{/if}
								<Badge variant={step.result === 'pass' ? 'default' : 'destructive'}
									>{step.result}</Badge
								>
							</div>
						</div>
						<p class="m-0 text-sm text-muted-foreground">Rule: {step.rule}</p>
						<p class="m-0 text-sm text-muted-foreground">Value: {step.value}</p>
					</li>
				{/each}
			</ul>
		{:else if steps.length > 0}
			<p class="m-0 p-4 text-sm text-muted-foreground">No passing decision steps to show.</p>
		{:else}
			<p class="m-0 p-4 text-sm text-muted-foreground">
				Select a device and run trace to see the decision steps.
			</p>
		{/if}
	</div>
	<div class="grid gap-3 rounded-md border border-border bg-muted/30 p-4">
		<div class="grid gap-3 sm:grid-cols-2">
			{#each summary as item (item.label)}
				<div class="grid gap-1">
					<span class="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
						>{item.label}</span
					>
					<span class="text-sm font-medium">{item.value}</span>
				</div>
			{/each}
		</div>
	</div>
</section>
