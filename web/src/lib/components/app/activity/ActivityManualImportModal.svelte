<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { activityDisplay, releaseGroupFromTitle } from './activityDisplay';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
		importing: boolean;
		error?: string;
		onImport: (_request: ManualImportRequest) => void;
		onClose: () => void;
	}

	let { activity, importing, error, onImport, onClose }: Props = $props();

	const summary = $derived(activityDisplay(activity));
	let initializedActivityId = $state('');
	let sourcePath = $state('');
	let targetFileName = $state('');
	let movieTitle = $state('');
	let year = $state<number | undefined>();
	let seasonNumber = $state<number | undefined>();
	let episodeNumber = $state<number | undefined>();
	let episodeTitle = $state('');
	let releaseGroup = $state('');
	let edition = $state('');
	let quality = $state('');
	let languagesText = $state('');

	$effect(() => {
		if (initializedActivityId === activity.id) return;
		initializedActivityId = activity.id;
		sourcePath = '';
		targetFileName = '';
		movieTitle = activity.mediaTitle;
		year = activity.mediaYear ?? undefined;
		seasonNumber = activity.mediaType === 'series' ? 1 : undefined;
		episodeNumber = activity.mediaType === 'series' ? 1 : undefined;
		episodeTitle = '';
		releaseGroup = releaseGroupFromTitle(activity.releaseTitle);
		edition = '';
		quality = summary.quality === '-' ? '' : summary.quality;
		languagesText = summary.languages.join(', ');
	});

	function submit() {
		onImport({
			sourcePath,
			targetFileName: optional(targetFileName),
			movieTitle: optional(movieTitle),
			year,
			seasonNumber,
			episodeNumber,
			episodeTitle: optional(episodeTitle),
			releaseGroup: optional(releaseGroup),
			edition: optional(edition),
			quality: optional(quality),
			languages: languagesText
				.split(',')
				.map((value) => value.trim())
				.filter(Boolean)
		});
	}

	function optional(value: string) {
		value = value.trim();
		return value === '' ? undefined : value;
	}
</script>

<SettingsFormModal
	title={activity.mediaTitle}
	modalClass="w-[min(1840px,calc(100vw-48px))]"
	{onClose}
>
	<div class="mb-4 grid gap-1">
		<p class="m-0 mb-1.5 text-xs font-extrabold text-muted-foreground uppercase">Manual import</p>
		<p class="m-0 text-sm leading-6 text-muted-foreground">{activity.releaseTitle}</p>
	</div>

	<form class="grid gap-4 md:grid-cols-2" onsubmit={(event) => (event.preventDefault(), submit())}>
		<label class="grid gap-1.5 md:col-span-2">
			<span class="text-sm font-bold text-muted-foreground">Source path</span>
			<Input bind:value={sourcePath} required placeholder="/downloads/release/file.mkv" />
		</label>
		<label class="grid gap-1.5 md:col-span-2">
			<span class="text-sm font-bold text-muted-foreground">Target filename override</span>
			<Input bind:value={targetFileName} placeholder="Leave empty to build from params" />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">
				{activity.mediaType === 'series' ? 'Series' : 'Movie'}
			</span>
			<Input bind:value={movieTitle} />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Year</span>
			<Input bind:value={year} min="0" type="number" />
		</label>
		{#if activity.mediaType === 'series'}
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Season</span>
				<Input bind:value={seasonNumber} min="0" type="number" required />
			</label>
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Episode</span>
				<Input bind:value={episodeNumber} min="0" type="number" required />
			</label>
			<label class="grid gap-1.5 md:col-span-2">
				<span class="text-sm font-bold text-muted-foreground">Episode title</span>
				<Input bind:value={episodeTitle} />
			</label>
		{/if}
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Release group</span>
			<Input bind:value={releaseGroup} />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Edition</span>
			<Input bind:value={edition} />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Quality</span>
			<Input bind:value={quality} />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Languages</span>
			<Input bind:value={languagesText} placeholder="English, German" />
		</label>
		{#if error}
			<p
				class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 font-bold text-destructive md:col-span-2"
			>
				{error}
			</p>
		{/if}
		<div class="flex items-center gap-3 md:col-span-2">
			<Button type="button" variant="outline" onclick={onClose}>Cancel</Button>
			<Button type="submit" disabled={importing}>{importing ? 'Importing' : 'Import'}</Button>
		</div>
	</form>
</SettingsFormModal>
