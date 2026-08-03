<script lang="ts">
	import { Progress } from '$lib/components/ui/progress';
	import { activityDisplay } from './activityDisplay';
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
	}

	let { activity }: Props = $props();
	const display = $derived(activityDisplay(activity));
	const showProgress = $derived(
		['queued', 'grabbed', 'downloading', 'completed'].includes(activity.status)
	);
	const progressValue = $derived(display.progressValue);
</script>

{#if showProgress}
	<div class="relative min-w-0">
		<Progress
			value={progressValue}
			class="h-6 w-full rounded-sm bg-muted [&_[data-slot=progress-indicator]]:bg-primary"
			aria-label="Download progress"
			aria-valuenow={display.progressValue}
			aria-valuetext={display.progressLabel}
		/>
		<span
			class="pointer-events-none absolute inset-0 grid place-items-center text-xs font-extrabold"
		>
			{display.progressLabel}
		</span>
	</div>
{/if}
