<script lang="ts">
	import { getSystemLogFileSettings, updateSystemLogFileSettings } from '$lib/settings/api';
	import type { SystemLogFileSettings } from '$lib/settings/types';

	let settings = $state<SystemLogFileSettings>();
	let enabled = $state(false);
	let directory = $state('.data/logs');
	let retentionDays = $state(7);
	let saving = $state(false);

	export async function load() {
		applySettings(await getSystemLogFileSettings());
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		try {
			applySettings(await updateSystemLogFileSettings({ enabled, directory, retentionDays }));
		} finally {
			saving = false;
		}
	}

	function applySettings(nextSettings: SystemLogFileSettings) {
		settings = nextSettings;
		enabled = nextSettings.enabled;
		directory = nextSettings.directory;
		retentionDays = nextSettings.retentionDays;
	}
</script>

<form class="settings-form compact-form log-file-form" onsubmit={save}>
	<label class="inline-check">
		<input type="checkbox" bind:checked={enabled} />
		<span>Write logs to files</span>
	</label>
	<label>
		<span>Log directory</span>
		<input bind:value={directory} />
	</label>
	<label>
		<span>Log file retention days</span>
		<input type="number" min="1" max="365" bind:value={retentionDays} />
	</label>
	<div class="form-actions">
		<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</button>
	</div>
</form>

{#if settings}
	<p class="muted">Effective log directory: {settings.effectiveDirectory}</p>
{/if}
