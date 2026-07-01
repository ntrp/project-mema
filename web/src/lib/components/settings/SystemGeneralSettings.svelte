<script lang="ts">
	import { onMount } from 'svelte';
	import { getSystemEventSettings, updateSystemEventSettings } from '$lib/settings/api';
	import SystemLogFileGeneralSettings from './SystemLogFileGeneralSettings.svelte';

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

<div class="page-heading">
	<p>Settings</p>
	<h1 id="settings-title">General</h1>
</div>

<section class="panel log-settings-panel" aria-label="General">
	<div class="section-heading">
		<button type="button" class="secondary compact-action" disabled={loading} onclick={load}>
			Refresh
		</button>
	</div>

	{#if errorMessage}
		<p class="inline-error">{errorMessage}</p>
	{/if}
	{#if message}
		<p class="muted">{message}</p>
	{/if}

	<form class="settings-form compact-form event-settings-form" onsubmit={save}>
		<label>
			<span>Event retention days</span>
			<input type="number" min="1" max="365" bind:value={retentionDays} />
		</label>
		<div class="form-actions">
			<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</button>
		</div>
	</form>

	<SystemLogFileGeneralSettings bind:this={logFileSettings} />
</section>
