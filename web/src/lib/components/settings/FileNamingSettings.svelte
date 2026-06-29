<script lang="ts">
	import { onMount } from 'svelte';

	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import { getFileNamingSettings, updateFileNamingSettings } from '$lib/settings/api';
	import type { FileNamingSettingsRequest } from '$lib/settings/types';

	type TemplateField = keyof FileNamingSettingsRequest;

	const defaultTemplates: FileNamingSettingsRequest = {
		movieFileFormat: '{Movie Title} ({Release Year}) {Quality Full}',
		movieFolderFormat: '{Movie Title} ({Release Year})',
		seriesEpisodeFormat:
			'{Series Title} - S{season:00}E{episode:00} - {Episode Title} {Quality Full}',
		dailyEpisodeFormat: '{Series Title} - {Air-Date} - {Episode Title} {Quality Full}',
		animeEpisodeFormat:
			'{Series Title} - S{season:00}E{episode:00} - {Episode Title} {Quality Full}',
		seriesFolderFormat: '{Series Title}',
		seasonFolderFormat: 'Season {season}',
		specialsFolderFormat: 'Specials'
	};

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
				{ key: 'dailyEpisodeFormat', label: 'Daily episode' }
			]
		},
		{
			id: 'anime-naming',
			title: 'Anime',
			fields: [{ key: 'animeEpisodeFormat', label: 'Episode' }]
		},
		{
			id: 'season-naming',
			title: 'Season',
			fields: [
				{ key: 'seasonFolderFormat', label: 'Season folder' },
				{ key: 'specialsFolderFormat', label: 'Specials folder' }
			]
		}
	];

	let templates = $state<FileNamingSettingsRequest>({ ...defaultTemplates });
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
		templates = { ...defaultTemplates };
		message = '';
		errorMessage = '';
	}

	function trimmedTemplates(): FileNamingSettingsRequest {
		return Object.fromEntries(
			Object.entries(templates).map(([key, value]) => [key, value.trim()])
		) as FileNamingSettingsRequest;
	}
</script>

<div class="panel file-naming-panel" aria-labelledby="file-naming-title">
	<form onsubmit={saveSettings}>
		<div class="section-heading">
			<div>
				<p class="section-kicker">Naming</p>
				<h2 id="file-naming-title">File naming</h2>
			</div>
			<div class="file-naming-actions">
				<button type="button" class="secondary" disabled={loading || saving} onclick={loadSettings}>
					Reload
				</button>
				<button
					type="button"
					class="secondary"
					disabled={loading || saving}
					onclick={resetDefaults}
				>
					Defaults
				</button>
				<button type="submit" disabled={loading || saving || hasValidationErrors}>
					{saving ? 'Saving' : 'Save templates'}
				</button>
			</div>
		</div>

		<NoticeStack {message} {errorMessage} />

		<div class="file-naming-grid">
			{#each templateSections as section (section.id)}
				<section class="file-naming-group" aria-labelledby={section.id}>
					<h3 id={section.id}>{section.title}</h3>
					{#each section.fields as field (field.key)}
						<label>
							<span>{field.label}</span>
							<textarea
								rows="2"
								required
								value={templates[field.key]}
								oninput={(event) => updateTemplate(field.key, event.currentTarget.value)}
							></textarea>
						</label>
					{/each}
				</section>
			{/each}
		</div>
	</form>
</div>
