<script lang="ts">
	import QualitySizeSettings from '$lib/components/settings/quality/QualitySizeSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Quality" titleId="settings-title" />
	<div class="space-y-4">
		<QualitySizeSettings />
	</div>
</SettingsRouteFrame>
