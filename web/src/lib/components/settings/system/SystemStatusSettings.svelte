<script lang="ts">
	import { onMount } from 'svelte';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
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

<Card.Root aria-label="About">
	<Card.Header class="border-b border-border">
		<Card.Title>About</Card.Title>
		<Card.Action>
			<Button type="button" variant="secondary" size="sm" disabled={loading} onclick={load}>
				<RefreshCwIcon class={loading ? 'animate-spin' : ''} />
				Refresh
			</Button>
		</Card.Action>
	</Card.Header>

	<Card.Content class="grid gap-4 pt-5">
		{#if errorMessage}
			<p class="text-sm font-medium text-destructive">{errorMessage}</p>
		{/if}

		{#if loading && !status}
			<p class="text-sm text-muted-foreground">Loading system status...</p>
		{:else if status}
			<dl class="m-0 grid grid-cols-1 gap-4 min-[761px]:grid-cols-2">
				{#each details as detail (detail.label)}
					<div class="grid min-w-0 gap-1">
						<dt class="text-xs font-medium tracking-wide text-muted-foreground uppercase">
							{detail.label}
						</dt>
						<dd class="m-0 break-words text-sm font-medium text-foreground">
							{#if detail.link}
								<!-- eslint-disable svelte/no-navigation-without-resolve -->
								<a
									class="text-primary underline-offset-4 hover:underline"
									href={detail.link}
									target="_blank"
									rel="noreferrer"
								>
									{detail.value}
								</a>
							{:else}
								{detail.value}
							{/if}
						</dd>
					</div>
				{/each}
			</dl>
		{/if}
	</Card.Content>
</Card.Root>
