<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { profileLanguageOptions } from '$lib/settings/languageCatalog';
	import type { Language } from '$lib/settings/types';

	interface Props {
		id: string;
		label: string;
		placeholder: string;
		languages: Language[];
		selectedIds: string[];
		onSelect: (_languageId: string) => void;
	}

	let { id, label, placeholder, languages, selectedIds, onSelect }: Props = $props();
	let query = $state('');
	let open = $state(false);
	const selected = $derived(new Set(selectedIds));
	const normalizedQuery = $derived(query.trim().toLowerCase());
	const options = $derived(
		profileLanguageOptions(languages, selectedIds)
			.filter((option) => !selected.has(option.id))
			.filter(
				(option) =>
					normalizedQuery === '' ||
					option.displayLabel.toLowerCase().includes(normalizedQuery) ||
					option.id.toLowerCase().includes(normalizedQuery)
			)
			.slice(0, 10)
	);

	function select(languageId: string) {
		onSelect(languageId);
		query = '';
		open = false;
	}
</script>

<div class="relative grid gap-1.5 text-sm">
	<Label for={id}>{label}</Label>
	<Input
		{id}
		bind:value={query}
		type="search"
		autocomplete="off"
		{placeholder}
		role="combobox"
		aria-expanded={open}
		aria-controls={`${id}-options`}
		onfocus={() => (open = true)}
		oninput={() => (open = true)}
		onblur={() => globalThis.setTimeout(() => (open = false), 120)}
	/>
	{#if open}
		<div
			id={`${id}-options`}
			role="listbox"
			class="absolute top-full z-20 mt-1 grid max-h-64 w-full overflow-auto rounded-md border border-border bg-popover p-1 shadow-lg"
		>
			{#each options as option (option.id)}
				<Button
					type="button"
					variant="ghost"
					class="h-8 justify-start px-2"
					onmousedown={(event) => event.preventDefault()}
					onclick={() => select(option.id)}
				>
					<PlusIcon aria-hidden="true" />
					<span class="truncate">{option.displayLabel}</span>
				</Button>
			{:else}
				<p class="m-0 px-2 py-1.5 text-sm text-muted-foreground">No languages found</p>
			{/each}
		</div>
	{/if}
</div>
