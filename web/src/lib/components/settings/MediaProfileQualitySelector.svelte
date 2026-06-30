<script lang="ts">
	import {
		groupQualitiesByResolution,
		type QualityResolutionGroup
	} from '$lib/settings/qualityGroups';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		qualities: QualitySizeSetting[];
		loading: boolean;
		error: string;
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, qualities, loading, error, onChange }: Props = $props();
	let qualityGroups = $derived(groupQualitiesByResolution(qualities));

	function updateQualityIds(qualityIds: string[]) {
		const upgradeUntilQualityId =
			form.upgradeUntilQualityId && qualityIds.includes(form.upgradeUntilQualityId)
				? form.upgradeUntilQualityId
				: undefined;
		onChange({ ...form, qualityIds, upgradeUntilQualityId });
	}

	function toggleQuality(qualityId: string) {
		if (form.qualityIds.includes(qualityId)) {
			updateQualityIds(form.qualityIds.filter((id) => id !== qualityId));
			return;
		}
		updateQualityIds([...form.qualityIds, qualityId]);
	}

	function selectAllQualities() {
		updateQualityIds(qualities.map((quality) => quality.qualityId));
	}

	function clearQualities() {
		updateQualityIds([]);
	}

	function selectQualityGroup(group: QualityResolutionGroup<QualitySizeSetting>) {
		updateQualityIds([
			...new Set([...form.qualityIds, ...group.qualities.map((quality) => quality.qualityId)])
		]);
	}

	function clearQualityGroup(group: QualityResolutionGroup<QualitySizeSetting>) {
		const groupIDs = new Set(group.qualities.map((quality) => quality.qualityId));
		updateQualityIds(form.qualityIds.filter((qualityId) => !groupIDs.has(qualityId)));
	}

	function selectedQualityCount(group: QualityResolutionGroup<QualitySizeSetting>) {
		return group.qualities.filter((quality) => form.qualityIds.includes(quality.qualityId)).length;
	}
</script>

<div class="wide profile-quality-header">
	<span>Qualities</span>
	<div>
		<button type="button" class="secondary" onclick={selectAllQualities}>Select all</button>
		<button type="button" class="secondary" onclick={clearQualities}>Clear</button>
	</div>
</div>

{#if error}
	<p class="form-status error wide">{error}</p>
{/if}

<div class="quality-group-stack wide" aria-label="Profile qualities">
	{#if loading}
		<p class="muted">Loading qualities</p>
	{:else}
		{#each qualityGroups as group (group.id)}
			<section class="quality-checkbox-group" aria-labelledby={`quality-group-${group.id}`}>
				<div class="quality-group-heading">
					<div>
						<h3 id={`quality-group-${group.id}`}>{group.label}</h3>
						<span>{selectedQualityCount(group)} / {group.qualities.length}</span>
					</div>
					<div>
						<button type="button" class="secondary" onclick={() => selectQualityGroup(group)}>
							Select
						</button>
						<button type="button" class="secondary" onclick={() => clearQualityGroup(group)}>
							Clear
						</button>
					</div>
				</div>
				<div class="quality-checkbox-grid">
					{#each group.qualities as quality (quality.qualityId)}
						<label class="quality-checkbox">
							<input
								type="checkbox"
								checked={form.qualityIds.includes(quality.qualityId)}
								onchange={() => toggleQuality(quality.qualityId)}
							/>
							<span>{quality.name}</span>
						</label>
					{/each}
				</div>
			</section>
		{/each}
	{/if}
</div>
