<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import RotateCwIcon from '@lucide/svelte/icons/rotate-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type { DLNASettingsRequest } from '$lib/settings/types';
	import { createDLNASettingsForm, allowedCidrsText } from './dlnaSettingsFormState';
	import { createDLNAResources } from './dlnaResources.svelte';
	import DLNADeviceProfilesPanel from './profiles/DLNADeviceProfilesPanel.svelte';
	import DLNASettingsForm from './DLNASettingsForm.svelte';

	const initialForm = createDLNASettingsForm();
	let form = $state(initialForm);
	let allowedText = $state(allowedCidrsText(initialForm.allowedCidrs));
	let errorMessage = $state('');
	const resources = createDLNAResources();
	let message = $state('');

	const loading = $derived(resources.settings.isFetching);
	const saving = $derived(resources.updateSettings.isPending);
	const status = $derived(resources.settings.data?.status);

	$effect(() => {
		const settings = resources.settings.data;
		if (settings) syncSettings(settings);
	});

	async function load() {
		errorMessage = '';
		try {
			const result = await resources.settings.refetch();
			syncSettings(result.data);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load DLNA settings';
		}
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		errorMessage = '';
		message = '';
		try {
			const saved = await resources.updateSettings.mutateAsync(normalizedForm());
			syncSettings(saved);
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

	function syncSettings(settings?: DLNASettingsRequest) {
		const next = createDLNASettingsForm(settings);
		form = next;
		allowedText = allowedCidrsText(next.allowedCidrs);
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
</script>

<div class="grid gap-6">
	<Card.Root aria-label="DLNA server configuration">
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
			<DLNASettingsForm
				bind:form
				bind:allowedText
				availableInterfaces={status?.availableInterfaces ?? []}
				profiles={resources.profiles.data ?? []}
				{saving}
				onSave={save}
			/>
		</Card.Content>
	</Card.Root>
	<DLNADeviceProfilesPanel />
</div>
