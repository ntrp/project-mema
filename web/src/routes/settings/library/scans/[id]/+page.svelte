<script lang="ts">
	import AppShell from '$lib/components/app/AppShell.svelte';

	interface Props {
		data: {
			id: string;
		};
	}

	let { data }: Props = $props();
</script>

<AppShell initialView="settings" initialSettingsSection="library" initialLibraryScanId={data.id} />
