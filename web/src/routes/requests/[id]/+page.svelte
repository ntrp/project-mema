<script lang="ts">
	import RequestsRoute from '$lib/features/requests/RequestsRoute.svelte';

	interface Props {
		data: {
			id: string;
		};
	}

	let { data }: Props = $props();
</script>

<RequestsRoute selectedRequestId={data.id} />
