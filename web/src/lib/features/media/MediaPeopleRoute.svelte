<script lang="ts">
	import MediaPeopleArea from '$lib/components/app/media/people/MediaPeopleArea.svelte';
	import type { PeopleSectionKind } from '$lib/components/app/shell/controller/index.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import type { MediaItem } from '$lib/settings/types';

	interface Props {
		kind: PeopleSectionKind;
	}

	let { kind }: Props = $props();
	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const detail = $derived(
		app.metadataDetail ??
			(library.data ?? []).find((item: MediaItem) => item.id === app.selectedMediaItemId)
	);
</script>

<MediaPeopleArea {detail} {kind} loading={app.loadingMetadataDetail && !detail} />
