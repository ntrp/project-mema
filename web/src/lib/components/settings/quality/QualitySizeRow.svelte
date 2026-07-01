<script lang="ts">
	import type { QualitySizeSetting } from '$lib/settings/types';
	import {
		activeTrackStyle,
		gibValue,
		mbPerMinuteTitle,
		nextSliderQuality,
		rowError,
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
	const sliderFields: SliderField[] = ['minimum', 'preferred', 'maximum'];
	const validationError = $derived(rowError(quality));
	const values = $derived(sliderValues(quality));

	function updateSlider(field: SliderField, rawValue: string) {
		onChange(nextSliderQuality(quality, field, rawValue));
	}

	function updateGibValue(field: SliderField, rawValue: string) {
		const parsed = Number.parseFloat(rawValue);
		if (Number.isFinite(parsed)) {
			updateSlider(field, String(parsed));
		}
	}

	function scrollGibValue(event: globalThis.WheelEvent, field: SliderField, value: number) {
		event.preventDefault();
		const wheelDelta = event.deltaY !== 0 ? event.deltaY : event.deltaX;
		const direction = wheelDelta < 0 ? 1 : -1;
		const step = event.shiftKey ? 0.1 : 0.01;
		updateSlider(field, String(Math.round((value + direction * step) * 100) / 100));
	}
</script>

<tr class:invalid-row={validationError !== ''}>
	<td>
		<strong>{quality.name}</strong>
		<span>{quality.qualityId}</span>
	</td>
	<td>
		<div class="quality-size-limit-row">
			<label>
				<span class="quality-size-field-label">Min</span>
				<span class="quality-size-unit-input">
					<input
						type="number"
						min="0"
						max={values.maximum - sliderHandleGap * 2}
						step="0.01"
						required
						aria-label={`${quality.name} minimum size GiB per hour`}
						title={mbPerMinuteTitle('Minimum', values.minimum)}
						value={gibValue(values.minimum)}
						oninput={(event) => updateGibValue('minimum', event.currentTarget.value)}
						onwheel={(event) => scrollGibValue(event, 'minimum', values.minimum)}
					/>
					<span>GiB/h</span>
				</span>
			</label>
			<div class="quality-size-range">
				<div class="quality-size-track" aria-hidden="true"></div>
				<div
					class="quality-size-track active"
					style={activeTrackStyle(values)}
					aria-hidden="true"
				></div>
				{#each sliderFields as field (field)}
					<input
						type="range"
						min="0"
						max={sliderMaxGibPerHour}
						step={sliderStepGibPerHour}
						value={values[field]}
						oninput={(event) => updateSlider(field, event.currentTarget.value)}
						aria-label={`${quality.name} ${field} size slider`}
						class={`quality-size-range-input ${field}`}
					/>
				{/each}
			</div>
			<label>
				<span class="quality-size-field-label">Preferred</span>
				<span class="quality-size-unit-input">
					<input
						type="number"
						min={values.minimum}
						max={values.maximum}
						step="0.01"
						aria-label={`${quality.name} preferred size GiB per hour`}
						title={mbPerMinuteTitle('Preferred', values.preferred)}
						value={gibValue(values.preferred)}
						oninput={(event) => updateGibValue('preferred', event.currentTarget.value)}
						onwheel={(event) => scrollGibValue(event, 'preferred', values.preferred)}
					/>
					<span>GiB/h</span>
				</span>
			</label>
			<label>
				<span class="quality-size-field-label">Max</span>
				<span class="quality-size-unit-input">
					<input
						type="number"
						min={values.minimum + sliderHandleGap * 2}
						step="0.01"
						aria-label={`${quality.name} maximum size GiB per hour`}
						title={mbPerMinuteTitle('Maximum', values.maximum)}
						value={gibValue(values.maximum)}
						oninput={(event) => updateGibValue('maximum', event.currentTarget.value)}
						onwheel={(event) => scrollGibValue(event, 'maximum', values.maximum)}
					/>
					<span>GiB/h</span>
				</span>
			</label>
		</div>
		{#if validationError}
			<span class="quality-size-inline-error">{validationError}</span>
		{/if}
	</td>
</tr>
