<script lang="ts">
	import type { SystemEventSeverity } from '$lib/settings/types';

	interface Props {
		severity: SystemEventSeverity;
	}

	let { severity }: Props = $props();

	const icons: Record<SystemEventSeverity, string> = {
		error: 'error',
		warning: 'warning',
		info: 'info'
	};
	const icon = $derived(icons[severity]);
	const label = $derived(severity[0].toUpperCase() + severity.slice(1));
</script>

<span class={`severity-icon severity-${severity}`} aria-label={label} title={label}>
	<span class="app-icon" aria-hidden="true">{icon}</span>
</span>
