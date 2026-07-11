<script lang="ts">
	import { onMount } from 'svelte';
	import FileNamingTemplateSection from '$lib/components/settings/library/FileNamingTemplateSection.svelte';
	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import {
		defaultFileNamingTemplates,
		fileNamingTemplateExample
	} from '$lib/settings/fileNamingTemplates';
	import { createFileNamingResource } from '$lib/features/settings/resources/filePolicies.svelte';
	import type { FileNamingSettingsRequest } from '$lib/settings/types';

	type TemplateField = keyof FileNamingSettingsRequest;

	const templateSections: {
		id: string;
		title: string;
		fields: { key: TemplateField; label: string }[];
	}[] = [
		{
			id: 'movie-naming',
			title: 'Movie',
			fields: [
				{ key: 'movieFolderFormat', label: 'Main folder' },
				{ key: 'movieFileFormat', label: 'File' }
			]
		},
		{
			id: 'series-naming',
			title: 'Series',
			fields: [
				{ key: 'seriesFolderFormat', label: 'Main folder' },
				{ key: 'seriesEpisodeFormat', label: 'Episode' },
				{ key: 'dailyEpisodeFormat', label: 'Daily episode' },
				{ key: 'animeEpisodeFormat', label: 'Anime episode' },
				{ key: 'seasonFolderFormat', label: 'Season folder' },
				{ key: 'specialsFolderFormat', label: 'Specials folder' }
			]
		}
	];

	let templates = $state<FileNamingSettingsRequest>({ ...defaultFileNamingTemplates });
	const resource = createFileNamingResource();
	let message = $state('');
	let errorMessage = $state('');

	const hasValidationErrors = $derived(
		Object.values(templates).some((value) => value.trim().length === 0)
	);

	const loading = $derived(resource.query.isPending || resource.query.isFetching);
	const saving = $derived(resource.save.isPending);
	onMount(() => {
		void loadSettings();
	});

	async function loadSettings() {
		message = '';
		errorMessage = '';
		try {
			const result = await resource.query.refetch();
			if (result.data) hydrate(result.data);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load file naming settings';
		} finally {
			/* query owns loading state */
		}
	}

	async function saveSettings(event: SubmitEvent) {
		event.preventDefault();
		message = '';
		errorMessage = '';
		if (hasValidationErrors) {
			errorMessage = 'All templates are required';
			return;
		}

		try {
			hydrate(await resource.save.mutateAsync(trimmedTemplates()));
			message = 'File naming settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save file naming settings';
		} finally {
			/* mutation owns saving state */
		}
	}

	function hydrate(settings: FileNamingSettingsRequest) {
		templates = { ...settings };
	}

	function updateTemplate(key: TemplateField, value: string) {
		templates = { ...templates, [key]: value };
		message = '';
	}

	function resetDefaults() {
		templates = { ...defaultFileNamingTemplates };
		message = '';
		errorMessage = '';
	}

	function trimmedTemplates(): FileNamingSettingsRequest {
		return Object.fromEntries(
			Object.entries(templates).map(([key, value]) => [key, value.trim()])
		) as FileNamingSettingsRequest;
	}

	function example(value: string) {
		return fileNamingTemplateExample(value);
	}
</script>

<Card class="overflow-visible p-5" aria-label="File naming">
	<form class="grid gap-4" onsubmit={saveSettings}>
		<SectionHeading title="File Naming">
			{#snippet actions()}
				<div class="flex flex-wrap justify-end gap-2.5">
					<Button
						type="button"
						variant="outline"
						disabled={loading || saving}
						onclick={loadSettings}
					>
						Reload
					</Button>
					<Button
						type="button"
						variant="outline"
						disabled={loading || saving}
						onclick={resetDefaults}
					>
						Defaults
					</Button>
					<Button type="submit" disabled={loading || saving || hasValidationErrors}>
						{saving ? 'Saving' : 'Save templates'}
					</Button>
				</div>
			{/snippet}
		</SectionHeading>

		<NoticeStack {message} {errorMessage} />

		<div class="grid gap-3.5">
			{#each templateSections as section (section.id)}
				<FileNamingTemplateSection
					id={section.id}
					title={section.title}
					fields={section.fields}
					{templates}
					onChange={updateTemplate}
					{example}
				/>
			{/each}
		</div>
	</form>
</Card>
