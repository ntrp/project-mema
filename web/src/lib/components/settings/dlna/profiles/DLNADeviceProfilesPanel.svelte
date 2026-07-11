<script lang="ts">
	import { onMount } from 'svelte';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
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

<Card.Root aria-label="DLNA device profiles">
	<Card.Header class="border-b border-border">
		<Card.Title>Device Profiles</Card.Title>
		<Card.Action>
			<Button type="button" variant="secondary" size="sm" disabled={vm.loading} onclick={vm.load}>
				<RefreshCwIcon class={vm.loading ? 'animate-spin' : ''} />
				Refresh
			</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content class="grid gap-6 pt-5">
		{#if vm.errorMessage}
			<p class="m-0 text-sm font-medium text-destructive">{vm.errorMessage}</p>
		{/if}
		{#if vm.message}
			<p class="m-0 text-sm text-muted-foreground">{vm.message}</p>
		{/if}
		<Card.Root aria-label="DLNA renderer profiles details" size="sm">
			<Card.Content class="pt-5">
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
			</Card.Content>
		</Card.Root>
	</Card.Content>
</Card.Root>

<Card.Root aria-label="DLNA renderer profiles">
	<Card.Content class="pt-5">
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
	</Card.Content>
</Card.Root>

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
