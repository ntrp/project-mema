<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { minimumAvailabilityOptions } from '$lib/components/settings/library/scan/libraryScanImport';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { providerDisplayName, providerPageUrl } from '$lib/settings/providerLinks';
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
	const metadataProviderLabel = $derived(providerDisplayName(item.externalProvider));
	const metadataProviderUrl = $derived(
		providerPageUrl(item.externalProvider, item.type, item.externalId, item.externalUrl)
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
	<div class="flex items-center">
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						variant="outline"
						size="sm"
						href={metadataProviderUrl}
						target="_blank"
						rel="noreferrer"
						disabled={!metadataProviderUrl}
						class="rounded-tr-none rounded-br-none border-r-0"
						aria-label={`Open ${metadataProviderLabel} page in a new tab`}
					>
						<span>{metadataProviderLabel}</span>
						<ExternalLinkIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Open {metadataProviderLabel} page</Tooltip.Content>
		</Tooltip.Root>
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
						class="rounded-tl-none rounded-bl-none"
						onclick={() => onRefreshMetadata(item)}
					>
						<RefreshCwIcon class={refreshing ? 'animate-spin' : undefined} aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Refresh metadata from {metadataProviderLabel}</Tooltip.Content>
		</Tooltip.Root>
	</div>
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
