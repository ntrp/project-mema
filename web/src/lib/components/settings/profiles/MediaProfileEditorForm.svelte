<script lang="ts">
	import MediaProfileCustomFormatScores from '$lib/components/settings/profiles/MediaProfileCustomFormatScores.svelte';
	import MediaProfileComponentTargets from '$lib/components/settings/profiles/component-targets/MediaProfileComponentTargets.svelte';
	import MediaProfileQualitySelector from '$lib/components/settings/profiles/MediaProfileQualitySelector.svelte';
	import MediaProfileRules from '$lib/components/settings/profiles/MediaProfileRules.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
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

	const canSave = $derived(form.name.trim() !== '' && form.qualityIds.length > 0 && !saving);
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
		<div
			class="grid min-w-0 items-start gap-4.5 min-[981px]:grid-cols-[minmax(0,2fr)_minmax(320px,1fr)]"
		>
			<div class="grid min-w-0 gap-3.5">
				<div class="grid gap-1.5">
					<Label>Name</Label>
					<Input bind:value={form.name} type="text" maxlength={200} required />
				</div>
				<Label class="flex items-center gap-2 text-sm">
					<Checkbox bind:checked={form.isDefault} />
					<span>Default profile</span>
				</Label>

				<MediaProfileRules {form} {qualities} {languages} onChange={(value) => (form = value)} />
				<MediaProfileComponentTargets {form} onChange={(value) => (form = value)} />
				<MediaProfileCustomFormatScores
					{form}
					{customFormats}
					onChange={(value) => (form = value)}
				/>
			</div>

			<aside class="grid min-w-0 gap-3.5 min-[981px]:sticky min-[981px]:top-0">
				<MediaProfileQualitySelector
					{form}
					{qualities}
					loading={loadingQualities}
					error={qualityError}
					onChange={(value) => (form = value)}
				/>
			</aside>
		</div>

		<div class="flex items-center justify-end gap-3">
			<Button type="button" variant="outline" href="/settings/profiles">Cancel</Button>
			<Button type="submit" disabled={!canSave}>
				{saving ? 'Saving' : submitLabel}
			</Button>
		</div>
	</form>
</Card>
