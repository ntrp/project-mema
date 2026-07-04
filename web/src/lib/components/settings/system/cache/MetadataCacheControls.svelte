<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	interface Props {
		pattern: string;
		placeholder?: string;
		clearing: boolean;
		showClearAll?: boolean;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
	}

	let {
		pattern = $bindable(),
		placeholder = 'Enter the pattern',
		clearing,
		showClearAll = true,
		onClearAll,
		onClearPattern
	}: Props = $props();

	let patternConfirmOpen = $state(false);
	let pendingPattern = $state('');
	let clearingPattern = $state(false);

	function requestClearPattern(event: SubmitEvent) {
		event.preventDefault();
		const nextPattern = pattern.trim();
		if (!nextPattern) return;
		pendingPattern = nextPattern;
		patternConfirmOpen = true;
	}

	function closePatternConfirmation() {
		patternConfirmOpen = false;
		pendingPattern = '';
	}

	async function confirmClearPattern() {
		if (!pendingPattern) return;
		clearingPattern = true;
		try {
			await onClearPattern(pendingPattern);
			pattern = '';
			closePatternConfirmation();
		} finally {
			clearingPattern = false;
		}
	}
</script>

<form class="grid items-end gap-3 md:grid-cols-[minmax(0,1fr)_auto]" onsubmit={requestClearPattern}>
	<div class="grid gap-1.5">
		<Label>Evict by regex</Label>
		<Input bind:value={pattern} {placeholder} autocomplete="off" />
	</div>
	<div class="flex flex-wrap justify-end gap-2">
		<Button
			type="submit"
			variant="destructive"
			disabled={clearing || clearingPattern || pattern.trim().length === 0}
		>
			{clearing || clearingPattern ? 'Evicting' : 'Evict'}
		</Button>
		{#if showClearAll}
			<ConfirmActionButton
				label="Reset all cache entries"
				title="Reset cache"
				description="Delete every cache entry?"
				confirmLabel="Reset all"
				confirmingLabel="Resetting"
				disabled={clearing}
				onConfirm={onClearAll}
			>
				Reset all
			</ConfirmActionButton>
		{/if}
	</div>
</form>

<Dialog.Root bind:open={patternConfirmOpen}>
	<Dialog.Content class="w-[min(460px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>Evict cache entries</Dialog.Title>
			<Dialog.Description>
				Delete cache entries matching the regex "{pendingPattern}"?
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button
				type="button"
				variant="outline"
				disabled={clearing || clearingPattern}
				onclick={closePatternConfirmation}
			>
				Cancel
			</Button>
			<Button
				type="button"
				variant="destructive"
				disabled={clearing || clearingPattern}
				onclick={confirmClearPattern}
			>
				{clearing || clearingPattern ? 'Evicting' : 'Evict'}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
