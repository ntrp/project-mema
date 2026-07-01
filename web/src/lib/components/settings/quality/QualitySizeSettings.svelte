<script lang="ts">
	import { onMount } from 'svelte';

	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import { listQualitySizeSettings, updateQualitySizeSettings } from '$lib/settings/api';
	import { groupQualitiesByResolution } from '$lib/settings/qualityGroups';
	import type { QualitySizeSetting } from '$lib/settings/types';
	import QualitySizeRow from './QualitySizeRow.svelte';
	import { qualityRequest, rowError } from './qualitySize';

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

	function updateQuality(nextQuality: QualitySizeSetting) {
		qualities = qualities.map((quality) =>
			quality.qualityId === nextQuality.qualityId ? nextQuality : quality
		);
		message = '';
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
								<QualitySizeRow {quality} onChange={updateQuality} />
							{/each}
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</form>
</div>
