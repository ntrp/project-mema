<script lang="ts">
	import LibraryFolderForm from '$lib/components/settings/LibraryFolderForm.svelte';
	import LibraryFolderAccordion from '$lib/components/settings/LibraryFolderAccordion.svelte';
	import FileNamingSettings from '$lib/components/settings/FileNamingSettings.svelte';
	import PathMappingSettings from '$lib/components/settings/PathMappingSettings.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Card } from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';
	import { emptyLibraryFolderForm } from '$lib/settings/forms';
	import type { LibraryScanImportRow } from '$lib/components/settings/libraryScanImport';
	import type {
		LibraryFolder,
		LibraryFolderForm as LibraryFolderFormValue,
		LibraryMediaKind,
		LibraryScan,
		PathMapping,
		PathMappingForm,
		MediaSearchResult,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		form: LibraryFolderFormValue;
		pathMappings: PathMapping[];
		pathMappingForm: PathMappingForm;
		scansByFolder: Record<string, LibraryScan>;
		openFolderId?: string;
		qualityProfiles: QualityProfileOption[];
		saving: boolean;
		scanningLibraryFolderId?: string;
		savingPathMapping: boolean;
		deletingPathMappingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onScan: (_id: string) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onSavePathMapping: (_event: SubmitEvent) => void | Promise<void>;
		onDeletePathMapping: (_id: string) => void | Promise<void>;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onImport: (_scan: LibraryScan, _rows: LibraryScanImportRow[]) => Promise<void>;
	}

	let {
		folders,
		form = $bindable(),
		pathMappings,
		pathMappingForm = $bindable(),
		scansByFolder,
		openFolderId,
		qualityProfiles,
		saving,
		scanningLibraryFolderId,
		savingPathMapping,
		deletingPathMappingId,
		onSave,
		onScan,
		onDelete,
		onSavePathMapping,
		onDeletePathMapping,
		onSearchMatch,
		onImport
	}: Props = $props();

	let modalOpen = $state(false);

	function openModal() {
		form = emptyLibraryFolderForm();
		modalOpen = true;
	}

	function closeModal() {
		form = emptyLibraryFolderForm();
		modalOpen = false;
	}

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (form.path.trim() === '') {
			modalOpen = false;
		}
	}
</script>

<PageHeading eyebrow="Settings" title="Library" titleId="settings-title" />
<div class="space-y-4">
	<Card class="overflow-visible p-5" aria-label="Root paths">
		<div class="grid gap-4">
			<SectionHeading title="Root Paths" />
			<section class="grid gap-3" aria-labelledby="root-paths-paths-title">
				<div class="flex items-center justify-between gap-3">
					<h3 id="root-paths-paths-title" class="m-0 text-lg text-foreground">Paths</h3>
					<div class="inline-flex items-center gap-2 text-xs font-black text-muted-foreground">
						<SettingsAddButton label="Add library folder" onclick={openModal} />
					</div>
				</div>
				<Separator />
				<LibraryFolderAccordion
					{folders}
					{scansByFolder}
					{openFolderId}
					scanningFolderId={scanningLibraryFolderId}
					{qualityProfiles}
					{onScan}
					{onDelete}
					{onSearchMatch}
					{onImport}
				/>
			</section>
			<PathMappingSettings
				mappings={pathMappings}
				bind:form={pathMappingForm}
				saving={savingPathMapping}
				deletingId={deletingPathMappingId}
				onSave={onSavePathMapping}
				onDelete={onDeletePathMapping}
			/>
		</div>
	</Card>
	<FileNamingSettings />
	{#if modalOpen}
		<SettingsFormModal title="Add library folder" onClose={closeModal}>
			<LibraryFolderForm bind:form {saving} onSave={save} />
		</SettingsFormModal>
	{/if}
</div>
