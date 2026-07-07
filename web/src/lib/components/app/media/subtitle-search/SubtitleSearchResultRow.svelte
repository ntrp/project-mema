<script lang="ts">
	import AlertTriangleIcon from '@lucide/svelte/icons/alert-triangle';
	import CheckCircleIcon from '@lucide/svelte/icons/check-circle-2';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import XCircleIcon from '@lucide/svelte/icons/x-circle';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { displayLanguage } from '$lib/settings/languageDisplay';
	import type { SubtitleCandidate } from '$lib/settings/types';

	interface Props {
		candidate: SubtitleCandidate;
		grabbing?: boolean;
		canManage: boolean;
		onGrab: (_candidate: SubtitleCandidate) => void;
	}

	let { candidate, grabbing = false, canManage, onGrab }: Props = $props();
	const matchTone = $derived(
		candidate.match.severity === 'success'
			? 'text-emerald-600'
			: candidate.match.severity === 'warning'
				? 'text-amber-600'
				: 'text-destructive'
	);
</script>

<Table.Row>
	<Table.Cell class="whitespace-nowrap">
		<Badge variant="outline">{candidate.protocol}</Badge>
	</Table.Cell>
	<Table.Cell class="max-w-48 truncate whitespace-nowrap">{candidate.providerName}</Table.Cell>
	<Table.Cell class="min-w-0 max-w-0">
		<span class="block truncate font-medium">{candidate.title}</span>
		<span class="block truncate text-xs text-muted-foreground"
			>{candidate.format.toUpperCase()}</span
		>
	</Table.Cell>
	<Table.Cell class="whitespace-nowrap">{displayLanguage(candidate.languageId)}</Table.Cell>
	<Table.Cell class="whitespace-nowrap">
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<span {...props} class={`inline-flex items-center gap-1 ${matchTone}`}>
						{#if candidate.match.severity === 'success'}
							<CheckCircleIcon class="size-4" />
						{:else if candidate.match.severity === 'warning'}
							<AlertTriangleIcon class="size-4" />
						{:else}
							<XCircleIcon class="size-4" />
						{/if}
						<span class="sr-only">{candidate.match.label}</span>
					</span>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content class="max-w-80">
				<div class="grid gap-1">
					<strong>{candidate.match.label}</strong>
					{#each candidate.match.details as detail (detail)}
						<span>{detail}</span>
					{/each}
				</div>
			</Tooltip.Content>
		</Tooltip.Root>
	</Table.Cell>
	<Table.Cell class="text-right">
		<Button
			type="button"
			size="icon-sm"
			class="bg-emerald-600 text-white hover:bg-emerald-700"
			disabled={!canManage || grabbing}
			aria-label="Grab subtitle"
			onclick={() => onGrab(candidate)}
		>
			<DownloadIcon aria-hidden="true" />
		</Button>
	</Table.Cell>
</Table.Row>
