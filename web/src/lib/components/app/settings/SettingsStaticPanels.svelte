<script lang="ts">
	import CustomFormatsSettings from '$lib/components/settings/CustomFormatsSettings.svelte';
	import FileNamingSettings from '$lib/components/settings/FileNamingSettings.svelte';
	import MetadataCacheSettings from '$lib/components/settings/MetadataCacheSettings.svelte';
	import MetadataProviderSettings from '$lib/components/settings/MetadataProviderSettings.svelte';
	import MediaProfilesSettings from '$lib/components/settings/MediaProfilesSettings.svelte';
	import QualitySizeSettings from '$lib/components/settings/quality/QualitySizeSettings.svelte';
	import TagSettings from '$lib/components/settings/TagSettings.svelte';
	import type {
		CustomFormat,
		CustomFormatForm as CustomFormatFormValue,
		IntegrationTestResults,
		MediaProfile,
		MediaProfileForm as MediaProfileFormValue,
		MetadataCacheResponse,
		MetadataProvider,
		MetadataProviderForm as MetadataProviderFormValue,
		SettingsSection,
		Tag,
		TagForm
	} from '$lib/settings/types';

	interface Props {
		activeSection: SettingsSection;
		metadataProviders: MetadataProvider[];
		metadataCache: MetadataCacheResponse;
		mediaProfiles: MediaProfile[];
		customFormats: CustomFormat[];
		tags: Tag[];
		metadataCachePattern: string;
		mediaProfileForm: MediaProfileFormValue;
		customFormatForm: CustomFormatFormValue;
		tagForm: TagForm;
		savingMetadataProviderId?: string;
		testingMetadataProviderId?: string;
		loadingMetadataCache: boolean;
		clearingMetadataCache: boolean;
		savingMediaProfile: boolean;
		deletingMediaProfileId?: string;
		savingCustomFormat: boolean;
		deletingCustomFormatId?: string;
		savingTag: boolean;
		deletingTagId?: string;
		metadataProviderTests: IntegrationTestResults;
		onSaveMetadataProvider: (_form: MetadataProviderFormValue) => void | Promise<void>;
		onTestMetadataProvider: (_id: string) => void | Promise<void>;
		onRefreshMetadataCache: () => void | Promise<void>;
		onClearMetadataCache: () => void | Promise<void>;
		onClearMetadataCachePattern: (_event: SubmitEvent) => void | Promise<void>;
		onSaveMediaProfile: (_event: SubmitEvent) => void | Promise<void>;
		onCancelMediaProfile: () => void;
		onEditMediaProfile: (_profile: MediaProfile) => void;
		onDeleteMediaProfile: (_id: string) => void | Promise<void>;
		onSaveCustomFormat: (_event: SubmitEvent) => void | Promise<void>;
		onImportCustomFormat: (_format: CustomFormatFormValue) => void | Promise<void>;
		onCancelCustomFormat: () => void;
		onEditCustomFormat: (_format: CustomFormat) => void;
		onDeleteCustomFormat: (_id: string) => void | Promise<void>;
		onSaveTag: (_event: SubmitEvent) => void | Promise<void>;
		onCancelTag: () => void;
		onEditTag: (_tag: Tag) => void;
		onDeleteTag: (_id: string) => void | Promise<void>;
	}

	let {
		activeSection,
		metadataProviders,
		metadataCache,
		mediaProfiles,
		customFormats,
		tags,
		metadataCachePattern = $bindable(),
		mediaProfileForm = $bindable(),
		customFormatForm = $bindable(),
		tagForm = $bindable(),
		savingMetadataProviderId,
		testingMetadataProviderId,
		loadingMetadataCache,
		clearingMetadataCache,
		savingMediaProfile,
		deletingMediaProfileId,
		savingCustomFormat,
		deletingCustomFormatId,
		savingTag,
		deletingTagId,
		metadataProviderTests,
		onSaveMetadataProvider,
		onTestMetadataProvider,
		onRefreshMetadataCache,
		onClearMetadataCache,
		onClearMetadataCachePattern,
		onSaveMediaProfile,
		onCancelMediaProfile,
		onEditMediaProfile,
		onDeleteMediaProfile,
		onSaveCustomFormat,
		onImportCustomFormat,
		onCancelCustomFormat,
		onEditCustomFormat,
		onDeleteCustomFormat,
		onSaveTag,
		onCancelTag,
		onEditTag,
		onDeleteTag
	}: Props = $props();
</script>

{#if activeSection === 'metadata'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">Metadata</h1>
	</div>
	<div class="settings-stack">
		<MetadataProviderSettings
			{metadataProviders}
			onSave={onSaveMetadataProvider}
			onTest={onTestMetadataProvider}
			testingId={testingMetadataProviderId}
			savingId={savingMetadataProviderId}
			testResults={metadataProviderTests}
		/>
		<MetadataCacheSettings
			cache={metadataCache}
			bind:pattern={metadataCachePattern}
			loading={loadingMetadataCache}
			clearing={clearingMetadataCache}
			onRefresh={onRefreshMetadataCache}
			onClearAll={onClearMetadataCache}
			onClearPattern={onClearMetadataCachePattern}
		/>
	</div>
{:else if activeSection === 'quality'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">Quality</h1>
	</div>
	<div class="settings-stack"><QualitySizeSettings /></div>
{:else if activeSection === 'profiles'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">Profiles</h1>
	</div>
	<div class="settings-stack">
		<MediaProfilesSettings
			profiles={mediaProfiles}
			{customFormats}
			bind:form={mediaProfileForm}
			saving={savingMediaProfile}
			deletingId={deletingMediaProfileId}
			onSave={onSaveMediaProfile}
			onCancel={onCancelMediaProfile}
			onEdit={onEditMediaProfile}
			onDelete={onDeleteMediaProfile}
		/>
	</div>
{:else if activeSection === 'custom-formats'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">Custom formats</h1>
	</div>
	<div class="settings-stack">
		<CustomFormatsSettings
			formats={customFormats}
			bind:form={customFormatForm}
			saving={savingCustomFormat}
			deletingId={deletingCustomFormatId}
			onSave={onSaveCustomFormat}
			onImport={onImportCustomFormat}
			onCancel={onCancelCustomFormat}
			onEdit={onEditCustomFormat}
			onDelete={onDeleteCustomFormat}
		/>
	</div>
{:else if activeSection === 'file-naming'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">File naming</h1>
	</div>
	<div class="settings-stack"><FileNamingSettings /></div>
{:else if activeSection === 'tags'}
	<div class="page-heading">
		<p>Settings</p>
		<h1 id="settings-title">Tags</h1>
	</div>
	<div class="settings-stack">
		<TagSettings
			{tags}
			bind:form={tagForm}
			saving={savingTag}
			deletingId={deletingTagId}
			onSave={onSaveTag}
			onCancel={onCancelTag}
			onEdit={onEditTag}
			onDelete={onDeleteTag}
		/>
	</div>
{/if}
