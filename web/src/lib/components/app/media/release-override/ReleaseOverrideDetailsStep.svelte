<script lang="ts">
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import { untrack } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type {
		Language,
		MediaItem,
		ReleaseCandidate,
		ReleaseOverrideDetails
	} from '$lib/settings/types';
	import {
		detailsFromOverrideDraft,
		overrideDraftFromRelease,
		type ReleaseOverrideDraft
	} from '$lib/components/app/media/release-override/releaseOverrideDetails';
	import ReleaseOverrideLanguageSelect from '$lib/components/app/media/release-override/ReleaseOverrideLanguageSelect.svelte';
	import ReleaseOverrideMovieField from '$lib/components/app/media/release-override/ReleaseOverrideMovieField.svelte';
	import ReleaseOverrideSeriesFields from '$lib/components/app/media/release-override/ReleaseOverrideSeriesFields.svelte';
	import ReleaseOverrideIcon from '$lib/components/app/media/release-override/ReleaseOverrideIcon.svelte';

	interface Props {
		item: MediaItem;
		release: ReleaseCandidate;
		languages: Language[];
		qualityOptions: string[];
		grabbing?: boolean;
		onBack: () => void;
		onConfirm: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch: boolean,
			_details: ReleaseOverrideDetails
		) => void;
	}

	let {
		item,
		release,
		languages,
		qualityOptions,
		grabbing = false,
		onBack,
		onConfirm
	}: Props = $props();
	let draft = $state<ReleaseOverrideDraft>(
		untrack(() => overrideDraftFromRelease(item, release, languages))
	);

	const isSeries = $derived(item.type === 'series');
	const qualities = $derived([...new Set([draft.quality, ...qualityOptions].filter(Boolean))]);

	function submit(event: SubmitEvent) {
		event.preventDefault();
		onConfirm(item, release, true, detailsFromOverrideDraft(draft, languages));
	}
</script>

<form class="grid gap-5" onsubmit={submit}>
	<div class="grid gap-3">
		<div class="grid gap-1.5">
			<Label for="override-title">Title</Label>
			<div
				id="override-title"
				role="textbox"
				aria-readonly="true"
				class="min-h-9 rounded-md border border-border bg-muted px-2.5 py-2 text-sm leading-5 break-all text-muted-foreground"
			>
				{release.title}
			</div>
		</div>
		{#if isSeries}
			<ReleaseOverrideSeriesFields {item} {draft} />
		{:else}
			<ReleaseOverrideMovieField bind:value={draft.movieTitle} />
		{/if}
		<div class="grid gap-3 md:grid-cols-2">
			<div class="grid gap-1.5">
				<Label for="override-quality">Quality</Label>
				<Select.Root type="single" bind:value={draft.quality}>
					<Select.Trigger id="override-quality" class="w-full"
						>{draft.quality || 'Quality'}</Select.Trigger
					>
					<Select.Content>
						{#each qualities as quality (quality)}
							<Select.Item value={quality} label={quality} />
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
			<ReleaseOverrideLanguageSelect bind:values={draft.languages} {languages} />
		</div>
		{#if isSeries}
			<div class="grid gap-1.5">
				<Label for="override-release-group">Release group</Label>
				<Input id="override-release-group" bind:value={draft.releaseGroup} />
			</div>
		{/if}
	</div>
	<div class="flex items-center justify-between gap-2">
		<Button type="button" variant="outline" onclick={onBack}>
			<ArrowLeftIcon aria-hidden="true" />
			<span>Back</span>
		</Button>
		<Button
			type="submit"
			class="bg-amber-400 text-amber-950 hover:bg-amber-500"
			disabled={grabbing}
		>
			<ReleaseOverrideIcon />
			<span>{grabbing ? 'Queueing' : 'Grab with override'}</span>
		</Button>
	</div>
</form>
