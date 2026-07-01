<script lang="ts">
	import MediaPeopleArea from '$lib/components/app/media/MediaPeopleArea.svelte';
	import type { PeopleSectionKind } from '$lib/components/app/shell/controller/index.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		kind: PeopleSectionKind;
	}

	let { kind }: Props = $props();
	const app = getAppShellContext();
</script>

<MediaPeopleArea
	detail={app.mediaPeopleDetail}
	{kind}
	loading={app.loadingMetadataDetail && !app.mediaPeopleDetail}
/>
