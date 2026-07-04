<script lang="ts">
	import { untrack } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';
	import { initialManualImportForm, manualImportRequestFromForm } from './activityManualImportForm';

	interface Props {
		activity: DownloadActivity;
		importing: boolean;
		error?: string;
		onImport: (_request: ManualImportRequest) => void;
		onClose: () => void;
	}

	let { activity, importing, error, onImport, onClose }: Props = $props();
	let initializedActivityId = $state(untrack(() => activity.id));
	let form = $state(untrack(() => initialManualImportForm(activity)));

	$effect(() => {
		if (initializedActivityId === activity.id) return;
		initializedActivityId = activity.id;
		form = initialManualImportForm(activity);
	});

	function submit() {
		onImport(manualImportRequestFromForm(form));
	}
</script>

<div class="mb-4 grid gap-1">
	<p class="m-0 mb-1.5 text-xs font-extrabold text-muted-foreground uppercase">Manual import</p>
	<p class="m-0 text-sm leading-6 text-muted-foreground">{activity.releaseTitle}</p>
</div>

<form class="grid gap-4 md:grid-cols-2" onsubmit={(event) => (event.preventDefault(), submit())}>
	<label class="grid gap-1.5 md:col-span-2">
		<span class="text-sm font-bold text-muted-foreground">Source path</span>
		<Input bind:value={form.sourcePath} required placeholder="/downloads/release/file.mkv" />
	</label>
	<label class="grid gap-1.5 md:col-span-2">
		<span class="text-sm font-bold text-muted-foreground">Target filename override</span>
		<Input bind:value={form.targetFileName} placeholder="Leave empty to build from params" />
	</label>
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">
			{activity.mediaType === 'serie' ? 'Series' : 'Movie'}
		</span>
		<Input bind:value={form.movieTitle} />
	</label>
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Year</span>
		<Input bind:value={form.year} min="0" type="number" />
	</label>
	{#if activity.mediaType === 'serie'}
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Season</span>
			<Input bind:value={form.seasonNumber} min="0" type="number" required />
		</label>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Episode</span>
			<Input bind:value={form.episodeNumber} min="0" type="number" required />
		</label>
		<label class="grid gap-1.5 md:col-span-2">
			<span class="text-sm font-bold text-muted-foreground">Episode title</span>
			<Input bind:value={form.episodeTitle} />
		</label>
	{/if}
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Release group</span>
		<Input bind:value={form.releaseGroup} />
	</label>
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Edition</span>
		<Input bind:value={form.edition} />
	</label>
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Quality</span>
		<Input bind:value={form.quality} />
	</label>
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Languages</span>
		<Input bind:value={form.languagesText} placeholder="English, German" />
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
