<script lang="ts">
	import { onMount } from 'svelte';
	import { Input } from '$lib/components/ui/input';
	import { Slider } from '$lib/components/ui/slider';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { QualitySizeSetting } from '$lib/settings/types';
	import {
		gibValue,
		mbPerMinuteTitle,
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

	function fieldInputClass() {
		return validationError ? 'min-w-0 border-destructive pr-12' : 'min-w-0 pr-12';
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
			<label class="relative block w-[124px]">
				<span
					class="absolute -top-2 left-2 z-10 bg-card px-1 text-xs font-black text-muted-foreground"
				>
					Min
				</span>
				<span
					class="relative block"
					onwheel={(event) => scrollGibValue(event, 'minimum', values.minimum)}
				>
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Input
									{...props}
									type="number"
									min="0"
									max={values.maximum - sliderHandleGap * 2}
									step={0.01}
									required
									aria-label={`${quality.name} minimum size GiB per hour`}
									class={fieldInputClass()}
									value={gibValue(values.minimum)}
									oninput={(event) => updateGibValue('minimum', event.currentTarget.value)}
									onpointerenter={focusInput}
									onwheelcapture={(event) => scrollGibValue(event, 'minimum', values.minimum)}
									onwheel={(event) => scrollGibValue(event, 'minimum', values.minimum)}
								/>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>{mbPerMinuteTitle('Minimum', values.minimum)}</Tooltip.Content>
					</Tooltip.Root>
					<span
						class="pointer-events-none absolute top-1/2 right-2 -translate-y-1/2 text-xs font-black text-muted-foreground"
					>
						GiB/h
					</span>
				</span>
			</label>
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
			<label class="relative block w-[124px]">
				<span
					class="absolute -top-2 left-2 z-10 bg-card px-1 text-xs font-black text-muted-foreground"
				>
					Preferred
				</span>
				<span
					class="relative block"
					onwheel={(event) => scrollGibValue(event, 'preferred', values.preferred)}
				>
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Input
									{...props}
									type="number"
									min={values.minimum}
									max={values.maximum}
									step={0.01}
									aria-label={`${quality.name} preferred size GiB per hour`}
									class={fieldInputClass()}
									value={gibValue(values.preferred)}
									oninput={(event) => updateGibValue('preferred', event.currentTarget.value)}
									onpointerenter={focusInput}
									onwheelcapture={(event) => scrollGibValue(event, 'preferred', values.preferred)}
									onwheel={(event) => scrollGibValue(event, 'preferred', values.preferred)}
								/>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>{mbPerMinuteTitle('Preferred', values.preferred)}</Tooltip.Content>
					</Tooltip.Root>
					<span
						class="pointer-events-none absolute top-1/2 right-2 -translate-y-1/2 text-xs font-black text-muted-foreground"
					>
						GiB/h
					</span>
				</span>
			</label>
			<label class="relative block w-[124px]">
				<span
					class="absolute -top-2 left-2 z-10 bg-card px-1 text-xs font-black text-muted-foreground"
				>
					Max
				</span>
				<span
					class="relative block"
					onwheel={(event) => scrollGibValue(event, 'maximum', values.maximum)}
				>
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Input
									{...props}
									type="number"
									min={values.minimum + sliderHandleGap * 2}
									step={0.01}
									aria-label={`${quality.name} maximum size GiB per hour`}
									class={fieldInputClass()}
									value={gibValue(values.maximum)}
									oninput={(event) => updateGibValue('maximum', event.currentTarget.value)}
									onpointerenter={focusInput}
									onwheelcapture={(event) => scrollGibValue(event, 'maximum', values.maximum)}
									onwheel={(event) => scrollGibValue(event, 'maximum', values.maximum)}
								/>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>{mbPerMinuteTitle('Maximum', values.maximum)}</Tooltip.Content>
					</Tooltip.Root>
					<span
						class="pointer-events-none absolute top-1/2 right-2 -translate-y-1/2 text-xs font-black text-muted-foreground"
					>
						GiB/h
					</span>
				</span>
			</label>
		</div>
		{#if validationError}
			<span class="mt-2 block text-xs font-black text-destructive">{validationError}</span>
		{/if}
	</Table.Cell>
</Table.Row>
