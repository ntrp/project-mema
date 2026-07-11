<script lang="ts">
	import MediaRequestArea from '$lib/components/app/requests/MediaRequestArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaRequestsQuery } from '$lib/features/library/queries.svelte';

	interface Props {
		selectedRequestId?: string;
	}

	let { selectedRequestId }: Props = $props();
	const app = getAppShellContext();
	const requests = createMediaRequestsQuery();
</script>

<MediaRequestArea
	requests={requests.data ?? []}
	{selectedRequestId}
	libraryFolders={app.libraryFolders}
	qualityProfiles={app.mediaProfiles}
	tags={app.tags}
	canManage={app.isAdmin}
	approvingRequestId={app.approvingRequestId}
	onApprove={app.approveMediaRequest}
/>
