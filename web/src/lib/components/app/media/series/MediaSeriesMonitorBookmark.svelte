<script lang="ts">
	import MediaMonitorBookmark from '$lib/components/app/media/detail/MediaMonitorBookmark.svelte';

	interface Props {
		name: string;
		monitored?: boolean;
		target: 'season' | 'episode';
		disabled?: boolean;
		onToggle: () => void;
	}

	let { name, monitored = false, target, disabled = false, onToggle }: Props = $props();

	const status = $derived(`${name} ${monitored ? 'monitored' : 'not monitored'}`);
	const hint = $derived(monitorHint());

	function monitorHint() {
		if (target === 'season') {
			return monitored
				? 'Click to stop monitoring all episodes in this season'
				: 'Click to monitor all episodes in this season';
		}
		return monitored ? 'Click to stop monitoring this episode' : 'Click to monitor this episode';
	}
</script>

<MediaMonitorBookmark {monitored} {status} {hint} {disabled} {onToggle} />
