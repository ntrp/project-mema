<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { minimumAvailabilityOptions } from '$lib/components/settings/libraryScanImport';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type {
		MediaItem,
		MediaItemUpdateRequest,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		qualityProfiles: QualityProfileOption[];
		refreshing?: boolean;
		savingOptions?: boolean;
		deleting?: boolean;
		onRefreshMetadata: (_item: MediaItem) => void;
		onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
		onDelete: (_item: MediaItem) => void;
	}

	let {
		item,
		qualityProfiles,
		refreshing = false,
		savingOptions = false,
		deleting = false,
		onRefreshMetadata,
		onSaveOptions,
		onDelete
	}: Props = $props();

	const qualityProfileOptions = $derived(
		qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	);

	function saveQualityProfile(value: string) {
		if (!value || value === item.qualityProfileId) {
			return;
		}
		onSaveOptions(item, {
			qualityProfileId: value,
			minimumAvailability: item.minimumAvailability
		});
	}

	function saveMinimumAvailability(value: string) {
		if (value === item.minimumAvailability || !item.qualityProfileId) {
			return;
		}
		onSaveOptions(item, {
			qualityProfileId: item.qualityProfileId,
			minimumAvailability: value as MinimumAvailability
		});
	}
</script>

<div class="ml-auto flex flex-wrap items-end gap-2.5">
	<div class="relative min-w-40 pt-2">
		<Label
			class="absolute top-0.5 left-2 z-10 px-1 text-[11px] leading-none font-extrabold text-muted-foreground"
			>Quality Profile</Label
		>
		<SettingsSelect
			value={item.qualityProfileId ?? ''}
			size="sm"
			options={qualityProfileOptions}
			disabled={savingOptions || qualityProfiles.length === 0}
			onValueChange={saveQualityProfile}
		/>
	</div>
	<div class="relative min-w-36 pt-2">
		<Label
			class="absolute top-0.5 left-2 z-10 px-1 text-[11px] leading-none font-extrabold text-muted-foreground"
			>Min. Availability</Label
		>
		<SettingsSelect
			value={item.minimumAvailability}
			size="sm"
			options={minimumAvailabilityOptions}
			disabled={savingOptions || !item.qualityProfileId}
			onValueChange={saveMinimumAvailability}
		/>
	</div>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon-sm"
					aria-label="Refresh metadata"
					disabled={refreshing}
					onclick={() => onRefreshMetadata(item)}
				>
					<RefreshCwIcon class={refreshing ? 'animate-spin' : undefined} aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Refresh metadata</Tooltip.Content>
	</Tooltip.Root>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="destructive"
					size="icon-sm"
					aria-label="Delete media"
					disabled={deleting}
					onclick={() => onDelete(item)}
				>
					<TrashIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Delete media</Tooltip.Content>
	</Tooltip.Root>
</div>
