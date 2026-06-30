<script lang="ts">
	interface Props {
		attributes: Record<string, unknown>;
	}

	let { attributes }: Props = $props();

	const formatted = $derived(JSON.stringify(attributes, null, 2));
</script>

<button type="button" class="log-json-trigger" aria-label="Log attributes">
	<span class="app-icon" aria-hidden="true">data_object</span>
	<code>{formatted}</code>
</button>
