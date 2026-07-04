<script lang="ts">
	interface Props {
		message: string;
	}

	let { message }: Props = $props();
</script>

<p
	class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
>
	{message}
</p>
