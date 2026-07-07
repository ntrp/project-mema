<script lang="ts">
	import DLNASettingsPanel from '$lib/components/settings/dlna/DLNASettingsPanel.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="DLNA" titleId="settings-title" />
	<DLNASettingsPanel />
</SettingsRouteFrame>
