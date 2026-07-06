<script lang="ts">
	import { resolve } from '$app/paths';
	import { Badge } from '$lib/components/ui/badge';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { languageLabel } from '$lib/settings/languageOptions';
	import { formatDate } from '$lib/settings/dateFormat';
	import SettingsRowActionButton from '../shared/SettingsRowActionButton.svelte';
	import type { MediaProfile, QualitySizeSetting } from '$lib/settings/types';

	interface Props {
		profiles: MediaProfile[];
		qualities: QualitySizeSetting[];
		deletingId?: string;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { profiles, qualities, deletingId, onDelete }: Props = $props();
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
		return (profile.audioTargets ?? [])
			.map((target) => languageLabel(target.languageId))
			.join(', ');
	}

	function scoreSummary(profile: MediaProfile) {
		return `${profile.minimumCustomFormatScore} min / ${profile.upgradeUntilCustomFormatScore} cutoff / ${profile.customFormatScores?.length ?? 0} scored`;
	}
</script>

<Card class="p-0">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Qualities</Table.Head>
				<Table.Head>Upgrade until</Table.Head>
				<Table.Head>Languages</Table.Head>
				<Table.Head>Score</Table.Head>
				<Table.Head>Updated</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each profiles as profile (profile.id)}
				<Table.Row>
					<Table.Cell>
						<div class="flex items-center gap-2">
							<strong>{profile.name}</strong>
							{#if profile.isDefault}
								<Badge>Default</Badge>
							{/if}
						</div>
					</Table.Cell>
					<Table.Cell>
						<div class="flex flex-wrap gap-1.5" aria-label={`${profile.name} qualities`}>
							{#each selectedQualityNames(profile).slice(0, 6) as name (name)}
								<Badge variant="secondary">{name}</Badge>
							{/each}
							{#if (profile.qualityIds?.length ?? 0) > 6}
								<Badge variant="outline">+{(profile.qualityIds?.length ?? 0) - 6}</Badge>
							{/if}
						</div>
					</Table.Cell>
					<Table.Cell>{upgradeUntilName(profile)}</Table.Cell>
					<Table.Cell>{languageSummary(profile)}</Table.Cell>
					<Table.Cell>{scoreSummary(profile)}</Table.Cell>
					<Table.Cell>{formatDate(profile.updatedAt)}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<SettingsRowActionButton
								label={`Edit ${profile.name}`}
								icon="edit"
								href={resolve('/settings/profiles/[id]', { id: profile.id })}
							/>
							<SettingsRowActionButton
								label={`${deletingId === profile.id ? 'Deleting' : 'Delete'} ${profile.name}`}
								icon="delete"
								variant="destructive"
								disabled={deletingId === profile.id}
								confirmTitle="Delete profile"
								confirmDescription={`Delete media profile "${profile.name}"?`}
								confirmLabel="Delete profile"
								onclick={() => onDelete(profile.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={7} class="py-8 text-center text-muted-foreground">
						No profiles configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
