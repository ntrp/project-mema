<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type { SeasonOption } from './releaseOverrideSeriesOptions';

	interface Props {
		value: string;
		label: string;
		seasons: SeasonOption[];
		onChange: (_value: string) => void;
	}

	let { value, label, seasons, onChange }: Props = $props();
</script>

<div class="grid gap-1.5">
	<Label for="override-season">Season</Label>
	{#if seasons.length > 0}
		<Select.Root type="single" {value} onValueChange={onChange}>
			<Select.Trigger id="override-season" class="w-full">
				{label || 'Season'}
			</Select.Trigger>
			<Select.Content>
				{#each seasons as option (option.value)}
					<Select.Item value={option.value} label={option.label} />
				{/each}
			</Select.Content>
		</Select.Root>
	{:else}
		<Input
			id="override-season"
			type="number"
			min="0"
			step="1"
			{value}
			oninput={(event) => onChange(event.currentTarget.value)}
		/>
	{/if}
</div>
