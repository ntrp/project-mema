<script lang="ts">
	import AppShell from '$lib/components/app/shell/AppShell.svelte';

	interface Props {
		data: {
			id: string;
		};
	}

	let { data }: Props = $props();
</script>

<AppShell initialHomeSection="requests" initialSelectedRequestId={data.id} />
