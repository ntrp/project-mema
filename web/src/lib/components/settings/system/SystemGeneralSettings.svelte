<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { getSystemEventSettings, updateSystemEventSettings } from './events/api';
	import SystemLogFileGeneralSettings from './logs/SystemLogFileGeneralSettings.svelte';

	let retentionDays = $state(7);
	let loading = $state(true);
	let saving = $state(false);
	let logFileSettings = $state<SystemLogFileGeneralSettings>();
	let errorMessage = $state('');
	let message = $state('');

	onMount(() => {
		void load();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			const eventSettings = await getSystemEventSettings();
			retentionDays = eventSettings.retentionDays;
			await logFileSettings?.load();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load general settings';
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
			retentionDays = (await updateSystemEventSettings({ retentionDays })).retentionDays;
			message = 'General settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save general settings';
		} finally {
			saving = false;
		}
	}
</script>

<PageHeading eyebrow="Settings" title="General" titleId="settings-title" />

<Card class="gap-4 p-5" aria-label="General">
	<SectionHeading>
		{#snippet actions()}
			<Button type="button" variant="outline" size="sm" disabled={loading} onclick={load}>
				Refresh
			</Button>
		{/snippet}
	</SectionHeading>

	{#if errorMessage}
		<p class="m-0 font-bold text-destructive">{errorMessage}</p>
	{/if}
	{#if message}
		<p class="m-0 text-sm leading-6 text-muted-foreground">{message}</p>
	{/if}

	<form class="grid gap-4 sm:grid-cols-2" onsubmit={save}>
		<div class="space-y-2">
			<Label for="event-retention-days">Event retention days</Label>
			<Input id="event-retention-days" type="number" min="1" max="365" bind:value={retentionDays} />
		</div>
		<div class="flex items-end justify-end">
			<Button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</Button>
		</div>
	</form>

	<SystemLogFileGeneralSettings bind:this={logFileSettings} />
</Card>
