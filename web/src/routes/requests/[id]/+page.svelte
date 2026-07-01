<script lang="ts">
	import HomeRoute from '$lib/features/home/HomeRoute.svelte';

	interface Props {
		data: {
			id: string;
		};
	}

	let { data }: Props = $props();
</script>

<HomeRoute section="requests" selectedRequestId={data.id} />
