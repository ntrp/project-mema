<script lang="ts">
	import { onMount } from 'svelte';
	import { getSystemStatus } from '$lib/settings/api';
	import type { SystemStatusResponse } from '$lib/settings/types';

	let status = $state<SystemStatusResponse>();
	let loading = $state(true);
	let errorMessage = $state('');

	const details = $derived(
		status
			? [
					{ label: 'Version', value: versionLabel(status) },
					{ label: 'Database', value: `${status.databaseType} ${status.databaseVersion}` },
					{ label: 'Licence', value: status.license },
					{ label: 'Source location', value: status.sourceLocation, link: sourceLink(status) }
				]
			: []
	);

	onMount(() => {
		void load();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			status = await getSystemStatus();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load system status';
		} finally {
			loading = false;
		}
	}

	function sourceLink(nextStatus: SystemStatusResponse) {
		if (/^https?:\/\//.test(nextStatus.sourceLocation)) {
			return nextStatus.sourceLocation;
		}
		return undefined;
	}

	function versionLabel(nextStatus: SystemStatusResponse) {
		if (!nextStatus.commit || nextStatus.commit === 'dev') {
			return nextStatus.version;
		}
		return `${nextStatus.version} (${nextStatus.commit})`;
	}
</script>

<section class="panel status-panel" aria-label="About">
	<div class="section-heading">
		<h2>About</h2>
		<button type="button" class="secondary compact-action" disabled={loading} onclick={load}>
			Refresh
		</button>
	</div>

	{#if errorMessage}
		<p class="inline-error">{errorMessage}</p>
	{/if}

	{#if loading && !status}
		<p class="muted">Loading system status...</p>
	{:else if status}
		<dl class="status-details">
			{#each details as detail (detail.label)}
				<div>
					<dt>{detail.label}</dt>
					<dd>
						{#if detail.link}
							<!-- eslint-disable svelte/no-navigation-without-resolve -->
							<a href={detail.link} target="_blank" rel="noreferrer">{detail.value}</a>
						{:else}
							{detail.value}
						{/if}
					</dd>
				</div>
			{/each}
		</dl>
	{/if}
</section>

<style>
	.status-panel {
		display: grid;
		gap: 16px;
	}

	.status-details {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 14px;
		margin: 0;
	}

	.status-details div {
		display: grid;
		gap: 5px;
		min-width: 0;
	}

	.status-details dt {
		color: #9aa7ba;
		font-size: 12px;
		font-weight: 800;
		text-transform: uppercase;
	}

	.status-details dd {
		margin: 0;
		overflow-wrap: anywhere;
		color: #e6edf7;
		font-weight: 800;
	}

	.status-details a {
		color: #67e8f9;
	}

	@media (width <= 760px) {
		.status-details {
			grid-template-columns: 1fr;
		}
	}
</style>
