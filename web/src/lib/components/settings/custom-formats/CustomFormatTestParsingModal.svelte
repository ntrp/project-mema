<script lang="ts">
	import ListChecksIcon from '@lucide/svelte/icons/list-checks';
	import XIcon from '@lucide/svelte/icons/x';
	import CustomFormatParsingResults from '$lib/components/settings/custom-formats/CustomFormatParsingResults.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
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

<SettingsFormModal
	title="Test parsing"
	modalClass="w-[min(1180px,calc(100vw-48px))] max-h-[min(860px,calc(100vh-48px))] max-sm:w-full max-sm:max-h-[calc(100vh-24px)]"
	{onClose}
>
	<div class="grid gap-3">
		<div
			class="grid grid-cols-[auto_minmax(0,1fr)_auto_auto] items-end gap-2.5 rounded-md border border-border bg-card p-2.5 max-sm:grid-cols-1"
		>
			<ListChecksIcon aria-hidden="true" />
			<label class="grid gap-1">
				<span>Release title</span>
				<Input bind:value={fileName} type="text" maxlength={500} />
			</label>
			<Button
				type="button"
				variant="outline"
				size="icon"
				aria-label="Clear file name"
				onclick={clearFileName}
			>
				<XIcon aria-hidden="true" />
			</Button>
		</div>
	</div>

	{#if error}
		<p
			class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 font-bold text-destructive"
		>
			{error}
		</p>
	{/if}

	{#if loading}
		<EmptyState class="grid min-h-[118px] place-items-center text-center" aria-live="polite">
			<InlineSpinner label="Parsing" />
		</EmptyState>
	{:else if !fileName.trim()}
		<EmptyState
			class="text-center"
			title="Enter a release title in the input above"
			description="mema will attempt to parse the title and show you details about it"
		/>
	{:else if result}
		<CustomFormatParsingResults {result} />
	{/if}
</SettingsFormModal>
