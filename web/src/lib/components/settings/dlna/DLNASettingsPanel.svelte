<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import RotateCwIcon from '@lucide/svelte/icons/rotate-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type { DLNASettingsRequest, DLNAStatus } from '$lib/settings/types';
	import { createDLNAResources } from './dlnaResources.svelte';
	import DLNADiagnosticsTables from './DLNADiagnosticsTables.svelte';
	import DLNADeviceProfilesPanel from './profiles/DLNADeviceProfilesPanel.svelte';
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

	const resources = createDLNAResources();
	let form = $derived(settingsForm(resources.settings.data));
	let allowedText = $derived(form.allowedCidrs.join('\n'));
	let errorMessage = $state('');
	let message = $state('');

	const loading = $derived(resources.settings.isFetching);
	const saving = $derived(resources.updateSettings.isPending);
	const status = $derived<DLNAStatus | undefined>(resources.settings.data?.status);
	const statusCells = $derived([
		{ label: 'State', value: status?.running ? 'Running' : 'Stopped' },
		{ label: 'SSDP', value: status?.lastSsdpEvent ?? 'None' },
		{ label: 'Last SOAP', value: status?.lastSoapAction ?? 'None' },
		{ label: 'Last error', value: status?.lastError ?? 'None' }
	]);

	async function load() {
		errorMessage = '';
		try {
			await resources.settings.refetch();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load DLNA settings';
		}
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		errorMessage = '';
		message = '';
		try {
			await resources.updateSettings.mutateAsync(normalizedForm());
			message = 'DLNA settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save DLNA settings';
		}
	}

	async function restart() {
		errorMessage = '';
		message = '';
		try {
			await resources.restart.mutateAsync();
			message = 'DLNA restarted';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not restart DLNA';
		}
	}

	function normalizedForm(): DLNASettingsRequest {
		return {
			...form,
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

	function settingsForm(settings = resources.settings.data): DLNASettingsRequest {
		if (!settings) return { ...defaultForm };
		return {
			enabled: settings.enabled,
			friendlyName: settings.friendlyName,
			interfaces: [...settings.interfaces],
			allowedCidrs: [...settings.allowedCidrs],
			announceIntervalSeconds: settings.announceIntervalSeconds,
			transcodeEnabled: settings.transcodeEnabled,
			thumbnailsEnabled: settings.thumbnailsEnabled,
			subtitlesEnabled: settings.subtitlesEnabled,
			defaultRendererProfile: settings.defaultRendererProfile
		};
	}
</script>

<div class="grid gap-6">
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
			<DLNASettingsForm
				bind:form
				bind:allowedText
				availableInterfaces={status?.availableInterfaces ?? []}
				{saving}
				onSave={save}
			/>
			<DLNADiagnosticsTables {status} />
		</Card.Content>
	</Card.Root>
	<DLNADeviceProfilesPanel />
</div>
