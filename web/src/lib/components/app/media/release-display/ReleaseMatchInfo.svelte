<script lang="ts">
	import CircleAlertIcon from '@lucide/svelte/icons/circle-alert';
	import CircleXIcon from '@lucide/svelte/icons/circle-x';
	import InfoIcon from '@lucide/svelte/icons/info';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaType } from '$lib/settings/types';
	import type { MatchInfo } from '$lib/components/app/media/release-display/releaseCandidateDisplay';
	import { parsedTooltipSections } from '$lib/components/app/media/release-display/releaseMatchTooltip';

	interface Props {
		info: MatchInfo;
		mediaType: MediaType;
	}

	let { info, mediaType }: Props = $props();
	let open = $state(false);

	const label = $derived(
		info.severity === 'error'
			? 'Release mismatch'
			: info.severity === 'warning'
				? 'Release warning'
				: 'Release match'
	);
	const iconClass = $derived(
		info.severity === 'error'
			? 'text-destructive'
			: info.severity === 'warning'
				? 'text-amber-500'
				: 'text-sky-500'
	);
	const parsedSections = $derived(parsedTooltipSections(info, mediaType));
	const matchedFormats = $derived(
		[...new Set((info.customFormatContributors ?? []).map((contributor) => contributor.label))]
			.filter(Boolean)
			.sort((left, right) => left.localeCompare(right))
	);
</script>

<Tooltip.Root bind:open>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<button
				{...props}
				type="button"
				class="inline-flex h-8 w-8 items-center justify-center rounded-md hover:bg-accent"
				aria-label={label}
			>
				{#if info.severity === 'error'}
					<CircleXIcon class={iconClass} aria-hidden="true" />
				{:else if info.severity === 'warning'}
					<CircleAlertIcon class={iconClass} aria-hidden="true" />
				{:else}
					<InfoIcon class={iconClass} aria-hidden="true" />
				{/if}
			</button>
		{/snippet}
	</Tooltip.Trigger>
	{#if open}
		<Tooltip.Content class="max-h-[min(520px,calc(100vh-96px))] max-w-112 overflow-auto">
			<div class="grid gap-3 text-left">
				{#if info.details.length > 0}
					<div class="grid gap-1">
						<span class="font-bold">Decision</span>
						{#each info.details as detail (detail)}
							<span>{detail}</span>
						{/each}
					</div>
				{/if}
				{#each parsedSections as section (section.label)}
					<div class="grid gap-1">
						<span class="font-bold">{section.label}</span>
						{#each section.fields as field (`${section.label}:${field.label}`)}
							<div class="grid grid-cols-[104px_minmax(0,1fr)] gap-3">
								<span class="text-muted-foreground">{field.label}</span>
								<span class="break-anywhere font-mono text-xs">{field.value}</span>
							</div>
						{/each}
					</div>
				{/each}
				{#if matchedFormats.length > 0}
					<div class="grid gap-1.5">
						<span class="font-bold">Matched formats</span>
						<span class="flex flex-wrap gap-1">
							{#each matchedFormats as format (format)}
								<Badge
									variant="outline"
									class="border-sky-500/35 bg-sky-500/10 text-sky-700 dark:text-sky-300"
								>
									{format}
								</Badge>
							{/each}
						</span>
					</div>
				{/if}
			</div>
		</Tooltip.Content>
	{/if}
</Tooltip.Root>
