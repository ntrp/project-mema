<script lang="ts">
	import MinusIcon from '@lucide/svelte/icons/minus';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import XIcon from '@lucide/svelte/icons/x';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { DiscoverMovieFacetOption } from '$lib/settings/types';

	interface Props {
		id: string;
		label: string;
		values: string[];
		excludedValues?: string[];
		placeholder: string;
		options: DiscoverMovieFacetOption[];
		loading?: boolean;
		onQuery: (_query: string) => void;
		onChange: (_values: string[]) => void;
		onExcludedChange?: (_values: string[]) => void;
		onSignedChange?: (_values: string[], _excludedValues: string[]) => void;
	}

	let {
		id,
		label,
		values,
		excludedValues = [],
		placeholder,
		options,
		loading = false,
		onQuery,
		onChange,
		onExcludedChange,
		onSignedChange
	}: Props = $props();
	let query = $state('');

	const supportsExclusion = $derived(Boolean(onExcludedChange || onSignedChange));

	function add(value: string, excluded = false) {
		const cleaned = value.trim();
		if (!cleaned) return;
		const nextValues = values.filter((item) => item !== cleaned);
		const nextExcluded = excludedValues.filter((item) => item !== cleaned);
		if (excluded && onExcludedChange) {
			updateValues(nextValues, [...nextExcluded, cleaned]);
		} else {
			updateValues([...nextValues, cleaned], nextExcluded);
		}
		query = '';
		onQuery('');
	}

	function updateValues(nextValues: string[], nextExcluded: string[]) {
		if (onSignedChange) {
			onSignedChange(nextValues, nextExcluded);
			return;
		}
		onExcludedChange?.(nextExcluded);
		onChange(nextValues);
	}

	function remove(value: string) {
		onChange(values.filter((item) => item !== value));
	}

	function removeExcluded(value: string) {
		onExcludedChange?.(excludedValues.filter((item) => item !== value));
	}

	function toggle(value: string, excluded: boolean) {
		if (!supportsExclusion) return;
		add(value, !excluded);
	}

	function handleInput(event: Event) {
		const target = event.currentTarget as HTMLInputElement;
		query = target.value;
		onQuery(query);
	}
</script>

<div class="grid gap-2">
	<Label for={id}>{label}</Label>
	<div class="flex flex-wrap gap-1.5">
		{#each values as value (value)}
			<Badge variant="default" class="gap-1 bg-accent text-accent-foreground hover:bg-accent/90">
				<button
					type="button"
					class="max-w-36 truncate"
					onclick={() => toggle(value, false)}
					aria-label={supportsExclusion ? `Exclude ${value}` : value}
				>
					{value}
				</button>
				<button type="button" class="grid size-4 place-items-center" onclick={() => remove(value)}>
					<XIcon aria-hidden="true" class="size-3" />
					<span class="sr-only">Remove {value}</span>
				</button>
			</Badge>
		{/each}
		{#each excludedValues as value (value)}
			<Badge
				variant="destructive"
				class="gap-1 border-destructive/20 bg-destructive/15 text-destructive hover:bg-destructive/20"
			>
				<button
					type="button"
					class="max-w-36 truncate"
					onclick={() => toggle(value, true)}
					aria-label={`Include ${value}`}
				>
					{value}
				</button>
				<button
					type="button"
					class="grid size-4 place-items-center"
					onclick={() => removeExcluded(value)}
				>
					<XIcon aria-hidden="true" class="size-3" />
					<span class="sr-only">Remove excluded {value}</span>
				</button>
			</Badge>
		{/each}
	</div>
	<div class="relative">
		<div class="flex gap-2">
			<Input {id} value={query} {placeholder} autocomplete="off" oninput={handleInput} />
			<Button
				type="button"
				variant="default"
				size="icon"
				class="bg-accent text-accent-foreground hover:bg-accent/90"
				onclick={() => add(query)}
				aria-label={`Add ${label}`}
			>
				<PlusIcon aria-hidden="true" />
			</Button>
		</div>
		{#if query.length >= 2 && (options.length > 0 || loading)}
			<div
				class="absolute z-20 mt-1 grid max-h-52 w-full overflow-auto rounded-md border border-border bg-popover p-1 shadow-md"
			>
				{#if loading}
					<span class="px-2 py-2 text-sm text-muted-foreground">Loading</span>
				{/if}
				{#each options as option (option.id)}
					<div class="flex items-center gap-1 rounded-md px-1 py-1 hover:bg-muted">
						<button
							type="button"
							class="min-w-0 flex-1 truncate px-1 py-1 text-left text-sm"
							onclick={() => add(option.name)}
						>
							{option.name}
						</button>
						<Button
							type="button"
							variant="default"
							size="icon-sm"
							class="bg-accent text-accent-foreground shadow-xs hover:bg-accent/90"
							onclick={() => add(option.name)}
							aria-label={`Include ${option.name}`}
						>
							<PlusIcon aria-hidden="true" />
						</Button>
						{#if supportsExclusion}
							<Button
								type="button"
								variant="default"
								size="icon-sm"
								class="bg-destructive text-destructive-foreground shadow-xs hover:bg-destructive/90"
								onclick={() => add(option.name, true)}
								aria-label={`Exclude ${option.name}`}
							>
								<MinusIcon aria-hidden="true" />
							</Button>
						{/if}
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
