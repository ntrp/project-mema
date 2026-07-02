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
	import { getFileNamingSettings, updateFileNamingSettings } from '$lib/settings/api';
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
	let loading = $state(true);
	let saving = $state(false);
	let message = $state('');
	let errorMessage = $state('');

	const hasValidationErrors = $derived(
		Object.values(templates).some((value) => value.trim().length === 0)
	);

	onMount(() => {
		void loadSettings();
	});

	async function loadSettings() {
		loading = true;
		message = '';
		errorMessage = '';
		try {
			const settings = await getFileNamingSettings();
			templates = {
				movieFileFormat: settings.movieFileFormat,
				movieFolderFormat: settings.movieFolderFormat,
				seriesEpisodeFormat: settings.seriesEpisodeFormat,
				dailyEpisodeFormat: settings.dailyEpisodeFormat,
				animeEpisodeFormat: settings.animeEpisodeFormat,
				seriesFolderFormat: settings.seriesFolderFormat,
				seasonFolderFormat: settings.seasonFolderFormat,
				specialsFolderFormat: settings.specialsFolderFormat
			};
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load file naming settings';
		} finally {
			loading = false;
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

		saving = true;
		try {
			const settings = await updateFileNamingSettings(trimmedTemplates());
			templates = {
				movieFileFormat: settings.movieFileFormat,
				movieFolderFormat: settings.movieFolderFormat,
				seriesEpisodeFormat: settings.seriesEpisodeFormat,
				dailyEpisodeFormat: settings.dailyEpisodeFormat,
				animeEpisodeFormat: settings.animeEpisodeFormat,
				seriesFolderFormat: settings.seriesFolderFormat,
				seasonFolderFormat: settings.seasonFolderFormat,
				specialsFolderFormat: settings.specialsFolderFormat
			};
			message = 'File naming settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save file naming settings';
		} finally {
			saving = false;
		}
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
