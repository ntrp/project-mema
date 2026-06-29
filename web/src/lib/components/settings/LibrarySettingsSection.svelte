<script lang="ts">
	import LibraryFolderForm from '$lib/components/settings/LibraryFolderForm.svelte';
	import LibraryFolderTable from '$lib/components/settings/LibraryFolderTable.svelte';
	import LibraryScanReview from '$lib/components/settings/LibraryScanReview.svelte';
	import PathMappingSettings from '$lib/components/settings/PathMappingSettings.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { emptyLibraryFolderForm } from '$lib/settings/forms';
	import type {
		LibraryFolder,
		LibraryFolderForm as LibraryFolderFormValue,
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		LibraryScanItemMatchRequest,
		PathMapping,
		PathMappingForm,
		MediaSearchResult
	} from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		form: LibraryFolderFormValue;
		pathMappings: PathMapping[];
		pathMappingForm: PathMappingForm;
		scan?: LibraryScan;
		saving: boolean;
		savingPathMapping: boolean;
		deletingPathMappingId?: string;
		loadingScan: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onSavePathMapping: (_event: SubmitEvent) => void | Promise<void>;
		onDeletePathMapping: (_id: string) => void | Promise<void>;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onMatch: (_item: LibraryScanItem, _request: LibraryScanItemMatchRequest) => Promise<void>;
	}

	let {
		folders,
		form = $bindable(),
		pathMappings,
		pathMappingForm = $bindable(),
		scan,
		saving,
		savingPathMapping,
		deletingPathMappingId,
		loadingScan,
		onSave,
		onDelete,
		onSavePathMapping,
		onDeletePathMapping,
		onSearchMatch,
		onMatch
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

<div class="page-heading">
	<p>Settings</p>
	<h1 id="settings-title">Library</h1>
</div>
<div class="settings-stack">
	<div class="settings-toolbar">
		<button type="button" onclick={openModal}>Add library folder</button>
	</div>
	<LibraryFolderTable {folders} {onDelete} />
	<PathMappingSettings
		mappings={pathMappings}
		bind:form={pathMappingForm}
		saving={savingPathMapping}
		deletingId={deletingPathMappingId}
		onSave={onSavePathMapping}
		onDelete={onDeletePathMapping}
	/>
	<LibraryScanReview {scan} loading={loadingScan} {onSearchMatch} {onMatch} />
	{#if modalOpen}
		<SettingsFormModal title="Add library folder" onClose={closeModal}>
			<LibraryFolderForm bind:form {saving} onSave={save} />
		</SettingsFormModal>
	{/if}
</div>
