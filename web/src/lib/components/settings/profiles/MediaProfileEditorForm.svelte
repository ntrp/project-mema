<script lang="ts">
	import MediaProfileCustomFormatScores from '$lib/components/settings/profiles/MediaProfileCustomFormatScores.svelte';
	import MediaProfileLanguageSelector from '$lib/components/settings/profiles/MediaProfileLanguageSelector.svelte';
	import MediaProfileRules from '$lib/components/settings/profiles/MediaProfileRules.svelte';
	import MediaProfileSubtitleSelector from '$lib/components/settings/profiles/MediaProfileSubtitleSelector.svelte';
	import MediaProfileVideoTarget from '$lib/components/settings/profiles/MediaProfileVideoTarget.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import type {
		CustomFormat,
		Language,
		MediaProfileForm,
		QualitySizeSetting
	} from '$lib/settings/types';

	interface Props {
		title: string;
		form: MediaProfileForm;
		customFormats: CustomFormat[];
		qualities: QualitySizeSetting[];
		languages: Language[];
		loadingQualities: boolean;
		qualityError: string;
		saving: boolean;
		saveError: string;
		submitLabel: string;
		onSubmit: (_event: SubmitEvent) => void | Promise<void>;
	}

	let {
		title,
		form = $bindable(),
		customFormats,
		qualities,
		languages,
		loadingQualities,
		qualityError,
		saving,
		saveError,
		submitLabel,
		onSubmit
	}: Props = $props();

	const canSave = $derived(
		form.name.trim() !== '' && form.qualityIds.length > 0 && form.audioTargets.length > 0 && !saving
	);
</script>

<PageHeading eyebrow="Settings" {title} titleId="settings-profile-title" />

<Card class="p-5">
	<form class="grid gap-4" onsubmit={onSubmit}>
		{#if saveError}
			<p
				class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 text-sm font-bold text-destructive"
			>
				{saveError}
			</p>
		{/if}
		<div class="grid min-w-0 grid-cols-1 gap-3.5">
			<div class="col-span-full min-w-0">
				<MediaProfileRules {form} {qualities} onChange={(value) => (form = value)} />
			</div>
			<div class="col-span-full min-w-0">
				<MediaProfileVideoTarget
					{form}
					{qualities}
					{loadingQualities}
					{qualityError}
					onChange={(value) => (form = value)}
				/>
			</div>
			<div class="col-span-full min-w-0">
				<MediaProfileLanguageSelector {form} {languages} onChange={(value) => (form = value)} />
			</div>
			<div class="col-span-full min-w-0">
				<MediaProfileSubtitleSelector {form} {languages} onChange={(value) => (form = value)} />
			</div>
			<div class="col-span-full min-w-0">
				<MediaProfileCustomFormatScores
					{form}
					{customFormats}
					onChange={(value) => (form = value)}
				/>
			</div>
		</div>

		<div class="flex items-center justify-end gap-3">
			<Button type="button" variant="outline" href="/settings/profiles">Cancel</Button>
			<Button type="submit" disabled={!canSave}>
				{saving ? 'Saving' : submitLabel}
			</Button>
		</div>
	</form>
</Card>
