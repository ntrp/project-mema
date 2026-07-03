<script lang="ts">
	import { Slider as SliderPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	type SliderProps = {
		value: number[];
		min?: number;
		max?: number;
		step?: number | number[];
		disabled?: boolean;
		ariaLabel?: string;
		onValueChange?: (_value: number[]) => void;
		onValueCommit?: (_value: number[]) => void;
		ref?: HTMLElement | null;
		class?: string;
		rangeClass?: string;
		rangeStyle?: string;
		thumbClass?: string | ((_index: number) => string);
	};

	let {
		ref = $bindable(null),
		class: className,
		rangeClass,
		rangeStyle,
		thumbClass,
		value,
		min = 0,
		max = 100,
		step = 1,
		disabled = false,
		ariaLabel,
		onValueChange,
		onValueCommit
	}: SliderProps = $props();

	function resolvedThumbClass(index: number) {
		return typeof thumbClass === 'function' ? thumbClass(index) : thumbClass;
	}

	function handleValueChange(nextValue: number[]) {
		onValueChange?.(Array.isArray(nextValue) ? nextValue : [nextValue]);
	}

	function handleValueCommit(nextValue: number[]) {
		onValueCommit?.(Array.isArray(nextValue) ? nextValue : [nextValue]);
	}
</script>

<SliderPrimitive.Root
	bind:ref
	type="multiple"
	data-slot="slider"
	class={cn(
		'relative flex w-full touch-none items-center select-none data-disabled:opacity-50',
		className
	)}
	{value}
	{min}
	{max}
	{step}
	{disabled}
	aria-label={ariaLabel}
	onValueChange={handleValueChange}
	onValueCommit={handleValueCommit}
>
	{#snippet children({ thumbItems })}
		<span data-slot="slider-track" class="bg-muted relative h-1.5 w-full grow rounded-full">
			<SliderPrimitive.Range
				data-slot="slider-range"
				class={cn('bg-primary absolute h-full rounded-full', rangeClass)}
			/>
			{#if rangeStyle}
				<span
					data-slot="slider-custom-range"
					class="absolute h-full rounded-full bg-primary"
					style={rangeStyle}
				></span>
			{/if}
		</span>
		{#each thumbItems as thumb (thumb.index)}
			<SliderPrimitive.Thumb
				index={thumb.index}
				data-slot="slider-thumb"
				class={cn(
					'border-primary bg-background ring-ring/50 focus-visible:ring-ring/50 block size-4 shrink-0 cursor-pointer rounded-full border shadow-sm transition-[color,box-shadow] outline-none focus-visible:ring-3 disabled:pointer-events-none disabled:opacity-50',
					resolvedThumbClass(thumb.index)
				)}
			/>
		{/each}
	{/snippet}
</SliderPrimitive.Root>
