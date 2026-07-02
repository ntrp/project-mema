<script lang="ts">
	import { onMount } from 'svelte';
	import { Slider } from '$lib/components/ui/slider';
	import * as Table from '$lib/components/ui/table';
	import type { QualitySizeSetting } from '$lib/settings/types';
	import QualitySizeField from './QualitySizeField.svelte';
	import {
		nextSliderQuality,
		rowError,
		activeTrackStyle,
		sliderHandleGap,
		sliderMaxGibPerHour,
		sliderStepGibPerHour,
		sliderValues,
		type SliderField
	} from './qualitySize';

	interface Props {
		quality: QualitySizeSetting;
		onChange: (_quality: QualitySizeSetting) => void;
	}

	let { quality, onChange }: Props = $props();
	const validationError = $derived(rowError(quality));
	const values = $derived(sliderValues(quality));
	const rangeValue = $derived([values.minimum, values.preferred, values.maximum]);

	onMount(() => {
		window.addEventListener('wheel', scrollInputField, { capture: true, passive: false });
		return () => window.removeEventListener('wheel', scrollInputField, { capture: true });
	});

	function updateSlider(field: SliderField, rawValue: string) {
		onChange(nextSliderQuality(quality, field, rawValue));
	}

	function updateGibValue(field: SliderField, rawValue: string) {
		const parsed = Number.parseFloat(rawValue);
		if (Number.isFinite(parsed)) {
			updateSlider(field, String(parsed));
		}
	}

	function updateRange(nextValue: number[]) {
		const [minimum, preferred, maximum] = nextValue;
		onChange(
			nextSliderQuality(
				nextSliderQuality(
					nextSliderQuality(quality, 'minimum', String(minimum)),
					'maximum',
					String(maximum)
				),
				'preferred',
				String(preferred)
			)
		);
	}

	function sliderThumbClass(index: number) {
		return index === 1
			? 'z-20 border-primary bg-primary text-primary-foreground focus-visible:ring-primary'
			: 'z-10';
	}

	function scrollGibValue(event: globalThis.WheelEvent, field: SliderField, value: number) {
		event.preventDefault();
		event.stopPropagation();
		const wheelDelta = event.deltaY !== 0 ? event.deltaY : event.deltaX;
		const direction = wheelDelta < 0 ? 1 : -1;
		const step = event.shiftKey ? 1 : 0.1;
		updateSlider(field, String(Math.round((value + direction * step) * 100) / 100));
	}

	function scrollInputField(event: globalThis.WheelEvent) {
		const target = event.target;
		if (!(target instanceof globalThis.HTMLInputElement)) {
			return;
		}
		const label = target.getAttribute('aria-label') ?? '';
		if (!label.startsWith(`${quality.name} `)) {
			return;
		}
		if (label.endsWith('minimum size GiB per hour')) {
			scrollGibValue(event, 'minimum', values.minimum);
		} else if (label.endsWith('preferred size GiB per hour')) {
			scrollGibValue(event, 'preferred', values.preferred);
		} else if (label.endsWith('maximum size GiB per hour')) {
			scrollGibValue(event, 'maximum', values.maximum);
		}
	}

	function focusInput(event: globalThis.PointerEvent) {
		if (event.currentTarget instanceof globalThis.HTMLInputElement) {
			event.currentTarget.focus();
		}
	}
</script>

<Table.Row>
	<Table.Cell class="w-[180px] min-w-40 align-top whitespace-normal">
		<strong class="block text-sm font-extrabold text-foreground">{quality.name}</strong>
		<span class="block text-xs font-extrabold text-muted-foreground">{quality.qualityId}</span>
	</Table.Cell>
	<Table.Cell class="align-top whitespace-normal">
		<div class="grid grid-cols-[124px_minmax(240px,1fr)_124px_124px] items-center gap-3">
			<QualitySizeField
				field="minimum"
				label="Min"
				qualityName={quality.name}
				value={values.minimum}
				min={0}
				max={values.maximum - sliderHandleGap * 2}
				{validationError}
				onChange={updateGibValue}
				onScroll={scrollGibValue}
				onFocus={focusInput}
			/>
			<Slider
				value={rangeValue}
				min={0}
				max={sliderMaxGibPerHour}
				step={sliderStepGibPerHour}
				ariaLabel={`${quality.name} size limits`}
				class="self-center px-3"
				rangeClass="hidden"
				rangeStyle={activeTrackStyle(values)}
				thumbClass={sliderThumbClass}
				onValueChange={updateRange}
			/>
			<QualitySizeField
				field="preferred"
				label="Preferred"
				qualityName={quality.name}
				value={values.preferred}
				min={values.minimum}
				max={values.maximum}
				{validationError}
				onChange={updateGibValue}
				onScroll={scrollGibValue}
				onFocus={focusInput}
			/>
			<QualitySizeField
				field="maximum"
				label="Max"
				qualityName={quality.name}
				value={values.maximum}
				min={values.minimum + sliderHandleGap * 2}
				{validationError}
				onChange={updateGibValue}
				onScroll={scrollGibValue}
				onFocus={focusInput}
			/>
		</div>
		{#if validationError}
			<span class="mt-2 block text-xs font-black text-destructive">{validationError}</span>
		{/if}
	</Table.Cell>
</Table.Row>
