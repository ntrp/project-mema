<script lang="ts">
	import { onMount } from 'svelte';
	import DLNADeviceOverrideTable from './DLNADeviceOverrideTable.svelte';
	import DLNADecisionTraceViewer from './decision-trace/DLNADecisionTraceViewer.svelte';
	import DLNAProfileEditor from './DLNAProfileEditor.svelte';
	import DLNAProfileModals from './DLNAProfileModals.svelte';
	import DLNAProfileTable from './DLNAProfileTable.svelte';
	import { DLNAProfilePanelState } from './dlnaProfilePanelState.svelte';

	const vm = new DLNAProfilePanelState();

	onMount(() => {
		void vm.load();
	});
</script>

<DLNADeviceOverrideTable
	devices={vm.devices}
	overrides={vm.overrides}
	profiles={vm.profiles}
	bind:overrideForm={vm.overrideForm}
	overrideJsonText={vm.overrideJsonText}
	saving={vm.saving}
	onOverrideJson={(value) => (vm.overrideJsonText = value)}
	onSave={vm.saveOverride}
	onDelete={vm.deleteOverride}
	onQuickAssign={vm.quickAssign}
/>

<DLNAProfileTable
	profiles={vm.filteredProfiles}
	search={vm.search}
	selectedId={vm.selectedId}
	onSearch={(value) => (vm.search = value)}
	onEdit={vm.openProfileEditor}
	onClone={vm.openClone}
	onReset={vm.resetProfile}
	onExport={vm.exportProfile}
	onDelete={vm.deleteProfile}
	onCreate={vm.newProfile}
	onImport={() => (vm.importOpen = true)}
	onTrace={vm.openTrace}
	onRestoreOriginals={vm.restoreOriginalProfiles}
/>

<DLNAProfileEditor
	open={vm.editorOpen}
	mode={vm.editorMode}
	profile={vm.selectedProfile}
	bind:form={vm.form}
	saving={vm.saving}
	errorMessage={vm.errorMessage}
	onSave={vm.saveProfile}
	onClose={vm.closeEditor}
/>

<DLNADecisionTraceViewer
	open={vm.traceOpen}
	devices={vm.devices}
	profiles={vm.profiles}
	selectedIp={vm.traceIp}
	mediaPath={vm.traceMediaPath}
	onClose={vm.closeTrace}
	onSelectedIp={(value) => (vm.traceIp = value)}
	onMediaPath={(value) => (vm.traceMediaPath = value)}
/>

<DLNAProfileModals
	cloneSource={vm.cloneSource}
	cloneId={vm.cloneId}
	cloneName={vm.cloneName}
	importOpen={vm.importOpen}
	importText={vm.importText}
	saving={vm.saving}
	onCloneId={(value) => (vm.cloneId = value)}
	onCloneName={(value) => (vm.cloneName = value)}
	onImportText={(value) => (vm.importText = value)}
	onCloseClone={() => (vm.cloneSource = undefined)}
	onCloseImport={() => (vm.importOpen = false)}
	onClone={vm.cloneProfile}
	onImport={vm.importProfile}
/>
