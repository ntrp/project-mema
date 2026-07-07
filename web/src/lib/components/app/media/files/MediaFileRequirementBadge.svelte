<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import {
		statusBadgeClass,
		statusBadgeVariant,
		type MediaFileSummaryStatus
	} from '$lib/components/app/media/files/mediaFileSummaryStatus';

	interface Props {
		status: MediaFileSummaryStatus;
	}

	let { status }: Props = $props();
	const detail = $derived(status.details.join(', '));
</script>

{#if status.state === 'ignored'}
	<span class="text-sm">{status.label}</span>
{:else}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Badge
					{...props}
					variant={statusBadgeVariant(status.state)}
					class={`justify-self-start ${statusBadgeClass(status.state) ?? ''}`}
					aria-label={detail}
				>
					{status.label}
				</Badge>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{detail}</Tooltip.Content>
	</Tooltip.Root>
{/if}
