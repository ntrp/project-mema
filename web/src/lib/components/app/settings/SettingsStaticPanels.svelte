<script lang="ts">
	import CustomFormatsSettings from '$lib/components/settings/CustomFormatsSettings.svelte';
	import MetadataProviderSettings from '$lib/components/settings/MetadataProviderSettings.svelte';
	import MediaProfilesSettings from '$lib/components/settings/MediaProfilesSettings.svelte';
	import QualitySizeSettings from '$lib/components/settings/quality/QualitySizeSettings.svelte';
	import TagSettings from '$lib/components/settings/TagSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import type {
		CustomFormat,
		CustomFormatForm as CustomFormatFormValue,
		IntegrationTestResults,
		MediaProfile,
		MetadataProvider,
		MetadataProviderForm as MetadataProviderFormValue,
		SettingsSection,
		Tag,
		TagForm
	} from '$lib/settings/types';

	interface Props {
		activeSection: SettingsSection;
		metadataProviders: MetadataProvider[];
		mediaProfiles: MediaProfile[];
		customFormats: CustomFormat[];
		tags: Tag[];
		customFormatForm: CustomFormatFormValue;
		tagForm: TagForm;
		savingMetadataProviderId?: string;
		testingMetadataProviderId?: string;
		deletingMediaProfileId?: string;
		savingCustomFormat: boolean;
		deletingCustomFormatId?: string;
		savingTag: boolean;
		deletingTagId?: string;
		metadataProviderTests: IntegrationTestResults;
		onSaveMetadataProvider: (_form: MetadataProviderFormValue) => void | Promise<void>;
		onTestMetadataProvider: (_id: string) => void | Promise<void>;
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
		mediaProfiles,
		customFormats,
		tags,
		customFormatForm = $bindable(),
		tagForm = $bindable(),
		savingMetadataProviderId,
		testingMetadataProviderId,
		deletingMediaProfileId,
		savingCustomFormat,
		deletingCustomFormatId,
		savingTag,
		deletingTagId,
		metadataProviderTests,
		onSaveMetadataProvider,
		onTestMetadataProvider,
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
	<PageHeading eyebrow="Settings" title="Metadata" titleId="settings-title" />
	<div class="space-y-4">
		<MetadataProviderSettings
			{metadataProviders}
			onSave={onSaveMetadataProvider}
			onTest={onTestMetadataProvider}
			testingId={testingMetadataProviderId}
			savingId={savingMetadataProviderId}
			testResults={metadataProviderTests}
		/>
	</div>
{:else if activeSection === 'quality'}
	<PageHeading eyebrow="Settings" title="Quality" titleId="settings-title" />
	<div class="space-y-4"><QualitySizeSettings /></div>
{:else if activeSection === 'profiles'}
	<PageHeading eyebrow="Settings" title="Profiles" titleId="settings-title" />
	<div class="space-y-4">
		<MediaProfilesSettings
			profiles={mediaProfiles}
			deletingId={deletingMediaProfileId}
			onDelete={onDeleteMediaProfile}
		/>
	</div>
{:else if activeSection === 'custom-formats'}
	<PageHeading eyebrow="Settings" title="Custom formats" titleId="settings-title" />
	<div class="space-y-4">
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
{:else if activeSection === 'tags'}
	<PageHeading eyebrow="Settings" title="Tags" titleId="settings-title" />
	<div class="space-y-4">
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
