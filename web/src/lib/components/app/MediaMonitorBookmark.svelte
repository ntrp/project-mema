<script lang="ts">
	interface Props {
		monitored?: boolean;
		label?: string;
	}

	let { monitored = false, label }: Props = $props();
	const title = $derived(label ?? (monitored ? 'Monitored' : 'Not monitored'));
</script>

<span class:monitored class="monitor-bookmark" {title} aria-label={title}>
	<span class="app-icon" aria-hidden="true">{monitored ? 'bookmark' : 'bookmark_border'}</span>
</span>
