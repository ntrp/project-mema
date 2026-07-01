<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import CustomFormatParsingResults from '$lib/components/settings/CustomFormatParsingResults.svelte';
	import { testCustomFormatParsing } from '$lib/settings/api';
	import type { CustomFormatParsingResponse } from '$lib/settings/types';

	interface Props {
		onClose: () => void;
	}

	let { onClose }: Props = $props();
	let fileName = $state('');
	let result = $state<CustomFormatParsingResponse | undefined>();
	let loading = $state(false);
	let error = $state('');
	let requestID = 0;

	$effect(() => {
		const value = fileName.trim();
		requestID += 1;
		const currentRequestID = requestID;
		if (!value) {
			result = undefined;
			error = '';
			loading = false;
			return;
		}
		loading = true;
		const timeout = window.setTimeout(() => {
			void runParsing(value, currentRequestID);
		}, 700);
		return () => window.clearTimeout(timeout);
	});

	async function runParsing(value: string, currentRequestID: number) {
		loading = true;
		error = '';
		try {
			const nextResult = await testCustomFormatParsing(value);
			if (currentRequestID === requestID) {
				result = nextResult;
			}
		} catch (caught) {
			if (currentRequestID === requestID) {
				error = caught instanceof Error ? caught.message : 'Could not test parsing';
				result = undefined;
			}
		} finally {
			if (currentRequestID === requestID) {
				loading = false;
			}
		}
	}

	function clearFileName() {
		fileName = '';
		result = undefined;
		error = '';
	}
</script>

<SettingsFormModal title="Test parsing" modalClass="test-parsing-modal" {onClose}>
	<div class="test-parsing-form">
		<div class="test-parsing-input">
			<span class="app-icon" aria-hidden="true">rule</span>
			<label>
				<span>Release title</span>
				<input bind:value={fileName} type="text" maxlength="500" />
			</label>
			<button
				type="button"
				class="secondary icon-button"
				aria-label="Clear file name"
				onclick={clearFileName}
			>
				<span class="app-icon" aria-hidden="true">close</span>
			</button>
		</div>
	</div>

	{#if error}
		<p class="form-status error">{error}</p>
	{/if}

	{#if loading}
		<section class="empty-state test-parsing-empty test-parsing-loading" aria-live="polite">
			<span class="test-parsing-spinner" aria-label="Parsing"></span>
		</section>
	{:else if !fileName.trim()}
		<section class="empty-state test-parsing-empty">
			<h3>Enter a release title in the input above</h3>
			<p>mema will attempt to parse the title and show you details about it</p>
		</section>
	{:else if result}
		<CustomFormatParsingResults {result} />
	{/if}
</SettingsFormModal>
