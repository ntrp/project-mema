<script lang="ts">
	import { tick } from 'svelte';

	import { Textarea } from '$lib/components/ui/textarea';
	import { fileNamingTemplateParameters } from '$lib/settings/fileNamingTemplates';

	interface Props {
		value: string;
		onChange: (_value: string) => void;
	}

	let { value, onChange }: Props = $props();
	let textarea = $state<globalThis.HTMLTextAreaElement | null>(null);
	let cursor = $state(0);
	let tokenStart = $state(-1);
	let query = $state('');
	let activeIndex = $state(0);

	const suggestions = $derived.by(() => {
		if (tokenStart < 0) {
			return [];
		}
		const normalized = query.toLowerCase();
		return fileNamingTemplateParameters
			.filter((parameter) => parameter.includes(normalized))
			.slice(0, 8);
	});
	const showSuggestions = $derived(suggestions.length > 0);

	function handleInput(event: Event) {
		const target = event.currentTarget as globalThis.HTMLTextAreaElement;
		onChange(target.value);
		updateTokenState(target);
	}

	function handleSelect() {
		if (textarea) {
			updateTokenState(textarea);
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (!showSuggestions) {
			return;
		}
		if (event.key === 'ArrowDown') {
			event.preventDefault();
			activeIndex = Math.min(activeIndex + 1, suggestions.length - 1);
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			activeIndex = Math.max(activeIndex - 1, 0);
		} else if (event.key === 'Enter' || event.key === 'Tab') {
			event.preventDefault();
			void insertParameter(suggestions[activeIndex] ?? suggestions[0]);
		} else if (event.key === 'Escape') {
			closeSuggestions();
		}
	}

	function updateTokenState(target: globalThis.HTMLTextAreaElement) {
		cursor = target.selectionStart ?? target.value.length;
		const beforeCursor = target.value.slice(0, cursor);
		const openIndex = beforeCursor.lastIndexOf('{');
		const closeIndex = beforeCursor.lastIndexOf('}');
		const nextQuery = beforeCursor.slice(openIndex + 1);
		if (openIndex <= closeIndex || !/^[a-zA-Z0-9_:-]*$/.test(nextQuery)) {
			closeSuggestions();
			return;
		}
		tokenStart = openIndex;
		query = nextQuery;
		activeIndex = 0;
	}

	async function insertParameter(parameter: string) {
		if (!textarea || tokenStart < 0) {
			return;
		}
		const nextValue = `${value.slice(0, tokenStart)}{${parameter}}${value.slice(cursor)}`;
		const nextCursor = tokenStart + parameter.length + 2;
		onChange(nextValue);
		closeSuggestions();
		await tick();
		textarea.focus();
		textarea.setSelectionRange(nextCursor, nextCursor);
	}

	function closeSuggestions() {
		tokenStart = -1;
		query = '';
		activeIndex = 0;
	}
</script>

<div class="relative">
	<Textarea
		bind:ref={textarea}
		class="min-h-14 resize-y font-mono leading-snug"
		rows={2}
		required
		{value}
		oninput={handleInput}
		onselect={handleSelect}
		onkeyup={handleSelect}
		onclick={handleSelect}
		onkeydown={handleKeydown}
	/>
	{#if showSuggestions}
		<div
			class="absolute z-30 mt-1 grid w-full max-w-sm gap-1 rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-lg"
		>
			{#each suggestions as parameter, index (parameter)}
				<button
					type="button"
					class="rounded-sm px-2 py-1.5 text-left font-mono text-sm hover:bg-accent hover:text-accent-foreground data-[active=true]:bg-accent data-[active=true]:text-accent-foreground"
					data-active={index === activeIndex}
					onmousedown={(event) => event.preventDefault()}
					onclick={() => insertParameter(parameter)}
				>
					{parameter}
				</button>
			{/each}
		</div>
	{/if}
</div>
