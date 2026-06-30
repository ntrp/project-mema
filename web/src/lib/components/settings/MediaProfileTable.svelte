<script lang="ts">
	import { languageLabel } from '$lib/settings/languageOptions';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { MediaProfile, QualitySizeSetting } from '$lib/settings/types';

	interface Props {
		profiles: MediaProfile[];
		qualities: QualitySizeSetting[];
		deletingId?: string;
		onEdit: (_profile: MediaProfile) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { profiles, qualities, deletingId, onEdit, onDelete }: Props = $props();
	let qualityNames = $derived(
		new Map(qualities.map((quality) => [quality.qualityId, quality.name]))
	);

	function selectedQualityNames(profile: MediaProfile) {
		return (profile.qualityIds ?? []).map((id) => qualityNames.get(id) ?? id).filter(Boolean);
	}

	function upgradeUntilName(profile: MediaProfile) {
		if (!profile.upgradesAllowed) {
			return 'Disabled';
		}
		return profile.upgradeUntilQualityId
			? (qualityNames.get(profile.upgradeUntilQualityId) ?? profile.upgradeUntilQualityId)
			: 'No cutoff';
	}

	function languageSummary(profile: MediaProfile) {
		return (profile.targetLanguages ?? []).map(languageLabel).join(', ') || 'Any';
	}

	function scoreSummary(profile: MediaProfile) {
		return `${profile.minimumCustomFormatScore} min / ${profile.upgradeUntilCustomFormatScore} cutoff / ${profile.customFormatScores?.length ?? 0} scored`;
	}
</script>

<div class="table-wrap">
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Qualities</th>
				<th>Upgrade until</th>
				<th>Languages</th>
				<th>Score</th>
				<th>Updated</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each profiles as profile (profile.id)}
				<tr>
					<td><strong>{profile.name}</strong></td>
					<td>
						<div class="quality-chip-list" aria-label={`${profile.name} qualities`}>
							{#each selectedQualityNames(profile).slice(0, 6) as name (name)}
								<span>{name}</span>
							{/each}
							{#if (profile.qualityIds?.length ?? 0) > 6}
								<span>+{(profile.qualityIds?.length ?? 0) - 6}</span>
							{/if}
						</div>
					</td>
					<td>{upgradeUntilName(profile)}</td>
					<td>{languageSummary(profile)}</td>
					<td>{scoreSummary(profile)}</td>
					<td>{formatDate(profile.updatedAt)}</td>
					<td class="row-actions">
						<button
							type="button"
							class="secondary icon-button"
							aria-label={`Edit ${profile.name}`}
							onclick={() => onEdit(profile)}
						>
							<span class="app-icon" aria-hidden="true">edit</span>
						</button>
						<button
							type="button"
							class="danger icon-button"
							disabled={deletingId === profile.id}
							aria-label={`${deletingId === profile.id ? 'Deleting' : 'Delete'} ${profile.name}`}
							onclick={() => onDelete(profile.id)}
						>
							<span class="app-icon" aria-hidden="true">delete</span>
						</button>
					</td>
				</tr>
			{:else}
				<tr>
					<td colspan="7" class="empty">No profiles configured</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
