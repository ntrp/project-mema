<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Textarea } from '$lib/components/ui/textarea';
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

<SettingsFormModal title="Import custom format" {onClose}>
	<form class="grid gap-4" onsubmit={importFormats}>
		<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
			<strong>Arr JSON</strong>
			<Button type="button" variant="outline" onclick={sampleJson}>Use sample</Button>
		</div>

		<label class="grid gap-2 text-sm">
			<span class="font-medium">Custom format JSON</span>
			<Textarea
				bind:value={rawJson}
				class="min-h-80 font-mono text-xs"
				rows={16}
				placeholder="Paste a Radarr or Sonarr custom format JSON export"
				oninput={parseImport}
				required
			/>
		</label>

		{#if error}
			<p
				class="m-0 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
			>
				{error}
			</p>
		{:else if parsed.length}
			<Card class="grid gap-1 p-3 text-sm">
				<strong>{parsed.length === 1 ? parsed[0].name : `${parsed.length} custom formats`}</strong>
				<span class="text-muted-foreground">
					{parsed.reduce((total, format) => total + format.includeSpecs.length, 0)} accepted /
					{parsed.reduce((total, format) => total + format.excludeSpecs.length, 0)} rejected conditions
				</span>
			</Card>
		{/if}

		<div class="flex justify-end gap-2">
			<Button type="button" variant="outline" onclick={onClose}>Cancel</Button>
			<Button type="submit" disabled={importing || rawJson.trim() === ''}>
				{importing ? 'Importing' : 'Import'}
			</Button>
		</div>
	</form>
</SettingsFormModal>
