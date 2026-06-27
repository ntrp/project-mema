<script lang="ts">
	import DownloadClientForm from '$lib/components/settings/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/DownloadClientTable.svelte';
	import IndexerForm from '$lib/components/settings/IndexerForm.svelte';
	import IndexerTable from '$lib/components/settings/IndexerTable.svelte';
	import SidebarMenu from './SidebarMenu.svelte';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults,
		SettingsSection
	} from '$lib/settings/types';

	interface Props {
		activeSection: SettingsSection;
		downloadClients: DownloadClient[];
		indexers: Indexer[];
		downloadForm: DownloadClientFormValue;
		indexerForm: IndexerFormValue;
		savingDownloadClient: boolean;
		savingIndexer: boolean;
		testingDownloadClientId?: string;
		testingIndexerId?: string;
		downloadClientTests: IntegrationTestResults;
		indexerTests: IntegrationTestResults;
		onSectionSelect: (_section: SettingsSection) => void;
		onSaveDownloadClient: (_event: SubmitEvent) => void | Promise<void>;
		onSaveIndexer: (_event: SubmitEvent) => void | Promise<void>;
		onCancelDownloadClient: () => void;
		onCancelIndexer: () => void;
		onEditDownloadClient: (_client: DownloadClient) => void;
		onEditIndexer: (_indexer: Indexer) => void;
		onDeleteDownloadClient: (_id: string) => void | Promise<void>;
		onDeleteIndexer: (_id: string) => void | Promise<void>;
		onTestDownloadClient: (_id: string) => void | Promise<void>;
		onTestIndexer: (_id: string) => void | Promise<void>;
	}

	let {
		activeSection,
		downloadClients,
		indexers,
		downloadForm = $bindable(),
		indexerForm = $bindable(),
		savingDownloadClient,
		savingIndexer,
		testingDownloadClientId,
		testingIndexerId,
		downloadClientTests,
		indexerTests,
		onSectionSelect,
		onSaveDownloadClient,
		onSaveIndexer,
		onCancelDownloadClient,
		onCancelIndexer,
		onEditDownloadClient,
		onEditIndexer,
		onDeleteDownloadClient,
		onDeleteIndexer,
		onTestDownloadClient,
		onTestIndexer
	}: Props = $props();

	const settingsItems = [
		{ value: 'download-clients', label: 'Download clients', meta: 'Torrent and NZB' },
		{ value: 'indexers', label: 'Indexers', meta: 'Torznab, Newznab, RSS' }
	] satisfies { value: SettingsSection; label: string; meta: string }[];
</script>

<div class="workspace-layout">
	<SidebarMenu
		title="Settings"
		items={settingsItems}
		active={activeSection}
		onSelect={(section) => onSectionSelect(section as SettingsSection)}
	/>

	<section class="workspace-main" aria-labelledby="settings-title">
		{#if activeSection === 'download-clients'}
			<div class="page-heading">
				<p>Settings</p>
				<h1 id="settings-title">Download clients</h1>
			</div>
			<div class="settings-stack">
				<DownloadClientForm
					bind:form={downloadForm}
					saving={savingDownloadClient}
					onSave={onSaveDownloadClient}
					onCancel={onCancelDownloadClient}
				/>
				<DownloadClientTable
					clients={downloadClients}
					onEdit={onEditDownloadClient}
					onDelete={onDeleteDownloadClient}
					onTest={onTestDownloadClient}
					testingId={testingDownloadClientId}
					testResults={downloadClientTests}
				/>
			</div>
		{:else}
			<div class="page-heading">
				<p>Settings</p>
				<h1 id="settings-title">Indexers</h1>
			</div>
			<div class="settings-stack">
				<IndexerForm
					bind:form={indexerForm}
					saving={savingIndexer}
					onSave={onSaveIndexer}
					onCancel={onCancelIndexer}
				/>
				<IndexerTable
					{indexers}
					onEdit={onEditIndexer}
					onDelete={onDeleteIndexer}
					onTest={onTestIndexer}
					testingId={testingIndexerId}
					testResults={indexerTests}
				/>
			</div>
		{/if}
	</section>
</div>
