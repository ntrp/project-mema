<script lang="ts">
	import { onMount } from 'svelte';

	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import { listQualitySizeSettings, updateQualitySizeSettings } from '$lib/settings/api';
	import { groupQualitiesByResolution } from '$lib/settings/qualityGroups';
	import type { QualitySizeSetting, QualitySizeSettingRequest } from '$lib/settings/types';

	type SliderField = 'minimum' | 'preferred' | 'maximum';

	const sliderMaxGibPerHour = 120;
	const sliderStepGibPerHour = 0.1;
	const sliderHandleGap = 0.1;

	let qualities = $state<QualitySizeSetting[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	let errorMessage = $state('');
	let message = $state('');

	const hasValidationErrors = $derived(qualities.some((quality) => rowError(quality) !== ''));
	const qualityGroups = $derived(groupQualitiesByResolution(qualities));

	onMount(() => {
		void loadQualitySizes();
	});

	async function loadQualitySizes() {
		loading = true;
		errorMessage = '';
		try {
			const response = await listQualitySizeSettings();
			qualities = response.qualities;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load quality sizes';
		} finally {
			loading = false;
		}
	}

	async function saveQualitySizes(event: SubmitEvent) {
		event.preventDefault();
		message = '';
		errorMessage = '';
		if (hasValidationErrors) {
			errorMessage = 'Fix invalid quality sizes before saving';
			return;
		}

		saving = true;
		try {
			const response = await updateQualitySizeSettings(qualities.map(qualityRequest));
			qualities = response.qualities;
			message = 'Quality sizes saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save quality sizes';
		} finally {
			saving = false;
		}
	}

	function rowError(quality: QualitySizeSetting) {
		const minimum = quality.minimumSizeMbPerMinute;
		const preferred = quality.preferredSizeMbPerMinute;
		const maximum = quality.maximumSizeMbPerMinute;
		if (minimum < 0) {
			return 'Minimum must be zero or greater';
		}
		if (preferred != null && preferred < minimum) {
			return 'Preferred must be at least minimum';
		}
		if (maximum != null && maximum < minimum) {
			return 'Maximum must be at least minimum';
		}
		if (preferred != null && maximum != null && preferred > maximum) {
			return 'Preferred must be at most maximum';
		}
		return '';
	}

	function qualityRequest(quality: QualitySizeSetting): QualitySizeSettingRequest {
		return {
			qualityId: quality.qualityId,
			minimumSizeMbPerMinute: quality.minimumSizeMbPerMinute,
			preferredSizeMbPerMinute: quality.preferredSizeMbPerMinute ?? null,
			maximumSizeMbPerMinute: quality.maximumSizeMbPerMinute ?? null
		};
	}

	function mbPerMinuteToGibPerHour(value: number | null | undefined) {
		if (value == null) {
			return sliderMaxGibPerHour;
		}
		return clamp(Math.round(((value * 60) / 1024) * 100) / 100, 0, sliderMaxGibPerHour);
	}

	function gibPerHourToMbPerMinute(value: number) {
		return Math.round(((value * 1024) / 60) * 100) / 100;
	}

	function sliderValues(quality: QualitySizeSetting) {
		const minimum = mbPerMinuteToGibPerHour(quality.minimumSizeMbPerMinute);
		const maximum = mbPerMinuteToGibPerHour(quality.maximumSizeMbPerMinute);
		const preferred = clamp(
			mbPerMinuteToGibPerHour(quality.preferredSizeMbPerMinute),
			minimum,
			maximum
		);
		return { minimum, preferred, maximum };
	}

	function updateSlider(qualityId: string, field: SliderField, rawValue: string) {
		const nextValue = clamp(Number.parseFloat(rawValue), 0, sliderMaxGibPerHour);
		if (!Number.isFinite(nextValue)) {
			return;
		}
		qualities = qualities.map((quality) => {
			if (quality.qualityId !== qualityId) {
				return quality;
			}
			let { minimum, preferred, maximum } = sliderValues(quality);
			if (field === 'minimum') {
				minimum = clamp(nextValue, 0, Math.max(0, maximum - sliderHandleGap * 2));
				preferred = Math.max(preferred, minimum);
			} else if (field === 'preferred') {
				preferred = clamp(nextValue, minimum, maximum);
			} else {
				maximum = clamp(nextValue, minimum + sliderHandleGap * 2, sliderMaxGibPerHour);
				preferred = Math.min(preferred, maximum);
			}
			return {
				...quality,
				minimumSizeMbPerMinute: gibPerHourToMbPerMinute(minimum),
				preferredSizeMbPerMinute:
					preferred >= sliderMaxGibPerHour ? null : gibPerHourToMbPerMinute(preferred),
				maximumSizeMbPerMinute:
					maximum >= sliderMaxGibPerHour ? null : gibPerHourToMbPerMinute(maximum)
			};
		});
		message = '';
	}

	function updateGibValue(qualityId: string, field: SliderField, rawValue: string) {
		const parsed = Number.parseFloat(rawValue);
		if (!Number.isFinite(parsed)) {
			return;
		}
		updateSlider(qualityId, field, String(parsed));
	}

	function scrollGibValue(
		event: globalThis.WheelEvent,
		qualityId: string,
		field: SliderField,
		value: number
	) {
		event.preventDefault();
		const wheelDelta = event.deltaY !== 0 ? event.deltaY : event.deltaX;
		const direction = wheelDelta < 0 ? 1 : -1;
		const step = event.shiftKey ? 0.1 : 0.01;
		updateSlider(qualityId, field, String(Math.round((value + direction * step) * 100) / 100));
	}

	function gibValue(value: number) {
		return value.toFixed(2);
	}

	function mbPerMinuteTitle(label: string, value: number) {
		return `${label}: ${gibPerHourToMbPerMinute(value).toFixed(2)} MB/m`;
	}

	function labelOffset(value: number) {
		return `${(value / sliderMaxGibPerHour) * 100}%`;
	}

	function activeTrackStyle(values: { minimum: number; maximum: number }) {
		return `--min-pct: ${labelOffset(values.minimum)}; --max-pct: ${labelOffset(values.maximum)}`;
	}

	function clamp(value: number, minimum: number, maximum: number) {
		return Math.min(Math.max(value, minimum), maximum);
	}
</script>

<div class="panel quality-size-panel" aria-labelledby="quality-size-title">
	<form onsubmit={saveQualitySizes}>
		<div class="section-heading">
			<div>
				<p class="section-kicker">Release scoring</p>
				<h2 id="quality-size-title">Quality sizes</h2>
			</div>
			<div class="quality-size-actions">
				<button
					type="button"
					class="secondary"
					disabled={loading || saving}
					onclick={loadQualitySizes}
				>
					Reload
				</button>
				<button type="submit" disabled={loading || saving || hasValidationErrors}>
					{saving ? 'Saving' : 'Save sizes'}
				</button>
			</div>
		</div>

		<NoticeStack {message} {errorMessage} />

		<div class="table-wrap">
			<table class="quality-size-table">
				<thead>
					<tr>
						<th>Quality</th>
						<th>Size limit</th>
					</tr>
				</thead>
				<tbody>
					{#if loading}
						<tr>
							<td colspan="2" class="empty">Loading quality sizes</td>
						</tr>
					{:else if qualities.length === 0}
						<tr>
							<td colspan="2" class="empty">No qualities loaded</td>
						</tr>
					{:else}
						{#each qualityGroups as group (group.id)}
							<tr class="quality-size-group-row">
								<th colspan="2">{group.label}</th>
							</tr>
							{#each group.qualities as quality (quality.qualityId)}
								{@const validationError = rowError(quality)}
								{@const values = sliderValues(quality)}
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
														oninput={(event) =>
															updateGibValue(
																quality.qualityId,
																'minimum',
																event.currentTarget.value
															)}
														onwheel={(event) =>
															scrollGibValue(event, quality.qualityId, 'minimum', values.minimum)}
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
												<input
													type="range"
													min="0"
													max={sliderMaxGibPerHour}
													step={sliderStepGibPerHour}
													value={values.minimum}
													oninput={(event) =>
														updateSlider(quality.qualityId, 'minimum', event.currentTarget.value)}
													aria-label={`${quality.name} minimum size slider`}
													class="quality-size-range-input minimum"
												/>
												<input
													type="range"
													min="0"
													max={sliderMaxGibPerHour}
													step={sliderStepGibPerHour}
													value={values.preferred}
													oninput={(event) =>
														updateSlider(quality.qualityId, 'preferred', event.currentTarget.value)}
													aria-label={`${quality.name} preferred size slider`}
													class="quality-size-range-input preferred"
												/>
												<input
													type="range"
													min="0"
													max={sliderMaxGibPerHour}
													step={sliderStepGibPerHour}
													value={values.maximum}
													oninput={(event) =>
														updateSlider(quality.qualityId, 'maximum', event.currentTarget.value)}
													aria-label={`${quality.name} maximum size slider`}
													class="quality-size-range-input maximum"
												/>
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
														oninput={(event) =>
															updateGibValue(
																quality.qualityId,
																'preferred',
																event.currentTarget.value
															)}
														onwheel={(event) =>
															scrollGibValue(
																event,
																quality.qualityId,
																'preferred',
																values.preferred
															)}
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
														oninput={(event) =>
															updateGibValue(
																quality.qualityId,
																'maximum',
																event.currentTarget.value
															)}
														onwheel={(event) =>
															scrollGibValue(event, quality.qualityId, 'maximum', values.maximum)}
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
							{/each}
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</form>
</div>
