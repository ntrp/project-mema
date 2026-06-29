<script lang="ts">
	import { onMount } from 'svelte';

	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import { listQualitySizeSettings, updateQualitySizeSettings } from '$lib/settings/api';
	import { groupQualitiesByResolution } from '$lib/settings/qualityGroups';
	import type { QualitySizeSetting, QualitySizeSettingRequest } from '$lib/settings/types';

	type SizeField = 'minimumSizeMbPerMinute' | 'preferredSizeMbPerMinute' | 'maximumSizeMbPerMinute';

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

	function updateSize(qualityId: string, field: SizeField, rawValue: string) {
		qualities = qualities.map((quality) => {
			if (quality.qualityId !== qualityId) {
				return quality;
			}
			return {
				...quality,
				[field]: parseSizeValue(rawValue, field === 'minimumSizeMbPerMinute')
			};
		});
		message = '';
	}

	function parseSizeValue(value: string, required: boolean) {
		const trimmed = value.trim();
		if (trimmed === '') {
			return required ? 0 : null;
		}
		const parsed = Number(trimmed);
		if (!Number.isFinite(parsed)) {
			return required ? 0 : null;
		}
		return Math.max(0, parsed);
	}

	function inputValue(value: number | null | undefined) {
		return value ?? '';
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
						<th>Minimum MB/min</th>
						<th>Preferred MB/min</th>
						<th>Maximum MB/min</th>
						<th>Status</th>
					</tr>
				</thead>
				<tbody>
					{#if loading}
						<tr>
							<td colspan="5" class="empty">Loading quality sizes</td>
						</tr>
					{:else if qualities.length === 0}
						<tr>
							<td colspan="5" class="empty">No qualities loaded</td>
						</tr>
					{:else}
						{#each qualityGroups as group (group.id)}
							<tr class="quality-size-group-row">
								<th colspan="5">{group.label}</th>
							</tr>
							{#each group.qualities as quality (quality.qualityId)}
								{@const validationError = rowError(quality)}
								<tr class:invalid-row={validationError !== ''}>
									<td>
										<strong>{quality.name}</strong>
										<span>{quality.qualityId}</span>
									</td>
									<td>
										<input
											type="number"
											min="0"
											step="0.1"
											required
											aria-label={`${quality.name} minimum size`}
											value={inputValue(quality.minimumSizeMbPerMinute)}
											oninput={(event) =>
												updateSize(
													quality.qualityId,
													'minimumSizeMbPerMinute',
													event.currentTarget.value
												)}
										/>
									</td>
									<td>
										<input
											type="number"
											min="0"
											step="0.1"
											aria-label={`${quality.name} preferred size`}
											value={inputValue(quality.preferredSizeMbPerMinute)}
											oninput={(event) =>
												updateSize(
													quality.qualityId,
													'preferredSizeMbPerMinute',
													event.currentTarget.value
												)}
										/>
									</td>
									<td>
										<input
											type="number"
											min="0"
											step="0.1"
											aria-label={`${quality.name} maximum size`}
											value={inputValue(quality.maximumSizeMbPerMinute)}
											oninput={(event) =>
												updateSize(
													quality.qualityId,
													'maximumSizeMbPerMinute',
													event.currentTarget.value
												)}
										/>
									</td>
									<td class="quality-size-status">
										{#if validationError}
											<span class="status-error">{validationError}</span>
										{:else}
											<span>Ready</span>
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
