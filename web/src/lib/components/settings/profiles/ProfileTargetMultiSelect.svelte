<script lang="ts">
	import XIcon from '@lucide/svelte/icons/x';
	import { selectedFirst } from '$lib/components/shared/multiSelectOrdering';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label';

	type Option = {
		value: string;
		label: string;
	};

	interface Props {
		id: string;
		label: string;
		values: string[];
		options: Option[];
		placeholder: string;
		labelClass?: string;
		onChange: (_values: string[]) => void;
	}

	let { id, label, values, options, placeholder, labelClass = '', onChange }: Props = $props();
	const selectedSet = $derived(new Set(values));
	const selectedOptions = $derived(options.filter((option) => selectedSet.has(option.value)));
	const sortedOptions = $derived(selectedFirst(options, selectedSet, (option) => option.value));

	function toggle(value: string, checked: boolean) {
		onChange(checked ? unique([...values, value]) : values.filter((item) => item !== value));
	}

	function clear() {
		onChange([]);
	}

	function unique(items: string[]) {
		return Array.from(new Set(items));
	}
</script>

<div class="grid gap-1.5">
	<Label for={id} class={labelClass}>{label}</Label>
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					{id}
					type="button"
					variant="outline"
					class="h-auto min-h-9 w-full justify-start gap-2 px-2 py-1.5"
				>
					{#if selectedOptions.length > 0}
						<span class="flex min-w-0 flex-wrap gap-1">
							{#each selectedOptions as option (option.value)}
								<Badge variant="secondary" class="max-w-36 truncate">{option.label}</Badge>
							{/each}
						</span>
					{:else}
						<span class="truncate text-muted-foreground">{placeholder}</span>
					{/if}
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="start" class="max-h-72 w-72">
			{#if values.length > 0}
				<DropdownMenu.Item onclick={clear}>
					<XIcon aria-hidden="true" />
					<span class="text-muted-foreground">Clear</span>
				</DropdownMenu.Item>
				<DropdownMenu.Separator />
			{/if}
			{#each sortedOptions as option (option.value)}
				<DropdownMenu.CheckboxItem
					checked={values.includes(option.value)}
					onCheckedChange={(checked) => toggle(option.value, checked === true)}
				>
					{option.label}
				</DropdownMenu.CheckboxItem>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
