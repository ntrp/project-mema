<script lang="ts">
	import * as Select from '$lib/components/ui/select';

	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		value: string;
		options: Option[];
		onValueChange: (_value: string) => void;
		disabled?: boolean;
		placeholder?: string;
		size?: 'sm' | 'default';
	}

	let {
		value,
		options,
		onValueChange,
		disabled = false,
		placeholder = 'Select option',
		size = 'default'
	}: Props = $props();

	let selectedLabel = $derived(
		options.find((option) => option.value === value)?.label ?? placeholder
	);
</script>

<Select.Root type="single" {value} {disabled} {onValueChange}>
	<Select.Trigger {size} class="w-full">{selectedLabel}</Select.Trigger>
	<Select.Content>
		{#each options as option (option.value)}
			<Select.Item value={option.value} label={option.label} />
		{/each}
	</Select.Content>
</Select.Root>
