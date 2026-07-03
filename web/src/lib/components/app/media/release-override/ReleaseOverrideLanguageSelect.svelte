<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { selectedFirst } from '$lib/components/shared/multiSelectOrdering';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label';
	import {
		languageLabelFromCatalog,
		languageOptionsFromCatalog
	} from '$lib/settings/languageCatalog';
	import type { Language } from '$lib/settings/types';

	interface Props {
		values: string[];
		languages: Language[];
	}

	let { values = $bindable(), languages }: Props = $props();

	const options = $derived(languageOptionsFromCatalog(languages));
	const selectedSet = $derived(new Set(values));
	const sortedOptions = $derived(selectedFirst(options, selectedSet, (option) => option.id));
	const selectedLabels = $derived(
		values.map((value) => languageLabelFromCatalog(value, languages))
	);

	function toggle(value: string, checked: boolean) {
		values = checked ? unique([...values, value]) : values.filter((item) => item !== value);
	}

	function clear() {
		values = [];
	}

	function unique(items: string[]) {
		return Array.from(new Set(items));
	}
</script>

<div class="grid gap-1.5">
	<Label for="override-languages">Languages</Label>
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					id="override-languages"
					type="button"
					variant="outline"
					class="min-h-9 w-full justify-between gap-2 py-1.5"
				>
					<span class="flex min-w-0 flex-1 flex-wrap gap-1">
						{#if selectedLabels.length > 0}
							{#each selectedLabels as label (label)}
								<Badge variant="secondary" class="max-w-32 truncate">{label}</Badge>
							{/each}
						{:else}
							<span class="truncate text-muted-foreground">Select languages</span>
						{/if}
					</span>
					<ChevronDownIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="start" class="max-h-72 w-72">
			<DropdownMenu.Item onclick={clear}>
				<span class="text-muted-foreground">Clear languages</span>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			{#each sortedOptions as option (option.id)}
				<DropdownMenu.CheckboxItem
					checked={values.includes(option.id)}
					onCheckedChange={(checked) => toggle(option.id, checked === true)}
				>
					<span class="truncate">{option.displayLabel}</span>
				</DropdownMenu.CheckboxItem>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
