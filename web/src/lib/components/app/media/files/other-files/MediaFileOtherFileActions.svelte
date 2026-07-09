<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import FileOutputIcon from '@lucide/svelte/icons/file-output';
	import SearchIcon from '@lucide/svelte/icons/search';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import MediaFileFulfillmentActions from '$lib/components/app/media/files/MediaFileFulfillmentActions.svelte';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type {
		MediaFulfillmentActionRequest,
		MediaItemSubtitle,
		MediaItemSubtitleSelectionRequest
	} from '$lib/settings/types';

	type OtherFile = MediaFileRow['otherFiles'][number];

	interface Props {
		file: OtherFile;
		subtitle?: MediaItemSubtitle;
		languageId?: string;
		subtitleMode: 'embedded' | 'external' | 'mixed';
		canManage: boolean;
		pendingFulfillmentActionKeys?: string[];
		canSearch: boolean;
		onSearch: (_languageId?: string) => void | Promise<void>;
		onManualSearch: (_languageId?: string) => void;
		onDelete: () => void | Promise<void>;
		onUpdateSubtitle: (
			_subtitle: MediaItemSubtitle,
			_request: MediaItemSubtitleSelectionRequest
		) => void | Promise<void>;
		onFulfillmentAction: (_request: MediaFulfillmentActionRequest) => void | Promise<void>;
	}

	let {
		file,
		subtitle,
		languageId,
		subtitleMode,
		canManage,
		pendingFulfillmentActionKeys = [],
		canSearch,
		onSearch,
		onManualSearch,
		onDelete,
		onUpdateSubtitle,
		onFulfillmentAction
	}: Props = $props();

	const fulfillmentRow = $derived.by((): MediaFileDetailRow | undefined =>
		file.type === 'subtitle'
			? {
					key: file.path,
					filePath: undefined,
					otherFileId: file.id,
					trackNumber: '-',
					type: 'subtitle',
					language: languageId ?? '-',
					description: file.path,
					missing: file.state?.visualState === 'missing_placeholder',
					unwanted: file.state?.visualState === 'unwanted',
					...file.state
				}
			: undefined
	);

	function updateRetention(retentionMode: 'external' | 'mux') {
		if (!subtitle) return;
		onUpdateSubtitle(subtitle, {
			selected: subtitle.selected,
			retentionMode
		});
	}
</script>

<span class="flex min-h-8 min-w-0 flex-wrap items-center gap-2 lg:justify-end">
	{#if file.type === 'subtitle'}
		{#if fulfillmentRow}
			<MediaFileFulfillmentActions
				row={fulfillmentRow}
				{canManage}
				{pendingFulfillmentActionKeys}
				{onFulfillmentAction}
			/>
		{/if}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="outline"
						size="icon-sm"
						aria-label="Auto search subtitle"
						disabled={!canManage || !canSearch}
						onclick={() => onSearch(languageId)}
					>
						<SearchIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Auto search subtitle</Tooltip.Content>
		</Tooltip.Root>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="outline"
						size="icon-sm"
						aria-label="Manual search subtitle"
						disabled={!canManage || !canSearch}
						onclick={() => onManualSearch(languageId)}
					>
						<UserIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Manual search subtitle</Tooltip.Content>
		</Tooltip.Root>
		{#if subtitle && subtitleMode !== 'external' && subtitle.retentionMode !== 'mux'}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Embed subtitle"
							disabled={!canManage}
							onclick={() => updateRetention('mux')}
						>
							<CaptionsIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>Embed subtitle</Tooltip.Content>
			</Tooltip.Root>
		{/if}
		{#if subtitle && subtitleMode === 'external' && subtitle.retentionMode !== 'external'}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Move subtitle out"
							disabled={!canManage}
							onclick={() => updateRetention('external')}
						>
							<FileOutputIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>Move subtitle out</Tooltip.Content>
			</Tooltip.Root>
		{/if}
	{/if}
	{#if file.status === 'available'}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="destructive"
						size="icon-sm"
						aria-label="Delete other file"
						disabled={!canManage}
						onclick={onDelete}
					>
						<TrashIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Delete other file</Tooltip.Content>
		</Tooltip.Root>
	{/if}
</span>
