<script lang="ts">
	import MovieIcon from '@lucide/svelte/icons/clapperboard';
	import SeriesIcon from '@lucide/svelte/icons/tv';
	import ActivityActions from '../ActivityActions.svelte';
	import ActivityProgressCell from '../ActivityProgressCell.svelte';
	import TruncatedText from './TruncatedText.svelte';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { formatBytes } from '$lib/components/app/media/file-data/mediaFileSize';
	import { Card } from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { activityDisplay, statusTone } from '../activityDisplay';
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activities: DownloadActivity[];
		canManage: boolean;
		cancellingId?: string;
		retryingId?: string;
		deletingId?: string;
		onManualImport: (_activity: DownloadActivity) => void;
		onCancel: (_activity: DownloadActivity) => void;
		onRetry: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
	}

	let {
		activities,
		canManage,
		cancellingId,
		retryingId,
		deletingId,
		onManualImport,
		onCancel,
		onRetry,
		onDelete
	}: Props = $props();
</script>

<Card class="overflow-hidden p-0">
	<Table.Root
		class="min-w-[1120px] table-fixed [&_td]:px-3 [&_td]:whitespace-nowrap [&_th]:px-3 [&_th]:whitespace-nowrap"
	>
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-10"><span class="sr-only">Select</span></Table.Head>
				<Table.Head>Media</Table.Head>
				<Table.Head class="w-18">Quality</Table.Head>
				<Table.Head class="w-24">Formats</Table.Head>
				<Table.Head class="w-20">Protocol</Table.Head>
				<Table.Head class="w-28">Client</Table.Head>
				<Table.Head class="w-18">Time left</Table.Head>
				<Table.Head class="w-56">Progress</Table.Head>
				<Table.Head class="w-24">Status</Table.Head>
				<Table.Head class="w-28 text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each activities as activity (activity.id)}
				{@const display = activityDisplay(activity)}
				<Table.Row>
					<Table.Cell class="w-10">
						<Checkbox aria-label={`Select ${activity.releaseTitle}`} />
					</Table.Cell>
					<Table.Cell class="min-w-0 whitespace-normal">
						<div class="grid min-w-0 gap-1">
							<div class="flex min-w-0 items-center gap-1.5">
								<span class="min-w-0 truncate font-bold">
									{activity.mediaTitle}{display.year
										? ` (${display.year})`
										: ''}{display.seriesEpisode
										? ` - ${display.seriesEpisode}`
										: ''}{display.episodeTitle ? ` - ${display.episodeTitle}` : ''}
								</span>
								<span
									role="img"
									aria-label={`${activity.mediaType === 'serie' ? 'Series' : 'Movie'} media type`}
									class="inline-grid shrink-0 text-muted-foreground"
								>
									{#if activity.mediaType === 'serie'}
										<SeriesIcon aria-hidden="true" class="size-4" />
									{:else}
										<MovieIcon aria-hidden="true" class="size-4" />
									{/if}
								</span>
							</div>
							<TruncatedText text={activity.releaseTitle} class="text-xs text-muted-foreground" />
						</div>
					</Table.Cell>
					<Table.Cell>{display.quality}</Table.Cell>
					<Table.Cell>
						<TruncatedText text={display.formats.length ? display.formats.join(', ') : '-'} />
					</Table.Cell>
					<Table.Cell class="capitalize">{activity.indexerProtocol || '-'}</Table.Cell>
					<Table.Cell><TruncatedText text={activity.downloadClientName} /></Table.Cell>
					<Table.Cell>{display.timeLeft}</Table.Cell>
					<Table.Cell><ActivityProgressCell {activity} /></Table.Cell>
					<Table.Cell>
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<span
										{...props}
										aria-label={`${activity.status} status${activity.error ? `: ${activity.error}` : ''}${activity.clientFiles?.length ? `: ${activity.clientFiles.length} client file(s)` : ''}`}
										class="inline-grid"
									>
										<StatusPill tone={statusTone(activity.status)} class="capitalize">
											{activity.status}
										</StatusPill>
									</span>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content class="max-w-100 break-words">
								{#if activity.error}
									<p>{activity.error}</p>
								{/if}
								{#if activity.clientFiles?.length}
									<div class="grid gap-2">
										{#each activity.clientFiles as file (file.path)}
											<div class="grid gap-0.5">
												<span class="break-all">{file.path}</span>
												<small class="text-muted-foreground">
													{formatBytes(file.sizeBytes)} · {file.complete
														? 'Complete'
														: 'Downloading'}
												</small>
											</div>
										{/each}
									</div>
								{:else if !activity.error}
									<span class="capitalize">{activity.status}</span>
								{/if}
							</Tooltip.Content>
						</Tooltip.Root>
					</Table.Cell>
					<Table.Cell class="min-w-28 text-right">
						<ActivityActions
							{activity}
							{canManage}
							{cancellingId}
							{retryingId}
							{deletingId}
							{onManualImport}
							{onCancel}
							{onRetry}
							{onDelete}
						/>
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
