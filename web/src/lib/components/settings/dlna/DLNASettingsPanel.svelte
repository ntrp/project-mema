<script lang="ts">
	import { onMount } from 'svelte';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import RotateCwIcon from '@lucide/svelte/icons/rotate-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import { getDLNASettings, restartDLNA, updateDLNASettings } from '$lib/settings/api';
	import type { DLNASettings, DLNASettingsRequest, DLNAStatus } from '$lib/settings/types';
	import DLNADiagnosticsTables from './DLNADiagnosticsTables.svelte';
	import DLNASettingsForm from './DLNASettingsForm.svelte';

	const defaultForm: DLNASettingsRequest = {
		enabled: false,
		friendlyName: 'Mema',
		interfaces: [],
		allowedCidrs: ['127.0.0.1/32', '::1/128'],
		announceIntervalSeconds: 1800,
		transcodeEnabled: true,
		thumbnailsEnabled: true,
		subtitlesEnabled: true,
		defaultRendererProfile: 'generic'
	};

	let settings = $state<DLNASettings>();
	let form = $state<DLNASettingsRequest>({ ...defaultForm });
	let interfacesText = $state('');
	let allowedText = $state(defaultForm.allowedCidrs.join('\n'));
	let loading = $state(true);
	let saving = $state(false);
	let errorMessage = $state('');
	let message = $state('');

	const status = $derived<DLNAStatus | undefined>(settings?.status);
	const statusCells = $derived([
		{ label: 'State', value: status?.running ? 'Running' : 'Stopped' },
		{ label: 'SSDP', value: status?.lastSsdpEvent ?? 'None' },
		{ label: 'Last SOAP', value: status?.lastSoapAction ?? 'None' },
		{ label: 'Last error', value: status?.lastError ?? 'None' }
	]);

	onMount(() => {
		void load();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			hydrate(await getDLNASettings());
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load DLNA settings';
		} finally {
			loading = false;
		}
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		errorMessage = '';
		message = '';
		try {
			hydrate(await updateDLNASettings(normalizedForm()));
			message = 'DLNA settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save DLNA settings';
		} finally {
			saving = false;
		}
	}

	async function restart() {
		errorMessage = '';
		message = '';
		try {
			hydrate(await restartDLNA());
			message = 'DLNA restarted';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not restart DLNA';
		}
	}

	function hydrate(next: DLNASettings) {
		settings = next;
		form = {
			enabled: next.enabled,
			friendlyName: next.friendlyName,
			interfaces: [...next.interfaces],
			allowedCidrs: [...next.allowedCidrs],
			announceIntervalSeconds: next.announceIntervalSeconds,
			transcodeEnabled: next.transcodeEnabled,
			thumbnailsEnabled: next.thumbnailsEnabled,
			subtitlesEnabled: next.subtitlesEnabled,
			defaultRendererProfile: next.defaultRendererProfile
		};
		interfacesText = form.interfaces.join('\n');
		allowedText = form.allowedCidrs.join('\n');
	}

	function normalizedForm(): DLNASettingsRequest {
		return {
			...form,
			interfaces: lines(interfacesText),
			allowedCidrs: lines(allowedText),
			announceIntervalSeconds: Number(form.announceIntervalSeconds)
		};
	}

	function lines(value: string) {
		return value
			.split('\n')
			.map((line) => line.trim())
			.filter(Boolean);
	}
</script>

<Card.Root aria-label="DLNA diagnostics">
	<Card.Header class="border-b border-border">
		<Card.Title>Server</Card.Title>
		<Card.Action class="flex gap-2">
			<Button type="button" variant="secondary" size="sm" disabled={loading} onclick={load}>
				<RefreshCwIcon class={loading ? 'animate-spin' : ''} />
				Refresh
			</Button>
			<ConfirmActionButton
				label="Restart DLNA"
				title="Restart DLNA"
				description="Active DLNA discovery and subscriptions will be restarted."
				confirmLabel="Restart"
				confirmingLabel="Restarting"
				variant="outline"
				size="sm"
				disabled={loading || saving}
				tooltip="Restart DLNA"
				onConfirm={restart}
			>
				<RotateCwIcon />
				Restart
			</ConfirmActionButton>
		</Card.Action>
	</Card.Header>
	<Card.Content class="grid gap-5 pt-5">
		{#if errorMessage}
			<p class="text-sm font-medium text-destructive">{errorMessage}</p>
		{/if}
		{#if message}
			<p class="text-sm text-muted-foreground">{message}</p>
		{/if}
		<div class="grid gap-3 sm:grid-cols-4">
			{#each statusCells as cell (cell.label)}
				<div class="grid gap-1 rounded-md border border-border p-3">
					<span class="text-xs font-medium text-muted-foreground uppercase">{cell.label}</span>
					<span class="break-words text-sm font-medium text-foreground">{cell.value}</span>
				</div>
			{/each}
		</div>
		<DLNASettingsForm bind:form bind:interfacesText bind:allowedText {saving} onSave={save} />
		<DLNADiagnosticsTables {status} />
	</Card.Content>
</Card.Root>
