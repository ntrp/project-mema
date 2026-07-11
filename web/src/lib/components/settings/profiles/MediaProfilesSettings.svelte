<script lang="ts">
	import { resolve } from '$app/paths';
	import MediaProfileTable from '$lib/components/settings/profiles/MediaProfileTable.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Card } from '$lib/components/ui/card';
	import { createQualitySizeResources } from '$lib/components/settings/quality/resources.svelte';
	import type { MediaProfile } from '$lib/settings/types';

	interface Props {
		profiles: MediaProfile[];
		deletingId?: string;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { profiles, deletingId, onDelete }: Props = $props();

	const qualitySizes = createQualitySizeResources();
	const qualities = $derived(qualitySizes.query.data ?? []);
	const qualityError = $derived(qualitySizes.query.error?.message ?? '');
</script>

<Card class="p-5" aria-label="Profiles">
	<SectionHeading>
		{#snippet actions()}
			<SettingsAddButton label="Add profile" href={resolve('/settings/profiles/new')} />
		{/snippet}
	</SectionHeading>

	{#if qualityError}
		<p
			class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 text-sm font-bold text-destructive"
		>
			{qualityError}
		</p>
	{/if}
	<MediaProfileTable {profiles} {qualities} {deletingId} {onDelete} />
</Card>
