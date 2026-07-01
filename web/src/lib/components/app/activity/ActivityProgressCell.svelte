<script lang="ts">
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Progress } from '$lib/components/ui/progress';
	import { cn } from '$lib/utils';
	import { activityDisplay, createdLabel } from './activityDisplay';
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
	}

	let { activity }: Props = $props();
	const display = $derived(activityDisplay(activity));
	const showProgress = $derived(
		['queued', 'grabbed', 'downloading', 'completed'].includes(activity.status)
	);
	const progressValue = $derived(display.progressValue ?? 42);
	const statusTone = $derived(
		activity.status === 'grabbed' || activity.status === 'completed'
			? 'success'
			: activity.status === 'queued' || activity.status === 'downloading'
				? 'pending'
				: activity.status === 'failed' || activity.status === 'cancelled'
					? 'error'
					: 'muted'
	);
</script>

<div class="grid min-w-45 gap-1.5">
	<StatusPill tone={statusTone}>{activity.status}</StatusPill>
	{#if showProgress}
		<Progress
			value={progressValue}
			class={cn(
				'h-2 w-full max-w-80 bg-muted [&_[data-slot=progress-indicator]]:bg-primary',
				display.progressValue === undefined && 'animate-pulse'
			)}
			aria-label="Download progress"
			aria-valuenow={display.progressValue}
		/>
		<small>{display.progressLabel}</small>
	{/if}
	<small>{createdLabel(activity.createdAt)}</small>
</div>
