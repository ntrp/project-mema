<script lang="ts">
	interface Props {
		clientCount: number;
		enabledClientCount: number;
		indexerCount: number;
		enabledIndexerCount: number;
	}

	let { clientCount, enabledClientCount, indexerCount, enabledIndexerCount }: Props = $props();
</script>

<header class="settings-header">
	<div>
		<p>Settings</p>
		<h1>Download clients and indexers</h1>
	</div>
	<div class="summary">
		<span>{clientCount} clients</span>
		<span>{enabledClientCount} enabled</span>
		<span>{indexerCount} indexers</span>
		<span>{enabledIndexerCount} enabled</span>
	</div>
</header>
