<script lang="ts">
	import AdvancedSearchArea from '$lib/components/app/search/AdvancedSearchArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		initialQuery: string;
	}

	let { initialQuery }: Props = $props();
	const app = getAppShellContext();
</script>

<AdvancedSearchArea
	{initialQuery}
	metadataProviders={app.metadataProviders}
	groups={app.advancedSearchGroups}
	searching={app.searchingAdvanced}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onSearch={app.advancedSearch}
	onAdd={app.addMedia}
/>
