<script lang="ts">
	import { onDestroy, tick } from 'svelte';

	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { fileNamingTemplateSuggestions } from '$lib/settings/fileNamingTemplates';

	interface Props {
		value: string;
		onChange: (_value: string) => void;
	}

	let { value, onChange }: Props = $props();
	let input = $state<globalThis.HTMLInputElement | null>(null);
	let cursor = $state(0);
	let tokenStart = $state(-1);
	let query = $state('');
	let activeIndex = $state(0);
	let tooltipOpenIndex = $state(-1);
	let tooltipTimer: ReturnType<typeof globalThis.setTimeout> | null = null;

	const suggestions = $derived.by(() => {
		if (tokenStart < 0) {
			return [];
		}
		return fileNamingTemplateSuggestions(query);
	});
	const showSuggestions = $derived(suggestions.length > 0);

	onDestroy(clearDescriptionTooltip);

	function handleInput(event: Event) {
		const target = event.currentTarget as globalThis.HTMLInputElement;
		clearDescriptionTooltip();
		onChange(target.value);
		updateTokenState(target);
	}

	function handleSelect() {
		if (input) {
			updateTokenState(input);
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (!showSuggestions) {
			return;
		}
		if (event.key === 'Escape') {
			event.preventDefault();
			event.stopPropagation();
			closeSuggestions();
		} else if (event.key === 'ArrowDown') {
			event.preventDefault();
			moveActiveIndex(Math.min(activeIndex + 1, suggestions.length - 1));
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			moveActiveIndex(Math.max(activeIndex - 1, 0));
		} else if (event.key === 'Enter' || event.key === 'Tab') {
			event.preventDefault();
			void insertParameter((suggestions[activeIndex] ?? suggestions[0]).param);
		}
	}

	function moveActiveIndex(nextIndex: number) {
		activeIndex = nextIndex;
		scheduleDescriptionTooltip(nextIndex);
	}

	function scheduleDescriptionTooltip(index: number) {
		clearDescriptionTooltip();
		tooltipTimer = globalThis.setTimeout(() => {
			if (showSuggestions && activeIndex === index) {
				tooltipOpenIndex = index;
			}
		}, 1000);
	}

	function clearDescriptionTooltip() {
		if (tooltipTimer) {
			globalThis.clearTimeout(tooltipTimer);
			tooltipTimer = null;
		}
		tooltipOpenIndex = -1;
	}

	function updateTokenState(target: globalThis.HTMLInputElement) {
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
		if (nextQuery !== query) {
			clearDescriptionTooltip();
			query = nextQuery;
			activeIndex = 0;
			return;
		}
		activeIndex = Math.min(activeIndex, Math.max(suggestions.length - 1, 0));
	}

	async function insertParameter(parameter: string) {
		if (!input || tokenStart < 0) {
			return;
		}
		const nextValue = `${value.slice(0, tokenStart)}{${parameter}}${value.slice(cursor)}`;
		const nextCursor = tokenStart + parameter.length + 2;
		onChange(nextValue);
		closeSuggestions();
		await tick();
		input.focus();
		input.setSelectionRange(nextCursor, nextCursor);
	}

	function closeSuggestions() {
		clearDescriptionTooltip();
		tokenStart = -1;
		query = '';
		activeIndex = 0;
	}
</script>

<div class="relative">
	<Input
		bind:ref={input}
		class="font-mono"
		required
		{value}
		oninput={handleInput}
		onselect={handleSelect}
		onclick={handleSelect}
		onkeydown={handleKeydown}
	/>
	{#if showSuggestions}
		<div
			class="absolute z-50 mt-1 grid w-[min(32rem,calc(100vw-2rem))] gap-1 rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-lg"
		>
			{#each suggestions as parameter, index (parameter.param)}
				<Tooltip.Root open={tooltipOpenIndex === index}>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<button
								type="button"
								class="grid min-w-0 grid-cols-[minmax(8rem,12rem)_minmax(10rem,1fr)] items-baseline gap-3 rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent hover:text-accent-foreground data-[active=true]:bg-accent data-[active=true]:text-accent-foreground"
								data-active={index === activeIndex}
								onmousedown={(event) => event.preventDefault()}
								onmousemove={clearDescriptionTooltip}
								onclick={() => insertParameter(parameter.param)}
								{...props}
							>
								<span class="truncate font-mono font-bold">{parameter.param}</span>
								<span class="truncate font-mono text-xs text-muted-foreground">{parameter.example}</span>
							</button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>{parameter.description}</Tooltip.Content>
				</Tooltip.Root>
			{/each}
		</div>
	{/if}
</div>
