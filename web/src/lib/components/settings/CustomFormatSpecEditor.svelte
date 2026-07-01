<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils';
	import type { CustomFormatSpec, CustomFormatSpecType } from '$lib/settings/types';

	interface Props {
		spec: CustomFormatSpec;
		labelPrefix: string;
		tone: 'include' | 'exclude';
		onChange: (_patch: Partial<CustomFormatSpec>) => void;
		onRemove: () => void;
	}

	const specTypes: { value: CustomFormatSpecType; label: string }[] = [
		{ value: 'releaseTitle', label: 'Release title regex' },
		{ value: 'source', label: 'Source' },
		{ value: 'resolution', label: 'Resolution' },
		{ value: 'quality', label: 'Quality' },
		{ value: 'videoCodec', label: 'Video codec' },
		{ value: 'audioCodec', label: 'Audio codec' },
		{ value: 'releaseGroup', label: 'Release group regex' },
		{ value: 'edition', label: 'Edition regex' },
		{ value: 'indexerFlag', label: 'Indexer flag' },
		{ value: 'language', label: 'Language' }
	];

	let { spec, labelPrefix, tone, onChange, onRemove }: Props = $props();

	function typeLabel() {
		return specTypes.find((type) => type.value === spec.type)?.label ?? 'Release title regex';
	}
</script>

<div
	class={cn(
		'grid gap-2 rounded-md border bg-muted/20 p-3',
		tone === 'include' ? 'border-primary/40' : 'border-destructive/40'
	)}
>
	<div class="grid gap-2 [grid-template-columns:repeat(auto-fit,minmax(min(100%,180px),1fr))]">
		<Input
			value={spec.name}
			type="text"
			placeholder="Label"
			aria-label={`${labelPrefix} condition label`}
			oninput={(event) => onChange({ name: event.currentTarget.value })}
		/>
		<Select.Root
			type="single"
			value={spec.type}
			onValueChange={(value: string) => onChange({ type: value as CustomFormatSpecType })}
		>
			<Select.Trigger class="w-full">{typeLabel()}</Select.Trigger>
			<Select.Content>
				{#each specTypes as type (type.value)}
					<Select.Item value={type.value} label={type.label} />
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	<Input
		value={spec.value}
		type="text"
		placeholder="Value or regex"
		aria-label={`${labelPrefix} condition value`}
		oninput={(event) => onChange({ value: event.currentTarget.value })}
	/>

	<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
		<label class="flex items-center gap-2 text-sm text-muted-foreground">
			<Checkbox
				checked={spec.required}
				onCheckedChange={(checked) => onChange({ required: checked === true })}
			/>
			<span>Required</span>
		</label>
		<Button
			type="button"
			variant="destructive"
			size="icon-sm"
			aria-label={`Remove ${labelPrefix.toLowerCase()} condition`}
			onclick={onRemove}
		>
			<TrashIcon aria-hidden="true" />
		</Button>
	</div>
</div>
