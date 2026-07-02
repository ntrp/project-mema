<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { gibValue, mbPerMinuteTitle, type SliderField } from './qualitySize';

	interface Props {
		field: SliderField;
		label: string;
		qualityName: string;
		value: number;
		min?: number;
		max?: number;
		validationError: string;
		onChange: (_field: SliderField, _rawValue: string) => void;
		onScroll: (_event: globalThis.WheelEvent, _field: SliderField, _value: number) => void;
		onFocus: (_event: globalThis.PointerEvent) => void;
	}

	let {
		field,
		label,
		qualityName,
		value,
		min,
		max,
		validationError,
		onChange,
		onScroll,
		onFocus
	}: Props = $props();
	const fieldClass = $derived(
		validationError ? 'min-w-0 border-destructive pr-12' : 'min-w-0 pr-12'
	);
</script>

<label class="relative block w-[124px]">
	<span class="absolute -top-2 left-2 z-10 bg-card px-1 text-xs font-black text-muted-foreground">
		{label}
	</span>
	<span class="relative block" onwheel={(event) => onScroll(event, field, value)}>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Input
						{...props}
						type="number"
						{min}
						{max}
						step={0.01}
						required={field === 'minimum'}
						aria-label={`${qualityName} ${field} size GiB per hour`}
						class={fieldClass}
						value={gibValue(value)}
						oninput={(event) => onChange(field, event.currentTarget.value)}
						onpointerenter={onFocus}
						onwheelcapture={(event) => onScroll(event, field, value)}
						onwheel={(event) => onScroll(event, field, value)}
					/>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>{mbPerMinuteTitle(label, value)}</Tooltip.Content>
		</Tooltip.Root>
		<span
			class="pointer-events-none absolute top-1/2 right-2 -translate-y-1/2 text-xs font-black text-muted-foreground"
		>
			GiB/h
		</span>
	</span>
</label>
