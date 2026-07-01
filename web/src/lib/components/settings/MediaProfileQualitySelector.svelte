<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
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

<div class="flex items-center justify-between gap-3">
	<span class="text-sm font-bold text-muted-foreground">Qualities</span>
	<div class="flex gap-2">
		<Button type="button" variant="outline" size="sm" onclick={selectAllQualities}
			>Select all</Button
		>
		<Button type="button" variant="outline" size="sm" onclick={clearQualities}>Clear</Button>
	</div>
</div>

{#if error}
	<p
		class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 font-bold text-destructive"
	>
		{error}
	</p>
{/if}

<div
	class="grid max-h-[min(640px,calc(100vh-260px))] gap-3 overflow-auto rounded-md border border-border bg-background p-2.5"
	aria-label="Profile qualities"
>
	{#if loading}
		<p class="m-0 text-sm leading-6 text-muted-foreground">Loading qualities</p>
	{:else}
		{#each qualityGroups as group (group.id)}
			<section
				class="grid gap-2 rounded-md border border-border bg-card p-2.5"
				aria-labelledby={`quality-group-${group.id}`}
			>
				<div class="flex items-center justify-between gap-3">
					<div class="flex min-w-0 items-baseline gap-2">
						<h3 id={`quality-group-${group.id}`} class="m-0 text-base text-foreground">
							{group.label}
						</h3>
						<span class="text-xs font-bold text-muted-foreground"
							>{selectedQualityCount(group)} / {group.qualities.length}</span
						>
					</div>
					<div class="flex gap-1.5">
						<Button
							type="button"
							variant="outline"
							size="sm"
							onclick={() => selectQualityGroup(group)}
						>
							Select
						</Button>
						<Button
							type="button"
							variant="outline"
							size="sm"
							onclick={() => clearQualityGroup(group)}
						>
							Clear
						</Button>
					</div>
				</div>
				<div class="grid gap-2 [grid-template-columns:repeat(auto-fill,minmax(170px,1fr))]">
					{#each group.qualities as quality (quality.qualityId)}
						<label
							class="grid grid-cols-[18px_minmax(0,1fr)] items-center gap-2 rounded-md border border-border bg-muted p-2"
						>
							<Checkbox
								checked={form.qualityIds.includes(quality.qualityId)}
								onclick={() => toggleQuality(quality.qualityId)}
							/>
							<span class="truncate">{quality.name}</span>
						</label>
					{/each}
				</div>
			</section>
		{/each}
	{/if}
</div>
