<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { parseArrCustomFormatImport } from '$lib/settings/arrCustomFormatImport';
	import type { CustomFormatForm } from '$lib/settings/types';

	interface Props {
		onClose: () => void;
		onImport: (_form: CustomFormatForm) => void | Promise<void>;
	}

	let { onClose, onImport }: Props = $props();
	let rawJson = $state('');
	let parsed = $state<CustomFormatForm[]>([]);
	let error = $state('');
	let importing = $state(false);

	function parseImport() {
		error = '';
		try {
			parsed = parseArrCustomFormatImport(rawJson);
		} catch (caught) {
			parsed = [];
			error = caught instanceof Error ? caught.message : 'Could not parse custom format JSON';
		}
	}

	async function importFormats(event: SubmitEvent) {
		event.preventDefault();
		parseImport();
		const formats = parseArrCustomFormatImport(rawJson);
		importing = true;
		error = '';
		try {
			for (const format of formats) {
				await onImport(format);
			}
			onClose();
		} catch (caught) {
			error = caught instanceof Error ? caught.message : 'Could not import custom format';
		} finally {
			importing = false;
		}
	}

	function sampleJson() {
		rawJson = JSON.stringify(
			{
				name: 'DSNP',
				specifications: [
					{
						name: 'Disney+',
						implementation: 'ReleaseTitleSpecification',
						negate: false,
						required: true,
						fields: [{ name: 'value', value: '\\\\b(dsnp|dsny|disney|Disney\\\\+)\\\\b' }]
					},
					{
						name: 'WEBDL',
						implementation: 'SourceSpecification',
						negate: false,
						required: false,
						fields: [{ name: 'value', value: 'WEBDL' }]
					}
				]
			},
			null,
			2
		);
		parseImport();
	}
</script>

<SettingsFormModal title="Import custom format" modalClass="custom-format-import-modal" {onClose}>
	<form class="settings-form custom-format-import-form" onsubmit={importFormats}>
		<div class="profile-quality-header">
			<strong>Arr JSON</strong>
			<button type="button" class="secondary" onclick={sampleJson}>Use sample</button>
		</div>

		<label>
			<span>Custom format JSON</span>
			<textarea
				bind:value={rawJson}
				rows="16"
				placeholder="Paste a Radarr or Sonarr custom format JSON export"
				oninput={parseImport}
				required></textarea>
		</label>

		{#if error}
			<p class="form-status error">{error}</p>
		{:else if parsed.length}
			<div class="custom-format-import-preview">
				<strong>{parsed.length === 1 ? parsed[0].name : `${parsed.length} custom formats`}</strong>
				<span>
					{parsed.reduce((total, format) => total + format.includeSpecs.length, 0)} accepted /
					{parsed.reduce((total, format) => total + format.excludeSpecs.length, 0)} rejected conditions
				</span>
			</div>
		{/if}

		<div class="form-actions">
			<button type="button" class="secondary" onclick={onClose}>Cancel</button>
			<button type="submit" disabled={importing || rawJson.trim() === ''}>
				{importing ? 'Importing' : 'Import'}
			</button>
		</div>
	</form>
</SettingsFormModal>
