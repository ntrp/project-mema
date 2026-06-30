<script lang="ts">
	import { onMount } from 'svelte';
	import {
		downloadSystemLogFile,
		getSystemLogFileSettings,
		listSystemLogFiles,
		updateSystemLogFileSettings
	} from '$lib/settings/api';
	import type { SystemLogFile, SystemLogFileSettings } from '$lib/settings/types';

	let settings = $state<SystemLogFileSettings>();
	let files = $state<SystemLogFile[]>([]);
	let enabled = $state(false);
	let directory = $state('.data/logs');
	let retentionDays = $state(7);
	let loading = $state(true);
	let saving = $state(false);
	let downloadingName = $state<string | undefined>();
	let errorMessage = $state('');
	let message = $state('');

	onMount(() => {
		void load();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			const [nextSettings, nextFiles] = await Promise.all([
				getSystemLogFileSettings(),
				listSystemLogFiles()
			]);
			applySettings(nextSettings);
			files = nextFiles;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load log files';
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
			applySettings(
				await updateSystemLogFileSettings({
					enabled,
					directory,
					retentionDays
				})
			);
			files = await listSystemLogFiles();
			message = 'Log file settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save log file settings';
		} finally {
			saving = false;
		}
	}

	async function download(file: SystemLogFile) {
		downloadingName = file.name;
		errorMessage = '';
		try {
			const blob = await downloadSystemLogFile(file.name);
			const href = globalThis.URL.createObjectURL(blob);
			const link = globalThis.document.createElement('a');
			link.href = href;
			link.download = file.name;
			link.click();
			globalThis.URL.revokeObjectURL(href);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not download log file';
		} finally {
			downloadingName = undefined;
		}
	}

	function applySettings(nextSettings: SystemLogFileSettings) {
		settings = nextSettings;
		enabled = nextSettings.enabled;
		directory = nextSettings.directory;
		retentionDays = nextSettings.retentionDays;
	}

	function formatTime(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			year: 'numeric',
			month: 'short',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit'
		}).format(new Date(value));
	}

	function formatSize(bytes: number) {
		if (bytes < 1024) {
			return `${bytes} B`;
		}
		if (bytes < 1024 * 1024) {
			return `${(bytes / 1024).toFixed(1)} KB`;
		}
		return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
	}
</script>

<section class="panel log-settings-panel" aria-label="Log files">
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

	<form class="settings-form compact-form log-file-form" onsubmit={save}>
		<label class="inline-check">
			<input type="checkbox" bind:checked={enabled} />
			<span>Write logs to files</span>
		</label>
		<label>
			<span>Directory</span>
			<input bind:value={directory} />
		</label>
		<label>
			<span>Retention days</span>
			<input type="number" min="1" max="365" bind:value={retentionDays} />
		</label>
		<div class="form-actions">
			<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</button>
		</div>
	</form>

	{#if settings}
		<p class="muted">Effective directory: {settings.effectiveDirectory}</p>
	{/if}

	<div class="table-wrap">
		<table class="data-table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Size</th>
					<th>Modified</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each files as file (file.name)}
					<tr>
						<td>{file.name}</td>
						<td>{formatSize(file.sizeBytes)}</td>
						<td>{formatTime(file.modifiedAt)}</td>
						<td>
							<button
								type="button"
								class="secondary compact-action"
								disabled={downloadingName === file.name}
								onclick={() => download(file)}
							>
								{downloadingName === file.name ? 'Downloading' : 'Download'}
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4">No log files retained.</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
