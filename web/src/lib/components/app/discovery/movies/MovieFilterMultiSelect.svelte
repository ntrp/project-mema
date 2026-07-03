<script lang="ts">
	import XIcon from '@lucide/svelte/icons/x';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
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
		onChange: (_values: string[]) => void;
	}

	let { id, label, values, options, placeholder, onChange }: Props = $props();

	function selectedLabel(value: string) {
		return options.find((option) => option.value === value)?.label ?? value;
	}

	function toggle(value: string, checked: boolean) {
		onChange(checked ? addValue(values, value) : values.filter((item) => item !== value));
	}

	function addValue(current: string[], value: string) {
		return current.includes(value) ? current : [...current, value];
	}
</script>

<div class="grid gap-2">
	<Label for={id}>{label}</Label>
	<details class="group rounded-md border border-border bg-background">
		<summary
			{id}
			class="flex min-h-10 cursor-pointer list-none items-center gap-2 px-3 py-2 text-sm text-foreground marker:hidden"
		>
			{#if values.length > 0}
				<span class="flex min-w-0 flex-wrap gap-1.5">
					{#each values as value (value)}
						<Badge variant="secondary" class="max-w-36 truncate">{selectedLabel(value)}</Badge>
					{/each}
				</span>
			{:else}
				<span class="text-muted-foreground">{placeholder}</span>
			{/if}
		</summary>
		<div class="grid max-h-56 gap-1 overflow-auto border-t border-border p-2">
			{#if values.length > 0}
				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="justify-start"
					onclick={() => onChange([])}
				>
					<XIcon aria-hidden="true" />
					Clear
				</Button>
			{/if}
			{#each options as option (option.value)}
				<label class="flex min-h-9 items-center gap-2 rounded-md px-2 text-sm hover:bg-muted">
					<Checkbox
						checked={values.includes(option.value)}
						onCheckedChange={(checked) => toggle(option.value, checked === true)}
					/>
					<span>{option.label}</span>
				</label>
			{/each}
		</div>
	</details>
</div>
