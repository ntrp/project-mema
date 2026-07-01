<script lang="ts">
	import CircleAlertIcon from '@lucide/svelte/icons/circle-alert';
	import InfoIcon from '@lucide/svelte/icons/info';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { SystemEventSeverity } from '$lib/settings/types';

	interface Props {
		severity: SystemEventSeverity;
	}

	let { severity }: Props = $props();

	const toneClasses: Record<SystemEventSeverity, string> = {
		error: 'text-destructive',
		warning: 'text-secondary-foreground',
		info: 'text-primary'
	};
	const Icon = $derived(
		severity === 'error' ? CircleAlertIcon : severity === 'warning' ? TriangleAlertIcon : InfoIcon
	);
	const toneClass = $derived(toneClasses[severity]);
	const label = $derived(severity[0].toUpperCase() + severity.slice(1));
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span
				class={`inline-flex size-6 items-center justify-center ${toneClass}`}
				aria-label={label}
				{...props}
			>
				<Icon aria-hidden="true" />
			</span>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content>{label}</Tooltip.Content>
</Tooltip.Root>
