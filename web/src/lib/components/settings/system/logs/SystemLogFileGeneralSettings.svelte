<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
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

<form class="grid gap-4 sm:grid-cols-2" onsubmit={save}>
	<div class="flex items-center gap-3 sm:col-span-2">
		<Switch id="write-log-files" bind:checked={enabled} />
		<Label for="write-log-files">Write logs to files</Label>
	</div>
	<div class="space-y-2">
		<Label for="log-directory">Log directory</Label>
		<Input id="log-directory" bind:value={directory} />
	</div>
	<div class="space-y-2">
		<Label for="log-retention-days">Log file retention days</Label>
		<Input id="log-retention-days" type="number" min="1" max="365" bind:value={retentionDays} />
	</div>
	<div class="flex justify-end sm:col-span-2">
		<Button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</Button>
	</div>
</form>

{#if settings}
	<p class="m-0 text-sm leading-6 text-muted-foreground">
		Effective log directory: {settings.effectiveDirectory}
	</p>
{/if}
