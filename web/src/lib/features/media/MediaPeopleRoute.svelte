<script lang="ts">
	import MediaPeopleArea from '$lib/components/app/media/people/MediaPeopleArea.svelte';
	import type { PeopleSectionKind } from '$lib/components/app/shell/controller/index.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import type { MediaItem, MediaMetadataDetails } from '$lib/settings/types';
	import { createMetadataDetailQuery } from './queries.svelte';

	interface Props {
		kind: PeopleSectionKind;
	}

	let { kind }: Props = $props();
	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const metadata = createMetadataDetailQuery({
		provider: () => app.route?.metadataProvider,
		type: () => app.route?.metadataType,
		id: () => app.route?.metadataExternalId
	});
	const detail = $derived(
		metadata.data ??
			((library.data ?? []).find(
				(item: MediaItem) => item.id === app.selectedMediaItemId
			) as unknown as MediaMetadataDetails | undefined)
	);
</script>

<MediaPeopleArea {detail} {kind} loading={metadata.isFetching && !detail} />
