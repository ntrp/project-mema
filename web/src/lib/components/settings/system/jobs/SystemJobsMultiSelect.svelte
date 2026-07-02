<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { selectedFirst } from '$lib/components/shared/multiSelectOrdering';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label';

	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		id: string;
		label: string;
		values: string[];
		options: Option[];
		placeholder: string;
	}

	let { id, label, values = $bindable(), options, placeholder }: Props = $props();

	const selectedLabel = $derived(summary(values, options, placeholder));
	const selectedSet = $derived(new Set(values));
	const sortedOptions = $derived(selectedFirst(options, selectedSet, (option) => option.value));

	function toggle(value: string, checked: boolean) {
		values = checked ? unique([...values, value]) : values.filter((item) => item !== value);
	}

	function clear() {
		values = [];
	}

	function unique(items: string[]) {
		return Array.from(new Set(items));
	}

	function summary(selected: string[], allOptions: Option[], emptyLabel: string) {
		if (selected.length === 0) return emptyLabel;
		if (selected.length === allOptions.length) return 'All';
		if (selected.length === 1) {
			return allOptions.find((option) => option.value === selected[0])?.label ?? selected[0];
		}
		return `${selected.length} selected`;
	}
</script>

<div class="grid gap-1.5">
	<Label for={id}>{label}</Label>
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					{id}
					type="button"
					variant="outline"
					class="w-full justify-between gap-2"
				>
					<span class="truncate">{selectedLabel}</span>
					<ChevronDownIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="start" class="max-h-72 w-64">
			<DropdownMenu.Item onclick={clear}>
				<span class="text-muted-foreground">All</span>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			{#each sortedOptions as option (option.value)}
				<DropdownMenu.CheckboxItem
					checked={values.includes(option.value)}
					onCheckedChange={(checked) => toggle(option.value, checked === true)}
				>
					<span class="truncate">{option.label}</span>
				</DropdownMenu.CheckboxItem>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
